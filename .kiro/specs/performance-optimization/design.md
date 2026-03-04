# Design Document

## Overview

本设计文档描述了 Tower-Go 系统的性能优化方案,涵盖数据库查询优化、缓存策略、内存管理、并发处理等多个维度。优化目标是在保持代码可维护性的前提下,显著提升系统的吞吐量、降低延迟并减少资源消耗。

核心设计原则:
- 渐进式优化: 优先优化热路径和高频操作
- 可测量性: 所有优化都应该有明确的性能指标
- 向后兼容: 优化不应破坏现有API和功能
- 可观测性: 提供完善的监控和诊断工具

## Architecture

### 整体架构

```
┌─────────────────────────────────────────────────────────────┐
│                      Application Layer                       │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ Controllers  │  │  Services    │  │  Middleware  │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
                            │
┌─────────────────────────────────────────────────────────────┐
│                   Performance Layer (NEW)                    │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ Query        │  │ Cache        │  │ Memory       │      │
│  │ Optimizer    │  │ Manager      │  │ Pool         │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ Batch        │  │ Worker       │  │ Metrics      │      │
│  │ Processor    │  │ Pool         │  │ Collector    │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
                            │
┌─────────────────────────────────────────────────────────────┐
│                      Data Access Layer                       │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ Enhanced     │  │ Connection   │  │ Index        │      │
│  │ Query        │  │ Pool         │  │ Strategy     │      │
│  │ Builder      │  │              │  │              │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
                            │
┌─────────────────────────────────────────────────────────────┐
│                      Storage Layer                           │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ MySQL        │  │ Redis        │  │ Local        │      │
│  │ (Optimized)  │  │ (Multi-tier) │  │ Cache        │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
```

### 性能优化层组件

1. **Query Optimizer**: 查询优化器,自动分析和优化SQL查询
2. **Cache Manager**: 多级缓存管理器,支持本地缓存和Redis
3. **Memory Pool**: 对象池管理器,复用高频对象
4. **Batch Processor**: 批处理器,优化批量操作
5. **Worker Pool**: 工作池,管理并发任务
6. **Metrics Collector**: 指标收集器,收集性能数据

## Components and Interfaces

### 1. Enhanced Query Builder

增强的查询构造器,在现有基础上添加智能优化功能。

```go
// QueryOptimizer 查询优化器接口
type QueryOptimizer interface {
    // OptimizeQuery 优化查询
    OptimizeQuery(query *gorm.DB) *gorm.DB
    
    // AnalyzeIndexUsage 分析索引使用情况
    AnalyzeIndexUsage(query string) (*IndexAnalysis, error)
    
    // SuggestIndex 建议索引
    SuggestIndex(tableName string, conditions []string) []string
}

// CursorPaginator 游标分页器
type CursorPaginator interface {
    // Paginate 执行游标分页
    Paginate(query *gorm.DB, cursor string, limit int) (*CursorResult, error)
    
    // EncodeCursor 编码游标
    EncodeCursor(lastID uint, lastValue interface{}) string
    
    // DecodeCursor 解码游标
    DecodeCursor(cursor string) (uint, interface{}, error)
}

// QueryCache 查询缓存
type QueryCache interface {
    // Get 获取缓存的查询结果
    Get(key string, dest interface{}) error
    
    // Set 设置查询结果缓存
    Set(key string, value interface{}, ttl time.Duration) error
    
    // Invalidate 失效缓存
    Invalidate(pattern string) error
}
```

### 2. Multi-tier Cache Manager

多级缓存管理器,支持本地缓存和Redis缓存。

```go
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
```

### 3. Memory Pool Manager

内存池管理器,复用高频对象。

```go
// MemoryPool 内存池接口
type MemoryPool interface {
    // Get 从池中获取对象
    Get() interface{}
    
    // Put 归还对象到池中
    Put(obj interface{})
    
    // Stats 获取池统计信息
    Stats() PoolStats
}

// ResponsePool 响应对象池
type ResponsePool interface {
    // GetResponse 获取响应对象
    GetResponse() *Response
    
    // PutResponse 归还响应对象
    PutResponse(resp *Response)
}

// BufferPool 缓冲区池
type BufferPool interface {
    // GetBuffer 获取缓冲区
    GetBuffer() *bytes.Buffer
    
    // PutBuffer 归还缓冲区
    PutBuffer(buf *bytes.Buffer)
}
```

