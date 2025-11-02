package module

import (
	"tower-go/model"

	"gorm.io/gorm"
)

type RoleMenuModule struct {
	db *gorm.DB
}

func NewRoleMenuModule(db *gorm.DB) *RoleMenuModule {
	return &RoleMenuModule{db: db}
}

// AssignMenusToRole 为角色分配菜单（覆盖式）
func (m *RoleMenuModule) AssignMenusToRole(roleID uint, menuIDs []uint) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		// 1. 删除该角色的所有菜单关联
		if err := tx.Where("role_id = ?", roleID).Delete(&model.RoleMenu{}).Error; err != nil {
			return err
		}

		// 2. 批量插入新的菜单关联
		if len(menuIDs) > 0 {
			roleMenus := make([]model.RoleMenu, len(menuIDs))
			for i, menuID := range menuIDs {
				roleMenus[i] = model.RoleMenu{
					RoleID: roleID,
					MenuID: menuID,
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

// DeleteByRoleID 删除角色的所有菜单关联
func (m *RoleMenuModule) DeleteByRoleID(roleID uint) error {
	return m.db.Where("role_id = ?", roleID).Delete(&model.RoleMenu{}).Error
}

// DeleteByMenuID 删除菜单的所有角色关联
func (m *RoleMenuModule) DeleteByMenuID(menuID uint) error {
	return m.db.Where("menu_id = ?", menuID).Delete(&model.RoleMenu{}).Error
}
