package repository

import "github.com/Kevin-Jii/tower-go/model"

// StoreSupplierRepository 门店供应商仓储接口
type StoreSupplierRepository interface {
	BindSuppliers(storeID uint, supplierIDs []uint) error
	UnbindSuppliers(storeID uint, supplierIDs []uint) error
	ListSuppliersByStoreID(storeID uint) ([]*model.StoreSupplier, error)
	ListProductsByStoreID(storeID, supplierID, categoryID uint, keyword string) ([]*model.SupplierProduct, error)
}

// PurchaseOrderRepository 采购单仓储接口
type PurchaseOrderRepository interface {
	Create(order *model.PurchaseOrder) error
	GetByID(id uint) (*model.PurchaseOrder, error)
	List(req *model.ListPurchaseOrderReq) ([]*model.PurchaseOrder, int64, error)
	UpdateByID(id uint, req *model.UpdatePurchaseOrderReq) error
	Delete(id uint) error
}

// SupplierRepository 供应商仓储接口
type SupplierRepository interface {
	Create(supplier *model.Supplier) error
	GetByID(id uint) (*model.Supplier, error)
	List(page, pageSize int, keyword string) ([]*model.Supplier, int64, error)
	Update(id uint, supplier *model.Supplier) error
	Delete(id uint) error
}

// SupplierProductRepository 供应商商品仓储接口
type SupplierProductRepository interface {
	Create(product *model.SupplierProduct) error
	GetByID(id uint) (*model.SupplierProduct, error)
	GetByIDs(ids []uint) ([]*model.SupplierProduct, error)
	List(req interface{}) ([]*model.SupplierProduct, int64, error)
	Update(id uint, product *model.SupplierProduct) error
	Delete(id uint) error
}
