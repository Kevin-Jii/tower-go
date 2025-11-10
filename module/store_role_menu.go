package module

import (
	"github.com/Kevin-Jii/tower-go/model"

	"gorm.io/gorm"
)

type StoreRoleMenuModule struct {
	db *gorm.DB
}

func NewStoreRoleMenuModule(db *gorm.DB) *StoreRoleMenuModule {
	return &StoreRoleMenuModule{db: db}
}

// AssignMenusToStoreRole 为门店角色分配菜单（覆盖式）
func (m *StoreRoleMenuModule) AssignMenusToStoreRole(storeID uint, roleID uint, menuIDs []uint) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		// 1. 删除该门店角色的所有菜单关联
		if err := tx.Where("store_id = ? AND role_id = ?", storeID, roleID).Delete(&model.StoreRoleMenu{}).Error; err != nil {
			return err
		}

		// 2. 批量插入新的菜单关联
		if len(menuIDs) > 0 {
			storeRoleMenus := make([]model.StoreRoleMenu, len(menuIDs))
			for i, menuID := range menuIDs {
				storeRoleMenus[i] = model.StoreRoleMenu{
					StoreID: storeID,
					RoleID:  roleID,
					MenuID:  menuID,
				}
			}
			if err := tx.Create(&storeRoleMenus).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// GetMenuIDsByStoreAndRole 获取门店角色的所有菜单ID
func (m *StoreRoleMenuModule) GetMenuIDsByStoreAndRole(storeID uint, roleID uint) ([]uint, error) {
	var menuIDs []uint
	err := m.db.Model(&model.StoreRoleMenu{}).
		Where("store_id = ? AND role_id = ?", storeID, roleID).
		Pluck("menu_id", &menuIDs).Error
	return menuIDs, err
}

// CopyStoreMenus 复制门店的菜单权限配置
func (m *StoreRoleMenuModule) CopyStoreMenus(fromStoreID uint, toStoreID uint, roleID uint) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		// 1. 查询源门店的菜单配置
		var menuIDs []uint
		if err := tx.Model(&model.StoreRoleMenu{}).
			Where("store_id = ? AND role_id = ?", fromStoreID, roleID).
			Pluck("menu_id", &menuIDs).Error; err != nil {
			return err
		}

		// 2. 删除目标门店的现有配置
		if err := tx.Where("store_id = ? AND role_id = ?", toStoreID, roleID).
			Delete(&model.StoreRoleMenu{}).Error; err != nil {
			return err
		}

		// 3. 为目标门店创建新配置
		if len(menuIDs) > 0 {
			storeRoleMenus := make([]model.StoreRoleMenu, len(menuIDs))
			for i, menuID := range menuIDs {
				storeRoleMenus[i] = model.StoreRoleMenu{
					StoreID: toStoreID,
					RoleID:  roleID,
					MenuID:  menuID,
				}
			}
			if err := tx.Create(&storeRoleMenus).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// DeleteByStoreID 删除门店的所有菜单关联
func (m *StoreRoleMenuModule) DeleteByStoreID(storeID uint) error {
	return m.db.Where("store_id = ?", storeID).Delete(&model.StoreRoleMenu{}).Error
}

// DeleteByRoleID 删除角色的所有门店菜单关联
func (m *StoreRoleMenuModule) DeleteByRoleID(roleID uint) error {
	return m.db.Where("role_id = ?", roleID).Delete(&model.StoreRoleMenu{}).Error
}

// DeleteByMenuID 删除菜单的所有门店角色关联
func (m *StoreRoleMenuModule) DeleteByMenuID(menuID uint) error {
	return m.db.Where("menu_id = ?", menuID).Delete(&model.StoreRoleMenu{}).Error
}
