# Multi-Tier Cache Manager

## Overview

The Multi-Tier Cache Manager provides a high-performance caching solution with automatic fallback between local memory cache and Redis. It implements intelligent cache access patterns, distributed locking to prevent cache stampede, and LRU eviction for memory efficiency.

## Features

### 1. Multi-Tier Cache Architecture

The cache manager automatically checks multiple cache layers in order:

1. **Local Cache (L1)**: In-memory LRU cache for ultra-fast access
2. **Redis Cache (L2)**: Distributed cache for shared data across instances
3. **Data Source**: Original data source (database, API, etc.)

### 2. LRU Cache Implementation

- **O(1) Time Complexity**: Both get and put operations are O(1)
- **Automatic Eviction**: Least Recently Used items are evicted when capacity is reached
- **TTL Support**: Items automatically expire after configured time
- **Thread-Safe**: Uses RWMutex for concurrent access

### 3. Cache Stampede Prevention

Uses distributed locks (Redis SETNX) to prevent multiple goroutines from loading the same data simultaneously when cache expires.

### 4. Automatic Cache Warming

When data is fetched from Redis, it's automatically written back to local cache for faster subsequent access.

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    Application Layer                     │
│                  (Your Business Logic)                   │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│                   CacheManager Interface                 │
│  • Get(key) → value                                     │
│  • Set(key, value, ttl)                                 │
│  • Delete(keys...)                                      │
│  • GetOrSet(key, fetch, ttl) → value                    │
└─────────────────────────────────────────────────────────┘
                          │
        ┌─────────────────┼─────────────────┐
        ▼                 ▼                 ▼
┌──────────────┐  ┌──────────────┐  ┌──────────────┐
│ Local Cache  │  │ Redis Cache  │  │ Distributed  │
│   (LRU)      │  │              │  │    Lock      │
│              │  │              │  │              │
│ • Fast       │  │ • Shared     │  │ • Prevents   │
│ • O(1) ops   │  │ • Persistent │  │   stampede   │
│ • Auto evict │  │ • Scalable   │  │ • SETNX      │
└──────────────┘  └──────────────┘  └──────────────┘
```

## Usage

### Basic Usage

```go
import (
    "context"
    "time"
    "github.com/Kevin-Jii/tower-go/pkg/performance"
)

// Create cache manager
cacheManager := performance.NewCacheManager()
ctx := context.Background()

// Set cache
err := cacheManager.Set(ctx, "user:123", userData, 5*time.Minute)

// Get cache
var user User
err = cacheManager.Get(ctx, "user:123", &user)
if err == performance.ErrCacheMiss {
    // Cache miss, load from database
}

// Delete cache
err = cacheManager.Delete(ctx, "user:123")
```

### GetOrSet Pattern (Recommended)

The `GetOrSet` method is the recommended way to use the cache manager. It automatically handles cache misses and prevents cache stampede:

```go
var user User
err := cacheManager.GetOrSet(ctx, "user:123", &user, 5*time.Minute, func() (interface{}, error) {
    // This function is only called on cache miss
    return loadUserFromDatabase(123)
})
```

### Using Local Cache Directly

For scenarios where you only need local caching:

```go
// Create LRU cache with capacity 1000 and 5-minute TTL
cache := performance.NewLRUCache(1000, 5*time.Minute)

// Set value
cache.Set("key", "value", 5*time.Minute)

// Get value
if value, ok := cache.Get("key"); ok {
    fmt.Println("Found:", value)
}

// Delete value
cache.Delete("key")

// Clear all
cache.Clear()

// Get size
size := cache.Len()
```

## Configuration

The cache manager is configured through environment variables:

```bash
# Local Cache Configuration
PERF_CACHE_LOCAL_ENABLED=true          # Enable local cache
PERF_CACHE_LOCAL_SIZE=10000            # Local cache capacity
PERF_CACHE_LOCAL_TTL_SECONDS=300       # Local cache TTL (5 minutes)

# Redis Configuration
PERF_CACHE_REDIS_ENABLED=true          # Enable Redis cache
PERF_CACHE_REDIS_PIPELINE_SIZE=100     # Pipeline batch size
PERF_CACHE_REDIS_POOL_SIZE=10          # Connection pool size

# Cache Strategy
PERF_CACHE_ENABLE_WARMUP=false         # Enable cache warmup
PERF_CACHE_ENABLE_LOCK=true            # Enable distributed lock
```

## Performance Characteristics

### Local Cache (LRU)

- **Get Operation**: O(1) time complexity
- **Set Operation**: O(1) time complexity
- **Memory Usage**: O(n) where n is capacity
- **Eviction**: Automatic LRU eviction when capacity reached
- **Concurrency**: Thread-safe with RWMutex

### Multi-Tier Access

| Scenario | Local Cache | Redis | Database | Total Time |
|----------|-------------|-------|----------|------------|
| Hot data (L1 hit) | ✓ | - | - | ~1μs |
| Warm data (L2 hit) | ✗ → ✓ | ✓ | - | ~1ms |
| Cold data (miss) | ✗ | ✗ | ✓ | ~10ms+ |

## Best Practices

### 1. Use GetOrSet for Cache-Aside Pattern

```go
// ✓ Good: Automatic cache management
err := cacheManager.GetOrSet(ctx, key, &result, ttl, fetchFunc)

