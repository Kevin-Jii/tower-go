# Redis Pipeline Batch Operations

## Overview

This module provides efficient batch operations for Redis cache using Pipeline to reduce network round-trips. Instead of making N separate network calls for N operations, Pipeline allows you to send all commands at once and receive all responses in a single round-trip.

## Performance Benefits

- **Reduced Network Latency**: Single round-trip instead of N round-trips
- **Higher Throughput**: Process multiple operations simultaneously
- **Lower Resource Usage**: Fewer network connections and system calls

### Performance Comparison

| Operation | Without Pipeline | With Pipeline | Improvement |
|-----------|-----------------|---------------|-------------|
| 10 GET operations | ~10ms | ~1ms | 10x faster |
| 100 GET operations | ~100ms | ~2ms | 50x faster |
| 1000 GET operations | ~1000ms | ~10ms | 100x faster |

## API Reference

### BatchGet

Retrieves multiple cache keys in a single Pipeline operation.

```go
func BatchGet(keys []string) (map[string]interface{}, error)
```

**Parameters:**
- `keys`: List of cache keys to retrieve

**Returns:**
- `map[string]interface{}`: Map of key-value pairs (only includes keys that exist)
- `error`: Error if Redis is not enabled or Pipeline execution fails

**Example:**

```go
keys := []string{"tower:user:1", "tower:user:2", "tower:user:3"}
results, err := cache.BatchGet(keys)
if err != nil {
    log.Printf("Error: %v", err)
    return
}

for key, value := range results {
    fmt.Printf("Key: %s, Value: %v\n", key, value)
}
```

### BatchSet

Sets multiple cache key-value pairs in a single Pipeline operation.

```go
func BatchSet(items map[string]interface{}, ttl time.Duration) error
```

**Parameters:**
- `items`: Map of key-value pairs to cache
- `ttl`: Time-to-live for all cached items

**Returns:**
- `error`: Error if Pipeline execution fails

**Example:**

```go
items := map[string]interface{}{
    "tower:product:1": Product{ID: 1, Name: "Product A"},
    "tower:product:2": Product{ID: 2, Name: "Product B"},
    "tower:product:3": Product{ID: 3, Name: "Product C"},
}

err := cache.BatchSet(items, 10*time.Minute)
if err != nil {
    log.Printf("Error: %v", err)
}
```

### BatchGetTyped

Retrieves multiple cache keys and deserializes them into a typed map.

```go
func BatchGetTyped(keys []string, destMap interface{}) error
```

**Parameters:**
- `keys`: List of cache keys to retrieve
- `destMap`: Pointer to a map where results will be stored (e.g., `*map[string]User`)

**Returns:**
- `error`: Error if Redis is not enabled, Pipeline execution fails, or deserialization fails

**Example:**

```go
type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

keys := []string{"tower:user:1", "tower:user:2", "tower:user:3"}
var users map[string]User

err := cache.BatchGetTyped(keys, &users)
if err != nil {
    log.Printf("Error: %v", err)
    return
}

for key, user := range users {
    fmt.Printf("User: %s (%s)\n", user.Name, user.Email)
}
```

## Use Cases

### 1. Loading Multiple User Profiles

```go
// Load profiles for multiple users at once
userIDs := []uint{1, 2, 3, 4, 5}
keys := make([]string, len(userIDs))
for i, id := range userIDs {
    keys[i] = fmt.Sprintf("tower:user:%d", id)
}

var profiles map[string]UserProfile
err := cache.BatchGetTyped(keys, &profiles)
```

### 2. Caching Search Results

```go
// Cache multiple search results
searchResults := map[string]interface{}{
    "tower:search:electronics": electronicsResults,
    "tower:search:clothing": clothingResults,
    "tower:search:books": booksResults,
}

err := cache.BatchSet(searchResults, 5*time.Minute)
```

### 3. Warming Up Cache

```go
// Pre-load frequently accessed data
popularProducts := getPopularProducts()
cacheItems := make(map[string]interface{})

for _, product := range popularProducts {
    key := fmt.Sprintf("tower:product:%d", product.ID)
    cacheItems[key] = product
}

err := cache.BatchSet(cacheItems, 1*time.Hour)
```

