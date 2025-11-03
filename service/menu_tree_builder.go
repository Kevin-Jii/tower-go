package service

import (
	"sync"
	"time"
	"tower-go/model"
)

// MenuTreeNode 菜单树节点（带缓存辅助信息）
type MenuTreeNode struct {
	*model.Menu
	Children []*MenuTreeNode `json:"children,omitempty"`
}

// MenuTreeCache 菜单树缓存
type MenuTreeCache struct {
	mu    sync.RWMutex
	cache map[string]*cacheEntry
	ttl   time.Duration
}

type cacheEntry struct {
	tree      []*model.Menu
	expiresAt time.Time
}

var globalMenuTreeCache *MenuTreeCache

func init() {
	globalMenuTreeCache = &MenuTreeCache{
		cache: make(map[string]*cacheEntry),
		ttl:   5 * time.Minute, // 菜单树缓存 5 分钟
	}
}

// GetMenuTreeCache 获取全局菜单树缓存
func GetMenuTreeCache() *MenuTreeCache {
	return globalMenuTreeCache
}

// Get 从缓存获取菜单树
func (c *MenuTreeCache) Get(key string) ([]*model.Menu, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.cache[key]
	if !ok {
		return nil, false
	}

	// 检查是否过期
	if time.Now().After(entry.expiresAt) {
		return nil, false
	}

	return entry.tree, true
}

// Set 设置缓存
func (c *MenuTreeCache) Set(key string, tree []*model.Menu) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[key] = &cacheEntry{
		tree:      tree,
		expiresAt: time.Now().Add(c.ttl),
	}
}

// Invalidate 使缓存失效
func (c *MenuTreeCache) Invalidate(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.cache, key)
}

// Clear 清空所有缓存
func (c *MenuTreeCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache = make(map[string]*cacheEntry)
}

// buildMenuTreeOptimized 优化的菜单树构建算法 O(n)
// 使用 Map 索引避免嵌套循环，一次遍历完成
func (s *MenuService) buildMenuTreeOptimized(menus []*model.Menu, parentID uint) []*model.Menu {
	if len(menus) == 0 {
		return []*model.Menu{}
	}

	// 第一步：创建 ID -> Menu 的映射表
	menuMap := make(map[uint]*model.Menu, len(menus))
	for _, menu := range menus {
		// 创建副本避免修改原始数据
		menuCopy := *menu
		menuCopy.Children = []*model.Menu{}
		menuMap[menu.ID] = &menuCopy
	}

	// 第二步：构建父子关系
	roots := []*model.Menu{}
	for _, menu := range menus {
		menuNode := menuMap[menu.ID]

		if menu.ParentID == parentID {
			// 顶层节点
			roots = append(roots, menuNode)
		} else if parent, ok := menuMap[menu.ParentID]; ok {
			// 添加到父节点的 children
			parent.Children = append(parent.Children, menuNode)
		}
	}

	return roots
}

// buildMenuTreeWithCache 带缓存的菜单树构建
func (s *MenuService) buildMenuTreeWithCache(menus []*model.Menu, parentID uint, cacheKey string) []*model.Menu {
	cache := GetMenuTreeCache()

	// 尝试从缓存获取
	if cached, ok := cache.Get(cacheKey); ok {
		return cached
	}

	// 构建菜单树
	tree := s.buildMenuTreeOptimized(menus, parentID)

	// 写入缓存
	cache.Set(cacheKey, tree)

	return tree
}

// InvalidateMenuCache 使菜单缓存失效（在菜单增删改时调用）
func (s *MenuService) InvalidateMenuCache() {
	GetMenuTreeCache().Clear()
}
