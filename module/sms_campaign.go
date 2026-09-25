package module

import (
	"time"

	"github.com/Kevin-Jii/tower-go/model"
	"gorm.io/gorm"
)

type SmsCampaignModule struct {
	db *gorm.DB
}

func NewSmsCampaignModule(db *gorm.DB) *SmsCampaignModule {
	return &SmsCampaignModule{db: db}
}

func (m *SmsCampaignModule) List() ([]*model.SmsCampaign, error) {
	var rows []*model.SmsCampaign
	if err := m.db.Order("id DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (m *SmsCampaignModule) GetByID(id uint) (*model.SmsCampaign, error) {
	var row model.SmsCampaign
	if err := m.db.First(&row, id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (m *SmsCampaignModule) Create(row *model.SmsCampaign) error {
	return m.db.Create(row).Error
}

func (m *SmsCampaignModule) Update(id uint, updates map[string]interface{}) error {
	return m.db.Model(&model.SmsCampaign{}).Where("id = ?", id).Updates(updates).Error
}

func (m *SmsCampaignModule) Delete(id uint) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("campaign_id = ?", id).Delete(&model.SmsSendRecord{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.SmsCampaign{}, id).Error
	})
}

func (m *SmsCampaignModule) ListDueScheduled(now time.Time, limit int) ([]*model.SmsCampaign, error) {
	var rows []*model.SmsCampaign
	q := m.db.Where("status = ? AND scheduled_at IS NOT NULL AND scheduled_at <= ?", model.SmsCampaignStatusScheduled, now).
		Order("scheduled_at ASC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	if err := q.Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (m *SmsCampaignModule) CreateSendRecords(records []*model.SmsSendRecord) error {
	if len(records) == 0 {
		return nil
	}
	return m.db.Create(&records).Error
}

func (m *SmsCampaignModule) ListSendRecords(campaignID uint, limit int) ([]*model.SmsSendRecord, error) {
	var rows []*model.SmsSendRecord
	q := m.db.Where("campaign_id = ?", campaignID).Order("id DESC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	if err := q.Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}
