# Hot Path Performance Optimizations

This package provides optimized utilities for hot path operations in the Tower-Go system, focusing on type assertions and concurrent map operations.

## Components

### 1. TypeConverter

Provides fast type conversions using type switches instead of reflection.

**Performance**: ~2.7 ns/op with 0 allocations

```go
converter := performance.GetTypeConverter()

// Convert to uint
val, err := converter.ToUint(someInterface)

// Convert to string
str := converter.ToString(someInterface)

// Convert to int
num, err := converter.ToInt(someInterface)
```

**Benefits**:
- 10-100x faster than reflection-based conversion
- Zero memory allocations
- Type-safe with proper error handling

### 2. ConcurrentCache

Thread-safe cache using `sync.Map` optimized for read-heavy workloads.

**Performance**: 
- Reads: ~3.1 ns/op (9x faster than mutex-protected map)
- Writes: ~40 ns/op (1.8x faster than mutex-protected map)

```go
cache := performance.NewConcurrentCache()

// Set value
cache.Set("key", value)

// Get value
val, ok := cache.Get("key")

// Atomic get or set
actual, loaded := cache.GetOrSet("key", value)

// Delete value
cache.Delete("key")

// Iterate over all entries
cache.Range(func(key, value interface{}) bool {
    // Process entry
    return true // continue iteration
})

// Clear all entries
cache.Clear()
```

**Benefits**:
- Optimized for concurrent reads (typical in caching scenarios)
- No lock contention for read operations
- Atomic operations for consistency

### 3. RegexCache

Cache for pre-compiled regular expressions to avoid repeated compilation.

```go
cache := performance.GetRegexCache()

// Get or compile regex
regex, err := cache.Get(`^\d{3}-\d{4}$`)

// Use the cached regex
if regex.MatchString("123-4567") {
    // Match found
}

// Pre-compile multiple patterns
patterns := []string{
    `^\d+$`,
    `^[a-z]+$`,
}
cache.Precompile(patterns)
```

**Benefits**:
- Regex compilation is expensive; caching saves CPU
- Thread-safe using sync.Map
- Automatic caching on first use

### 4. OptimizedValidator

Validation utilities with pre-compiled regex patterns.

```go
validator := performance.GetOptimizedValidator()

// Validate phone number
if validator.ValidatePhone("13800138000") {
    // Valid phone
}

// Validate email
if validator.ValidateEmail("user@example.com") {
    // Valid email
}

// Validate password strength
ok, msg := validator.ValidatePasswordStrength("MyPass123")

// Sanitize input
clean := validator.SanitizeInput(userInput)
```

**Benefits**:
- Pre-compiled regex patterns for common validations
- Faster than compiling regex on each validation
- Consistent validation logic

### 5. OptimizedSessionManager

WebSocket session manager using sync.Map for better concurrent performance.

```go
manager := performance.NewOptimizedSessionManager("multi", 3)

// Create session
session, kicked := manager.CreateSession(userID, deviceID, token, expiresAt, conn)

// Remove session
manager.RemoveSession(sessionID)

// Kick user (all sessions)
count := manager.KickUser(userID, "reason")

// Kick single session
ok := manager.KickSession(sessionID, "reason")

// List user sessions
sessions := manager.ListUserSessions(userID)

// Broadcast to user
count := manager.Broadcast(userID, message)
```

**Benefits**:
- Better concurrent performance for session lookups
- Supports single sign-on and multi-device strategies
- Thread-safe operations

### 6. ContextExtractor

Optimized context value extraction using type switches.

```go
extractor := performance.GetContextExtractor()

// Extract StoreID
storeID, err := extractor.GetStoreID(ctx)

// Extract UserID
userID, err := extractor.GetUserID(ctx)

// Extract string value
str, err := extractor.GetString(ctx, "key")

// Extract int value
num, err := extractor.GetInt(ctx, "key")

// Extract bool value
flag, err := extractor.GetBool(ctx, "key")
```

**Benefits**:
- Fast type assertions using type switches
- Proper error handling
- Consistent API

## Usage Guidelines

### When to Use Type Switches

Use type switches instead of reflection when:
- You know the possible types at compile time
- The operation is in a hot path (executed frequently)
- You need maximum performance

### When to Use sync.Map

Use `sync.Map` instead of mutex-protected maps when:
- The workload is read-heavy (more reads than writes)
- Keys are written once and read many times
- You need lock-free reads

**Note**: For write-heavy workloads, a regular map with `sync.RWMutex` may perform better.

### When to Pre-compile Regex

Pre-compile regex patterns when:
- The pattern is used multiple times
- The pattern is used in a hot path
- You want to avoid compilation overhead

## Performance Characteristics

| Operation | Time | Allocations | vs Alternative |
|-----------|------|-------------|----------------|
| Type Switch (ToUint) | ~2.7 ns | 0 | 10-100x faster than reflection |
| sync.Map Read | ~3.1 ns | 0 | 9x faster than mutex map |
| sync.Map Write | ~40 ns | 2 | 1.8x faster than mutex map |
| Regex Cache | Amortized O(1) | 0 (after first) | Avoids repeated compilation |

## Testing

The package includes comprehensive tests:

```bash
# Run all tests
go test ./pkg/performance/...

# Run property-based tests
go test -v ./pkg/performance/... -run TestProperty

# Run benchmarks
go test -bench=. -benchmem ./pkg/performance/...
```

## Property-Based Tests

The package includes property-based tests that verify:

1. **Property 34**: Type switch performance - Type switches correctly convert all numeric types
2. **Property 35**: sync.Map concurrent performance - Concurrent operations are thread-safe and performant

These tests run 100+ iterations with random inputs to ensure correctness across all scenarios.

## Integration

To use these optimizations in your code:

```go
import "github.com/Kevin-Jii/tower-go/pkg/performance"

// Use global instances
converter := performance.GetTypeConverter()
cache := performance.NewConcurrentCache()
validator := performance.GetOptimizedValidator()
extractor := performance.GetContextExtractor()
```

## Best Practices

1. **Reuse instances**: Create cache and validator instances once and reuse them
2. **Pre-compile patterns**: Use `Precompile()` for known regex patterns at startup
3. **Profile first**: Use benchmarks to verify optimizations help your specific use case
4. **Test thoroughly**: Property-based tests help catch edge cases

## References

- Requirements: 8.1 (Type switch optimization), 8.2 (sync.Map usage)
- Design: Hot Path Optimization section
- Properties: 34 (Type switch performance), 35 (sync.Map concurrent performance)
