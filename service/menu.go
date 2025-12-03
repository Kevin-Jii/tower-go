package service

import (
	"errors"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
	"github.com/Kevin-Jii/tower-go/pkg/composite"
)

type MenuService struct {
	menuModule          *module.MenuModule
	roleMenuModule      *module.RoleMenuModule
	storeRoleMenuModule *module.StoreRoleMenuModule
}

func NewMenuService(
	menuModule *module.MenuModule,
	roleMenuModule *module.RoleMenuModule,
	storeRoleMenuModule *module.StoreRoleMenuModule,
) *MenuService {
	return &MenuService{
		menuModule:          menuModule,
		roleMenuModule:      roleMenuModule,
		storeRoleMenuModule: storeRoleMenuModule,
	}
}

// CreateMenu 创建菜单
func (s *MenuService) CreateMenu(req *model.CreateMenuReq) error {
	menu := &model.Menu{
		ParentID:   req.ParentID,
		Name:       req.Name,
		Title:      req.Title,
		Icon:       req.Icon,
		Path:       req.Path,
		Component:  req.Component,
		Type:       req.Type,
		Sort:       req.Sort,
		Permission: req.Permission,
		Visible:    req.Visible,
		Status:     req.Status,
		Remark:     req.Remark,
	}

	// 创建菜单后使缓存失效
	defer s.InvalidateMenuCache()

	return s.menuModule.Create(menu)
}

// GetMenu 获取菜单详情
func (s *MenuService) GetMenu(id uint) (*model.Menu, error) {
	return s.menuModule.GetByID(id)
}

// ListMenus 获取所有菜单（平铺列表）
func (s *MenuService) ListMenus() ([]*model.Menu, error) {
	return s.menuModule.List()
}

// GetMenuTree 获取菜单树形结构
func (s *MenuService) GetMenuTree() ([]*model.Menu, error) {
	// 获取所有菜单
	menus, err := s.menuModule.List()
	if err != nil {
		return nil, err
	}

	// 使用组合模式构建树
	tree := composite.NewMenuTree().Build(menus)
	return tree.ToMenus(), nil
}

// GetMenuTreeComposite 获取组合模式的菜单树（提供更多操作能力）
func (s *MenuService) GetMenuTreeComposite() (*composite.MenuTree, error) {
	menus, err := s.menuModule.List()
	if err != nil {
		return nil, err
	}
	return composite.NewMenuTree().Build(menus), nil
}

// buildMenuTree 构建树形菜单结构（保留旧方法用于兼容，但内部调用优化版本）
func (s *MenuService) buildMenuTree(menus []*model.Menu, parentID uint) []*model.Menu {
	tree := composite.NewMenuTree().Build(menus)
	return tree.ToMenus()
}

// UpdateMenu 更新菜单
func (s *MenuService) UpdateMenu(id uint, req *model.UpdateMenuReq) error {
	// 检查菜单是否存在
	_, err := s.menuModule.GetByID(id)
	if err != nil {
		return err
	}

	// 如果修改了父级，检查是否会形成循环引用
	if req.ParentID != nil && *req.ParentID != 0 {
		if err := s.checkCircularReference(id, *req.ParentID); err != nil {
			return err
		}
	}

	// 更新菜单后使缓存失效
	defer s.InvalidateMenuCache()

	return s.menuModule.Update(id, req)
}

// checkCircularReference 检查是否会形成循环引用
func (s *MenuService) checkCircularReference(menuID uint, parentID uint) error {
	if menuID == parentID {
		return errors.New("父级菜单不能是自己")
	}

	// 向上查找，检查parentID的所有祖先节点
	currentID := parentID
	for currentID != 0 {
		if currentID == menuID {
			return errors.New("不能将菜单移动到自己的子菜单下")
		}

		menu, err := s.menuModule.GetByID(currentID)
		if err != nil {
			return err
		}
		currentID = menu.ParentID
	}

	return nil
}

// DeleteMenu 删除菜单
func (s *MenuService) DeleteMenu(id uint) error {
	// 检查是否有子菜单
	children, err := s.menuModule.ListByParentID(id)
	if err != nil {
		return err
	}
	if len(children) > 0 {
		return errors.New("存在子菜单，无法删除")
	}

	// 删除菜单相关的权限关联
	_ = s.roleMenuModule.DeleteByMenuID(id)
	_ = s.storeRoleMenuModule.DeleteByMenuID(id)

	// 删除菜单后使缓存失效
	defer s.InvalidateMenuCache()

	return s.menuModule.Delete(id)
}

// AssignMenusToRole 为角色分配菜单（支持权限位）
func (s *MenuService) AssignMenusToRole(req *model.AssignMenusToRoleReq) error {
	return s.roleMenuModule.AssignMenusToRole(req.RoleID, req.MenuIDs, req.Perms)
}

