package utils

import (
	"fmt"
	"time"
)

// CacheDomain 缓存域抽象：统一 Key 格式与 TTL
type CacheDomain struct {
	KeyPattern string        // 例如 "dish:categories:store:%d"
	TTL        time.Duration // 过期时间
}

// Key 根据参数格式化实际键（目前支持单个或多个整型参数）
func (d CacheDomain) Key(args ...interface{}) string { return fmt.Sprintf(d.KeyPattern, args...) }

// GetOrSet 执行查缓存或加载数据
func (d CacheDomain) GetOrSet(dest interface{}, fetch func() (interface{}, error), args ...interface{}) error {
	key := d.Key(args...)
	return CacheGetOrSet(key, dest, d.TTL, fetch)
}

// Invalidate 失效该域指定参数的 Key（一次性）
func (d CacheDomain) Invalidate(args ...interface{}) { CacheDelete(d.Key(args...)) }

// InvalidateMulti 支持失效多个主键集合
func (d CacheDomain) InvalidateMulti(argLists ...[]interface{}) {
	if len(argLists) == 0 {
		return
	}
	keys := make([]string, 0, len(argLists))
	for _, lst := range argLists {
		keys = append(keys, d.Key(lst...))
	}
	CacheDelete(keys...)
}

// 具体域定义（按既有 TTL 需求）
var (
	DishCategoriesDomain           = CacheDomain{KeyPattern: CacheKeyDishCategories, TTL: 10 * time.Minute}
	DishCategoriesWithDishesDomain = CacheDomain{KeyPattern: CacheKeyDishCategoriesWithDishes, TTL: 5 * time.Minute}
)
