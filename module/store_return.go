package module

import (
	"fmt"
	"strings"
	"time"

	"github.com/Kevin-Jii/tower-go/model"
	"gorm.io/gorm"
)

type StoreReturnModule struct {
	db *gorm.DB
}

func NewStoreReturnModule(db *gorm.DB) *StoreReturnModule {
	return &StoreReturnModule{db: db}
}

func (m *StoreReturnModule) Create(record *model.StoreReturn) error {
	return m.db.Create(record).Error
}

func IsDuplicateKeyError(err error) bool {
	return err != nil && (strings.Contains(err.Error(), "Error 1062") || strings.Contains(strings.ToLower(err.Error()), "duplicate entry"))
}

func (m *StoreReturnModule) GetByClientReqIDScoped(clientReqID string, storeID uint, hqUnbound bool) (*model.StoreReturn, error) {
	clientReqID = strings.TrimSpace(clientReqID)
	if clientReqID == "" {
		return nil, gorm.ErrRecordNotFound
	}
	var record model.StoreReturn
	query := m.db.Preload("Items").Preload("Store").Where("client_req_id = ?", clientReqID)
	if !hqUnbound {
		query = query.Where("store_id = ?", storeID)
	}
	if err := query.First(&record).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

func (m *StoreReturnModule) GetByIDScoped(id, storeID uint, hqUnbound bool) (*model.StoreReturn, error) {
	var record model.StoreReturn
	query := m.db.Preload("Items").Preload("Store").Where("id = ?", id)
	if !hqUnbound {
		query = query.Where("store_id = ?", storeID)
	}
	if err := query.First(&record).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

func (m *StoreReturnModule) List(req *model.ListStoreReturnReq) ([]*model.StoreReturn, int64, error) {
	records := make([]*model.StoreReturn, 0)
	var total int64

	query := m.db.Model(&model.StoreReturn{}).Preload("Store")
	if req.StoreID > 0 {
		query = query.Where("store_id = ?", req.StoreID)
	}
	if keyword := strings.TrimSpace(req.Keyword); keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where(
			"return_no LIKE ? OR remark LIKE ? OR operator_name LIKE ? OR EXISTS (SELECT 1 FROM store_return_items sri WHERE sri.return_id = store_returns.id AND sri.deleted_at IS NULL AND sri.product_name LIKE ?)",
			like, like, like, like,
		)
	}
	if req.StartDate != "" {
		query = query.Where("return_date >= ?", req.StartDate)
	}
	if req.EndDate != "" {
		query = query.Where("return_date <= ?", req.EndDate)
	}

	if err := query.Count(&total).Error; err != nil {
		return records, 0, err
	}

	offset := (req.Page - 1) * req.PageSize
	if err := query.Preload("Items").Order("id DESC").Offset(offset).Limit(req.PageSize).Find(&records).Error; err != nil {
		return records, 0, err
	}
	return records, total, nil
}

func (m *StoreReturnModule) Update(record *model.StoreReturn) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("return_id = ?", record.ID).Delete(&model.StoreReturnItem{}).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.StoreReturn{}).Where("id = ?", record.ID).Updates(map[string]interface{}{
			"store_id":      record.StoreID,
			"client_req_id": record.ClientReqID,
			"return_date":   record.ReturnDate,
			"logistics_fee": record.LogisticsFee,
			"total_deposit": record.TotalDeposit,
			"item_count":    record.ItemCount,
			"remark":        record.Remark,
			"operator_id":   record.OperatorID,
			"operator_name": record.OperatorName,
		}).Error; err != nil {
			return err
		}
		for i := range record.Items {
			record.Items[i].ReturnID = record.ID
		}
		if len(record.Items) > 0 {
			return tx.Create(&record.Items).Error
		}
		return nil
	})
}

func (m *StoreReturnModule) Delete(id, storeID uint, hqUnbound bool) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		query := tx.Model(&model.StoreReturn{}).Where("id = ?", id)
		if !hqUnbound {
			query = query.Where("store_id = ?", storeID)
		}
		var existing model.StoreReturn
		if err := query.First(&existing).Error; err != nil {
			return err
		}
		if err := tx.Where("return_id = ?", existing.ID).Delete(&model.StoreReturnItem{}).Error; err != nil {
			return err
		}
		return tx.Delete(&existing).Error
	})
}