// GetRoleMenus 获取角色的菜单列表
func (s *MenuService) GetRoleMenus(roleID uint) ([]*model.Menu, error) {
	return s.menuModule.GetMenusByRoleID(roleID)
}

// GetRoleMenuTree 获取角色的菜单树
func (s *MenuService) GetRoleMenuTree(roleID uint) ([]*model.Menu, error) {
	menus, err := s.menuModule.GetMenusByRoleID(roleID)
	if err != nil {
		return nil, err
	}

	tree := composite.NewMenuTree().Build(menus)
	return tree.ToMenus(), nil
}

// GetRoleMenuIDs 获取角色的所有菜单ID（用于前端回显）
func (s *MenuService) GetRoleMenuIDs(roleID uint) ([]uint, error) {
	return s.roleMenuModule.GetMenuIDsByRoleID(roleID)
}

// AssignMenusToStoreRole 为门店角色分配菜单（支持权限位）
func (s *MenuService) AssignMenusToStoreRole(req *model.AssignStoreMenusReq) error {
	return s.storeRoleMenuModule.AssignMenusToStoreRole(req.StoreID, req.RoleID, req.MenuIDs, req.Perms)
}

// GetStoreRoleMenus 获取门店角色的菜单列表
func (s *MenuService) GetStoreRoleMenus(storeID uint, roleID uint) ([]*model.Menu, error) {
	return s.menuModule.GetMenusByStoreAndRole(storeID, roleID)
}

// GetStoreRoleMenuTree 获取门店角色的菜单树
func (s *MenuService) GetStoreRoleMenuTree(storeID uint, roleID uint) ([]*model.Menu, error) {
	menus, err := s.menuModule.GetMenusByStoreAndRole(storeID, roleID)
	if err != nil {
		return nil, err
	}

	tree := composite.NewMenuTree().Build(menus)
	return tree.ToMenus(), nil
}

// GetStoreRoleMenuIDs 获取门店角色的所有菜单ID
func (s *MenuService) GetStoreRoleMenuIDs(storeID uint, roleID uint) ([]uint, error) {
	return s.storeRoleMenuModule.GetMenuIDsByStoreAndRole(storeID, roleID)
}

// CopyStoreMenus 复制门店菜单权限
func (s *MenuService) CopyStoreMenus(req *model.CopyStoreMenusReq) error {
	return s.storeRoleMenuModule.CopyStoreMenus(req.FromStoreID, req.ToStoreID, req.RoleID)
}

// GetUserMenus 获取用户的菜单（根据用户的门店和角色）
func (s *MenuService) GetUserMenus(storeID uint, roleID uint) ([]*model.Menu, error) {
	return s.menuModule.GetMenusByStoreAndRole(storeID, roleID)
}

// GetUserMenuTree 获取用户的菜单树（用于前端渲染侧边栏）
func (s *MenuService) GetUserMenuTree(storeID uint, roleID uint) ([]*model.Menu, error) {
	menus, err := s.menuModule.GetMenusByStoreAndRole(storeID, roleID)
	if err != nil {
		return nil, err
	}

	// 使用组合模式构建树，只返回可见菜单
	tree := composite.NewMenuTree().Build(menus)
	visibleTree := tree.GetVisibleMenus()
	return visibleTree.ToMenus(), nil
}

// GetUserPermissions 获取用户的所有权限标识（用于前端按钮级权限控制）
func (s *MenuService) GetUserPermissions(storeID uint, roleID uint) ([]string, error) {
	menus, err := s.menuModule.GetMenusByStoreAndRole(storeID, roleID)
	if err != nil {
		return nil, err
	}

	// 使用组合模式收集权限
	tree := composite.NewMenuTree().Build(menus)
	return tree.GetAllPermissions(), nil
}

// GetAllPermissions 获取所有菜单权限标识（总部管理员）
func (s *MenuService) GetAllPermissions() ([]string, error) {
	menus, err := s.menuModule.List()
	if err != nil {
		return nil, err
	}

	// 使用组合模式 + 访问者模式收集权限
	tree := composite.NewMenuTree().Build(menus)
	enabledTree := tree.GetEnabledMenus()
	return enabledTree.GetAllPermissions(), nil
}

// GetRoleMenuPermissions 获取角色的菜单权限映射
func (s *MenuService) GetRoleMenuPermissions(roleID uint) (map[uint]uint8, error) {
	return s.roleMenuModule.GetMenuPermissionsByRoleID(roleID)
}

// GetStoreRoleMenuPermissions 获取门店角色的菜单权限映射
func (s *MenuService) GetStoreRoleMenuPermissions(storeID uint, roleID uint) (map[uint]uint8, error) {
	return s.storeRoleMenuModule.GetMenuPermissionsByStoreAndRole(storeID, roleID)
}