// ✗ Avoid: Manual cache management
result, err := cacheManager.Get(ctx, key)
if err == ErrCacheMiss {
    result = fetchFromDB()
    cacheManager.Set(ctx, key, result, ttl)
}
```

### 2. Choose Appropriate TTL

- **Hot data**: 5-15 minutes (frequently accessed)
- **Warm data**: 30-60 minutes (moderately accessed)
- **Cold data**: 2-24 hours (rarely accessed)

### 3. Use Structured Keys

```go
// ✓ Good: Structured, easy to invalidate
"user:123"
"store:456:products"
"order:789:details"

// ✗ Avoid: Unstructured keys
"u123"
"data_456"
```

### 4. Batch Operations

```go
// ✓ Good: Batch delete
cacheManager.Delete(ctx, "key1", "key2", "key3")

// ✗ Avoid: Multiple single deletes
cacheManager.Delete(ctx, "key1")
cacheManager.Delete(ctx, "key2")
cacheManager.Delete(ctx, "key3")
```

### 5. Handle Cache Errors Gracefully

```go
var user User
err := cacheManager.Get(ctx, key, &user)
if err != nil {
    // Log error but continue with database query
    log.Warn("Cache error", zap.Error(err))
    user = loadFromDatabase()
}
```

## Cache Invalidation Strategies

### 1. Time-Based Expiration (TTL)

Automatic expiration after configured time:

```go
cacheManager.Set(ctx, key, value, 5*time.Minute)
```

### 2. Explicit Invalidation

Delete cache when data changes:

```go
func UpdateUser(user User) error {
    // Update database
    err := db.Save(&user)
    if err != nil {
        return err
    }
    
    // Invalidate cache
    cacheManager.Delete(ctx, fmt.Sprintf("user:%d", user.ID))
    return nil
}
```

### 3. Pattern-Based Invalidation

For invalidating multiple related keys, use a consistent naming pattern and delete in batch:

```go
// Delete all user-related caches
keys := []string{
    fmt.Sprintf("user:%d", userID),
    fmt.Sprintf("user:%d:profile", userID),
    fmt.Sprintf("user:%d:orders", userID),
}
cacheManager.Delete(ctx, keys...)
```

## Monitoring

### Cache Statistics

```go
// Get local cache size
size := cache.Len()

// Monitor cache hit rate
hits := 0
misses := 0

err := cacheManager.Get(ctx, key, &result)
if err == nil {
    hits++
} else if err == ErrCacheMiss {
    misses++
}

hitRate := float64(hits) / float64(hits + misses)
```

## Testing

The cache manager includes comprehensive tests:

```bash
# Run all cache manager tests
go test -v ./pkg/performance -run "TestLRUCache|TestMultiTierCacheManager"

# Run specific test
go test -v ./pkg/performance -run TestLRUCache_Eviction

# Run with race detection
go test -race ./pkg/performance
```

## Implementation Details

### LRU Cache Structure

The LRU cache uses a doubly-linked list combined with a hash map:

- **Hash Map**: O(1) key lookup
- **Doubly-Linked List**: O(1) insertion/deletion and LRU ordering
- **Sentinel Nodes**: Simplify edge cases (head/tail)

### Distributed Lock

Uses Redis SETNX for distributed locking:

```go
// Acquire lock
locked, err := lock.Lock(ctx, "lock:key", 10*time.Second)
if locked {
    defer lock.Unlock(ctx, "lock:key")
    // Critical section
}
```

## Troubleshooting

### High Memory Usage

If local cache is consuming too much memory:

1. Reduce `PERF_CACHE_LOCAL_SIZE`
2. Reduce `PERF_CACHE_LOCAL_TTL_SECONDS`
3. Monitor eviction rate

### Low Cache Hit Rate

If cache hit rate is low:

1. Increase cache capacity
2. Increase TTL for stable data
3. Review cache key patterns
4. Check if data is actually cacheable

### Redis Connection Issues

If Redis is unavailable:

1. Cache manager automatically falls back to local cache only
2. Distributed lock is disabled
3. System continues to function (degraded mode)

## Related Documentation

- [Performance Configuration](../../config/performance.go)
- [Query Cache](../../utils/database/QUERY_CACHE_README.md)
- [Performance Optimization Guide](../../docs/PERFORMANCE_TUNING_GUIDE.md)
