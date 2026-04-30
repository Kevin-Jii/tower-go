package module

import (
	"database/sql"
	"time"

	"github.com/Kevin-Jii/tower-go/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ThirdPartyOrderModule struct {
	db *gorm.DB
}

func NewThirdPartyOrderModule(db *gorm.DB) *ThirdPartyOrderModule {
	return &ThirdPartyOrderModule{db: db}
}

func (m *ThirdPartyOrderModule) UpsertBatch(rows []model.ThirdPartyOrder) error {
	if len(rows) == 0 {
		return nil
	}
	return m.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "order_no"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"account_id",
			"platform_name",
			"place_time",
			"place_date",
			"order_trade_status",
			"status_name",
			"pay_amount",
			"total_amount",
			"total_item_num",
			"raw_json",
			"synced_at",
			"updated_at",
		}),
	}).Create(&rows).Error
}

// GetLatestPlaceTimeByAccount 获取账号已同步订单的最新提报时间
func (m *ThirdPartyOrderModule) GetLatestPlaceTimeByAccount(accountID uint) (*time.Time, error) {
	var t sql.NullTime
	err := m.db.Model(&model.ThirdPartyOrder{}).
		Where("account_id = ?", accountID).
		Select("MAX(place_time)").
		Scan(&t).Error
	if err != nil {
		return nil, err
	}
	if !t.Valid {
		return nil, nil
	}
	v := t.Time
	return &v, nil
}

// ListByAccount 分页查询账号已同步订单（按提报时间倒序）
func (m *ThirdPartyOrderModule) ListByAccount(accountID uint, page, pageSize int) ([]*model.ThirdPartyOrder, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	var total int64
	query := m.db.Model(&model.ThirdPartyOrder{}).Where("account_id = ?", accountID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	rows := make([]*model.ThirdPartyOrder, 0)
	offset := (page - 1) * pageSize
	if err := query.Order("place_time DESC, id DESC").Offset(offset).Limit(pageSize).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func (m *ThirdPartyOrderModule) ListByAccountsAndDate(accountIDs []uint, placeDate string) ([]*model.ThirdPartyOrder, error) {
	if len(accountIDs) == 0 {
		return []*model.ThirdPartyOrder{}, nil
	}
	rows := make([]*model.ThirdPartyOrder, 0)
	if err := m.db.Model(&model.ThirdPartyOrder{}).
		Where("account_id IN ? AND place_date = ?", accountIDs, placeDate).
		Order("place_time DESC, id DESC").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (m *ThirdPartyOrderModule) ListByAccountsAndDateRange(accountIDs []uint, startDate, endDate string) ([]*model.ThirdPartyOrder, error) {
	if len(accountIDs) == 0 {
		return []*model.ThirdPartyOrder{}, nil
	}
	rows := make([]*model.ThirdPartyOrder, 0)
	if err := m.db.Model(&model.ThirdPartyOrder{}).
		Where("account_id IN ? AND place_date >= ? AND place_date <= ?", accountIDs, startDate, endDate).
		Order("place_date DESC, place_time DESC, id DESC").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}