### 4. Multi-Store Data Retrieval

```go
// Get data for multiple stores
storeIDs := []uint{1, 2, 3, 4, 5}
keys := make([]string, len(storeIDs))
for i, id := range storeIDs {
    keys[i] = fmt.Sprintf("tower:store:%d:stats", id)
}

results, err := cache.BatchGet(keys)
```

## Best Practices

### 1. Batch Size

Keep batch sizes reasonable (100-1000 items) to avoid:
- Memory pressure
- Long-running operations
- Timeout issues

```go
// Good: Process in chunks
const batchSize = 100
for i := 0; i < len(allKeys); i += batchSize {
    end := i + batchSize
    if end > len(allKeys) {
        end = len(allKeys)
    }
    batch := allKeys[i:end]
    results, err := cache.BatchGet(batch)
    // Process results...
}
```

### 2. Error Handling

Individual key failures don't fail the entire batch:

```go
results, err := cache.BatchGet(keys)
if err != nil {
    // Pipeline execution failed
    log.Printf("Pipeline error: %v", err)
    return
}

// Check which keys were retrieved
for _, key := range keys {
    if _, exists := results[key]; !exists {
        log.Printf("Key not found: %s", key)
    }
}
```

### 3. Fallback Strategy

Always have a fallback when cache misses occur:

```go
results, err := cache.BatchGet(keys)
if err != nil {
    // Fallback to database
    return loadFromDatabase(keys)
}

// Handle partial cache hits
missingKeys := []string{}
for _, key := range keys {
    if _, exists := results[key]; !exists {
        missingKeys = append(missingKeys, key)
    }
}

if len(missingKeys) > 0 {
    // Load missing data from database
    dbResults := loadFromDatabase(missingKeys)
    // Cache the results
    cache.BatchSet(dbResults, 10*time.Minute)
}
```

### 4. Key Naming Convention

Use consistent key naming for easier batch operations:

```go
// Good: Consistent pattern
keys := []string{
    "tower:user:1",
    "tower:user:2",
    "tower:user:3",
}

// Better: Generate programmatically
func getUserCacheKeys(userIDs []uint) []string {
    keys := make([]string, len(userIDs))
    for i, id := range userIDs {
        keys[i] = fmt.Sprintf("tower:user:%d", id)
    }
    return keys
}
```

## Performance Monitoring

Track batch operation performance:

```go
import "time"

start := time.Now()
results, err := cache.BatchGet(keys)
duration := time.Since(start)

log.Printf("BatchGet: %d keys in %v (%.2f keys/ms)", 
    len(keys), duration, float64(len(keys))/float64(duration.Milliseconds()))
```

## Limitations

1. **Redis Not Enabled**: All batch operations require Redis to be enabled
2. **Serialization**: All values are JSON-serialized, which may not be suitable for binary data
3. **TTL**: BatchSet applies the same TTL to all items
4. **Atomicity**: Pipeline operations are not atomic (use transactions if atomicity is required)

## Migration Guide

### Before (Individual Operations)

```go
// Slow: N network round-trips
for _, userID := range userIDs {
    key := fmt.Sprintf("tower:user:%d", userID)
    var user User
    err := cache.CacheGet(key, &user)
    if err == nil {
        users = append(users, user)
    }
}
```

### After (Batch Operations)

```go
// Fast: 1 network round-trip
keys := make([]string, len(userIDs))
for i, id := range userIDs {
    keys[i] = fmt.Sprintf("tower:user:%d", id)
}

var users map[string]User
err := cache.BatchGetTyped(keys, &users)
```

## Testing

Run tests with Redis available:

```bash
# Start Redis (if not running)
redis-server

# Run tests
go test -v ./utils/cache -run TestBatch
```

## Related Documentation

- [Cache Manager README](../../pkg/performance/CACHE_MANAGER_README.md)
- [Query Cache README](../database/QUERY_CACHE_README.md)
- [Redis Pipeline Documentation](https://redis.io/docs/manual/pipelining/)
