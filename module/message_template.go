package module

import (
	"github.com/Kevin-Jii/tower-go/model"
	"gorm.io/gorm"
)

type MessageTemplateModule struct {
	db *gorm.DB
}

func NewMessageTemplateModule(db *gorm.DB) *MessageTemplateModule {
	return &MessageTemplateModule{db: db}
}

// GetByCode 根据编码获取模板
func (m *MessageTemplateModule) GetByCode(code string) (*model.MessageTemplate, error) {
	var template model.MessageTemplate
	if err := m.db.Where("code = ? AND is_enabled = ?", code, true).First(&template).Error; err != nil {
		return nil, err
	}
	return &template, nil
}

// GetByID 根据ID获取模板
func (m *MessageTemplateModule) GetByID(id uint) (*model.MessageTemplate, error) {
	var template model.MessageTemplate
	if err := m.db.First(&template, id).Error; err != nil {
		return nil, err
	}
	return &template, nil
}

// List 获取模板列表
func (m *MessageTemplateModule) List() ([]*model.MessageTemplate, error) {
	var templates []*model.MessageTemplate
	if err := m.db.Order("id ASC").Find(&templates).Error; err != nil {
		return nil, err
	}
	return templates, nil
}

// Create 创建模板
func (m *MessageTemplateModule) Create(template *model.MessageTemplate) error {
	return m.db.Create(template).Error
}

// Update 更新模板
func (m *MessageTemplateModule) Update(id uint, updates map[string]interface{}) error {
	return m.db.Model(&model.MessageTemplate{}).Where("id = ?", id).Updates(updates).Error
}

// Delete 删除模板
func (m *MessageTemplateModule) Delete(id uint) error {
	return m.db.Delete(&model.MessageTemplate{}, id).Error
}

// Upsert 创建或更新模板（根据code）
func (m *MessageTemplateModule) Upsert(template *model.MessageTemplate) error {
	var existing model.MessageTemplate
	err := m.db.Where("code = ?", template.Code).First(&existing).Error
	if err == gorm.ErrRecordNotFound {
		return m.db.Create(template).Error
	}
	if err != nil {
		return err
	}
	// 更新现有记录
	return m.db.Model(&existing).Updates(map[string]interface{}{
		"name":        template.Name,
		"title":       template.Title,
		"content":     template.Content,
		"description": template.Description,
		"variables":   template.Variables,
	}).Error
}
