package module

import (
	"fmt"

	"github.com/Kevin-Jii/tower-go/model"

	"gorm.io/gorm"
)

type SupplierModule struct {
	db *gorm.DB
}

func NewSupplierModule(db *gorm.DB) *SupplierModule {
	return &SupplierModule{db: db}
}

func (m *SupplierModule) GetDB() *gorm.DB {
	return m.db
}

func (m *SupplierModule) Create(supplier *model.Supplier) error {
	return m.db.Create(supplier).Error
}

func (m *SupplierModule) GetByID(id uint) (*model.Supplier, error) {
	var supplier model.Supplier
	if err := m.db.First(&supplier, id).Error; err != nil {
		return nil, err
	}
	return &supplier, nil
}

func (m *SupplierModule) GetByCode(code string) (*model.Supplier, error) {
	var supplier model.Supplier
	if err := m.db.Where("supplier_code = ?", code).First(&supplier).Error; err != nil {
		return nil, err
	}
	return &supplier, nil
}

func (m *SupplierModule) ExistsByCode(code string) (bool, error) {
	var count int64
	if err := m.db.Model(&model.Supplier{}).Where("supplier_code = ?", code).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (m *SupplierModule) List(req *model.ListSupplierReq) ([]*model.Supplier, int64, error) {
	var suppliers []*model.Supplier
	var total int64

	query := m.db.Model(&model.Supplier{})

	if req.Keyword != "" {
		keyword := "%" + req.Keyword + "%"
		query = query.Where("supplier_name LIKE ? OR supplier_code LIKE ?", keyword, keyword)
	}

	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (req.Page - 1) * req.PageSize
	if err := query.Order("id DESC").Offset(offset).Limit(req.PageSize).Find(&suppliers).Error; err != nil {
		return nil, 0, err
	}

	return suppliers, total, nil
}

func (m *SupplierModule) UpdateByID(id uint, req *model.UpdateSupplierReq) error {
	updateMap := make(map[string]interface{})

	if req.SupplierName != nil {
		updateMap["supplier_name"] = *req.SupplierName
	}
	if req.ContactPerson != nil {
		updateMap["contact_person"] = *req.ContactPerson
	}
	if req.ContactPhone != nil {
		updateMap["contact_phone"] = *req.ContactPhone
	}
	if req.ContactEmail != nil {
		updateMap["contact_email"] = *req.ContactEmail
	}
	if req.SupplierAddress != nil {
		updateMap["supplier_address"] = *req.SupplierAddress
	}
	if req.Remark != nil {
		updateMap["remark"] = *req.Remark
	}
	if req.Status != nil {
		updateMap["status"] = *req.Status
	}

	if len(updateMap) == 0 {
		return nil
	}
	return m.db.Model(&model.Supplier{}).Where("id = ?", id).Updates(updateMap).Error
}

func (m *SupplierModule) Delete(id uint) error {
	return m.db.Delete(&model.Supplier{}, id).Error
}

// GetMaxSeqByStorePrefix 获取指定门店前缀的最大序号
func (m *SupplierModule) GetMaxSeqByStorePrefix(storeID uint) (int, error) {
	var maxCode string
	prefix := fmt.Sprintf("%d", storeID)
	// 查找以门店ID开头的最大编码
	err := m.db.Model(&model.Supplier{}).
		Where("supplier_code LIKE ?", prefix+"%").
		Order("supplier_code DESC").
		Limit(1).
		Pluck("supplier_code", &maxCode).Error

	if err != nil || maxCode == "" {
		return 0, nil
	}

	// 提取序号部分（去掉门店ID前缀）
	if len(maxCode) > len(prefix) {
		seqStr := maxCode[len(prefix):]
		var seq int
		fmt.Sscanf(seqStr, "%d", &seq)
		return seq, nil
	}
	return 0, nil
}
