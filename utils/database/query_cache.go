package database

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Kevin-Jii/tower-go/utils/cache"
	"github.com/redis/go-redis/v9"
)

// redisQueryCache Redis查询缓存实现
type redisQueryCache struct {
	client *redis.Client
	prefix string
}

// NewRedisQueryCache 创建Redis查询缓存
func NewRedisQueryCache(prefix string) QueryCache {
	client := cache.GetRedisClient()
	if client == nil {
		return &noOpQueryCache{}
	}

	if prefix == "" {
		prefix = "query:"
	}

	return &redisQueryCache{
		client: client,
		prefix: prefix,
	}
}

// Get 获取缓存的查询结果
func (r *redisQueryCache) Get(ctx context.Context, key string, dest interface{}) error {
	if r.client == nil {
		return redis.Nil
	}

	fullKey := r.prefix + key
	data, err := r.client.Get(ctx, fullKey).Bytes()
	if err != nil {
		return err
	}

	return json.Unmarshal(data, dest)
}

// Set 设置查询结果缓存
func (r *redisQueryCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if r.client == nil {
		return nil
	}

	fullKey := r.prefix + key
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, fullKey, jsonData, ttl).Err()
}

// Invalidate 失效缓存（支持模式匹配）
func (r *redisQueryCache) Invalidate(ctx context.Context, pattern string) error {
	if r.client == nil {
		return nil
	}

	fullPattern := r.prefix + pattern
	keys, err := r.client.Keys(ctx, fullPattern).Result()
	if err != nil {
		return err
	}

	if len(keys) == 0 {
		return nil
	}

	return r.client.Del(ctx, keys...).Err()
}

// noOpQueryCache 空操作查询缓存（当Redis未启用时使用）
type noOpQueryCache struct{}

func (n *noOpQueryCache) Get(ctx context.Context, key string, dest interface{}) error {
	return redis.Nil
}

func (n *noOpQueryCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return nil
}

func (n *noOpQueryCache) Invalidate(ctx context.Context, pattern string) error {
	return nil
}
