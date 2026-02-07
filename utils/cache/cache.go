package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Kevin-Jii/tower-go/config"
	"github.com/Kevin-Jii/tower-go/utils/logging"
	"go.uber.org/zap"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

// InitCache 初始化缓存
func InitCache() {
	cfg := config.GetRedisConfig()
	if !cfg.Enabled {
		logging.LogInfo("Redis is disabled")
		return
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		logging.LogError("Failed to connect to Redis", zap.Error(err))
		return
	}

	logging.LogInfo("Redis initialized successfully")
}

// GetRedisClient 获取 Redis 客户端
func GetRedisClient() *redis.Client {
	return redisClient
}

// IsRedisEnabled 检查 Redis 是否启用
func IsRedisEnabled() bool {
	return config.GetRedisConfig().Enabled && redisClient != nil
}

// CachePrefix 缓存前缀
const CachePrefix = "tower:"

// CacheGetOrSet 获取或设置缓存
func CacheGetOrSet(key string, dest interface{}, ttl time.Duration, fetch func() (interface{}, error)) error {
	if !IsRedisEnabled() {
		// Redis 未启用，直接获取数据
		data, err := fetch()
		if err != nil {
			return err
		}
		// 赋值给 dest
		if jsonData, err := json.Marshal(data); err == nil {
			json.Unmarshal(jsonData, dest)
		}
		return nil
	}

	ctx := context.Background()

	// 尝试从缓存获取
	data, err := redisClient.Get(ctx, key).Bytes()
	if err == nil {
		// 缓存命中，反序列化
		return json.Unmarshal(data, dest)
	}

	// 缓存未命中，获取数据
	result, err := fetch()
	if err != nil {
		return err
	}

	// 序列化并写入缓存
	jsonData, err := json.Marshal(result)
	if err != nil {
		return err
	}

	if err := redisClient.Set(ctx, key, jsonData, ttl).Err(); err != nil {
		logging.LogError("Failed to set cache", zap.String("key", key), zap.Error(err))
	}

	// 赋值给 dest
	return json.Unmarshal(jsonData, dest)
}

// CacheDelete 删除缓存
func CacheDelete(keys ...string) error {
	if !IsRedisEnabled() {
		return nil
	}

	ctx := context.Background()
	return redisClient.Del(ctx, keys...).Err()
}

// CacheDeleteByPattern 按模式删除缓存
func CacheDeleteByPattern(pattern string) error {
	if !IsRedisEnabled() {
		return nil
	}

	ctx := context.Background()

	// 查找匹配的键
	keys, err := redisClient.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) == 0 {
		return nil
	}

	// 删除匹配的键
	return redisClient.Del(ctx, keys...).Err()
}

// CacheTTL 获取缓存剩余时间
func CacheTTL(key string) (time.Duration, error) {
	if !IsRedisEnabled() {
		return 0, nil
	}

	ctx := context.Background()
	return redisClient.TTL(ctx, key).Result()
}

// 缓存键定义
var (
	CacheKeyDishCategories           = fmt.Sprintf("%sdish:categories:store", CachePrefix)
	CacheKeyDishCategoriesWithDishes = fmt.Sprintf("%sdish:categories:with:dishes:store", CachePrefix)
	CacheKeyMenuTree                 = fmt.Sprintf("%smenu:tree:role", CachePrefix)
	CacheKeyRoleMenus                = fmt.Sprintf("%srole:menus:role", CachePrefix)
	CacheKeyStoreRoleMenus           = fmt.Sprintf("%sstore:role:menus:store", CachePrefix)
	MenuTreeTTL                      = 30 * time.Minute
	RoleMenusTTL                     = 30 * time.Minute
	PermissionsTTL                   = 30 * time.Minute
)

// CacheGet 获取缓存
func CacheGet(key string, dest interface{}) error {
	if !IsRedisEnabled() {
		return nil
	}

	ctx := context.Background()
	data, err := redisClient.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}

	return json.Unmarshal(data, dest)
}

// CacheSet 设置缓存
func CacheSet(key string, value interface{}, ttl time.Duration) error {
	if !IsRedisEnabled() {
		return nil
	}

	ctx := context.Background()
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return redisClient.Set(ctx, key, jsonData, ttl).Err()
}

// InvalidateMenuCache 清除菜单缓存
func InvalidateMenuCache() {
	CacheDeleteByPattern(fmt.Sprintf("%smenu:*", CachePrefix))
}
