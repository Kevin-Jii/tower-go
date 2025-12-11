package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
	"github.com/Kevin-Jii/tower-go/utils/logging"
)

// parseDate 解析日期字符串
func parseDate(s string) (time.Time, error) {
	return time.Parse("2006-01-02", s)
}

type InventoryService struct {
	inventoryModule *module.InventoryModule
	userModule      *module.UserModule
	storeModule     *module.StoreModule
	productModule   *module.SupplierProductModule
	dingTalkService *DingTalkService
	botModule       *module.DingTalkBotModule
	templateService *MessageTemplateService
}

func NewInventoryService(
	inventoryModule *module.InventoryModule,
	userModule *module.UserModule,
	storeModule *module.StoreModule,
	productModule *module.SupplierProductModule,
	dingTalkService *DingTalkService,
	botModule *module.DingTalkBotModule,
	templateService *MessageTemplateService,
) *InventoryService {
	return &InventoryService{
		inventoryModule: inventoryModule,
		userModule:      userModule,
		storeModule:     storeModule,
		productModule:   productModule,
		dingTalkService: dingTalkService,
		botModule:       botModule,
		templateService: templateService,
	}
}

// GetInventory 获取库存
func (s *InventoryService) GetInventory(storeID, productID uint) (*model.Inventory, error) {
	return s.inventoryModule.GetByStoreAndProduct(storeID, productID)
}

// ListInventory 库存列表
func (s *InventoryService) ListInventory(req *model.ListInventoryReq) ([]*model.InventoryWithProduct, int64, error) {
	return s.inventoryModule.List(req)
}

// CreateOrder 创建出入库单
func (s *InventoryService) CreateOrder(storeID, operatorID uint, req *model.CreateInventoryOrderReq) (*model.InventoryOrder, error) {
	// 出库时校验库存
	if req.Type == model.InventoryTypeOut {
		for _, item := range req.Items {
			inv, err := s.inventoryModule.GetByStoreAndProduct(storeID, item.ProductID)
			if err != nil {
				// 获取商品名称用于错误提示
				productName := "未知商品"
				if product, _ := s.productModule.GetByID(item.ProductID); product != nil {
					productName = product.Name
				}
				return nil, fmt.Errorf("商品【%s】不在库存中，无法出库", productName)
			}
			if inv.Quantity < item.Quantity {
				productName := "未知商品"
				if product, _ := s.productModule.GetByID(item.ProductID); product != nil {
					productName = product.Name
				}
				return nil, fmt.Errorf("商品【%s】库存不足，当前库存: %.2f，出库数量: %.2f", productName, inv.Quantity, item.Quantity)
			}
		}
	}

	// 生成单据编号
	orderNo := s.inventoryModule.GenerateOrderNo(req.Type)

	// 获取门店信息
	storeName := ""
	if store, err := s.storeModule.GetByID(storeID); err == nil && store != nil {
		storeName = store.Name
	}

	// 获取操作人信息
	operatorName := ""
	operatorPhone := ""
	if user, err := s.userModule.GetByID(operatorID); err == nil && user != nil {
		operatorName = user.Nickname
		if operatorName == "" {
			operatorName = user.Username
		}
		operatorPhone = user.Phone
	}

	// 构建明细
	var items []model.InventoryOrderItem
	var totalQuantity float64

	for _, item := range req.Items {
		// 获取商品信息
		productName := ""
		unit := ""
		if product, err := s.productModule.GetByID(item.ProductID); err == nil && product != nil {
			productName = product.Name
			unit = product.Unit
		}

		orderItem := model.InventoryOrderItem{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Unit:        unit,
			Remark:      item.Remark,
		}

		// 入库时解析生产日期和截止日期
		if req.Type == model.InventoryTypeIn {
			if item.ProductionDate != "" {
				if t, err := parseDate(item.ProductionDate); err == nil {
					orderItem.ProductionDate = &t
				}
			}
			if item.ExpiryDate != "" {
				if t, err := parseDate(item.ExpiryDate); err == nil {
					orderItem.ExpiryDate = &t
				}
			}
		}

		items = append(items, orderItem)
		totalQuantity += item.Quantity
	}

	// 创建出入库单
	order := &model.InventoryOrder{
		OrderNo:       orderNo,
		Type:          req.Type,
		StoreID:       storeID,
		StoreName:     storeName,
		Reason:        req.Reason,
		Remark:        req.Remark,
		TotalQuantity: totalQuantity,
		ItemCount:     len(items),
		OperatorID:    operatorID,
		OperatorName:  operatorName,
		OperatorPhone: operatorPhone,
		Items:         items,
	}

	if err := s.inventoryModule.CreateOrder(order); err != nil {
		return nil, err
	}

	// 更新库存
	for _, item := range req.Items {
		unit := ""
		if product, err := s.productModule.GetByID(item.ProductID); err == nil && product != nil {
			unit = product.Unit
		}

		if req.Type == model.InventoryTypeIn {
			if err := s.inventoryModule.AddQuantity(storeID, item.ProductID, item.Quantity, unit); err != nil {
				return nil, err
			}
		} else {
			if err := s.inventoryModule.SubQuantity(storeID, item.ProductID, item.Quantity); err != nil {
				return nil, err
			}
		}
	}

	// 异步发送钉钉通知（仅入库）
	if req.Type == model.InventoryTypeIn {
		go s.sendDingTalkNotification(order, storeID)
	}

	return order, nil
}

