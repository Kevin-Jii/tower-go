package module

import (
	"errors"
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

func (m *ProductUnitSpecModule) ListEnabledByProductIDWithConsumables(productID, storeID uint) ([]*model.ProductUnitSpec, error) {
	var specs []*model.ProductUnitSpec
	query := m.db.Where("product_id = ? AND is_enabled = 1", productID)
	if storeID > 0 {
		query = query.
			Preload("Consumables", "store_id = ?", storeID).
			Preload("Consumables.ConsumableProduct")
	}
	if err := query.Order("factor_to_base asc, id asc").Find(&specs).Error; err != nil {
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

func (m *ProductUnitSpecModule) ListByProductIDsWithConsumables(productIDs []uint, storeID uint) ([]*model.ProductUnitSpec, error) {
	if len(productIDs) == 0 {
		return nil, nil
	}
	var specs []*model.ProductUnitSpec
	query := m.db.Where("product_id IN ? AND is_enabled = 1", productIDs)
	if storeID > 0 {
		query = query.
			Preload("Consumables", "store_id = ?", storeID).
			Preload("Consumables.ConsumableProduct")
	}
	if err := query.Order("product_id asc, factor_to_base asc, id asc").Find(&specs).Error; err != nil {
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
	return m.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("unit_spec_id = ?", id).Delete(&model.ProductUnitSpecConsumable{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.ProductUnitSpec{}, id).Error
	})
}

func (m *ProductUnitSpecModule) UpsertByProductAndUnit(spec *model.ProductUnitSpec) error {
	return m.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "product_id"}, {Name: "unit_code"}, {Name: "unit_name"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"factor_to_base", "precision", "cost_price", "sale_price", "is_saleable", "is_enabled",
		}),
	}).Create(spec).Error
}

func (m *ProductUnitSpecModule) GetByProductAndUnitName(productID uint, unitCode, unitName string) (*model.ProductUnitSpec, error) {
	var spec model.ProductUnitSpec
	if err := m.db.Where("product_id = ? AND unit_code = ? AND unit_name = ?", productID, unitCode, unitName).First(&spec).Error; err != nil {
		return nil, err
	}
	return &spec, nil
}

func (m *ProductUnitSpecModule) ReplaceConsumables(unitSpecID, storeID uint, items []model.ProductUnitSpecConsumableItemReq) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("unit_spec_id = ? AND store_id = ?", unitSpecID, storeID).
			Delete(&model.ProductUnitSpecConsumable{}).Error; err != nil {
			return err
		}
		if len(items) == 0 {
			return nil
		}

		ids := make([]uint, 0, len(items))
		seen := make(map[uint]struct{}, len(items))
		for _, item := range items {
			if item.ConsumableProductID == 0 || item.Quantity <= 0 {
				return errors.New("invalid product unit spec consumable")
			}
			if _, exists := seen[item.ConsumableProductID]; exists {
				return errors.New("duplicate product unit spec consumable")
			}
			seen[item.ConsumableProductID] = struct{}{}
			ids = append(ids, item.ConsumableProductID)
		}
		var count int64
		if err := tx.Model(&model.StoreAccountConsumableProduct{}).
			Where("store_id = ? AND id IN ?", storeID, ids).Count(&count).Error; err != nil {
			return err
		}
		if count != int64(len(ids)) {
			return gorm.ErrRecordNotFound
		}

		rows := make([]model.ProductUnitSpecConsumable, 0, len(items))
		for _, item := range items {
			rows = append(rows, model.ProductUnitSpecConsumable{
				UnitSpecID:          unitSpecID,
				StoreID:             storeID,
				ConsumableProductID: item.ConsumableProductID,
				Quantity:            item.Quantity,
			})
		}
		return tx.Create(&rows).Error
	})
}
