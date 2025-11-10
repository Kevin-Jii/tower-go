package module

import (
	"fmt"
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/utils/cache"
	updatesPkg "github.com/Kevin-Jii/tower-go/utils/updates"

	"gorm.io/gorm"
)

type MenuModule struct {
	db *gorm.DB
}

func NewMenuModule(db *gorm.DB) *MenuModule {
	return &MenuModule{db: db}
}

// Create 创建菜单（清除缓存）
func (m *MenuModule) Create(menu *model.Menu) error {
	err := m.db.Create(menu).Error
	if err == nil {
		cache.InvalidateMenuCache()
	}
	return err
}

// GetByID 根据ID获取菜单
func (m *MenuModule) GetByID(id uint) (*model.Menu, error) {
	var menu model.Menu
	if err := m.db.First(&menu, id).Error; err != nil {
		return nil, err
	}
	return &menu, nil
}

// List 获取所有菜单（带缓存）
func (m *MenuModule) List() ([]*model.Menu, error) {
	var menus []*model.Menu

	// 尝试从缓存获取
	cacheKey := cache.CacheKeyMenuTree
	err := cache.CacheGet(cacheKey, &menus)
	if err == nil && len(menus) > 0 {
		return menus, nil
	}

	// 缓存未命中，从数据库查询
	if err := m.db.Order("sort ASC, id ASC").Find(&menus).Error; err != nil {
		return nil, err
	}

	// 保存到缓存
	cache.CacheSet(cacheKey, menus, cache.MenuTreeTTL)
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

// Update 更新菜单（清除缓存）
func (m *MenuModule) Update(id uint, req *model.UpdateMenuReq) error {
	updateMap := updatesPkg.BuildUpdatesFromReq(req)
	if len(updateMap) == 0 {
		return nil
	}
	err := m.db.Model(&model.Menu{}).Where("id = ?", id).Updates(updateMap).Error
	if err == nil {
		cache.InvalidateMenuCache()
	}
	return err
}

// Delete 删除菜单（清除缓存）
func (m *MenuModule) Delete(id uint) error {
	err := m.db.Delete(&model.Menu{}, id).Error
	if err == nil {
		cache.InvalidateMenuCache()
	}
	return err
}

// GetMenusByRoleID 根据角色ID获取菜单列表（带缓存）
func (m *MenuModule) GetMenusByRoleID(roleID uint) ([]*model.Menu, error) {
	var menus []*model.Menu

	// 尝试从缓存获取
	cacheKey := fmt.Sprintf(cache.CacheKeyRoleMenus, roleID)
	err := cache.CacheGet(cacheKey, &menus)
	if err == nil && len(menus) > 0 {
		return menus, nil
	}

	// 缓存未命中，从数据库查询
	err = m.db.Table("menus").
		Joins("INNER JOIN role_menus ON menus.id = role_menus.menu_id").
		Where("role_menus.role_id = ? AND menus.status = 1", roleID).
		Order("menus.sort ASC").
		Find(&menus).Error

	if err != nil {
		return nil, err
	}

	// 保存到缓存
	cache.CacheSet(cacheKey, menus, cache.PermissionsTTL)
	return menus, err
}

// GetMenusByStoreAndRole 根据门店ID和角色ID获取菜单（门店定制权限，带缓存）
func (m *MenuModule) GetMenusByStoreAndRole(storeID uint, roleID uint) ([]*model.Menu, error) {
	var menus []*model.Menu

	// 尝试从缓存获取
	cacheKey := fmt.Sprintf(cache.CacheKeyStoreRoleMenus, storeID, roleID)
	err := cache.CacheGet(cacheKey, &menus)
	if err == nil && len(menus) > 0 {
		return menus, nil
	}

	// 缓存未命中，优先查询门店定制权限
	err = m.db.Table("menus").
		Joins("INNER JOIN store_role_menus ON menus.id = store_role_menus.menu_id").
		Where("store_role_menus.store_id = ? AND store_role_menus.role_id = ? AND menus.status = 1", storeID, roleID).
		Order("menus.sort ASC").
		Find(&menus).Error

	// 如果门店没有定制权限，则使用角色默认权限
	if err == nil && len(menus) == 0 {
		return m.GetMenusByRoleID(roleID)
	}

	if err != nil {
		return nil, err
	}

	// 保存到缓存
	cache.CacheSet(cacheKey, menus, cache.PermissionsTTL)
	return menus, err
}
