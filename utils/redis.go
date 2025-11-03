package utils

import (
	"context"
	"fmt"
	"tower-go/config"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var RedisClient *redis.Client
var redisEnabled bool

// InitRedis 初始化 Redis 连接
func InitRedis(cfg config.RedisConfig) error {
	if !cfg.Enabled {
		LogInfo("Redis 缓存已禁用")
		redisEnabled = false
		return nil
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// 测试连接
	ctx := context.Background()
	if err := RedisClient.Ping(ctx).Err(); err != nil {
		LogError("Redis 连接失败", zap.Error(err))
		redisEnabled = false
		return err
	}

	redisEnabled = true
	LogInfo("Redis 连接成功", zap.String("addr", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)))
	return nil
}

// IsRedisEnabled 检查 Redis 是否启用
func IsRedisEnabled() bool {
	return redisEnabled && RedisClient != nil
}

// CloseRedis 关闭 Redis 连接
func CloseRedis() error {
	if RedisClient != nil {
		return RedisClient.Close()
	}
	return nil
}
