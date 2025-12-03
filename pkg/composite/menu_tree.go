package composite

import (
	"sort"

	"github.com/Kevin-Jii/tower-go/model"
)

// MenuComponent 菜单组件接口（组合模式核心）
type MenuComponent interface {
	GetID() uint
	GetParentID() uint
	GetName() string
	GetType() int
	IsLeaf() bool
	Add(child MenuComponent)
	Remove(id uint)
	GetChild(id uint) MenuComponent
	GetChildren() []MenuComponent
	Accept(visitor MenuVisitor)
	ToMenu() *model.Menu
}

// MenuVisitor 菜单访问者接口（配合访问者模式）
type MenuVisitor interface {
	VisitDirectory(node *MenuNode)
	VisitMenu(node *MenuNode)
	VisitButton(node *MenuNode)
}

// MenuNode 菜单节点（组合模式的具体实现）
type MenuNode struct {
	menu     *model.Menu
	children []MenuComponent
}

// NewMenuNode 创建菜单节点
func NewMenuNode(menu *model.Menu) *MenuNode {
	return &MenuNode{
		menu:     menu,
		children: make([]MenuComponent, 0),
	}
}

// GetID 获取ID
func (n *MenuNode) GetID() uint {
	return n.menu.ID
}

// GetParentID 获取父ID
func (n *MenuNode) GetParentID() uint {
	return n.menu.ParentID
}

// GetName 获取名称
func (n *MenuNode) GetName() string {
	return n.menu.Name
}

// GetType 获取类型
func (n *MenuNode) GetType() int {
	return n.menu.Type
}

// IsLeaf 是否叶子节点
func (n *MenuNode) IsLeaf() bool {
	return len(n.children) == 0
}

// Add 添加子节点
func (n *MenuNode) Add(child MenuComponent) {
	n.children = append(n.children, child)
}

// Remove 移除子节点
func (n *MenuNode) Remove(id uint) {
	for i, child := range n.children {
		if child.GetID() == id {
			n.children = append(n.children[:i], n.children[i+1:]...)
			return
		}
	}
}

// GetChild 获取指定子节点
func (n *MenuNode) GetChild(id uint) MenuComponent {
	for _, child := range n.children {
		if child.GetID() == id {
			return child
		}
	}
	return nil
}

// GetChildren 获取所有子节点
func (n *MenuNode) GetChildren() []MenuComponent {
	return n.children
}

// GetMenu 获取原始菜单数据
func (n *MenuNode) GetMenu() *model.Menu {
	return n.menu
}

// Accept 接受访问者
func (n *MenuNode) Accept(visitor MenuVisitor) {
	switch n.menu.Type {
	case 1: // 目录
		visitor.VisitDirectory(n)
	case 2: // 菜单
		visitor.VisitMenu(n)
	case 3: // 按钮
		visitor.VisitButton(n)
	}

	// 递归访问子节点
	for _, child := range n.children {
		child.Accept(visitor)
	}
}

// ToMenu 转换为 model.Menu（包含子节点）
func (n *MenuNode) ToMenu() *model.Menu {
	menu := *n.menu
	menu.Children = make([]*model.Menu, 0, len(n.children))

	for _, child := range n.children {
		menu.Children = append(menu.Children, child.ToMenu())
	}

	return &menu
}

// MenuTree 菜单树（组合模式的容器）
type MenuTree struct {
	roots    []MenuComponent
	nodeMap  map[uint]MenuComponent
	sortFunc func(a, b MenuComponent) bool
}

// NewMenuTree 创建菜单树
func NewMenuTree() *MenuTree {
	return &MenuTree{
		roots:   make([]MenuComponent, 0),
		nodeMap: make(map[uint]MenuComponent),
		sortFunc: func(a, b MenuComponent) bool {
			// 默认按 Sort 字段排序
			aNode := a.(*MenuNode)
			bNode := b.(*MenuNode)
			return aNode.menu.Sort < bNode.menu.Sort
		},
	}
}

// SetSortFunc 设置排序函数
func (t *MenuTree) SetSortFunc(fn func(a, b MenuComponent) bool) {
	t.sortFunc = fn
}

// Build 从菜单列表构建树（O(n) 复杂度）
func (t *MenuTree) Build(menus []*model.Menu) *MenuTree {
	if len(menus) == 0 {
		return t
	}

	// 第一步：创建所有节点
	for _, menu := range menus {
		node := NewMenuNode(menu)
		t.nodeMap[menu.ID] = node
	}

	// 第二步：建立父子关系
	for _, menu := range menus {
		node := t.nodeMap[menu.ID]
		if menu.ParentID == 0 {
			t.roots = append(t.roots, node)
		} else if parent, ok := t.nodeMap[menu.ParentID]; ok {
			parent.Add(node)
		}
	}

	// 第三步：排序
	t.sortAll()

	return t
}

