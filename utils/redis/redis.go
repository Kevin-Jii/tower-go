package redis

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

var client *redis.Client

// InitRedis 初始化 Redis 连接
func InitRedis(cfg config.RedisConfig) error {
	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: 100,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	logging.LogInfo("Redis connected successfully",
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
	)

	return nil
}

// GetClient 获取 Redis 客户端
func GetClient() *redis.Client {
	return client
}

// Set 设置键值
func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if client == nil {
		return nil
	}

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return client.Set(ctx, key, data, expiration).Err()
}

// Get 获取值
func Get(ctx context.Context, key string, result interface{}) error {
	if client == nil {
		return nil
	}

	data, err := client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}

	return json.Unmarshal(data, result)
}

// Delete 删除键
func Delete(ctx context.Context, keys ...string) error {
	if client == nil {
		return nil
	}

	return client.Del(ctx, keys...).Err()
}

// Exists 检查键是否存在
func Exists(ctx context.Context, key string) (bool, error) {
	if client == nil {
		return false, nil
	}

	n, err := client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return n > 0, nil
}

// TTL 获取键的剩余生存时间
func TTL(ctx context.Context, key string) (time.Duration, error) {
	if client == nil {
		return 0, nil
	}

	return client.TTL(ctx, key).Result()
}

// Incr 原子递增
func Incr(ctx context.Context, key string) (int64, error) {
	if client == nil {
		return 0, nil
	}

	return client.Incr(ctx, key).Result()
}

// SetNX 分布式锁
func SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	if client == nil {
		return false, nil
	}

	return client.SetNX(ctx, key, value, expiration).Result()
}

// CachePrefix 缓存前缀
const CachePrefix = "tower:"

// UserCacheKeys 用户相关缓存键
var UserCacheKeys = struct {
	UserByID      func(id uint) string
	UserByPhone   func(phone string) string
	UserList      func(storeID uint, page int) string
	StoreByID     func(id uint) string
	StoreList     string
	DictByType    func(typeCode string) string
	MenuByRoleID  func(roleID uint) string
}{
	UserByID:      func(id uint) string { return fmt.Sprintf("%suser:%d", CachePrefix, id) },
	UserByPhone:   func(phone string) string { return fmt.Sprintf("%suser:phone:%s", CachePrefix, phone) },
	UserList:      func(storeID uint, page int) string { return fmt.Sprintf("%suser:list:%d:%d", CachePrefix, storeID, page) },
	StoreByID:     func(id uint) string { return fmt.Sprintf("%sstore:%d", CachePrefix, id) },
	StoreList:     fmt.Sprintf("%sstore:list", CachePrefix),
	DictByType:    func(typeCode string) string { return fmt.Sprintf("%sdict:%s", CachePrefix, typeCode) },
	MenuByRoleID:  func(roleID uint) string { return fmt.Sprintf("%smenu:%d", CachePrefix, roleID) },
}

// CacheUser 缓存用户信息
func CacheUser(ctx context.Context, user interface{}) error {
	// 这里简化处理，实际应该从对象中提取 ID
	return nil
}

// GetCachedUser 获取缓存的用户信息
func GetCachedUser(ctx context.Context, userID uint, result interface{}) error {
	return Get(ctx, UserCacheKeys.UserByID(userID), result)
}

// IsRedisEnabled 检查 Redis 是否启用
func IsRedisEnabled() bool {
	return client != nil
}

// CloseRedis 关闭 Redis 连接
func CloseRedis() error {
	if client != nil {
		return client.Close()
	}
	return nil
}
