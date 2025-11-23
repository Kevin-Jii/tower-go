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

// AssignMenusToStoreRole 为门店角色分配菜单（覆盖式，支持权限位）
func (m *StoreRoleMenuModule) AssignMenusToStoreRole(storeID uint, roleID uint, menuIDs []uint, perms map[uint]uint8) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		// 1. 删除该门店角色的所有菜单关联
		if err := tx.Where("store_id = ? AND role_id = ?", storeID, roleID).Delete(&model.StoreRoleMenu{}).Error; err != nil {
			return err
		}

		// 2. 批量插入新的菜单关联（带权限位）
		if len(menuIDs) > 0 {
			storeRoleMenus := make([]model.StoreRoleMenu, len(menuIDs))
			for i, menuID := range menuIDs {
				// 获取该菜单的权限位，默认为全部权限
				perm := model.PermAll
				if perms != nil {
					if p, ok := perms[menuID]; ok {
						perm = p
					}
				}
				storeRoleMenus[i] = model.StoreRoleMenu{
					StoreID:     storeID,
					RoleID:      roleID,
					MenuID:      menuID,
					Permissions: perm,
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

// GetMenuPermissionsByStoreAndRole 获取门店角色的菜单权限映射
func (m *StoreRoleMenuModule) GetMenuPermissionsByStoreAndRole(storeID uint, roleID uint) (map[uint]uint8, error) {
	var storeRoleMenus []model.StoreRoleMenu
	err := m.db.Where("store_id = ? AND role_id = ?", storeID, roleID).Find(&storeRoleMenus).Error
	if err != nil {
		return nil, err
	}

	perms := make(map[uint]uint8)
	for _, srm := range storeRoleMenus {
		perms[srm.MenuID] = srm.Permissions
	}
	return perms, nil
}

// CopyStoreMenus 复制门店的菜单权限配置（包含权限位）
func (m *StoreRoleMenuModule) CopyStoreMenus(fromStoreID uint, toStoreID uint, roleID uint) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		// 1. 查询源门店的菜单配置（包含权限位）
		var sourceMenus []model.StoreRoleMenu
		if err := tx.Where("store_id = ? AND role_id = ?", fromStoreID, roleID).
			Find(&sourceMenus).Error; err != nil {
			return err
		}

		// 2. 删除目标门店的现有配置
		if err := tx.Where("store_id = ? AND role_id = ?", toStoreID, roleID).
			Delete(&model.StoreRoleMenu{}).Error; err != nil {
			return err
		}

		// 3. 为目标门店创建新配置（复制权限位）
		if len(sourceMenus) > 0 {
			storeRoleMenus := make([]model.StoreRoleMenu, len(sourceMenus))
			for i, sm := range sourceMenus {
				storeRoleMenus[i] = model.StoreRoleMenu{
					StoreID:     toStoreID,
					RoleID:      roleID,
					MenuID:      sm.MenuID,
					Permissions: sm.Permissions,
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