### 4. Batch Processor

批处理器,优化批量操作。

```go
// BatchProcessor 批处理器接口
type BatchProcessor interface {
    // BatchInsert 批量插入
    BatchInsert(ctx context.Context, records interface{}, batchSize int) error
    
    // BatchUpdate 批量更新
    BatchUpdate(ctx context.Context, updates []BatchUpdate) error
    
    // BatchDelete 批量软删除
    BatchDelete(ctx context.Context, tableName string, ids []uint) error
    
    // Process 通用批处理
    Process(ctx context.Context, items []interface{}, 
            processor func(batch []interface{}) error, batchSize int) error
}

// BatchUpdate 批量更新项
type BatchUpdate struct {
    ID     uint
    Fields map[string]interface{}
}

// ChunkProcessor 分块处理器
type ChunkProcessor interface {
    // ProcessInChunks 分块处理
    ProcessInChunks(ctx context.Context, totalCount int, chunkSize int,
                    processor func(offset, limit int) error) error
}
```

### 5. Worker Pool

工作池,管理并发任务。

```go
// WorkerPool 工作池接口
type WorkerPool interface {
    // Submit 提交任务
    Submit(task func()) error
    
    // SubmitWithContext 提交带上下文的任务
    SubmitWithContext(ctx context.Context, task func(context.Context)) error
    
    // Wait 等待所有任务完成
    Wait()
    
    // Stop 停止工作池
    Stop()
    
    // Stats 获取工作池统计
    Stats() WorkerPoolStats
}

// Semaphore 信号量接口
type Semaphore interface {
    // Acquire 获取信号量
    Acquire(ctx context.Context) error
    
    // Release 释放信号量
    Release()
    
    // TryAcquire 尝试获取信号量
    TryAcquire() bool
}
```

### 6. Metrics Collector

指标收集器,收集性能数据。

```go
// MetricsCollector 指标收集器接口
type MetricsCollector interface {
    // RecordQuery 记录查询
    RecordQuery(query string, duration time.Duration, err error)
    
    // RecordCache 记录缓存操作
    RecordCache(operation string, hit bool, duration time.Duration)
    
    // RecordRequest 记录请求
    RecordRequest(method, path string, status int, duration time.Duration)
    
    // GetMetrics 获取指标
    GetMetrics() *Metrics
}

// SlowQueryLogger 慢查询日志器
type SlowQueryLogger interface {
    // LogSlowQuery 记录慢查询
    LogSlowQuery(query string, duration time.Duration, explain string)
    
    // GetSlowQueries 获取慢查询列表
    GetSlowQueries(limit int) []SlowQuery
}
```

## Data Models

### 性能配置模型

```go
// PerformanceConfig 性能配置
type PerformanceConfig struct {
    // 数据库配置
    Database DatabasePerformanceConfig
    
    // 缓存配置
    Cache CachePerformanceConfig
    
    // 并发配置
    Concurrency ConcurrencyConfig
    
    // 监控配置
    Monitoring MonitoringConfig
}

// DatabasePerformanceConfig 数据库性能配置
type DatabasePerformanceConfig struct {
    // 连接池配置
    MaxOpenConns    int           // 最大打开连接数
    MaxIdleConns    int           // 最大空闲连接数
    ConnMaxLifetime time.Duration // 连接最大生命周期
    ConnMaxIdleTime time.Duration // 连接最大空闲时间
    
    // 查询配置
    SlowQueryThreshold time.Duration // 慢查询阈值
    EnableQueryLog     bool          // 是否启用查询日志
    EnableExplain      bool          // 是否启用EXPLAIN
    
    // 批处理配置
    DefaultBatchSize int // 默认批处理大小
    MaxBatchSize     int // 最大批处理大小
}

// CachePerformanceConfig 缓存性能配置
type CachePerformanceConfig struct {
    // 本地缓存配置
    LocalCacheEnabled bool          // 是否启用本地缓存
    LocalCacheSize    int           // 本地缓存大小
    LocalCacheTTL     time.Duration // 本地缓存TTL
    
    // Redis配置
    RedisEnabled      bool          // 是否启用Redis
    RedisPipelineSize int           // Pipeline批量大小
    RedisPoolSize     int           // 连接池大小
    
    // 缓存策略
    EnableCacheWarmup bool // 是否启用缓存预热
    EnableCacheLock   bool // 是否启用缓存锁
}

// ConcurrencyConfig 并发配置
type ConcurrencyConfig struct {
    // 工作池配置
    WorkerPoolSize    int // 工作池大小
    WorkerQueueSize   int // 工作队列大小
    WorkerIdleTimeout time.Duration
    
    // 并发限制
    MaxConcurrentRequests int // 最大并发请求数
    MaxConcurrentQueries  int // 最大并发查询数
}

// MonitoringConfig 监控配置
type MonitoringConfig struct {
    // 指标配置
    EnableMetrics     bool          // 是否启用指标
    MetricsInterval   time.Duration // 指标收集间隔
    MetricsRetention  time.Duration // 指标保留时间
    
    // 慢查询配置
    SlowQueryEnabled  bool // 是否启用慢查询日志
    SlowQueryLimit    int  // 慢查询保留数量
    
    // Pprof配置
    PprofEnabled bool   // 是否启用pprof
    PprofPort    int    // pprof端口
}
```

