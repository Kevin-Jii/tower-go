package module

import (
	"fmt"
	"time"

	"github.com/Kevin-Jii/tower-go/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SmsCampaignModule struct{ db *gorm.DB }

func NewSmsCampaignModule(db *gorm.DB) *SmsCampaignModule { return &SmsCampaignModule{db: db} }

func (m *SmsCampaignModule) scoped(storeID uint, allStores bool) *gorm.DB {
	q := m.db.Model(&model.SmsCampaign{})
	if !allStores || storeID > 0 {
		q = q.Where("owner_store_id = ?", storeID)
	}
	return q
}

func (m *SmsCampaignModule) List(storeID uint, allStores bool) ([]*model.SmsCampaign, error) {
	var rows []*model.SmsCampaign
	if err := m.scoped(storeID, allStores).Preload("Segments", func(db *gorm.DB) *gorm.DB { return db.Order("position ASC, id ASC") }).Order("id DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (m *SmsCampaignModule) GetByID(id, storeID uint, allStores bool) (*model.SmsCampaign, error) {
	var row model.SmsCampaign
	if err := m.scoped(storeID, allStores).Preload("Segments", func(db *gorm.DB) *gorm.DB { return db.Order("position ASC, id ASC") }).Where("id = ?", id).First(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (m *SmsCampaignModule) GetUnscopedByID(id uint) (*model.SmsCampaign, error) {
	return m.GetByID(id, 0, true)
}

func (m *SmsCampaignModule) Create(row *model.SmsCampaign) error {
	return m.db.Transaction(func(tx *gorm.DB) error { return tx.Create(row).Error })
}

func (m *SmsCampaignModule) Update(id uint, updates map[string]interface{}, segments *[]model.SmsCampaignSegment) error {
	return m.update(id, nil, updates, segments)
}

// UpdateEditable serializes edits against the send claim. Once a worker changes
// the status to sending, the conditional parent update affects no rows and no
// segments are replaced.
func (m *SmsCampaignModule) UpdateEditable(id uint, updates map[string]interface{}, segments *[]model.SmsCampaignSegment) (bool, error) {
	errNotEditable := fmt.Errorf("campaign is not editable")
	err := m.db.Transaction(func(tx *gorm.DB) error {
		if len(updates) == 0 {
			updates = map[string]interface{}{"updated_at": time.Now()}
		}
		res := tx.Model(&model.SmsCampaign{}).
			Where("id = ? AND status IN ?", id, []string{model.SmsCampaignStatusDraft, model.SmsCampaignStatusScheduled, model.SmsCampaignStatusCancelled, model.SmsCampaignStatusFailed}).
			Updates(updates)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected != 1 {
			return errNotEditable
		}
		return replaceCampaignSegments(tx, id, segments)
	})
	if err != nil {
		if err == errNotEditable {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (m *SmsCampaignModule) update(id uint, statuses []string, updates map[string]interface{}, segments *[]model.SmsCampaignSegment) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		if len(updates) > 0 {
			q := tx.Model(&model.SmsCampaign{}).Where("id = ?", id)
			if len(statuses) > 0 {
				q = q.Where("status IN ?", statuses)
			}
			if err := q.Updates(updates).Error; err != nil {
				return err
			}
		}
		return replaceCampaignSegments(tx, id, segments)
	})
}

func replaceCampaignSegments(tx *gorm.DB, id uint, segments *[]model.SmsCampaignSegment) error {
	if segments == nil {
		return nil
	}
	if err := tx.Where("campaign_id = ?", id).Delete(&model.SmsCampaignSegment{}).Error; err != nil {
		return err
	}
	if len(*segments) == 0 {
		return nil
	}
	return tx.Create(segments).Error
}

func (m *SmsCampaignModule) DeleteEditable(id uint) (bool, error) {
	deleted := false
	err := m.db.Transaction(func(tx *gorm.DB) error {
		res := tx.Where("id = ? AND status <> ?", id, model.SmsCampaignStatusSending).Delete(&model.SmsCampaign{})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected != 1 {
			return nil
		}
		deleted = true
		if err := tx.Where("campaign_id = ?", id).Delete(&model.SmsSendRecord{}).Error; err != nil {
			return err
		}
		return tx.Where("campaign_id = ?", id).Delete(&model.SmsCampaignSegment{}).Error
	})
	return deleted, err
}

func (m *SmsCampaignModule) CancelEditable(id uint) (bool, error) {
	res := m.db.Model(&model.SmsCampaign{}).
		Where("id = ? AND status IN ?", id, []string{model.SmsCampaignStatusDraft, model.SmsCampaignStatusScheduled}).
		Updates(map[string]interface{}{"status": model.SmsCampaignStatusCancelled, "scheduled_at": nil})
	return res.RowsAffected == 1, res.Error
}

func (m *SmsCampaignModule) ListDueScheduled(now time.Time, limit int) ([]*model.SmsCampaign, error) {
	var rows []*model.SmsCampaign
	q := m.db.Preload("Segments", func(db *gorm.DB) *gorm.DB { return db.Order("position ASC, id ASC") }).Where("status = ? AND scheduled_at IS NOT NULL AND scheduled_at <= ?", model.SmsCampaignStatusScheduled, now.UTC()).Order("scheduled_at ASC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	if err := q.Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// Claim atomically moves a draft/scheduled campaign to sending. RowsAffected=0
// means another worker already claimed it or its state is no longer sendable.
func (m *SmsCampaignModule) Claim(id uint) (bool, error) {
	res := m.db.Model(&model.SmsCampaign{}).Where("id = ? AND status IN ?", id, []string{model.SmsCampaignStatusDraft, model.SmsCampaignStatusScheduled}).Updates(map[string]interface{}{"status": model.SmsCampaignStatusSending, "last_error": ""})
	return res.RowsAffected == 1, res.Error
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

func (m *SmsCampaignModule) GetSendRecord(campaignID, recordID uint) (*model.SmsSendRecord, error) {
	var row model.SmsSendRecord
	if err := m.db.Where("campaign_id = ? AND id = ?", campaignID, recordID).First(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// ClaimFailedSendRecord prevents concurrent clicks from sending the same failed
// record more than once. A failed retry must call FinishSendRecord to release it.
func (m *SmsCampaignModule) ClaimFailedSendRecord(campaignID, recordID uint) (bool, error) {
	res := m.db.Model(&model.SmsSendRecord{}).
		Where("campaign_id = ? AND id = ? AND status = ?", campaignID, recordID, model.SmsSendRecordFailed).
		Updates(map[string]interface{}{"status": model.SmsSendRecordRetrying, "error_message": ""})
	return res.RowsAffected == 1, res.Error
}

func (m *SmsCampaignModule) FinishSendRecord(campaignID, recordID uint, status, bizID, errorMessage string) error {
	return m.db.Model(&model.SmsSendRecord{}).
		Where("campaign_id = ? AND id = ? AND status = ?", campaignID, recordID, model.SmsSendRecordRetrying).
		Updates(map[string]interface{}{"status": status, "biz_id": bizID, "error_message": errorMessage, "created_at": time.Now().UTC()}).Error
}

func (m *SmsCampaignModule) MarkRetrySuccess(campaignID, recordID uint, bizID string) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&model.SmsSendRecord{}).
			Where("campaign_id = ? AND id = ? AND status = ?", campaignID, recordID, model.SmsSendRecordRetrying).
			Updates(map[string]interface{}{"status": model.SmsSendRecordSuccess, "biz_id": bizID, "error_message": "", "created_at": time.Now().UTC()})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected != 1 {
			return fmt.Errorf("send record retry state changed")
		}
		var campaign model.SmsCampaign
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Select("id", "status", "success_count", "fail_count", "last_error").First(&campaign, campaignID).Error; err != nil {
			return err
		}
		campaign.SuccessCount++
		if campaign.FailCount > 0 {
			campaign.FailCount--
		}
		updates := map[string]interface{}{
			"success_count": campaign.SuccessCount,
			"fail_count":    campaign.FailCount,
		}
		if campaign.FailCount == 0 {
			updates["last_error"] = ""
			updates["status"] = model.SmsCampaignStatusSent
		}
		return tx.Model(&model.SmsCampaign{}).Where("id = ?", campaignID).Updates(updates).Error
	})
}
