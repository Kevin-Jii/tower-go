package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
	"github.com/Kevin-Jii/tower-go/pkg/statemachine"
	"github.com/Kevin-Jii/tower-go/utils"
	"github.com/Kevin-Jii/tower-go/utils/logging"
)

type PurchaseOrderService struct {
	orderModule         *module.PurchaseOrderModule
	productModule       *module.SupplierProductModule
	storeSupplierModule *module.StoreSupplierModule
	storeModule         *module.StoreModule
	botModule           *module.DingTalkBotModule
	dingTalkService     *DingTalkService
	stateMachine        *statemachine.StateMachine
}

func NewPurchaseOrderService(
	orderModule *module.PurchaseOrderModule,
	productModule *module.SupplierProductModule,
	storeSupplierModule *module.StoreSupplierModule,
	storeModule *module.StoreModule,
	botModule *module.DingTalkBotModule,
	dingTalkService *DingTalkService,
) *PurchaseOrderService {
	// 创建状态机并注册钩子
	sm := statemachine.NewOrderStateMachine()

	// 确认订单时发布事件
	sm.OnAction(statemachine.ActionConfirm, func(from, to statemachine.State) error {
		utils.GlobalEventBus.Publish(utils.EventOrderConfirmed, map[string]interface{}{
			"from": from,
			"to":   to,
		})
		return nil
	})

	// 完成订单时发布事件
	sm.OnAction(statemachine.ActionComplete, func(from, to statemachine.State) error {
		utils.GlobalEventBus.Publish(utils.EventOrderCompleted, map[string]interface{}{
			"from": from,
			"to":   to,
		})
		return nil
	})

	// 取消订单时发布事件
	sm.OnAction(statemachine.ActionCancel, func(from, to statemachine.State) error {
		utils.GlobalEventBus.Publish(utils.EventOrderCancelled, map[string]interface{}{
			"from": from,
			"to":   to,
		})
		return nil
	})

	return &PurchaseOrderService{
		orderModule:         orderModule,
		productModule:       productModule,
		storeSupplierModule: storeSupplierModule,
		storeModule:         storeModule,
		botModule:           botModule,
		dingTalkService:     dingTalkService,
		stateMachine:        sm,
	}
}

// CreateOrder 创建采购单
func (s *PurchaseOrderService) CreateOrder(storeID, userID uint, req *model.CreatePurchaseOrderReq) (*model.PurchaseOrder, error) {
	// 解析报菜日期：前端不传时自动取当天
	now := time.Now()
	orderDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	if strings.TrimSpace(req.OrderDate) != "" {
		parsedDate, err := time.Parse("2006-01-02", req.OrderDate)
		if err != nil {
			return nil, errors.New("invalid order date format")
		}
		orderDate = parsedDate
	}

	// 获取商品信息
	var productIDs []uint
	for _, item := range req.Items {
		productIDs = append(productIDs, item.ProductID)
	}

	// 验证商品是否已绑定到门店
	unboundIDs, err := s.storeSupplierModule.ValidateStoreProducts(storeID, productIDs)
	if err != nil {
		return nil, err
	}
	if len(unboundIDs) > 0 {
		return nil, errors.New("some products are not bound to the store")
	}

	products, err := s.productModule.GetByIDs(productIDs)
	if err != nil {
		return nil, err
	}

	// 构建商品ID到商品的映射
	productMap := make(map[uint]*model.SupplierProduct)
	for _, p := range products {
		productMap[p.ID] = p
	}

	// 创建采购单
	order := &model.PurchaseOrder{
		OrderNo:   s.orderModule.GenerateOrderNo(),
		StoreID:   storeID,
		Status:    model.PurchaseStatusPending,
		Remark:    req.Remark,
		OrderDate: orderDate,
		CreatedBy: userID,
	}

	if err := s.orderModule.Create(order); err != nil {
		return nil, err
	}

	// 创建明细
	var totalAmount float64
	var items []model.PurchaseOrderItem
	for _, reqItem := range req.Items {
		product, ok := productMap[reqItem.ProductID]
		if !ok {
			continue
		}

		amount := product.Price * reqItem.Quantity
		totalAmount += amount

		items = append(items, model.PurchaseOrderItem{
			OrderID:    order.ID,
			SupplierID: product.SupplierID,
			ProductID:  reqItem.ProductID,
			Quantity:   reqItem.Quantity,
			UnitPrice:  product.Price,
			Amount:     amount,
			Remark:     reqItem.Remark,
		})
	}

	if len(items) > 0 {
		if err := s.orderModule.CreateItems(items); err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("no valid purchase items")
	}

	// 更新总金额
	s.orderModule.GetDB().Model(&model.PurchaseOrder{}).Where("id = ?", order.ID).Update("total_amount", totalAmount)
	order.TotalAmount = totalAmount

	// 发布订单创建事件
	utils.GlobalEventBus.Publish(utils.EventOrderCreated, order)

	// 异步推送钉钉通知
	go s.sendDingTalkNotification(order)

	return order, nil
}

