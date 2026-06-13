package module

import (
	"time"

	"github.com/Kevin-Jii/tower-go/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MeituanAIModule struct {
	db *gorm.DB
}

func NewMeituanAIModule(db *gorm.DB) *MeituanAIModule {
	return &MeituanAIModule{db: db}
}

func (m *MeituanAIModule) CreateAccount(row *model.MeituanAIOperatorAccount) error {
	return m.db.Create(row).Error
}

func (m *MeituanAIModule) UpdateAccount(id, storeID uint, hqUnbound bool, updates map[string]interface{}) error {
	q := m.db.Model(&model.MeituanAIOperatorAccount{}).Where("id = ?", id)
	if !hqUnbound {
		q = q.Where("store_id = ?", storeID)
	}
	return q.Updates(updates).Error
}

func (m *MeituanAIModule) GetAccount(id, storeID uint, hqUnbound bool) (*model.MeituanAIOperatorAccount, error) {
	var row model.MeituanAIOperatorAccount
	q := m.db.Where("id = ?", id)
	if !hqUnbound {
		q = q.Where("store_id = ?", storeID)
	}
	if err := q.First(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (m *MeituanAIModule) ListAccounts(storeID uint, hqUnbound bool) ([]*model.MeituanAIOperatorAccount, error) {
	var rows []*model.MeituanAIOperatorAccount
	q := m.db.Model(&model.MeituanAIOperatorAccount{})
	if !hqUnbound || storeID > 0 {
		q = q.Where("store_id = ?", storeID)
	}
	if err := q.Order("id DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (m *MeituanAIModule) UpsertOrders(account *model.MeituanAIOperatorAccount, orders []model.MeituanAIOrder) (int, error) {
	if len(orders) == 0 {
		return 0, nil
	}
	now := time.Now()
	return len(orders), m.db.Transaction(func(tx *gorm.DB) error {
		for i := range orders {
			orders[i].StoreID = account.StoreID
			orders[i].AccountID = account.ID
			orders[i].ImportedAt = now
			if err := tx.Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "account_id"}, {Name: "order_no"}},
				DoUpdates: clause.AssignmentColumns([]string{
					"order_time", "customer_name", "product_summary", "original_amount", "actual_amount",
					"delivery_fee", "platform_fee", "refund_amount", "status", "raw_json", "imported_at", "updated_at",
				}),
			}).Create(&orders[i]).Error; err != nil {
				return err
			}
		}
		return tx.Model(&model.MeituanAIOperatorAccount{}).Where("id = ?", account.ID).Updates(map[string]interface{}{
			"last_imported_at": now,
		}).Error
	})
}

func (m *MeituanAIModule) UpsertReviews(account *model.MeituanAIOperatorAccount, reviews []model.MeituanAIReview) (int, error) {
	if len(reviews) == 0 {
		return 0, nil
	}
	now := time.Now()
	return len(reviews), m.db.Transaction(func(tx *gorm.DB) error {
		for i := range reviews {
			reviews[i].StoreID = account.StoreID
			reviews[i].AccountID = account.ID
			reviews[i].ImportedAt = now
			if err := tx.Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "account_id"}, {Name: "review_id"}},
				DoUpdates: clause.AssignmentColumns([]string{
					"order_no", "rating", "content", "sentiment", "tags", "suggested_reply", "review_time", "reply_status", "imported_at", "updated_at",
				}),
			}).Create(&reviews[i]).Error; err != nil {
				return err
			}
		}
		return tx.Model(&model.MeituanAIOperatorAccount{}).Where("id = ?", account.ID).Updates(map[string]interface{}{
			"last_imported_at": now,
		}).Error
	})
}

func (m *MeituanAIModule) ListOrders(req *model.ListMeituanAIReq) ([]*model.MeituanAIOrder, int64, error) {
	rows := make([]*model.MeituanAIOrder, 0)
	var total int64
	q := m.db.Model(&model.MeituanAIOrder{}).Where("store_id = ?", req.StoreID)
	if req.AccountID > 0 {
		q = q.Where("account_id = ?", req.AccountID)
	}
	if req.StartDate != "" {
		q = q.Where("order_time >= ?", req.StartDate+" 00:00:00")
	}
	if req.EndDate != "" {
		q = q.Where("order_time <= ?", req.EndDate+" 23:59:59")
	}
	if req.Keyword != "" {
		kw := "%" + req.Keyword + "%"
		q = q.Where("order_no LIKE ? OR product_summary LIKE ? OR customer_name LIKE ?", kw, kw, kw)
	}
	if err := q.Count(&total).Error; err != nil {
		return rows, 0, err
	}
	if err := q.Order("order_time DESC, id DESC").Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&rows).Error; err != nil {
		return rows, 0, err
	}
	return rows, total, nil
}

