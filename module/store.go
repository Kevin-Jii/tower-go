package module

import (
	"errors"
	"strings"

	"github.com/Kevin-Jii/tower-go/model"
	updatesPkg "github.com/Kevin-Jii/tower-go/utils/updates"

	"gorm.io/gorm"
)

type StoreModule struct {
	db *gorm.DB
}

func NewStoreModule(db *gorm.DB) *StoreModule {
	return &StoreModule{db: db}
}

// GetDB 返回底层数据库实例
func (m *StoreModule) GetDB() *gorm.DB {
	return m.db
}

// Create 创建门店
func (m *StoreModule) Create(store *model.Store) error {
	return m.db.Create(store).Error
}

// GetByID 根据ID获取门店
func (m *StoreModule) GetByID(id uint) (*model.Store, error) {
	var store model.Store
	if err := m.db.Preload("ThirdPartyAccount").First(&store, id).Error; err != nil {
		return nil, err
	}
	return &store, nil
}

// GetIDByStoreCode 按门店编码（store_code）解析门店主键
func (m *StoreModule) GetIDByStoreCode(code string) (uint, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return 0, errors.New("门店编码不能为空")
	}
	var store model.Store
	if err := m.db.Where("store_code = ?", code).First(&store).Error; err != nil {
		return 0, err
	}
	return store.ID, nil
}

// List 获取门店列表（全部数据，不分页；排除系统总部门店 StoreCodeHQ，不返回给前端）
func (m *StoreModule) List() ([]*model.Store, int64, error) {
	var stores []*model.Store
	var total int64

	cond := m.db.Model(&model.Store{}).Where("store_code IS NULL OR store_code <> ?", model.StoreCodeHQ)
	if err := cond.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := m.db.Model(&model.Store{}).Where("store_code IS NULL OR store_code <> ?", model.StoreCodeHQ).
		Preload("ThirdPartyAccount").Find(&stores).Error; err != nil {
		return nil, 0, err
	}

	return stores, total, nil
}

// Update 更新门店信息
func (m *StoreModule) Update(store *model.Store) error {
	return m.db.Save(store).Error
}

// UpdateByID 根据ID更新门店信息（动态更新，避免整行覆盖）
func (m *StoreModule) UpdateByID(id uint, req *model.UpdateStoreReq) error {
	updateMap := updatesPkg.BuildUpdatesFromReq(req)
	if len(updateMap) == 0 {
		return nil
	}
	return m.db.Model(&model.Store{}).Where("id = ?", id).Updates(updateMap).Error
}

// Delete 删除门店
func (m *StoreModule) Delete(id uint) error {
	return m.db.Delete(&model.Store{}, id).Error
}

// BindThirdPartyAccount 绑定门店第三方账号（accountID=nil 表示解绑）
func (m *StoreModule) BindThirdPartyAccount(storeID uint, accountID *uint) error {
	if accountID != nil {
		var count int64
		if err := m.db.Model(&model.Store{}).
			Where("third_party_account_id = ? AND id <> ?", *accountID, storeID).
			Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errors.New("该第三方账号已绑定其他门店")
		}
	}
	return m.db.Model(&model.Store{}).Where("id = ?", storeID).Updates(map[string]interface{}{
		"third_party_account_id": accountID,
	}).Error
}
