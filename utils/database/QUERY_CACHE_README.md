# Query Cache Implementation

## Overview

The query cache implementation provides automatic caching of database query results using Redis. This helps reduce database load and improve query performance for frequently accessed data.

## Features

1. **Automatic Cache Key Generation**: Generates cache keys based on query SQL and parameters using MD5 hashing
2. **Custom Cache Keys**: Supports custom cache keys for specific use cases
3. **Empty Condition Short-Circuit**: Queries without conditions are not cached (Requirement 2.5)
4. **Aggregation Query Caching**: Supports caching of aggregation queries (Requirement 2.4)
5. **Cache Invalidation**: Pattern-based cache invalidation support
6. **Redis Integration**: Uses existing Redis infrastructure

## Usage

### Basic Usage

```go
import (
    "context"
    "time"
    "github.com/Kevin-Jii/tower-go/utils/database"
)

// Create Redis query cache
cache := database.NewRedisQueryCache("myapp:")

// Create query builder with cache enabled
qb := database.NewQueryBuilder(db)
qb.Where("store_id = ?", 1).
   Where("status = ?", 1).
   OrderBy("created_at", "DESC").
   Limit(10).
   WithCache(cache, 5*time.Minute)

// Execute query (first time: queries database and caches result)
var products []Product
err := qb.FindWithCache(context.Background(), &products)

// Execute same query again (second time: retrieves from cache)
err = qb.FindWithCache(context.Background(), &products)
```

### Custom Cache Key

```go
// Use custom cache key for specific queries
qb := database.NewQueryBuilder(db)
qb.Where("store_id = ?", 1).
   WithCache(cache, 5*time.Minute).
   WithCacheKey("store:1:products")

var products []Product
err := qb.FindWithCache(context.Background(), &products)
```

### Cache Invalidation

```go
// Invalidate all product-related caches
qb := database.NewQueryBuilder(db)
qb.WithCache(cache, 5*time.Minute)
err := qb.InvalidateCache(context.Background(), "query:*product*")
```

### Aggregation Query Caching

```go
// Cache aggregation queries
qb := database.NewQueryBuilder(db)
qb.Where("store_id = ?", 1).
   Where("created_at >= ?", time.Now().AddDate(0, -1, 0)).
   WithCache(cache, 10*time.Minute)

var count int64
err := qb.Count(&count)
```

## Implementation Details

### Cache Key Generation

Cache keys are generated using MD5 hash of:
- WHERE conditions
- Query arguments
- ORDER BY clauses
- LIMIT and OFFSET values

Format: `query:<md5_hash>`

### Empty Condition Short-Circuit

Queries without any conditions, orders, limits, or offsets will NOT be cached. This prevents caching of full table scans which could:
- Consume excessive cache memory
- Return stale data for frequently updated tables
- Provide minimal performance benefit

### Redis Integration

The implementation uses the existing Redis client from `utils/cache` package. If Redis is not enabled, the cache operations become no-ops and queries execute normally.

## Requirements Validation

This implementation satisfies the following requirements:

- **Requirement 2.4**: Aggregation query caching - Queries can be cached with configurable TTL
- **Requirement 2.5**: Empty condition short-circuit - Queries without conditions are not cached

## Testing

Run tests with:

```bash
go test -v ./utils/database -run "TestGenerateCacheKey|TestWithCache|TestNoOpQueryCache"
```

## Performance Considerations

1. **Cache TTL**: Choose appropriate TTL based on data update frequency
2. **Cache Key Size**: MD5 hash ensures consistent key size regardless of query complexity
3. **Memory Usage**: Monitor Redis memory usage for cached query results
4. **Cache Invalidation**: Implement proper cache invalidation on data updates

## Future Enhancements

- Local cache layer for hot data (multi-tier caching)
- Cache warming strategies
- Distributed lock for cache stampede prevention
- Query result compression for large datasets
