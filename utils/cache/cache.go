package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Kevin-Jii/tower-go/utils/logging"
	"github.com/Kevin-Jii/tower-go/utils/redis"

	"go.uber.org/zap"
)

// CacheKey 缓存键常量
const (
	CacheKeyRole                     = "role:%d"                              // 角色缓存 role:1
	CacheKeyRoleMenus                = "role:menus:%d"                        // 角色菜单缓存 role:menus:1
	CacheKeyStoreRoleMenus           = "store:role:menus:%d:%d"               // 门店角色菜单缓存 store:role:menus:1:2
	CacheKeyMenuTree                 = "menu:tree"                            // 菜单树缓存
	CacheKeyUserPermissions          = "user:permissions:%d"                  // 用户权限缓存 user:permissions:1
	CacheKeyUser                     = "user:%d"                              // 用户缓存 user:1
	CacheKeyStore                    = "store:%d"                             // 门店缓存 store:1
	CacheKeyDish                     = "dish:%d"                              // 菜品缓存 dish:1
	CacheKeyDishList                 = "dish:list:store:%d"                   // 门店菜品列表缓存 dish:list:store:1
	CacheKeyDishCategories           = "dish:categories:store:%d"             // 门店分类列表
	CacheKeyDishCategoriesWithDishes = "dish:categories:with:dishes:store:%d" // 门店分类+菜品聚合
)

// 默认缓存过期时间
const (
	DefaultCacheTTL = 1 * time.Hour    // 默认1小时
	ShortCacheTTL   = 5 * time.Minute  // 短期5分钟
	LongCacheTTL    = 24 * time.Hour   // 长期24小时
	MenuTreeTTL     = 30 * time.Minute // 菜单树30分钟
	PermissionsTTL  = 15 * time.Minute // 权限15分钟
)

// CacheGet 从缓存获取数据（自动 JSON 反序列化）
func CacheGet(key string, dest interface{}) error {
	if !redis.IsRedisEnabled() {
		return fmt.Errorf("redis not enabled")
	}

	ctx := context.Background()
	val, err := redis.RedisClient.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(val), dest); err != nil {
		logging.LogError("缓存反序列化失败", zap.String("key", key), zap.Error(err))
		return err
	}

	logging.LogDebug("缓存命中", zap.String("key", key))
	return nil
}

// CacheSet 设置缓存（自动 JSON 序列化）
func CacheSet(key string, value interface{}, ttl time.Duration) error {
	if !redis.IsRedisEnabled() {
		return nil // Redis 未启用时静默失败
	}

	ctx := context.Background()
	jsonData, err := json.Marshal(value)
	if err != nil {
		logging.LogError("缓存序列化失败", zap.String("key", key), zap.Error(err))
		return err
	}

	if err := redis.RedisClient.Set(ctx, key, jsonData, ttl).Err(); err != nil {
		logging.LogError("缓存设置失败", zap.String("key", key), zap.Error(err))
		return err
	}

	logging.LogDebug("缓存设置成功", zap.String("key", key), zap.Duration("ttl", ttl))
	return nil
}

// CacheDelete 删除缓存
func CacheDelete(keys ...string) error {
	if !redis.IsRedisEnabled() {
		return nil
	}

	ctx := context.Background()
	if err := redis.RedisClient.Del(ctx, keys...).Err(); err != nil {
		logging.LogError("缓存删除失败", zap.Strings("keys", keys), zap.Error(err))
		return err
	}

	logging.LogDebug("缓存删除成功", zap.Strings("keys", keys))
	return nil
}

// CacheDeletePattern 按模式删除缓存（慎用，性能较低）
func CacheDeletePattern(pattern string) error {
	if !redis.IsRedisEnabled() {
		return nil
	}

	ctx := context.Background()
	iter := redis.RedisClient.Scan(ctx, 0, pattern, 0).Iterator()

	var keys []string
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}

	if err := iter.Err(); err != nil {
		logging.LogError("缓存模式扫描失败", zap.String("pattern", pattern), zap.Error(err))
		return err
	}

	if len(keys) > 0 {
		if err := redis.RedisClient.Del(ctx, keys...).Err(); err != nil {
			logging.LogError("缓存批量删除失败", zap.String("pattern", pattern), zap.Error(err))
			return err
		}
		logging.LogInfo("缓存批量删除成功", zap.String("pattern", pattern), zap.Int("count", len(keys)))
	}

	return nil
}

