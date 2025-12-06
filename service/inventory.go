package service

import (
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
)

type InventoryService struct {
	inventoryModule *module.InventoryModule
}

func NewInventoryService(inventoryModule *module.InventoryModule) *InventoryService {
	return &InventoryService{inventoryModule: inventoryModule}
}

// GetInventory 获取库存
func (s *InventoryService) GetInventory(storeID, productID uint) (*model.Inventory, error) {
	return s.inventoryModule.GetByStoreAndProduct(storeID, productID)
}

// ListInventory 库存列表
func (s *InventoryService) ListInventory(req *model.ListInventoryReq) ([]*model.InventoryWithProduct, int64, error) {
	return s.inventoryModule.List(req)
}

// CreateRecord 创建出入库记录
func (s *InventoryService) CreateRecord(storeID, operatorID uint, req *model.CreateInventoryRecordReq) error {
	// 生成单据编号
	recordNo := s.inventoryModule.GenerateRecordNo(req.Type)

	// 获取商品单位（这里简化处理，实际可以从商品表获取）
	unit := ""

	// 创建记录
	record := &model.InventoryRecord{
		RecordNo:   recordNo,
		StoreID:    storeID,
		ProductID:  req.ProductID,
		Type:       req.Type,
		Quantity:   req.Quantity,
		Unit:       unit,
		Reason:     req.Reason,
		Remark:     req.Remark,
		OperatorID: operatorID,
	}

	if err := s.inventoryModule.CreateRecord(record); err != nil {
		return err
	}

	// 更新库存
	if req.Type == model.InventoryTypeIn {
		return s.inventoryModule.AddQuantity(storeID, req.ProductID, req.Quantity, unit)
	}
	return s.inventoryModule.SubQuantity(storeID, req.ProductID, req.Quantity)
}

// ListRecords 出入库记录列表
func (s *InventoryService) ListRecords(req *model.ListInventoryRecordReq) ([]*model.InventoryRecord, int64, error) {
	return s.inventoryModule.ListRecords(req)
}
