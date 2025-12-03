package service

import (
	"errors"
	"time"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
	"github.com/Kevin-Jii/tower-go/pkg/statemachine"
	"github.com/Kevin-Jii/tower-go/utils"
)

type PurchaseOrderService struct {
	orderModule         *module.PurchaseOrderModule
	productModule       *module.SupplierProductModule
	storeSupplierModule *module.StoreSupplierModule
	stateMachine        *statemachine.StateMachine
}

func NewPurchaseOrderService(
	orderModule *module.PurchaseOrderModule,
	productModule *module.SupplierProductModule,
	storeSupplierModule *module.StoreSupplierModule,
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
		stateMachine:        sm,
	}
}

// CreateOrder 创建采购单
func (s *PurchaseOrderService) CreateOrder(storeID, userID uint, req *model.CreatePurchaseOrderReq) (*model.PurchaseOrder, error) {
	// 解析日期
	orderDate, err := time.Parse("2006-01-02", req.OrderDate)
	if err != nil {
		return nil, errors.New("invalid order date format")
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
	}

	// 更新总金额
	s.orderModule.GetDB().Model(&model.PurchaseOrder{}).Where("id = ?", order.ID).Update("total_amount", totalAmount)
	order.TotalAmount = totalAmount

	// 发布订单创建事件
	utils.GlobalEventBus.Publish(utils.EventOrderCreated, order)

	return order, nil
}

// GetOrder 获取采购单详情
func (s *PurchaseOrderService) GetOrder(id uint) (*model.PurchaseOrder, error) {
	return s.orderModule.GetByID(id)
}

// ListOrders 获取采购单列表
func (s *PurchaseOrderService) ListOrders(req *model.ListPurchaseOrderReq) ([]*model.PurchaseOrder, int64, error) {
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
