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
	// B2B 供货单独以供货单为准：兼容尚未生成记账单的历史订单，同时避免新订单与关联记账单重复统计。
	err := m.db.WithContext(ctx).Raw(`
SELECT revenue.store_id,
       COALESCE(SUM(revenue.total_amount), 0) AS total_amount,
       COUNT(*) AS order_count
FROM (
    SELECT store_id, total_amount
    FROM store_accounts
    WHERE deleted_at IS NULL
      AND account_date = ?
      AND is_canceled = ?
      AND source_type <> ?
    UNION ALL
    SELECT store_id, total_amount
    FROM b2b_supply_orders
    WHERE deleted_at IS NULL
      AND order_date = ?
      AND delivery_status <> ?
) AS revenue
GROUP BY revenue.store_id`,
		businessDate, false, model.StoreAccountSourceB2BSupplyOrder,
		businessDate, model.B2BDeliveryCancel,
	).Scan(&summaries).Error
	return summaries, err
}

func (m *DailyTurnoverModule) ListChannels(ctx context.Context, businessDate string) ([]DailyTurnoverChannelRow, error) {
	channels := make([]DailyTurnoverChannelRow, 0)
	err := m.db.WithContext(ctx).Raw(`
SELECT revenue.store_id,
       revenue.channel,
       COALESCE(SUM(revenue.total_amount), 0) AS amount,
       COUNT(*) AS order_count
FROM (
    SELECT store_id, channel, total_amount
    FROM store_accounts
    WHERE deleted_at IS NULL
      AND account_date = ?
      AND is_canceled = ?
      AND source_type <> ?
    UNION ALL
    SELECT store_id, ? AS channel, total_amount
    FROM b2b_supply_orders
    WHERE deleted_at IS NULL
      AND order_date = ?
      AND delivery_status <> ?
) AS revenue
GROUP BY revenue.store_id, revenue.channel
ORDER BY revenue.store_id ASC, amount DESC`,
		businessDate, false, model.StoreAccountSourceB2BSupplyOrder,
		"B2B供货", businessDate, model.B2BDeliveryCancel,
	).Scan(&channels).Error
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
