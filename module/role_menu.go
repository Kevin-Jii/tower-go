package module

import (
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/utils/cache"

	"gorm.io/gorm"
)

type RoleMenuModule struct {
	db *gorm.DB
}

func NewRoleMenuModule(db *gorm.DB) *RoleMenuModule {
	return &RoleMenuModule{db: db}
}

// AssignMenusToRole 为角色分配菜单（覆盖式，支持权限位）
func (m *RoleMenuModule) AssignMenusToRole(roleID uint, menuIDs []uint, perms map[uint]uint8) error {
	// 分配完成后清除缓存
	defer cache.InvalidateMenuCache()
	
	return m.db.Transaction(func(tx *gorm.DB) error {
		// 1. 删除该角色的所有菜单关联
		if err := tx.Where("role_id = ?", roleID).Delete(&model.RoleMenu{}).Error; err != nil {
			return err
		}

		// 2. 批量插入新的菜单关联（带权限位）
		if len(menuIDs) > 0 {
			roleMenus := make([]model.RoleMenu, len(menuIDs))
			for i, menuID := range menuIDs {
				// 获取该菜单的权限位，默认为全部权限
				perm := model.PermAll
				if perms != nil {
					if p, ok := perms[menuID]; ok {
						perm = p
					}
				}
				roleMenus[i] = model.RoleMenu{
					RoleID:      roleID,
					MenuID:      menuID,
					Permissions: perm,
				}
			}
			if err := tx.Create(&roleMenus).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// GetMenuIDsByRoleID 获取角色的所有菜单ID
func (m *RoleMenuModule) GetMenuIDsByRoleID(roleID uint) ([]uint, error) {
	var menuIDs []uint
	err := m.db.Model(&model.RoleMenu{}).
		Where("role_id = ?", roleID).
		Pluck("menu_id", &menuIDs).Error
	return menuIDs, err
}

// GetMenuPermissionsByRoleID 获取角色的菜单权限映射
func (m *RoleMenuModule) GetMenuPermissionsByRoleID(roleID uint) (map[uint]uint8, error) {
	var roleMenus []model.RoleMenu
	err := m.db.Where("role_id = ?", roleID).Find(&roleMenus).Error
	if err != nil {
		return nil, err
	}

	perms := make(map[uint]uint8)
	for _, rm := range roleMenus {
		perms[rm.MenuID] = rm.Permissions
	}
	return perms, nil
}

// DeleteByRoleID 删除角色的所有菜单关联
func (m *RoleMenuModule) DeleteByRoleID(roleID uint) error {
	return m.db.Where("role_id = ?", roleID).Delete(&model.RoleMenu{}).Error
}

// DeleteByMenuID 删除菜单的所有角色关联
func (m *RoleMenuModule) DeleteByMenuID(menuID uint) error {
	return m.db.Where("menu_id = ?", menuID).Delete(&model.RoleMenu{}).Error
}