// GetOrder 获取采购单详情（包含关联数据）
func (s *PurchaseOrderService) GetOrder(id uint) (*model.PurchaseOrder, error) {
	return s.orderModule.GetByIDWithDetails(id)
}

// ListOrders 获取采购单列表（ctx 须含 AuthContext，见 middleware.AttachAuthContextToHTTPRequest）
func (s *PurchaseOrderService) ListOrders(ctx context.Context, req *model.ListPurchaseOrderReq) ([]*model.PurchaseOrder, int64, error) {
	applyListRBACFromContextToPurchaseOrder(ctx, req)
	return s.orderModule.List(req)
}

// UpdateOrder 更新采购单
func (s *PurchaseOrderService) UpdateOrder(id uint, req *model.UpdatePurchaseOrderReq) error {
	order, err := s.orderModule.GetByID(id)
	if err != nil {
		return errors.New("order not found")
	}

	// 使用状态机验证状态转换
	if req.Status != nil {
		action := s.getActionForStatus(order.Status, *req.Status)
		if action == "" {
			return errors.New("invalid status transition")
		}

		// 执行状态机转换（会触发钩子和事件）
		_, err := s.stateMachine.Execute(statemachine.State(order.Status), action)
		if err != nil {
			return err
		}
	}

	return s.orderModule.UpdateByID(id, req)
}

// getActionForStatus 根据目标状态获取对应的动作
func (s *PurchaseOrderService) getActionForStatus(currentStatus, newStatus int8) statemachine.Action {
	switch newStatus {
	case model.PurchaseStatusConfirmed:
		return statemachine.ActionConfirm
	case model.PurchaseStatusCompleted:
		return statemachine.ActionComplete
	case model.PurchaseStatusCancelled:
		return statemachine.ActionCancel
	default:
		return ""
	}
}

// GetAvailableActions 获取订单可用的操作
func (s *PurchaseOrderService) GetAvailableActions(id uint) ([]string, error) {
	order, err := s.orderModule.GetByID(id)
	if err != nil {
		return nil, errors.New("order not found")
	}

	actions := s.stateMachine.GetAvailableActions(statemachine.State(order.Status))
	result := make([]string, len(actions))
	for i, a := range actions {
		result[i] = string(a)
	}
	return result, nil
}

// ConfirmOrder 确认采购单（便捷方法）
func (s *PurchaseOrderService) ConfirmOrder(id uint) error {
	status := model.PurchaseStatusConfirmed
	return s.UpdateOrder(id, &model.UpdatePurchaseOrderReq{Status: &status})
}

// CompleteOrder 完成采购单（便捷方法）
func (s *PurchaseOrderService) CompleteOrder(id uint) error {
	status := model.PurchaseStatusCompleted
	return s.UpdateOrder(id, &model.UpdatePurchaseOrderReq{Status: &status})
}

// CancelOrder 取消采购单（便捷方法）
func (s *PurchaseOrderService) CancelOrder(id uint, reason string) error {
	status := model.PurchaseStatusCancelled
	return s.UpdateOrder(id, &model.UpdatePurchaseOrderReq{Status: &status, Remark: reason})
}

// DeleteOrder 删除采购单
func (s *PurchaseOrderService) DeleteOrder(id uint) error {
	order, err := s.orderModule.GetByID(id)
	if err != nil {
		return errors.New("order not found")
	}
	if order.Status != model.PurchaseStatusPending && order.Status != model.PurchaseStatusCancelled {
		return errors.New("order cannot be deleted: only pending or cancelled orders can be deleted, current status is " + getStatusName(order.Status))
	}
	return s.orderModule.Delete(id)
}

// getStatusName 获取状态名称
func getStatusName(status int8) string {
	return statemachine.GetStateName(statemachine.State(status))
}