### 性能指标模型

```go
// Metrics 性能指标
type Metrics struct {
    // 请求指标
    RequestCount  int64         // 请求总数
    RequestQPS    float64       // 每秒请求数
    AvgLatency    time.Duration // 平均延迟
    P95Latency    time.Duration // P95延迟
    P99Latency    time.Duration // P99延迟
    ErrorRate     float64       // 错误率
    
    // 数据库指标
    QueryCount    int64         // 查询总数
    QueryQPS      float64       // 每秒查询数
    SlowQueryCount int64        // 慢查询数
    AvgQueryTime  time.Duration // 平均查询时间
    
    // 缓存指标
    CacheHitCount  int64   // 缓存命中数
    CacheMissCount int64   // 缓存未命中数
    CacheHitRate   float64 // 缓存命中率
    
    // 资源指标
    GoroutineCount int     // Goroutine数量
    MemoryUsage    uint64  // 内存使用量
    CPUUsage       float64 // CPU使用率
}

// SlowQuery 慢查询记录
type SlowQuery struct {
    Query     string        // 查询语句
    Duration  time.Duration // 执行时间
    Explain   string        // 执行计划
    Timestamp time.Time     // 时间戳
}

// IndexAnalysis 索引分析结果
type IndexAnalysis struct {
    TableName    string   // 表名
    UsedIndexes  []string // 使用的索引
    MissingIndexes []string // 缺失的索引
    Suggestions  []string // 优化建议
}
```

## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system-essentially, a formal statement about what the system should do. 
Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*

### Database Query Optimization Properties

Property 1: Index performance improvement
*For any* store accounting statistics query with date range and store_id filters, when executed with composite indexes, the query time should be less than 100ms
**Validates: Requirements 1.1**

Property 2: Inventory uniqueness and performance
*For any* inventory query with store_id and product_id combination, the unique composite index should ensure uniqueness and query time should be faster than without index
**Validates: Requirements 1.2**

Property 3: N+1 query elimination
*For any* purchase order list query with N orders, the total number of SQL queries executed should be constant (typically 2-3) regardless of N, not N+1
**Validates: Requirements 1.3**

Property 4: Prefix index optimization
*For any* LIKE search query with prefix pattern (e.g., "abc%"), the query should use the prefix index and perform faster than full table scan
**Validates: Requirements 1.5**

### Query Builder Enhancement Properties

Property 5: Automatic index strategy
*For any* query built with multiple conditions, the Query Builder should generate SQL that uses the most optimal index available
**Validates: Requirements 2.1**

Property 6: Cursor pagination performance
*For any* large dataset pagination (>10000 records), cursor-based pagination should perform significantly faster than OFFSET-based pagination
**Validates: Requirements 2.2**

Property 7: JOIN deduplication
*For any* complex query that joins the same table multiple times, the Query Builder should automatically merge redundant JOINs
**Validates: Requirements 2.3**

Property 8: Aggregation caching
*For any* aggregation query, when executed twice with the same parameters, the second execution should retrieve results from cache
**Validates: Requirements 2.4**

Property 9: Empty condition short-circuit
*For any* query where all conditions are empty or nil, the Query Builder should return early without executing a database query
**Validates: Requirements 2.5**

### Cache Strategy Properties

