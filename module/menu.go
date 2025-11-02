package module

import (
	"tower-go/model"

	"gorm.io/gorm"
)

type MenuModule struct {
	db *gorm.DB
}

func NewMenuModule(db *gorm.DB) *MenuModule {
	return &MenuModule{db: db}
}

// Create 创建菜单
func (m *MenuModule) Create(menu *model.Menu) error {
	return m.db.Create(menu).Error
}

// GetByID 根据ID获取菜单
func (m *MenuModule) GetByID(id uint) (*model.Menu, error) {
	var menu model.Menu
	if err := m.db.First(&menu, id).Error; err != nil {
		return nil, err
	}
	return &menu, nil
}

// List 获取所有菜单
func (m *MenuModule) List() ([]*model.Menu, error) {
	var menus []*model.Menu
	if err := m.db.Order("sort ASC, id ASC").Find(&menus).Error; err != nil {
		return nil, err
	}
	return menus, nil
}

// ListByParentID 根据父ID获取子菜单
func (m *MenuModule) ListByParentID(parentID uint) ([]*model.Menu, error) {
	var menus []*model.Menu
	if err := m.db.Where("parent_id = ?", parentID).Order("sort ASC").Find(&menus).Error; err != nil {
		return nil, err
	}
	return menus, nil
}

// Update 更新菜单
func (m *MenuModule) Update(id uint, req *model.UpdateMenuReq) error {
	updates := make(map[string]interface{})

	if req.ParentID != nil {
		updates["parent_id"] = *req.ParentID
	}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Icon != "" {
		updates["icon"] = req.Icon
	}
	if req.Path != "" {
		updates["path"] = req.Path
	}
	if req.Component != "" {
		updates["component"] = req.Component
	}
	if req.Type != nil {
		updates["type"] = *req.Type
	}
	if req.Sort != nil {
		updates["sort"] = *req.Sort
	}
	if req.Permission != "" {
		updates["permission"] = req.Permission
	}
	if req.Visible != nil {
		updates["visible"] = *req.Visible
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Remark != "" {
		updates["remark"] = req.Remark
	}

	return m.db.Model(&model.Menu{}).Where("id = ?", id).Updates(updates).Error
}

// Delete 删除菜单
func (m *MenuModule) Delete(id uint) error {
	return m.db.Delete(&model.Menu{}, id).Error
}

// GetMenusByRoleID 根据角色ID获取菜单列表
func (m *MenuModule) GetMenusByRoleID(roleID uint) ([]*model.Menu, error) {
	var menus []*model.Menu
	err := m.db.Table("menus").
		Joins("INNER JOIN role_menus ON menus.id = role_menus.menu_id").
		Where("role_menus.role_id = ? AND menus.status = 1", roleID).
		Order("menus.sort ASC").
		Find(&menus).Error
	return menus, err
}

// GetMenusByStoreAndRole 根据门店ID和角色ID获取菜单（门店定制权限）
func (m *MenuModule) GetMenusByStoreAndRole(storeID uint, roleID uint) ([]*model.Menu, error) {
	var menus []*model.Menu

	// 优先查询门店定制权限
	err := m.db.Table("menus").
		Joins("INNER JOIN store_role_menus ON menus.id = store_role_menus.menu_id").
		Where("store_role_menus.store_id = ? AND store_role_menus.role_id = ? AND menus.status = 1", storeID, roleID).
		Order("menus.sort ASC").
		Find(&menus).Error

	// 如果门店没有定制权限，则使用角色默认权限
	if err == nil && len(menus) == 0 {
		return m.GetMenusByRoleID(roleID)
	}

	return menus, err
}