// GetOrdersBySupplier 按供应商分组获取采购单明细
func (s *PurchaseOrderService) GetOrdersBySupplier(orderID uint) ([]model.SupplierGroupedItems, error) {
	itemsMap, err := s.orderModule.GetOrdersBySupplier(orderID)
	if err != nil {
		return nil, err
	}

	var result []model.SupplierGroupedItems
	for supplierID, items := range itemsMap {
		var supplierName string
		var groupedItems []model.SupplierGroupedItemDetail
		var subTotal float64

		for _, item := range items {
			if item.Supplier != nil {
				supplierName = item.Supplier.SupplierName
			}

			productName := ""
			unit := ""
			if item.Product != nil {
				productName = item.Product.Name
				unit = item.Product.Unit
			}

			groupedItems = append(groupedItems, model.SupplierGroupedItemDetail{
				ID:          item.ID,
				ProductID:   item.ProductID,
				ProductName: productName,
				Unit:        unit,
				Quantity:    item.Quantity,
				UnitPrice:   item.UnitPrice,
				Amount:      item.Amount,
				Remark:      item.Remark,
			})
			subTotal += item.Amount
		}

		result = append(result, model.SupplierGroupedItems{
			SupplierID:   supplierID,
			SupplierName: supplierName,
			Items:        groupedItems,
			SubTotal:     subTotal,
		})
	}

	return result, nil
}

// sendDingTalkNotification 采购单创建后异步推送钉钉通知
func (s *PurchaseOrderService) sendDingTalkNotification(order *model.PurchaseOrder) {
	if s.dingTalkService == nil || s.storeModule == nil || s.botModule == nil {
		return
	}

	// 获取门店信息
	store, err := s.storeModule.GetByID(order.StoreID)
	if err != nil || store == nil {
		logging.LogWarn(fmt.Sprintf("[PurchaseOrder] 推送钉钉失败，获取门店信息错误: storeID=%d, err=%v", order.StoreID, err))
		return
	}

	if store.Phone == "" {
		logging.LogWarn(fmt.Sprintf("[PurchaseOrder] 门店无联系电话，跳过推送: storeID=%d", order.StoreID))
		return
	}

	// 获取门店绑定的机器人
	bot, err := s.botModule.GetByStoreID(order.StoreID)
	if err != nil || bot == nil {
		logging.LogWarn(fmt.Sprintf("[PurchaseOrder] 未找到门店绑定的机器人: storeID=%d, err=%v", order.StoreID, err))
		return
	}

	if !bot.IsEnabled || bot.BotType != "stream" {
		return
	}

	// 获取完整采购单（含明细）
	fullOrder, err := s.orderModule.GetByIDWithDetails(order.ID)
	if err != nil || fullOrder == nil {
		fullOrder = order
	}

	// 构建商品明细（按供应商分组）
	supplierItems := make(map[uint][]string)
	supplierNames := make(map[uint]string)
	for _, item := range fullOrder.Items {
		supplierName := "未知供应商"
		if item.Supplier != nil {
			supplierName = item.Supplier.SupplierName
		}
		productName := "未知商品"
		unit := ""
		if item.Product != nil {
			productName = item.Product.Name
			unit = item.Product.Unit
		}
		qty := strconv.FormatFloat(item.Quantity, 'f', -1, 64)
		line := fmt.Sprintf("- %s &nbsp; %s%s &nbsp; ¥%.2f", productName, qty, unit, item.UnitPrice)
		supplierItems[item.SupplierID] = append(supplierItems[item.SupplierID], line)
		supplierNames[item.SupplierID] = supplierName
	}

	var supplierBlocks []string
	for sid, lines := range supplierItems {
		block := fmt.Sprintf("**【%s】**\n\n%s", supplierNames[sid], strings.Join(lines, "\n\n"))
		supplierBlocks = append(supplierBlocks, block)
	}

	creatorName := ""
	if fullOrder.Creator != nil {
		creatorName = fullOrder.Creator.Username
	}

	title := fmt.Sprintf("📋 新采购单 - %s", store.Name)
	text := fmt.Sprintf("## %s\n\n"+
		"**单号：** %s\n\n"+
		"**门店：** %s\n\n"+
		"**报菜日期：** %s\n\n"+
		"**制单人：** %s\n\n"+
		"### 采购明细\n\n"+
		"%s\n\n"+
		"---\n\n"+
		"**合计：** ¥%.2f\n\n"+
		"%s",
		title,
		fullOrder.OrderNo,
		store.Name,
		fullOrder.OrderDate.Format("2006-01-02"),
		creatorName,
		strings.Join(supplierBlocks, "\n\n"),
		fullOrder.TotalAmount,
		time.Now().Format("2006-01-02 15:04:05"),
	)

	if err := s.dingTalkService.SendStreamMarkdownToMobile(bot, title, text, store.Phone); err != nil {
		logging.LogWarn(fmt.Sprintf("[PurchaseOrder] 钉钉推送失败: orderNo=%s, err=%v", order.OrderNo, err))
	} else {
		logging.LogInfo(fmt.Sprintf("[PurchaseOrder] 钉钉推送成功: orderNo=%s, store=%s", order.OrderNo, store.Name))
	}
}
