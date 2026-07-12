package module

import (
	"strings"

	"github.com/Kevin-Jii/tower-go/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductUnitSpecModule struct {
	db *gorm.DB
}

func NewProductUnitSpecModule(db *gorm.DB) *ProductUnitSpecModule {
	return &ProductUnitSpecModule{db: db}
}

func (m *ProductUnitSpecModule) GetByProductAndUnit(productID uint, unit string) (*model.ProductUnitSpec, error) {
	var spec model.ProductUnitSpec
	u := strings.TrimSpace(unit)
	err := m.db.Where("product_id = ? AND is_enabled = 1 AND unit_name = ?", productID, u).
		Order("id asc").
		First(&spec).Error
	if err == nil {
		return &spec, nil
	}
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	var specs []model.ProductUnitSpec
	if err := m.db.Where("product_id = ? AND is_enabled = 1 AND unit_code = ?", productID, u).
		Order("id ASC").Limit(2).Find(&specs).Error; err != nil {
		return nil, err
	}
	if len(specs) != 1 {
		return nil, gorm.ErrRecordNotFound
	}
	spec = specs[0]
	return &spec, nil
}

func (m *ProductUnitSpecModule) Create(spec *model.ProductUnitSpec) error {
	return m.db.Create(spec).Error
}

func (m *ProductUnitSpecModule) GetByID(id uint) (*model.ProductUnitSpec, error) {
	var spec model.ProductUnitSpec
	if err := m.db.First(&spec, id).Error; err != nil {
		return nil, err
	}
	return &spec, nil
}

func (m *ProductUnitSpecModule) ListByProductID(productID uint) ([]*model.ProductUnitSpec, error) {
	var specs []*model.ProductUnitSpec
	if err := m.db.Where("product_id = ?", productID).Order("id asc").Find(&specs).Error; err != nil {
		return nil, err
	}
	return specs, nil
}

func (m *ProductUnitSpecModule) ListEnabledByProductID(productID uint) ([]*model.ProductUnitSpec, error) {
	var specs []*model.ProductUnitSpec
	if err := m.db.Where("product_id = ? AND is_enabled = 1", productID).Order("factor_to_base asc, id asc").Find(&specs).Error; err != nil {
		return nil, err
	}
	return specs, nil
}

// ListByProductIDs 批量查询多个商品的单位规格（按 product_id、id 排序）
func (m *ProductUnitSpecModule) ListByProductIDs(productIDs []uint) ([]*model.ProductUnitSpec, error) {
	if len(productIDs) == 0 {
		return nil, nil
	}
	var specs []*model.ProductUnitSpec
	if err := m.db.Where("product_id IN ? AND is_enabled = 1", productIDs).Order("product_id asc, factor_to_base asc, id asc").Find(&specs).Error; err != nil {
		return nil, err
	}
	return specs, nil
}

func (m *ProductUnitSpecModule) UpdateByID(id uint, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}
	return m.db.Model(&model.ProductUnitSpec{}).Where("id = ?", id).Updates(updates).Error
}

func (m *ProductUnitSpecModule) DeleteByID(id uint) error {
	return m.db.Delete(&model.ProductUnitSpec{}, id).Error
}

func (m *ProductUnitSpecModule) UpsertByProductAndUnit(spec *model.ProductUnitSpec) error {
	return m.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "product_id"}, {Name: "unit_code"}, {Name: "unit_name"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"factor_to_base", "precision", "cost_price", "sale_price", "is_enabled",
		}),
	}).Create(spec).Error
}