func (m *StoreReturnModule) Stats(req *model.ListStoreReturnReq) (*model.StoreReturnStats, error) {
	stats := &model.StoreReturnStats{}
	query := m.db.Model(&model.StoreReturn{})
	if req.StoreID > 0 {
		query = query.Where("store_id = ?", req.StoreID)
	}
	if req.StartDate != "" {
		query = query.Where("return_date >= ?", req.StartDate)
	}
	if req.EndDate != "" {
		query = query.Where("return_date <= ?", req.EndDate)
	}
	if err := query.Select("COALESCE(SUM(total_deposit),0) AS total_deposit, COALESCE(SUM(logistics_fee),0) AS logistics_fee, COUNT(*) AS return_count, COALESCE(SUM(item_count),0) AS item_count").Scan(stats).Error; err != nil {
		return nil, err
	}
	return stats, nil
}

func (m *StoreReturnModule) CreateProduct(product *model.StoreReturnProduct) error {
	return m.db.Create(product).Error
}

func (m *StoreReturnModule) GetProductByIDScoped(id, storeID uint, hqUnbound bool) (*model.StoreReturnProduct, error) {
	var product model.StoreReturnProduct
	query := m.db.Where("id = ?", id)
	if !hqUnbound {
		query = query.Where("store_id = ?", storeID)
	}
	if err := query.First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (m *StoreReturnModule) GetProductMap(ids []uint, storeID uint, hqUnbound bool) (map[uint]*model.StoreReturnProduct, error) {
	result := make(map[uint]*model.StoreReturnProduct)
	if len(ids) == 0 {
		return result, nil
	}
	var products []*model.StoreReturnProduct
	query := m.db.Where("id IN ?", ids)
	if !hqUnbound {
		query = query.Where("store_id = ?", storeID)
	}
	if err := query.Find(&products).Error; err != nil {
		return nil, err
	}
	for _, product := range products {
		result[product.ID] = product
	}
	return result, nil
}

func (m *StoreReturnModule) ListProducts(req *model.ListStoreReturnProductReq) ([]*model.StoreReturnProduct, int64, error) {
	products := make([]*model.StoreReturnProduct, 0)
	var total int64

	query := m.db.Model(&model.StoreReturnProduct{}).Preload("Store")
	if req.StoreID > 0 {
		query = query.Where("store_id = ?", req.StoreID)
	}
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}
	if keyword := strings.TrimSpace(req.Keyword); keyword != "" {
		query = query.Where("product_name LIKE ? OR remark LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if err := query.Count(&total).Error; err != nil {
		return products, 0, err
	}
	offset := (req.Page - 1) * req.PageSize
	if err := query.Order("id DESC").Offset(offset).Limit(req.PageSize).Find(&products).Error; err != nil {
		return products, 0, err
	}
	return products, total, nil
}

func (m *StoreReturnModule) UpdateProduct(product *model.StoreReturnProduct) error {
	return m.db.Model(&model.StoreReturnProduct{}).Where("id = ?", product.ID).Updates(map[string]interface{}{
		"store_id":     product.StoreID,
		"product_name": product.ProductName,
		"deposit":      product.Deposit,
		"remark":       product.Remark,
		"status":       product.Status,
	}).Error
}

func (m *StoreReturnModule) DeleteProduct(id, storeID uint, hqUnbound bool) error {
	query := m.db.Where("id = ?", id)
	if !hqUnbound {
		query = query.Where("store_id = ?", storeID)
	}
	return query.Delete(&model.StoreReturnProduct{}).Error
}

func (m *StoreReturnModule) GenerateReturnNo() string {
	prefix := "FC"
	today := time.Now().Format("20060102")
	pattern := prefix + today + "%"

	var maxNo string
	m.db.Model(&model.StoreReturn{}).
		Where("return_no LIKE ?", pattern).
		Order("return_no DESC").
		Limit(1).
		Pluck("return_no", &maxNo)

	seq := 1
	if maxNo != "" && len(maxNo) >= 14 {
		fmt.Sscanf(maxNo[len(maxNo)-4:], "%d", &seq)
		seq++
	}
	return fmt.Sprintf("%s%s%04d", prefix, today, seq)
}
