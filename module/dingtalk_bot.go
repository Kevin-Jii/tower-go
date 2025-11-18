package module

import (
	"github.com/Kevin-Jii/tower-go/model"

	"gorm.io/gorm"
)

type DingTalkBotModule struct {
	db *gorm.DB
}

func NewDingTalkBotModule(db *gorm.DB) *DingTalkBotModule {
	return &DingTalkBotModule{db: db}
}

// Create 创建钉钉机器人配置
func (m *DingTalkBotModule) Create(bot *model.DingTalkBot) error {
	return m.db.Create(bot).Error
}

// GetByID 根据ID获取配置
func (m *DingTalkBotModule) GetByID(id uint) (*model.DingTalkBot, error) {
	var bot model.DingTalkBot
	if err := m.db.Preload("Store").First(&bot, id).Error; err != nil {
		return nil, err
	}
	return &bot, nil
}

// List 获取机器人列表
func (m *DingTalkBotModule) List(page, pageSize int) ([]*model.DingTalkBot, int64, error) {
	var bots []*model.DingTalkBot
	var total int64

	query := m.db.Model(&model.DingTalkBot{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Store").
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&bots).Error; err != nil {
		return nil, 0, err
	}

	return bots, total, nil
}

// ListEnabledByStoreID 获取指定门店启用的机器人（包含全局机器人）
func (m *DingTalkBotModule) ListEnabledByStoreID(storeID uint) ([]*model.DingTalkBot, error) {
	var bots []*model.DingTalkBot
	if err := m.db.Where("is_enabled = ? AND (store_id = ? OR store_id IS NULL)", true, storeID).
		Find(&bots).Error; err != nil {
		return nil, err
	}
	return bots, nil
}

// ListAllEnabled 获取所有启用的机器人
func (m *DingTalkBotModule) ListAllEnabled() ([]*model.DingTalkBot, error) {
	var bots []*model.DingTalkBot
	if err := m.db.Where("is_enabled = ?", true).Find(&bots).Error; err != nil {
		return nil, err
	}
	return bots, nil
}

// ListEnabledStreamBots 获取所有启用的 Stream 类型机器人
func (m *DingTalkBotModule) ListEnabledStreamBots() ([]*model.DingTalkBot, error) {
	var bots []*model.DingTalkBot
	if err := m.db.Where("is_enabled = ? AND bot_type = ?", true, "stream").
		Find(&bots).Error; err != nil {
		return nil, err
	}
	return bots, nil
}

// Update 更新机器人配置
func (m *DingTalkBotModule) Update(id uint, updates map[string]interface{}) error {
	return m.db.Model(&model.DingTalkBot{}).Where("id = ?", id).Updates(updates).Error
}

// Delete 删除机器人配置
func (m *DingTalkBotModule) Delete(id uint) error {
	return m.db.Delete(&model.DingTalkBot{}, id).Error
}

// ExistsByWebhook 检查 Webhook 是否已存在
func (m *DingTalkBotModule) ExistsByWebhook(webhook string, excludeID uint) (bool, error) {
	var count int64
	query := m.db.Model(&model.DingTalkBot{}).Where("webhook = ?", webhook)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// FindByClientID 根据 ClientID 查找机器人
func (m *DingTalkBotModule) FindByClientID(clientID string) (*model.DingTalkBot, error) {
	var bot model.DingTalkBot
	if err := m.db.Where("client_id = ?", clientID).First(&bot).Error; err != nil {
		return nil, err
	}
	return &bot, nil
}

// UpdateName 更新机器人名称
func (m *DingTalkBotModule) UpdateName(id uint, name string) error {
	return m.db.Model(&model.DingTalkBot{}).Where("id = ?", id).Update("name", name).Error
}

// GetStoreByID 获取门店信息
func (m *DingTalkBotModule) GetStoreByID(storeID uint) (*model.Store, error) {
	var store model.Store
	if err := m.db.First(&store, storeID).Error; err != nil {
		return nil, err
	}
	return &store, nil
}
