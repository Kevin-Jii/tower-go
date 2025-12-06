package service

import (
	"time"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
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
}

func NewInventoryService(inventoryModule *module.InventoryModule, userModule *module.UserModule, storeModule *module.StoreModule, productModule *module.SupplierProductModule) *InventoryService {
	return &InventoryService{
		inventoryModule: inventoryModule,
		userModule:      userModule,
		storeModule:     storeModule,
		productModule:   productModule,
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

	return order, nil
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
