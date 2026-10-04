package module

import (
	"github.com/Kevin-Jii/tower-go/model"
	"gorm.io/gorm"
)

type AliyunSmsTemplateModule struct{ db *gorm.DB }

func NewAliyunSmsTemplateModule(db *gorm.DB) *AliyunSmsTemplateModule {
	return &AliyunSmsTemplateModule{db: db}
}

func (m *AliyunSmsTemplateModule) List(keyword, auditStatus string) ([]model.AliyunSmsTemplate, error) {
	q := m.db.Model(&model.AliyunSmsTemplate{}).Order("id DESC")
	if keyword = trim(keyword); keyword != "" {
		like := "%" + keyword + "%"
		q = q.Where("name LIKE ? OR template_code LIKE ? OR content LIKE ?", like, like, like)
	}
	if auditStatus = trim(auditStatus); auditStatus != "" {
		q = q.Where("audit_status = ?", auditStatus)
	}
	var rows []model.AliyunSmsTemplate
	if err := q.Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (m *AliyunSmsTemplateModule) GetByCode(code string) (*model.AliyunSmsTemplate, error) {
	var row model.AliyunSmsTemplate
	if err := m.db.Where("template_code = ?", code).First(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (m *AliyunSmsTemplateModule) Upsert(row *model.AliyunSmsTemplate) error {
	return m.db.Save(row).Error
}

func (m *AliyunSmsTemplateModule) DeleteByCode(code string) error {
	return m.db.Where("template_code = ?", code).Delete(&model.AliyunSmsTemplate{}).Error
}

func trim(s string) string {
	for len(s) > 0 && (s[0] == ' ' || s[0] == '\t' || s[0] == '\n' || s[0] == '\r') {
		s = s[1:]
	}
	for len(s) > 0 && (s[len(s)-1] == ' ' || s[len(s)-1] == '\t' || s[len(s)-1] == '\n' || s[len(s)-1] == '\r') {
		s = s[:len(s)-1]
	}
	return s
}
