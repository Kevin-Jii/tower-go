package performance

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Kevin-Jii/tower-go/config"
	"github.com/Kevin-Jii/tower-go/utils/cache"
	"github.com/Kevin-Jii/tower-go/utils/logging"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var (
	// ErrCacheMiss 缓存未命中错误
	ErrCacheMiss = errors.New("cache miss")
	// ErrCacheLockFailed 获取缓存锁失败
	ErrCacheLockFailed = errors.New("failed to acquire cache lock")
)

// CacheManager 缓存管理器接口
type CacheManager interface {
	// Get 获取缓存(自动选择最优层级)
	Get(ctx context.Context, key string, dest interface{}) error

	// Set 设置缓存(写入所有层级)
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error

	// Delete 删除缓存(删除所有层级)
	Delete(ctx context.Context, keys ...string) error

	// GetOrSet 获取或设置缓存
	GetOrSet(ctx context.Context, key string, dest interface{}, ttl time.Duration,
		fetch func() (interface{}, error)) error
}

// LocalCache 本地缓存接口
type LocalCache interface {
	// Get 获取本地缓存
	Get(key string) (interface{}, bool)

	// Set 设置本地缓存
	Set(key string, value interface{}, ttl time.Duration)

	// Delete 删除本地缓存
	Delete(key string)

	// Clear 清空本地缓存
	Clear()

	// Len 返回缓存项数量
	Len() int
}

// DistributedLock 分布式锁接口
type DistributedLock interface {
	// Lock 获取锁
	Lock(ctx context.Context, key string, ttl time.Duration) (bool, error)

	// Unlock 释放锁
	Unlock(ctx context.Context, key string) error

	// Extend 延长锁时间
	Extend(ctx context.Context, key string, ttl time.Duration) error
}

// multiTierCacheManager 多级缓存管理器实现
type multiTierCacheManager struct {
	localCache  LocalCache
	redisClient *redis.Client
	lock        DistributedLock
	config      config.CachePerformanceConfig
}

// NewCacheManager 创建缓存管理器
func NewCacheManager() CacheManager {
	cfg := config.GetPerformanceConfig().Cache

	var localCache LocalCache
	if cfg.LocalCacheEnabled {
		localCache = NewLRUCache(cfg.LocalCacheSize, cfg.LocalCacheTTL)
	}

	var lock DistributedLock
	if cfg.EnableCacheLock && cache.IsRedisEnabled() {
		lock = NewRedisLock(cache.GetRedisClient())
	}

	return &multiTierCacheManager{
		localCache:  localCache,
		redisClient: cache.GetRedisClient(),
		lock:        lock,
		config:      cfg,
	}
}

// Get 获取缓存(自动选择最优层级)
func (m *multiTierCacheManager) Get(ctx context.Context, key string, dest interface{}) error {
	// 1. 尝试从本地缓存获取
	if m.localCache != nil {
		if value, ok := m.localCache.Get(key); ok {
			// 本地缓存命中
			return copyValue(value, dest)
		}
	}

	// 2. 尝试从Redis获取
	if m.redisClient != nil && m.config.RedisEnabled {
		data, err := m.redisClient.Get(ctx, key).Bytes()
		if err == nil {
			// Redis缓存命中
			if err := json.Unmarshal(data, dest); err != nil {
				return err
			}

			// 回写到本地缓存
			if m.localCache != nil {
				m.localCache.Set(key, dest, m.config.LocalCacheTTL)
			}

			return nil
		}

		if err != redis.Nil {
			// Redis错误(非未命中)
			logging.LogError("Redis get error", zap.String("key", key), zap.Error(err))
		}
	}

	return ErrCacheMiss
}