func (m *MeituanAIModule) ListReviews(req *model.ListMeituanAIReq) ([]*model.MeituanAIReview, int64, error) {
	rows := make([]*model.MeituanAIReview, 0)
	var total int64
	q := m.db.Model(&model.MeituanAIReview{}).Where("store_id = ?", req.StoreID)
	if req.AccountID > 0 {
		q = q.Where("account_id = ?", req.AccountID)
	}
	if req.StartDate != "" {
		q = q.Where("review_time >= ?", req.StartDate+" 00:00:00")
	}
	if req.EndDate != "" {
		q = q.Where("review_time <= ?", req.EndDate+" 23:59:59")
	}
	if req.Keyword != "" {
		kw := "%" + req.Keyword + "%"
		q = q.Where("order_no LIKE ? OR content LIKE ? OR tags LIKE ?", kw, kw, kw)
	}
	if err := q.Count(&total).Error; err != nil {
		return rows, 0, err
	}
	if err := q.Order("review_time DESC, id DESC").Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&rows).Error; err != nil {
		return rows, 0, err
	}
	return rows, total, nil
}

func (m *MeituanAIModule) CreateSuggestions(rows []model.MeituanAISuggestion) error {
	if len(rows) == 0 {
		return nil
	}
	return m.db.Create(&rows).Error
}

func (m *MeituanAIModule) ClearPendingSuggestions(storeID, accountID uint) error {
	return m.db.Where("store_id = ? AND account_id = ? AND status = ?", storeID, accountID, model.MeituanSuggestionStatusPending).
		Delete(&model.MeituanAISuggestion{}).Error
}

func (m *MeituanAIModule) ListSuggestions(req *model.ListMeituanAIReq) ([]*model.MeituanAISuggestion, int64, error) {
	rows := make([]*model.MeituanAISuggestion, 0)
	var total int64
	q := m.db.Model(&model.MeituanAISuggestion{}).Where("store_id = ?", req.StoreID)
	if req.AccountID > 0 {
		q = q.Where("account_id = ?", req.AccountID)
	}
	if req.Keyword != "" {
		kw := "%" + req.Keyword + "%"
		q = q.Where("title LIKE ? OR content LIKE ? OR reason LIKE ?", kw, kw, kw)
	}
	if err := q.Count(&total).Error; err != nil {
		return rows, 0, err
	}
	if err := q.Order("FIELD(status,'pending','approved','done','ignored'), impact_score DESC, id DESC").
		Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&rows).Error; err != nil {
		return rows, 0, err
	}
	return rows, total, nil
}

func (m *MeituanAIModule) UpdateSuggestionStatus(id, storeID uint, status string) error {
	updates := map[string]interface{}{"status": status}
	now := time.Now()
	if status == model.MeituanSuggestionStatusApproved {
		updates["approved_at"] = &now
	}
	if status == model.MeituanSuggestionStatusDone {
		updates["done_at"] = &now
	}
	return m.db.Model(&model.MeituanAISuggestion{}).Where("id = ? AND store_id = ?", id, storeID).Updates(updates).Error
}

func (m *MeituanAIModule) Dashboard(storeID, accountID uint, startDate, endDate string) (*model.MeituanAIDashboard, error) {
	stats := &model.MeituanAIDashboard{}
	orders := m.db.Model(&model.MeituanAIOrder{}).Where("store_id = ?", storeID)
	if accountID > 0 {
		orders = orders.Where("account_id = ?", accountID)
	}
	if startDate != "" {
		orders = orders.Where("order_time >= ?", startDate+" 00:00:00")
	}
	if endDate != "" {
		orders = orders.Where("order_time <= ?", endDate+" 23:59:59")
	}
	if err := orders.Count(&stats.OrderCount).Error; err != nil {
		return nil, err
	}
	if err := orders.Select("COALESCE(SUM(actual_amount),0), COALESCE(SUM(refund_amount),0), COALESCE(SUM(platform_fee),0)").
		Row().Scan(&stats.SalesAmount, &stats.RefundAmount, &stats.PlatformFee); err != nil {
		return nil, err
	}
	if stats.OrderCount > 0 {
		stats.AvgOrderAmount = stats.SalesAmount / float64(stats.OrderCount)
	}

	reviews := m.db.Model(&model.MeituanAIReview{}).Where("store_id = ?", storeID)
	if accountID > 0 {
		reviews = reviews.Where("account_id = ?", accountID)
	}
	if startDate != "" {
		reviews = reviews.Where("review_time >= ?", startDate+" 00:00:00")
	}
	if endDate != "" {
		reviews = reviews.Where("review_time <= ?", endDate+" 23:59:59")
	}
	if err := reviews.Count(&stats.ReviewCount).Error; err != nil {
		return nil, err
	}
	if err := reviews.Where("rating <= 3 OR sentiment = ?", "negative").Count(&stats.NegativeCount).Error; err != nil {
		return nil, err
	}
	if stats.ReviewCount > 0 {
		stats.NegativeRate = float64(stats.NegativeCount) / float64(stats.ReviewCount) * 100
	}

	if err := m.db.Model(&model.MeituanAISuggestion{}).
		Where("store_id = ? AND status = ?", storeID, model.MeituanSuggestionStatusPending).
		Count(&stats.PendingSuggestions).Error; err != nil {
		return nil, err
	}
	return stats, nil
}