Property 10: Multi-tier cache access
*For any* hot data access, the system should check local cache first, then Redis, then database, in that order
**Validates: Requirements 3.1**

Property 11: Cache stampede prevention
*For any* cache expiration under high concurrency, only one request should load data from database while others wait
**Validates: Requirements 3.2**

Property 12: Pipeline batch operations
*For any* batch cache query of N keys, using Pipeline should result in 1 network round-trip instead of N
**Validates: Requirements 3.3**

Property 13: Delayed double deletion
*For any* data update operation, the cache should be deleted both before and after the database update
**Validates: Requirements 3.4**

Property 14: LRU eviction
*For any* cache at full capacity, when adding new items, the least recently used items should be evicted first
**Validates: Requirements 3.5**

### Batch Processing Properties

Property 15: Batch insert efficiency
*For any* batch insert of N records, the system should execute 1 SQL statement instead of N separate INSERT statements
**Validates: Requirements 4.1**

Property 16: Batch update with CASE WHEN
*For any* batch update of N records with different values, the system should use a single UPDATE with CASE WHEN instead of N UPDATE statements
**Validates: Requirements 4.2**

Property 17: Chunked processing memory safety
*For any* large dataset processing, memory usage should remain bounded regardless of total dataset size
**Validates: Requirements 4.3**

Property 18: Partial success and resume
*For any* batch operation that fails mid-way, the system should record progress and allow resuming from the failure point
**Validates: Requirements 4.4**

Property 19: Soft delete batch operations
*For any* batch delete operation, records should be marked with deleted_at timestamp instead of being physically removed
**Validates: Requirements 4.5**

### Connection Pool Properties

Property 20: Idle connection reclamation
*For any* database connection idle beyond the configured timeout, the connection should be automatically closed and removed from the pool
**Validates: Requirements 5.2**

Property 21: Connection pool queue
*For any* request when connection pool is exhausted, the request should wait in queue instead of immediately failing
**Validates: Requirements 5.3**

Property 22: Unhealthy connection handling
*For any* database connection that encounters an error, the connection should be marked as unhealthy and automatically reconnected
**Validates: Requirements 5.4**

Property 23: Long transaction monitoring
*For any* database transaction exceeding the configured timeout threshold, an alert should be logged
**Validates: Requirements 5.5**

### Memory Optimization Properties

Property 24: Object pool reuse
*For any* high-frequency request handling, response objects should be reused from pool, reducing memory allocations
**Validates: Requirements 6.1**

Property 25: JSON encoder pooling
*For any* JSON serialization operation, encoders and buffers should be reused from pool
**Validates: Requirements 6.2**

Property 26: Streaming large files
*For any* large file processing, memory usage should remain constant regardless of file size
**Validates: Requirements 6.3**

Property 27: String builder efficiency
*For any* string concatenation operation, using strings.Builder should result in fewer memory allocations than using + operator
**Validates: Requirements 6.4**

Property 28: Slice pre-allocation
*For any* slice operation where final size is known, pre-allocating capacity should reduce the number of reallocations
**Validates: Requirements 6.5**

### Concurrency Properties

Property 29: Worker pool parallelism
*For any* set of N independent tasks, using worker pool should complete them faster than sequential processing
**Validates: Requirements 7.1**

Property 30: Read-write lock efficiency
*For any* shared resource with concurrent access, using RWMutex should allow multiple concurrent reads
**Validates: Requirements 7.2**

Property 31: Async task non-blocking
*For any* async task submission, the calling goroutine should return immediately without waiting for task completion
**Validates: Requirements 7.3**

Property 32: Semaphore concurrency limit
*For any* high concurrency scenario, the number of concurrent operations should not exceed the configured semaphore limit
**Validates: Requirements 7.4**

Property 33: Context-based goroutine lifecycle
*For any* goroutine spawned with context, when context is cancelled, the goroutine should terminate gracefully
**Validates: Requirements 7.5**

### Hot Path Optimization Properties

Property 34: Type switch performance
*For any* frequent type assertion operation, using type switch should be faster than using reflection
**Validates: Requirements 8.1**

Property 35: sync.Map concurrent performance
*For any* concurrent map operations with read-heavy workload, sync.Map should perform better than mutex-protected map
**Validates: Requirements 8.2**

Property 36: Time format caching
*For any* frequent time formatting operation, caching results should reduce redundant computations
**Validates: Requirements 8.3**