// Set 设置缓存(写入所有层级)
func (m *multiTierCacheManager) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	// 1. 写入本地缓存
	if m.localCache != nil {
		localTTL := ttl
		if localTTL > m.config.LocalCacheTTL {
			localTTL = m.config.LocalCacheTTL
		}
		m.localCache.Set(key, value, localTTL)
	}

	// 2. 写入Redis
	if m.redisClient != nil && m.config.RedisEnabled {
		data, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value: %w", err)
		}

		if err := m.redisClient.Set(ctx, key, data, ttl).Err(); err != nil {
			logging.LogError("Redis set error", zap.String("key", key), zap.Error(err))
			return fmt.Errorf("failed to set redis cache: %w", err)
		}
	}

	return nil
}

// Delete 删除缓存(删除所有层级)
func (m *multiTierCacheManager) Delete(ctx context.Context, keys ...string) error {
	// 1. 删除本地缓存
	if m.localCache != nil {
		for _, key := range keys {
			m.localCache.Delete(key)
		}
	}

	// 2. 删除Redis缓存
	if m.redisClient != nil && m.config.RedisEnabled && len(keys) > 0 {
		if err := m.redisClient.Del(ctx, keys...).Err(); err != nil {
			logging.LogError("Redis delete error", zap.Strings("keys", keys), zap.Error(err))
			return fmt.Errorf("failed to delete redis cache: %w", err)
		}
	}

	return nil
}

// GetOrSet 获取或设置缓存
func (m *multiTierCacheManager) GetOrSet(ctx context.Context, key string, dest interface{}, ttl time.Duration,
	fetch func() (interface{}, error)) error {

	// 1. 尝试获取缓存
	err := m.Get(ctx, key, dest)
	if err == nil {
		return nil
	}

	if err != ErrCacheMiss {
		// 非缓存未命中错误
		logging.LogError("Cache get error", zap.String("key", key), zap.Error(err))
	}

	// 2. 缓存未命中,使用分布式锁防止缓存击穿
	if m.lock != nil && m.config.EnableCacheLock {
		lockKey := fmt.Sprintf("lock:%s", key)
		locked, err := m.lock.Lock(ctx, lockKey, 10*time.Second)
		if err != nil {
			logging.LogError("Failed to acquire lock", zap.String("key", lockKey), zap.Error(err))
			// 锁获取失败,继续执行(降级)
		} else if locked {
			// 获取锁成功
			defer m.lock.Unlock(ctx, lockKey)

			// 再次尝试获取缓存(可能其他goroutine已经加载)
			err := m.Get(ctx, key, dest)
			if err == nil {
				return nil
			}
		} else {
			// 未获取到锁,等待一段时间后重试
			time.Sleep(100 * time.Millisecond)
			err := m.Get(ctx, key, dest)
			if err == nil {
				return nil
			}
		}
	}

	// 3. 从数据源获取数据
	value, err := fetch()
	if err != nil {
		return fmt.Errorf("failed to fetch data: %w", err)
	}

	// 4. 设置缓存
	if err := m.Set(ctx, key, value, ttl); err != nil {
		logging.LogError("Failed to set cache", zap.String("key", key), zap.Error(err))
		// 设置缓存失败不影响返回结果
	}

	// 5. 复制结果到dest
	return copyValue(value, dest)
}

// copyValue 复制值到目标
func copyValue(src, dest interface{}) error {
	// 使用JSON序列化/反序列化进行深拷贝
	data, err := json.Marshal(src)
	if err != nil {
		return fmt.Errorf("failed to marshal source: %w", err)
	}

	if err := json.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("failed to unmarshal to dest: %w", err)
	}

	return nil
}

// cacheItem 缓存项
type cacheItem struct {
	value      interface{}
	expireTime time.Time
}

// isExpired 检查是否过期
func (item *cacheItem) isExpired() bool {
	return time.Now().After(item.expireTime)
}

// lruCache LRU缓存实现
type lruCache struct {
	mu         sync.RWMutex
	capacity   int
	items      map[string]*lruNode
	head       *lruNode // 最近使用的节点
	tail       *lruNode // 最久未使用的节点
	defaultTTL time.Duration
}

// lruNode LRU双向链表节点
type lruNode struct {
	key  string
	item *cacheItem
	prev *lruNode
	next *lruNode
}