// sortAll 递归排序所有节点
func (t *MenuTree) sortAll() {
	t.sortNodes(t.roots)
	for _, root := range t.roots {
		t.sortChildren(root)
	}
}

func (t *MenuTree) sortNodes(nodes []MenuComponent) {
	sort.Slice(nodes, func(i, j int) bool {
		return t.sortFunc(nodes[i], nodes[j])
	})
}

func (t *MenuTree) sortChildren(node MenuComponent) {
	children := node.GetChildren()
	if len(children) > 0 {
		t.sortNodes(children)
		for _, child := range children {
			t.sortChildren(child)
		}
	}
}

// GetRoots 获取根节点
func (t *MenuTree) GetRoots() []MenuComponent {
	return t.roots
}

// GetNode 获取指定节点
func (t *MenuTree) GetNode(id uint) MenuComponent {
	return t.nodeMap[id]
}

// ToMenus 转换为 model.Menu 切片
func (t *MenuTree) ToMenus() []*model.Menu {
	result := make([]*model.Menu, 0, len(t.roots))
	for _, root := range t.roots {
		result = append(result, root.ToMenu())
	}
	return result
}

// Filter 过滤菜单树（返回新树）
func (t *MenuTree) Filter(predicate func(MenuComponent) bool) *MenuTree {
	newTree := NewMenuTree()
	newTree.sortFunc = t.sortFunc

	// 收集符合条件的菜单ID及其所有祖先
	keepIDs := make(map[uint]bool)
	for _, node := range t.nodeMap {
		if predicate(node) {
			// 标记该节点及其所有祖先
			t.markAncestors(node, keepIDs)
		}
	}

	// 重建树
	var filteredMenus []*model.Menu
	for id := range keepIDs {
		if node, ok := t.nodeMap[id]; ok {
			filteredMenus = append(filteredMenus, node.(*MenuNode).menu)
		}
	}

	return newTree.Build(filteredMenus)
}

func (t *MenuTree) markAncestors(node MenuComponent, keepIDs map[uint]bool) {
	keepIDs[node.GetID()] = true
	if node.GetParentID() != 0 {
		if parent, ok := t.nodeMap[node.GetParentID()]; ok {
			t.markAncestors(parent, keepIDs)
		}
	}
}

// Walk 遍历树（深度优先）
func (t *MenuTree) Walk(fn func(node MenuComponent, depth int) bool) {
	for _, root := range t.roots {
		if !t.walkNode(root, 0, fn) {
			return
		}
	}
}

func (t *MenuTree) walkNode(node MenuComponent, depth int, fn func(MenuComponent, int) bool) bool {
	if !fn(node, depth) {
		return false
	}
	for _, child := range node.GetChildren() {
		if !t.walkNode(child, depth+1, fn) {
			return false
		}
	}
	return true
}

// FindByPath 根据路径查找节点
func (t *MenuTree) FindByPath(path string) MenuComponent {
	for _, node := range t.nodeMap {
		if node.(*MenuNode).menu.Path == path {
			return node
		}
	}
	return nil
}

// FindByPermission 根据权限标识查找节点
func (t *MenuTree) FindByPermission(permission string) MenuComponent {
	for _, node := range t.nodeMap {
		if node.(*MenuNode).menu.Permission == permission {
			return node
		}
	}
	return nil
}

// GetAllPermissions 获取所有权限标识
func (t *MenuTree) GetAllPermissions() []string {
	var permissions []string
	for _, node := range t.nodeMap {
		perm := node.(*MenuNode).menu.Permission
		if perm != "" {
			permissions = append(permissions, perm)
		}
	}
	return permissions
}

// GetVisibleMenus 获取可见菜单
func (t *MenuTree) GetVisibleMenus() *MenuTree {
	return t.Filter(func(node MenuComponent) bool {
		return node.(*MenuNode).menu.Visible == 1
	})
}

// GetEnabledMenus 获取启用的菜单
func (t *MenuTree) GetEnabledMenus() *MenuTree {
	return t.Filter(func(node MenuComponent) bool {
		return node.(*MenuNode).menu.Status == 1
	})
}

// Count 统计节点数量
func (t *MenuTree) Count() int {
	return len(t.nodeMap)
}

// Depth 获取树的最大深度
func (t *MenuTree) Depth() int {
	maxDepth := 0
	t.Walk(func(node MenuComponent, depth int) bool {
		if depth > maxDepth {
			maxDepth = depth
		}
		return true
	})
	return maxDepth
}