Property 37: Regex pre-compilation
*For any* frequent regex matching, pre-compiled regex should perform faster than compiling on each use
**Validates: Requirements 8.4**

Property 38: Fast JSON library
*For any* frequent JSON parsing operation, using sonic library should be faster than standard library
**Validates: Requirements 8.5**

### Monitoring Properties

Property 39: Slow query logging
*For any* database query exceeding 100ms, the query should be logged with execution time and query text
**Validates: Requirements 9.1**

Property 40: Query plan recording
*For any* database operation when explain is enabled, the execution plan should be recorded for analysis
**Validates: Requirements 9.2**

### Data Structure Properties

Property 41: Map vs slice lookup
*For any* lookup operation, using map should have O(1) time complexity while slice traversal has O(n)
**Validates: Requirements 10.1**

Property 42: Ordered traversal
*For any* ordered map or skip list, iteration should return elements in sorted order
**Validates: Requirements 10.2**

Property 43: Set memory efficiency
*For any* set implementation, using map[T]struct{} should consume less memory than map[T]bool
**Validates: Requirements 10.3**

Property 44: Heap priority queue
*For any* priority queue operation, heap-based implementation should have O(log n) complexity while sorted slice has O(n log n)
**Validates: Requirements 10.4**

Property 45: LRU cache O(1) operations
*For any* LRU cache using doubly-linked list and hash map, both get and put operations should have O(1) time complexity
**Validates: Requirements 10.5**

## Error Handling

### 错误处理策略

1. **查询优化错误**
   - 索引创建失败: 记录错误但不影响系统运行
   - 查询超时: 返回超时错误,记录慢查询日志
   - 连接池耗尽: 等待可用连接或返回服务繁忙错误

2. **缓存错误**
   - 缓存未命中: 降级到数据库查询
   - Redis连接失败: 降级到仅使用本地缓存
   - 缓存序列化失败: 记录错误,跳过缓存

3. **批处理错误**
   - 部分记录失败: 记录失败记录,继续处理其他记录
   - 事务失败: 回滚整个批次,记录错误
   - 内存不足: 减小批次大小重试

4. **并发错误**
   - Goroutine panic: 使用recover捕获,记录错误
   - 死锁检测: 设置超时,超时后返回错误
   - 资源耗尽: 拒绝新请求,返回服务繁忙

### 错误恢复机制

```go
// ErrorRecovery 错误恢复接口
type ErrorRecovery interface {
    // Retry 重试机制
    Retry(ctx context.Context, operation func() error, maxRetries int) error
    
    // CircuitBreaker 熔断器
    CircuitBreaker(operation func() error) error
    
    // Fallback 降级处理
    Fallback(primary func() error, fallback func() error) error
}
```

## Testing Strategy

### 测试方法

本项目采用双重测试策略:

1. **单元测试**: 验证具体实现的正确性
   - 测试各个组件的基本功能
   - 测试边界条件和错误处理
   - 测试与现有代码的集成

2. **属性测试**: 验证通用属性在所有输入下都成立
   - 使用 [go-fuzz](https://github.com/dvyukov/go-fuzz) 或 [gopter](https://github.com/leanovate/gopter) 进行属性测试
   - 每个属性测试至少运行100次迭代
   - 生成随机输入验证性能属性

### 性能测试库

使用 [gopter](https://github.com/leanovate/gopter) 作为属性测试库,原因:
- 原生Go实现,无需外部依赖
- 支持自定义生成器
- 良好的错误报告
- 支持收缩(shrinking)找到最小失败用例

### 基准测试

使用Go标准库的testing.B进行基准测试:
- 对比优化前后的性能差异
- 测量内存分配和GC压力
- 使用benchstat分析结果

### 测试标注规范

每个属性测试必须使用以下格式标注:

```go
// **Feature: performance-optimization, Property 1: Index performance improvement**
func TestProperty_IndexPerformance(t *testing.T) {
    // 测试实现
}
```

### 测试覆盖目标

- 单元测试覆盖率: >80%
- 属性测试: 覆盖所有45个correctness properties
- 基准测试: 覆盖所有性能关键路径
- 集成测试: 验证端到端性能提升

### 测试环境

- 开发环境: 本地MySQL + Redis
- CI环境: Docker容器化测试
- 性能测试: 独立的性能测试环境,模拟生产负载
