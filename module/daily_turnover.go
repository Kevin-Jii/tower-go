package module

import (
	"context"

	"github.com/Kevin-Jii/tower-go/model"
	"gorm.io/gorm"
)

type DailyTurnoverStoreRow struct {
	StoreID   uint
	StoreName string
}

type DailyTurnoverSummaryRow struct {
	StoreID     uint
	TotalAmount float64
	OrderCount  int64
}

type DailyTurnoverChannelRow struct {
	StoreID    uint
	Channel    string
	Amount     float64
	OrderCount int64
}

type DailyTurnoverAdminRow struct {
	StoreID uint
	UserID  uint
	OpenID  string
}

type DailyTurnoverModule struct {
	db *gorm.DB
}

func NewDailyTurnoverModule(db *gorm.DB) *DailyTurnoverModule {
	return &DailyTurnoverModule{db: db}
}

func (m *DailyTurnoverModule) ListActiveStores(ctx context.Context) ([]DailyTurnoverStoreRow, error) {
	stores := make([]DailyTurnoverStoreRow, 0)
	err := m.db.WithContext(ctx).
		Table("stores").
		Select("id AS store_id, name AS store_name").
		Where("status = ?", 1).
		Where("store_code IS NULL OR store_code <> ?", model.StoreCodeHQ).
		Order("id ASC").
		Scan(&stores).Error
	return stores, err
}

func (m *DailyTurnoverModule) ListSummaries(ctx context.Context, businessDate string) ([]DailyTurnoverSummaryRow, error) {
	summaries := make([]DailyTurnoverSummaryRow, 0)
	err := m.db.WithContext(ctx).
		Model(&model.StoreAccount{}).
		Select("store_id, COALESCE(SUM(total_amount), 0) AS total_amount, COUNT(*) AS order_count").
		Where("account_date = ? AND is_canceled = ?", businessDate, false).
		Group("store_id").
		Scan(&summaries).Error
	return summaries, err
}

func (m *DailyTurnoverModule) ListChannels(ctx context.Context, businessDate string) ([]DailyTurnoverChannelRow, error) {
	channels := make([]DailyTurnoverChannelRow, 0)
	err := m.db.WithContext(ctx).
		Model(&model.StoreAccount{}).
		Select("store_id, channel, COALESCE(SUM(total_amount), 0) AS amount, COUNT(*) AS order_count").
		Where("account_date = ? AND is_canceled = ?", businessDate, false).
		Group("store_id, channel").
		Order("store_id ASC, amount DESC").
		Scan(&channels).Error
	return channels, err
}

func (m *DailyTurnoverModule) ListAdministrators(ctx context.Context, storeIDs []uint) ([]DailyTurnoverAdminRow, error) {
	admins := make([]DailyTurnoverAdminRow, 0)
	if len(storeIDs) == 0 {
		return admins, nil
	}

	err := m.db.WithContext(ctx).
		Table("users").
		Select("users.store_id, users.id AS user_id, users.wechat_open_id AS open_id").
		Joins("JOIN roles ON roles.id = users.role_id").
		Where("users.store_id IN ?", storeIDs).
		Where("users.status = ? AND roles.status = ?", 1, 1).
		Where("roles.code IN ?", []string{model.RoleCodeStoreAdmin, model.RoleCodeAdmin}).
		Where("users.wechat_open_id IS NOT NULL AND TRIM(users.wechat_open_id) <> ''").
		Order("users.store_id ASC, users.id ASC").
		Scan(&admins).Error
	return admins, err
}