// sendDingTalkNotification 发送入库通知到门店负责人
func (s *InventoryService) sendDingTalkNotification(order *model.InventoryOrder, storeID uint) {
	if s.dingTalkService == nil || s.storeModule == nil || s.botModule == nil {
		return
	}

	// 获取门店信息
	store, err := s.storeModule.GetByID(storeID)
	if err != nil || store == nil {
		if logging.SugaredLogger != nil {
			logging.SugaredLogger.Warnw("Failed to get store for notification", "storeID", storeID, "error", err)
		}
		return
	}

	// 检查门店是否有联系电话
	if store.Phone == "" {
		if logging.SugaredLogger != nil {
			logging.SugaredLogger.Warnw("Store has no phone, skip notification", "storeID", storeID)
		}
		return
	}

	// 获取门店绑定的机器人
	bot, err := s.botModule.GetByStoreID(storeID)
	if err != nil || bot == nil {
		if logging.SugaredLogger != nil {
			logging.SugaredLogger.Warnw("No bot found for store", "storeID", storeID, "error", err)
		}
		return
	}

	if !bot.IsEnabled || bot.BotType != "stream" {
		return
	}

	// 构建商品明细
	var itemLines []string
	for i, item := range order.Items {
		line := fmt.Sprintf("%d. %s x%.2f%s", i+1, item.ProductName, item.Quantity, item.Unit)
		itemLines = append(itemLines, line)
	}

	// 入库类型显示
	orderType := order.Reason
	if orderType == "" {
		orderType = "入库"
	}

	var title, text string

	// 尝试使用模板
	if s.templateService != nil {
		data := map[string]interface{}{
			"StoreName":    store.Name,
			"OrderNo":      order.OrderNo,
			"OrderType":    orderType,
			"OrderDate":    order.CreatedAt.Format("2006-01-02"),
			"OperatorName": order.OperatorName,
			"ItemList":     strings.Join(itemLines, "\n\n"),
			"TotalAmount":  fmt.Sprintf("%.2f", order.TotalQuantity),
			"ItemCount":    order.ItemCount,
			"CreateTime":   time.Now().Format("2006-01-02 15:04:05"),
		}
		var err error
		title, text, err = s.templateService.RenderTemplate(model.TemplateInventoryCreated, data)
		if err != nil {
			if logging.SugaredLogger != nil {
				logging.SugaredLogger.Warnw("Failed to render template, using default", "error", err)
			}
		}
	}

	// 如果模板渲染失败，使用默认格式
	if text == "" {
		title = fmt.Sprintf("📦 新入库通知 - %s", store.Name)
		text = fmt.Sprintf("## %s\n\n"+
			"**入库单号：** %s\n\n"+
			"**入库类型：** %s\n\n"+
			"**入库日期：** %s\n\n"+
			"**操作人：** %s\n\n"+
			"### 入库明细\n\n"+
			"%s\n\n"+
			"**总数量：** %.2f\n\n"+
			"**商品种类：** %d 项\n\n"+
			"---\n\n"+
			"%s",
			title,
			order.OrderNo,
			orderType,
			order.CreatedAt.Format("2006-01-02"),
			order.OperatorName,
			strings.Join(itemLines, "\n\n"),
			order.TotalQuantity,
			order.ItemCount,
			time.Now().Format("2006-01-02 15:04:05"),
		)
	}

	// 发送通知
	if err := s.dingTalkService.SendStreamMarkdownToMobile(bot, title, text, store.Phone); err != nil {
		if logging.SugaredLogger != nil {
			logging.SugaredLogger.Errorw("Failed to send inventory notification",
				"storeID", storeID,
				"orderNo", order.OrderNo,
				"error", err,
			)
		}
	} else {
		if logging.SugaredLogger != nil {
			logging.SugaredLogger.Infow("Inventory notification sent",
				"storeID", storeID,
				"orderNo", order.OrderNo,
				"mobile", store.Phone,
			)
		}
	}
}

// GetOrderByNo 根据单号获取出入库单详情
func (s *InventoryService) GetOrderByNo(orderNo string) (*model.InventoryOrder, error) {
	return s.inventoryModule.GetOrderByNo(orderNo)
}

// GetOrderByID 根据ID获取出入库单详情
func (s *InventoryService) GetOrderByID(id uint) (*model.InventoryOrder, error) {
	return s.inventoryModule.GetOrderByID(id)
}

// ListOrders 出入库单列表
func (s *InventoryService) ListOrders(req *model.ListInventoryOrderReq) ([]*model.InventoryOrder, int64, error) {
	return s.inventoryModule.ListOrders(req)
}

// UpdateInventory 修改库存数量
func (s *InventoryService) UpdateInventory(id uint, quantity float64) error {
	return s.inventoryModule.UpdateQuantity(id, quantity)
}

// GetInventoryByID 根据ID获取库存
func (s *InventoryService) GetInventoryByID(id uint) (*model.Inventory, error) {
	return s.inventoryModule.GetByID(id)
}
