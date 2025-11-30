package service

import (
	"errors"
	"time"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
)

type PurchaseOrderService struct {
	orderModule         *module.PurchaseOrderModule
	productModule       *module.SupplierProductModule
	storeSupplierModule *module.StoreSupplierModule
}

func NewPurchaseOrderService(
	orderModule *module.PurchaseOrderModule,
	productModule *module.SupplierProductModule,
	storeSupplierModule *module.StoreSupplierModule,
) *PurchaseOrderService {
	return &PurchaseOrderService{
		orderModule:         orderModule,
		productModule:       productModule,
		storeSupplierModule: storeSupplierModule,
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

	// 验证状态转换
	if req.Status != nil {
		if !isValidStatusTransition(order.Status, *req.Status) {
			return errors.New("invalid status transition")
		}
	}

	return s.orderModule.UpdateByID(id, req)
}

// isValidStatusTransition 验证状态转换是否有效
// 有效的状态转换：
// 待确认(1) -> 已确认(2), 已取消(4)
// 已确认(2) -> 已完成(3), 已取消(4)
// 已完成(3) -> 无法转换
// 已取消(4) -> 无法转换
func isValidStatusTransition(currentStatus, newStatus int8) bool {
	validTransitions := map[int8][]int8{
		model.PurchaseStatusPending:   {model.PurchaseStatusConfirmed, model.PurchaseStatusCancelled},
		model.PurchaseStatusConfirmed: {model.PurchaseStatusCompleted, model.PurchaseStatusCancelled},
		model.PurchaseStatusCompleted: {},
		model.PurchaseStatusCancelled: {},
	}

	allowedStatuses, ok := validTransitions[currentStatus]
	if !ok {
		return false
	}

	for _, allowed := range allowedStatuses {
		if newStatus == allowed {
			return true
		}
	}
	return false
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
	switch status {
	case model.PurchaseStatusPending:
		return "pending"
	case model.PurchaseStatusConfirmed:
		return "confirmed"
	case model.PurchaseStatusCompleted:
		return "completed"
	case model.PurchaseStatusCancelled:
		return "cancelled"
	default:
		return "unknown"
	}
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
