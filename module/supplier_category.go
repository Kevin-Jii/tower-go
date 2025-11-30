package module

import (
	"github.com/Kevin-Jii/tower-go/model"
	updatesPkg "github.com/Kevin-Jii/tower-go/utils/updates"
	"gorm.io/gorm"
)

type SupplierCategoryModule struct {
	db *gorm.DB
}

func NewSupplierCategoryModule(db *gorm.DB) *SupplierCategoryModule {
	return &SupplierCategoryModule{db: db}
}

func (m *SupplierCategoryModule) Create(category *model.SupplierCategory) error {
	return m.db.Create(category).Error
}

func (m *SupplierCategoryModule) GetByID(id uint) (*model.SupplierCategory, error) {
	var category model.SupplierCategory
	if err := m.db.Preload("Supplier").First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (m *SupplierCategoryModule) ListBySupplierID(supplierID uint) ([]*model.SupplierCategory, error) {
	var categories []*model.SupplierCategory
	if err := m.db.Where("supplier_id = ? AND status = 1", supplierID).Order("sort ASC, id ASC").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (m *SupplierCategoryModule) UpdateByID(id uint, req *model.UpdateSupplierCategoryReq) error {
	updateMap := updatesPkg.BuildUpdatesFromReq(req)
	if len(updateMap) == 0 {
		return nil
	}
	return m.db.Model(&model.SupplierCategory{}).Where("id = ?", id).Updates(updateMap).Error
}

func (m *SupplierCategoryModule) Delete(id uint) error {
	return m.db.Delete(&model.SupplierCategory{}, id).Error
}
