package module

import (
	"github.com/Kevin-Jii/tower-go/model"
	updatesPkg "github.com/Kevin-Jii/tower-go/utils/updates"
	"gorm.io/gorm"
)

type SupplierProductModule struct {
	db *gorm.DB
}

func NewSupplierProductModule(db *gorm.DB) *SupplierProductModule {
	return &SupplierProductModule{db: db}
}

func (m *SupplierProductModule) Create(product *model.SupplierProduct) error {
	return m.db.Create(product).Error
}

func (m *SupplierProductModule) GetByID(id uint) (*model.SupplierProduct, error) {
	var product model.SupplierProduct
	if err := m.db.Preload("Supplier").Preload("Category").First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (m *SupplierProductModule) List(req *model.ListSupplierProductReq) ([]*model.SupplierProduct, int64, error) {
	var products []*model.SupplierProduct
	var total int64

	query := m.db.Model(&model.SupplierProduct{})

	if req.SupplierID > 0 {
		query = query.Where("supplier_id = ?", req.SupplierID)
	}
	if req.CategoryID > 0 {
		query = query.Where("category_id = ?", req.CategoryID)
	}
	if req.Keyword != "" {
		keyword := "%" + req.Keyword + "%"
		query = query.Where("name LIKE ?", keyword)
	}
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (req.Page - 1) * req.PageSize
	if err := query.Preload("Supplier").Preload("Category").Order("id DESC").Offset(offset).Limit(req.PageSize).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (m *SupplierProductModule) UpdateByID(id uint, req *model.UpdateSupplierProductReq) error {
	updateMap := updatesPkg.BuildUpdatesFromReq(req)
	if len(updateMap) == 0 {
		return nil
	}
	return m.db.Model(&model.SupplierProduct{}).Where("id = ?", id).Updates(updateMap).Error
}

func (m *SupplierProductModule) Delete(id uint) error {
	return m.db.Delete(&model.SupplierProduct{}, id).Error
}

// GetByIDs 批量获取商品
func (m *SupplierProductModule) GetByIDs(ids []uint) ([]*model.SupplierProduct, error) {
	var products []*model.SupplierProduct
	if err := m.db.Preload("Supplier").Preload("Category").Where("id IN ?", ids).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}