// CacheExists 检查缓存是否存在
func CacheExists(key string) bool {
	if !redis.IsRedisEnabled() {
		return false
	}

	ctx := context.Background()
	exists, err := redis.RedisClient.Exists(ctx, key).Result()
	if err != nil {
		return false
	}
	return exists > 0
}

// CacheExpire 设置缓存过期时间
func CacheExpire(key string, ttl time.Duration) error {
	if !redis.IsRedisEnabled() {
		return nil
	}

	ctx := context.Background()
	return redis.RedisClient.Expire(ctx, key, ttl).Err()
}

// CacheClear 清空所有缓存（危险操作）
func CacheClear() error {
	if !redis.IsRedisEnabled() {
		return nil
	}

	ctx := context.Background()
	if err := redis.RedisClient.FlushDB(ctx).Err(); err != nil {
		logging.LogError("清空缓存失败", zap.Error(err))
		return err
	}

	logging.LogWarn("所有缓存已清空")
	return nil
}

// CacheGetOrSet 获取缓存，若不存在则执行回调函数并缓存结果
func CacheGetOrSet(key string, dest interface{}, ttl time.Duration, fetchFunc func() (interface{}, error)) error {
	// 先尝试从缓存获取
	err := CacheGet(key, dest)
	if err == nil {
		return nil // 缓存命中
	}

	// 缓存未命中，执行回调函数获取数据
	data, err := fetchFunc()
	if err != nil {
		return err
	}

	// 保存到缓存
	if err := CacheSet(key, data, ttl); err != nil {
		logging.LogWarn("缓存设置失败，但数据已获取", zap.String("key", key), zap.Error(err))
	}

	// 将数据赋值给 dest
	jsonData, _ := json.Marshal(data)
	return json.Unmarshal(jsonData, dest)
}

// InvalidateUserCache 清除用户相关缓存
func InvalidateUserCache(userID uint) {
	keys := []string{
		fmt.Sprintf(CacheKeyUser, userID),
		fmt.Sprintf(CacheKeyUserPermissions, userID),
	}
	CacheDelete(keys...)
}

// InvalidateRoleCache 清除角色相关缓存
func InvalidateRoleCache(roleID uint) {
	keys := []string{
		fmt.Sprintf(CacheKeyRole, roleID),
		fmt.Sprintf(CacheKeyRoleMenus, roleID),
	}
	CacheDelete(keys...)
	// 同时清除所有用户权限缓存（因为角色变更影响权限）
	CacheDeletePattern("user:permissions:*")
}

// InvalidateMenuCache 清除菜单相关缓存
func InvalidateMenuCache() {
	CacheDelete(CacheKeyMenuTree)
	CacheDeletePattern("role:menus:*")
	CacheDeletePattern("store:role:menus:*")
	CacheDeletePattern("user:permissions:*")
}

// InvalidateStoreCache 清除门店相关缓存
func InvalidateStoreCache(storeID uint) {
	keys := []string{
		fmt.Sprintf(CacheKeyStore, storeID),
		fmt.Sprintf(CacheKeyDishList, storeID),
	}
	CacheDelete(keys...)
}

// InvalidateDishCache 清除菜品相关缓存
func InvalidateDishCache(dishID uint, storeID uint) {
	keys := []string{
		fmt.Sprintf(CacheKeyDish, dishID),
		fmt.Sprintf(CacheKeyDishList, storeID),
		fmt.Sprintf(CacheKeyDishCategoriesWithDishes, storeID),
		fmt.Sprintf(CacheKeyDishCategories, storeID),
	}
	CacheDelete(keys...)
}

// InvalidateDishCategoryCache 清除分类相关缓存
func InvalidateDishCategoryCache(storeID uint) {
	keys := []string{
		fmt.Sprintf(CacheKeyDishCategories, storeID),
		fmt.Sprintf(CacheKeyDishCategoriesWithDishes, storeID),
	}
	CacheDelete(keys...)
}
