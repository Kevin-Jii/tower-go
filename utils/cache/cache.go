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
	CacheKeyStoreRoleMenus           = fmt.Sprintf("%sstore:role:menus:v2:store", CachePrefix)
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

// BatchGet 批量获取缓存
// 使用 Redis Pipeline 批量获取多个键的值，减少网络往返次数
func BatchGet(keys []string) (map[string]interface{}, error) {
	if !IsRedisEnabled() {
		return nil, fmt.Errorf("redis is not enabled")
	}

	if len(keys) == 0 {
		return make(map[string]interface{}), nil
	}

	ctx := context.Background()
	pipe := redisClient.Pipeline()

	// 批量添加 GET 命令到 Pipeline
	cmds := make(map[string]*redis.StringCmd, len(keys))
	for _, key := range keys {
		cmds[key] = pipe.Get(ctx, key)
	}

	// 执行 Pipeline
	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		// 如果不是 redis.Nil 错误，记录日志
		logging.LogError("Failed to execute pipeline batch get", zap.Error(err))
	}

	// 收集结果
	result := make(map[string]interface{}, len(keys))
	for key, cmd := range cmds {
		data, err := cmd.Bytes()
		if err == redis.Nil {
			// 键不存在，跳过
			continue
		}
		if err != nil {
			logging.LogError("Failed to get value from pipeline", zap.String("key", key), zap.Error(err))
			continue
		}

		// 尝试反序列化 JSON
		var value interface{}
		if err := json.Unmarshal(data, &value); err != nil {
			// 如果不是 JSON，直接存储字符串
			result[key] = string(data)
		} else {
			result[key] = value
		}
	}

	return result, nil
}

// BatchSet 批量设置缓存
// 使用 Redis Pipeline 批量设置多个键值对，减少网络往返次数
func BatchSet(items map[string]interface{}, ttl time.Duration) error {
	if !IsRedisEnabled() {
		return nil
	}

	if len(items) == 0 {
		return nil
	}

	ctx := context.Background()
	pipe := redisClient.Pipeline()

	// 批量添加 SET 命令到 Pipeline
	for key, value := range items {
		jsonData, err := json.Marshal(value)
		if err != nil {
			logging.LogError("Failed to marshal value for batch set", zap.String("key", key), zap.Error(err))
			continue
		}
		pipe.Set(ctx, key, jsonData, ttl)
	}

	// 执行 Pipeline
	_, err := pipe.Exec(ctx)
	if err != nil {
		logging.LogError("Failed to execute pipeline batch set", zap.Error(err))
		return err
	}

	return nil
}

// BatchGetTyped 批量获取缓存并反序列化到指定类型
// 使用 Redis Pipeline 批量获取多个键的值，并反序列化到 map[string]T
func BatchGetTyped(keys []string, destMap interface{}) error {
	if !IsRedisEnabled() {
		return fmt.Errorf("redis is not enabled")
	}

	if len(keys) == 0 {
		return nil
	}

	ctx := context.Background()
	pipe := redisClient.Pipeline()

	// 批量添加 GET 命令到 Pipeline
	cmds := make(map[string]*redis.StringCmd, len(keys))
	for _, key := range keys {
		cmds[key] = pipe.Get(ctx, key)
	}

	// 执行 Pipeline
	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		logging.LogError("Failed to execute pipeline batch get typed", zap.Error(err))
	}

	// 收集结果并反序列化
	results := make(map[string]json.RawMessage, len(keys))
	for key, cmd := range cmds {
		data, err := cmd.Bytes()
		if err == redis.Nil {
			// 键不存在，跳过
			continue
		}
		if err != nil {
			logging.LogError("Failed to get value from pipeline", zap.String("key", key), zap.Error(err))
			continue
		}
		results[key] = data
	}

	// 将结果序列化为 JSON 并反序列化到目标 map
	if len(results) > 0 {
		jsonData, err := json.Marshal(results)
		if err != nil {
			return err
		}
		return json.Unmarshal(jsonData, destMap)
	}

	return nil
}
