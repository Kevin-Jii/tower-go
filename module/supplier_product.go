package module

import (
	"errors"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/pkg/apicode"
	updatesPkg "github.com/Kevin-Jii/tower-go/utils/updates"
	"github.com/go-sql-driver/mysql"
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

func (m *SupplierProductModule) List(req *model.ListSupplierProductReq) ([]*model.SupplierProduct, error) {
	var products []*model.SupplierProduct

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

	if err := query.Preload("Supplier").Preload("Category").Order("id DESC").Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (m *SupplierProductModule) UpdateByID(id uint, req *model.UpdateSupplierProductReq) error {
	updateMap := updatesPkg.BuildUpdatesFromReq(req)
	if len(updateMap) == 0 {
		return nil
	}
	return m.db.Model(&model.SupplierProduct{}).Where("id = ?", id).Updates(updateMap).Error
}

func (m *SupplierProductModule) Delete(id uint) error {
	err := m.db.Delete(&model.SupplierProduct{}, id).Error
	if isForeignKeyReferenceError(err) {
		return apicode.Wrap(
			apicode.ResourceInUse.WithMessage("商品已有库存、单据或其他业务记录，无法删除，请改为停用"),
			err,
		)
	}
	return err
}

func isForeignKeyReferenceError(err error) bool {
	var mysqlErr *mysql.MySQLError
	return errors.As(err, &mysqlErr) && mysqlErr.Number == 1451
}

// GetByIDs 批量获取商品
func (m *SupplierProductModule) GetByIDs(ids []uint) ([]*model.SupplierProduct, error) {
	var products []*model.SupplierProduct
	if err := m.db.Preload("Supplier").Preload("Category").Where("id IN ?", ids).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}