// NewLRUCache 创建LRU缓存
func NewLRUCache(capacity int, defaultTTL time.Duration) LocalCache {
	if capacity <= 0 {
		capacity = 1000 // 默认容量
	}

	cache := &lruCache{
		capacity:   capacity,
		items:      make(map[string]*lruNode, capacity),
		defaultTTL: defaultTTL,
	}

	// 初始化哨兵节点
	cache.head = &lruNode{}
	cache.tail = &lruNode{}
	cache.head.next = cache.tail
	cache.tail.prev = cache.head

	return cache
}

// Get 获取本地缓存
func (c *lruCache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	node, ok := c.items[key]
	if !ok {
		return nil, false
	}

	// 检查是否过期
	if node.item.isExpired() {
		c.removeNode(node)
		delete(c.items, key)
		return nil, false
	}

	// 移动到链表头部(最近使用)
	c.moveToHead(node)

	return node.item.value, true
}

// Set 设置本地缓存
func (c *lruCache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ttl == 0 {
		ttl = c.defaultTTL
	}

	expireTime := time.Now().Add(ttl)

	// 如果key已存在,更新值并移到头部
	if node, ok := c.items[key]; ok {
		node.item.value = value
		node.item.expireTime = expireTime
		c.moveToHead(node)
		return
	}

	// 创建新节点
	newNode := &lruNode{
		key: key,
		item: &cacheItem{
			value:      value,
			expireTime: expireTime,
		},
	}

	// 添加到map和链表头部
	c.items[key] = newNode
	c.addToHead(newNode)

	// 检查容量,如果超出则淘汰最久未使用的项
	if len(c.items) > c.capacity {
		// 移除尾部节点
		tailNode := c.tail.prev
		c.removeNode(tailNode)
		delete(c.items, tailNode.key)
	}
}

// Delete 删除本地缓存
func (c *lruCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if node, ok := c.items[key]; ok {
		c.removeNode(node)
		delete(c.items, key)
	}
}

// Clear 清空本地缓存
func (c *lruCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]*lruNode, c.capacity)
	c.head.next = c.tail
	c.tail.prev = c.head
}

// Len 返回缓存项数量
func (c *lruCache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.items)
}

// addToHead 添加节点到链表头部
func (c *lruCache) addToHead(node *lruNode) {
	node.prev = c.head
	node.next = c.head.next
	c.head.next.prev = node
	c.head.next = node
}

// removeNode 从链表中移除节点
func (c *lruCache) removeNode(node *lruNode) {
	node.prev.next = node.next
	node.next.prev = node.prev
}

// moveToHead 移动节点到链表头部
func (c *lruCache) moveToHead(node *lruNode) {
	c.removeNode(node)
	c.addToHead(node)
}

// redisLock Redis分布式锁实现
type redisLock struct {
	client *redis.Client
}

// NewRedisLock 创建Redis分布式锁
func NewRedisLock(client *redis.Client) DistributedLock {
	return &redisLock{
		client: client,
	}
}

// Lock 获取锁
func (l *redisLock) Lock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	if l.client == nil {
		return false, errors.New("redis client is nil")
	}

	// 使用SETNX实现分布式锁
	locked, err := l.client.SetNX(ctx, key, "1", ttl).Result()
	if err != nil {
		return false, fmt.Errorf("failed to acquire lock: %w", err)
	}

	return locked, nil
}

// Unlock 释放锁
func (l *redisLock) Unlock(ctx context.Context, key string) error {
	if l.client == nil {
		return errors.New("redis client is nil")
	}

	err := l.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to release lock: %w", err)
	}

	return nil
}

// Extend 延长锁时间
func (l *redisLock) Extend(ctx context.Context, key string, ttl time.Duration) error {
	if l.client == nil {
		return errors.New("redis client is nil")
	}

	err := l.client.Expire(ctx, key, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to extend lock: %w", err)
	}

	return nil
}
