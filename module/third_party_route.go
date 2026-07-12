package module

import (
	"time"

	"github.com/Kevin-Jii/tower-go/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ThirdPartyRouteModule struct {
	db *gorm.DB
}

func NewThirdPartyRouteModule(db *gorm.DB) *ThirdPartyRouteModule {
	return &ThirdPartyRouteModule{db: db}
}

func (m *ThirdPartyRouteModule) List() ([]*model.ThirdPartyRoute, error) {
	var rows []*model.ThirdPartyRoute
	err := m.db.Preload("Stores").Preload("Stores.Store").Order("id DESC").Find(&rows).Error
	return rows, err
}

func (m *ThirdPartyRouteModule) GetByID(id uint) (*model.ThirdPartyRoute, error) {
	var row model.ThirdPartyRoute
	if err := m.db.Preload("Stores").Preload("Stores.Store").First(&row, id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (m *ThirdPartyRouteModule) Create(row *model.ThirdPartyRoute) error {
	return m.db.Create(row).Error
}

func (m *ThirdPartyRouteModule) Update(id uint, updates map[string]interface{}) error {
	return m.db.Model(&model.ThirdPartyRoute{}).Where("id = ?", id).Updates(updates).Error
}

func (m *ThirdPartyRouteModule) Delete(id uint) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("route_id = ?", id).Delete(&model.ThirdPartyRouteStore{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.ThirdPartyRoute{}, id).Error
	})
}

func (m *ThirdPartyRouteModule) ReplaceStores(routeID uint, storeIDs []uint) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("route_id = ?", routeID).Delete(&model.ThirdPartyRouteStore{}).Error; err != nil {
			return err
		}
		if len(storeIDs) == 0 {
			return nil
		}
		rows := make([]model.ThirdPartyRouteStore, 0, len(storeIDs))
		for i, storeID := range storeIDs {
			rows = append(rows, model.ThirdPartyRouteStore{
				RouteID: routeID,
				StoreID: storeID,
				Sort:    i + 1,
			})
		}
		return tx.Create(&rows).Error
	})
}

func (m *ThirdPartyRouteModule) SaveLogisticsSheet(row *model.ThirdPartyLogisticsSheet) error {
	row.UpdatedAt = time.Now()
	return m.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "route_id"}, {Name: "start_date"}, {Name: "end_date"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"sheet_date", "headers_json", "rows_json", "products_json", "updated_at",
		}),
	}).Create(row).Error
}

// GetLogisticsSheet 按路线 + 导入日期区间取单条（用于保存后回读）
func (m *ThirdPartyRouteModule) GetLogisticsSheet(routeID uint, startDate, endDate string) (*model.ThirdPartyLogisticsSheet, error) {
	var row model.ThirdPartyLogisticsSheet
	if err := m.db.Where("route_id = ? AND start_date = ? AND end_date = ?", routeID, startDate, endDate).
		First(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (m *ThirdPartyRouteModule) ListLogisticsSheets(routeID uint) ([]*model.ThirdPartyLogisticsSheet, error) {
	rows := make([]*model.ThirdPartyLogisticsSheet, 0)
	if err := m.db.Model(&model.ThirdPartyLogisticsSheet{}).
		Where("route_id = ?", routeID).
		Order("updated_at DESC, id DESC").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}
