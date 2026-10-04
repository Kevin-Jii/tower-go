package module

import (
	"strings"

	"github.com/Kevin-Jii/tower-go/model"
	"gorm.io/gorm"
)

type AliyunSmsTemplateModule struct{ db *gorm.DB }

func NewAliyunSmsTemplateModule(db *gorm.DB) *AliyunSmsTemplateModule {
	return &AliyunSmsTemplateModule{db: db}
}

func (m *AliyunSmsTemplateModule) scoped(storeID uint, allStores bool) *gorm.DB {
	q := m.db.Model(&model.AliyunSmsTemplate{})
	if !allStores || storeID > 0 {
		q = q.Where("owner_store_id = ?", storeID)
	}
	return q
}

func (m *AliyunSmsTemplateModule) List(storeID uint, allStores bool, keyword, auditStatus string) ([]model.AliyunSmsTemplate, error) {
	q := m.scoped(storeID, allStores).Order("id DESC")
	if keyword = strings.TrimSpace(keyword); keyword != "" {
		like := "%" + keyword + "%"
		q = q.Where("name LIKE ? OR template_code LIKE ? OR content LIKE ?", like, like, like)
	}
	if auditStatus = strings.TrimSpace(auditStatus); auditStatus != "" {
		q = q.Where("audit_status = ?", auditStatus)
	}
	var rows []model.AliyunSmsTemplate
	if err := q.Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (m *AliyunSmsTemplateModule) GetByCode(code string, storeID uint, allStores bool) (*model.AliyunSmsTemplate, error) {
	var row model.AliyunSmsTemplate
	if err := m.scoped(storeID, allStores).Where("template_code = ?", code).First(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (m *AliyunSmsTemplateModule) ExistsByName(name string, storeID uint) (bool, error) {
	var count int64
	err := m.db.Model(&model.AliyunSmsTemplate{}).
		Where("owner_store_id = ? AND name = ?", storeID, strings.TrimSpace(name)).
		Count(&count).Error
	return count > 0, err
}

func (m *AliyunSmsTemplateModule) Upsert(row *model.AliyunSmsTemplate) error {
	return m.db.Save(row).Error
}

func (m *AliyunSmsTemplateModule) Delete(row *model.AliyunSmsTemplate) error {
	return m.db.Delete(row).Error
}
