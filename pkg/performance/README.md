# Hot Path Performance Optimizations

This package provides optimized utilities for hot path operations in the Tower-Go system, focusing on type assertions and concurrent map operations.

## Components

### 1. QueryOptimizer

Intelligent query analyzer and optimizer that detects performance issues and suggests improvements.

**Features**:
- Automatic index usage analysis
- Detection of N+1 query problems
- OFFSET pagination warnings
- Duplicate JOIN detection
- Full table scan identification
- LIKE prefix optimization suggestions

```go
optimizer := performance.NewQueryOptimizer(db)

// Analyze a SQL query
result, err := optimizer.Analyze("SELECT * FROM store_accounts WHERE store_id = ?")

// Check for issues
for _, issue := range result.Issues {
    fmt.Printf("[%s] %s: %s\n", issue.Severity, issue.Type, issue.Message)
    fmt.Printf("Suggestion: %s\n", issue.Suggestion)
}

// Get recommendations
for _, rec := range result.Recommendations {
    fmt.Println("Recommendation:", rec)
}

// Check index usage
fmt.Printf("Used indexes: %v\n", result.IndexUsage.UsedIndexes)
fmt.Printf("Missing indexes: %v\n", result.IndexUsage.MissingIndexes)
fmt.Printf("Table scan: %v\n", result.IndexUsage.TableScan)
```

**Issue Detection**:
- `missing_index`: Fields that may benefit from indexes
- `offset_pagination`: OFFSET usage in large datasets
- `full_table_scan`: Queries without WHERE clauses
- `like_prefix`: LIKE patterns starting with wildcards
- `duplicate_join`: Same table joined multiple times
- `select_all`: SELECT * usage

**Benefits**:
- Proactive performance issue detection
- Actionable optimization suggestions
- Index strategy recommendations
- Query complexity analysis

### 2. IndexAnalyzer

Analyzes index usage patterns and suggests optimal indexing strategies.

```go
analyzer := performance.NewIndexAnalyzer(db)

// Analyze index usage for a table
usage := analyzer.AnalyzeIndexUsage("store_accounts", []string{"store_id", "account_date"})

fmt.Printf("Potential indexes: %v\n", usage.PotentialIndexes)
fmt.Printf("Table scan detected: %v\n", usage.TableScan)
```

**Benefits**:
- Identifies missing indexes
- Suggests composite indexes
- Detects full table scans
- Provides index recommendations

### 3. JoinDeduplicator

Detects and prevents duplicate JOIN operations in queries.

```go
deduplicator := performance.NewJoinDeduplicator()

// Add tables to track
if deduplicator.Add("users") {
    // First time joining users table
}

if !deduplicator.Add("users") {
    // Duplicate JOIN detected!
}

// Reset for next query
deduplicator.Reset()
```

**Benefits**:
- Prevents redundant JOINs
- Reduces query complexity
- Improves query performance

### 4. TypeConverter

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

### 5. ConcurrentCache

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

### 6. RegexCache

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

### 7. OptimizedValidator

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

### 8. OptimizedSessionManager

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

### 9. ContextExtractor

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
| Query Analysis | ~100 μs | Minimal | Prevents slow queries |
| Index Analysis | O(n) | Minimal | Identifies missing indexes |
| JOIN Deduplication | O(1) | 0 | Prevents redundant JOINs |
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

1. **Property 5**: Automatic index strategy - Query optimizer detects performance issues
2. **Property 7**: JOIN deduplication - Duplicate JOINs are detected and prevented
3. **Property 34**: Type switch performance - Type switches correctly convert all numeric types
4. **Property 35**: sync.Map concurrent performance - Concurrent operations are thread-safe and performant

These tests run 100+ iterations with random inputs to ensure correctness across all scenarios.

## Integration

To use these optimizations in your code:

```go
import "github.com/Kevin-Jii/tower-go/pkg/performance"

// Query optimization
optimizer := performance.NewQueryOptimizer(db)
result, _ := optimizer.Analyze(sqlQuery)

// Use global instances for hot path operations
converter := performance.GetTypeConverter()
cache := performance.NewConcurrentCache()
validator := performance.GetOptimizedValidator()
extractor := performance.GetContextExtractor()
```

## Best Practices

1. **Analyze queries early**: Use QueryOptimizer during development to catch issues
2. **Monitor index usage**: Regularly check IndexAnalyzer recommendations
3. **Reuse instances**: Create cache and validator instances once and reuse them
4. **Pre-compile patterns**: Use `Precompile()` for known regex patterns at startup
5. **Profile first**: Use benchmarks to verify optimizations help your specific use case
6. **Test thoroughly**: Property-based tests help catch edge cases

## References

- Requirements: 2.1 (Automatic index strategy), 2.3 (JOIN deduplication), 8.1 (Type switch optimization), 8.2 (sync.Map usage)
- Design: Query Builder Enhancement, Hot Path Optimization sections
- Properties: 5 (Automatic index strategy), 7 (JOIN deduplication), 34 (Type switch performance), 35 (sync.Map concurrent performance)
