package module

import (
	"github.com/Kevin-Jii/tower-go/model"
	"gorm.io/gorm"
)

type AIAssistantModule struct{ db *gorm.DB }

func NewAIAssistantModule(db *gorm.DB) *AIAssistantModule { return &AIAssistantModule{db: db} }
func (m *AIAssistantModule) Config() (*model.AIAssistantConfig, error) {
	var row model.AIAssistantConfig
	err := m.db.First(&row, 1).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &row, err
}
func (m *AIAssistantModule) SaveConfig(row *model.AIAssistantConfig) error {
	row.ID = 1
	return m.db.Save(row).Error
}
func (m *AIAssistantModule) CreateConversation(row *model.AIAssistantConversation) error {
	return m.db.Create(row).Error
}
func (m *AIAssistantModule) Conversation(id, storeID uint) (*model.AIAssistantConversation, error) {
	var row model.AIAssistantConversation
	err := m.db.Where("id = ? AND store_id = ?", id, storeID).First(&row).Error
	return &row, err
}
func (m *AIAssistantModule) Conversations(storeID uint) ([]model.AIAssistantConversation, error) {
	var rows []model.AIAssistantConversation
	err := m.db.Where("store_id = ?", storeID).Order("updated_at DESC").Limit(100).Find(&rows).Error
	return rows, err
}
func (m *AIAssistantModule) RenameConversation(id, storeID uint, title string) error {
	return m.db.Model(&model.AIAssistantConversation{}).Where("id = ? AND store_id = ?", id, storeID).Update("title", title).Error
}
func (m *AIAssistantModule) DeleteConversation(id, storeID uint) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND store_id = ?", id, storeID).Delete(&model.AIAssistantConversation{}).Error; err != nil {
			return err
		}
		return tx.Where("conversation_id = ?", id).Delete(&model.AIAssistantMessage{}).Error
	})
}
func (m *AIAssistantModule) Messages(id uint) ([]model.AIAssistantMessage, error) {
	var rows []model.AIAssistantMessage
	err := m.db.Where("conversation_id = ?", id).Order("id ASC").Limit(200).Find(&rows).Error
	return rows, err
}
func (m *AIAssistantModule) AddMessage(row *model.AIAssistantMessage) error {
	return m.db.Create(row).Error
}
func (m *AIAssistantModule) TouchConversation(id uint) error {
	return m.db.Model(&model.AIAssistantConversation{}).Where("id = ?", id).UpdateColumn("updated_at", gorm.Expr("CURRENT_TIMESTAMP")).Error
}
func (m *AIAssistantModule) RenameIfUntitled(id uint, title string) error {
	return m.db.Model(&model.AIAssistantConversation{}).Where("id = ? AND title = ?", id, "新对话").Update("title", title).Error
}

func (m *AIAssistantModule) PaymentSummary(storeID uint, startDate, endDate string) (map[string]any, error) {
	query := m.db.Model(&model.StoreAccount{}).Where("deleted_at IS NULL AND is_canceled = 0 AND account_date >= ? AND account_date <= ?", startDate, endDate)
	if storeID > 0 {
		query = query.Where("store_id = ?", storeID)
	}
	var rows []struct {
		PaymentStatus int     `json:"payment_status"`
		Orders        int64   `json:"orders"`
		Amount        float64 `json:"amount"`
	}
	if err := query.Select("payment_status, COUNT(*) AS orders, COALESCE(SUM(total_amount), 0) AS amount").Group("payment_status").Scan(&rows).Error; err != nil {
		return nil, err
	}
	summary := map[string]any{"paid_orders": int64(0), "paid_amount": float64(0), "unpaid_orders": int64(0), "unpaid_amount": float64(0)}
	for _, row := range rows {
		if row.PaymentStatus == model.StoreAccountPaymentPaid {
			summary["paid_orders"], summary["paid_amount"] = row.Orders, row.Amount
		} else if row.PaymentStatus == model.StoreAccountPaymentUnpaid {
			summary["unpaid_orders"], summary["unpaid_amount"] = row.Orders, row.Amount
		}
	}
	return summary, nil
}
