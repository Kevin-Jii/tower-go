package config

import (
	"runtime"
	"time"
)

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
	RedisEnabled      bool // 是否启用Redis
	RedisPipelineSize int  // Pipeline批量大小
	RedisPoolSize     int  // 连接池大小

	// 缓存策略
	EnableCacheWarmup bool // 是否启用缓存预热
	EnableCacheLock   bool // 是否启用缓存锁
}

// ConcurrencyConfig 并发配置
type ConcurrencyConfig struct {
	// 工作池配置
	WorkerPoolSize    int           // 工作池大小
	WorkerQueueSize   int           // 工作队列大小
	WorkerIdleTimeout time.Duration // 工作器空闲超时

	// 并发限制
	MaxConcurrentRequests int // 最大并发请求数
	MaxConcurrentQueries  int // 最大并发查询数
}

// MonitoringConfig 监控配置
type MonitoringConfig struct {
	// 指标配置
	EnableMetrics    bool          // 是否启用指标
	MetricsInterval  time.Duration // 指标收集间隔
	MetricsRetention time.Duration // 指标保留时间

	// 慢查询配置
	SlowQueryEnabled bool // 是否启用慢查询日志
	SlowQueryLimit   int  // 慢查询保留数量

	// Pprof配置
	PprofEnabled bool // 是否启用pprof
	PprofPort    int  // pprof端口
}

// loadPerformanceConfig 加载性能配置
func loadPerformanceConfig() PerformanceConfig {
	return PerformanceConfig{
		Database:    loadDatabasePerformanceConfig(),
		Cache:       loadCachePerformanceConfig(),
		Concurrency: loadConcurrencyConfig(),
		Monitoring:  loadMonitoringConfig(),
	}
}

// loadDatabasePerformanceConfig 加载数据库性能配置
func loadDatabasePerformanceConfig() DatabasePerformanceConfig {
	// 根据CPU核心数计算默认连接数
	cpuCount := runtime.NumCPU()
	defaultMaxOpenConns := cpuCount * 2
	defaultMaxIdleConns := cpuCount

	return DatabasePerformanceConfig{
		// 连接池配置
		MaxOpenConns:    getAppInt("PERF_DB_MAX_OPEN_CONNS", defaultMaxOpenConns),
		MaxIdleConns:    getAppInt("PERF_DB_MAX_IDLE_CONNS", defaultMaxIdleConns),
		ConnMaxLifetime: time.Duration(getAppInt("PERF_DB_CONN_MAX_LIFETIME_MINUTES", 60)) * time.Minute,
		ConnMaxIdleTime: time.Duration(getAppInt("PERF_DB_CONN_MAX_IDLE_TIME_MINUTES", 10)) * time.Minute,

		// 查询配置
		SlowQueryThreshold: time.Duration(getAppInt("PERF_DB_SLOW_QUERY_MS", 100)) * time.Millisecond,
		EnableQueryLog:     getAppBool("PERF_DB_ENABLE_QUERY_LOG", false),
		EnableExplain:      getAppBool("PERF_DB_ENABLE_EXPLAIN", false),

		// 批处理配置
		DefaultBatchSize: getAppInt("PERF_DB_DEFAULT_BATCH_SIZE", 100),
		MaxBatchSize:     getAppInt("PERF_DB_MAX_BATCH_SIZE", 1000),
	}
}

// loadCachePerformanceConfig 加载缓存性能配置
func loadCachePerformanceConfig() CachePerformanceConfig {
	return CachePerformanceConfig{
		// 本地缓存配置
		LocalCacheEnabled: getAppBool("PERF_CACHE_LOCAL_ENABLED", true),
		LocalCacheSize:    getAppInt("PERF_CACHE_LOCAL_SIZE", 10000),
		LocalCacheTTL:     time.Duration(getAppInt("PERF_CACHE_LOCAL_TTL_SECONDS", 300)) * time.Second,

		// Redis配置
		RedisEnabled:      getAppBool("PERF_CACHE_REDIS_ENABLED", true),
		RedisPipelineSize: getAppInt("PERF_CACHE_REDIS_PIPELINE_SIZE", 100),
		RedisPoolSize:     getAppInt("PERF_CACHE_REDIS_POOL_SIZE", 10),

		// 缓存策略
		EnableCacheWarmup: getAppBool("PERF_CACHE_ENABLE_WARMUP", false),
		EnableCacheLock:   getAppBool("PERF_CACHE_ENABLE_LOCK", true),
	}
}

// loadConcurrencyConfig 加载并发配置
func loadConcurrencyConfig() ConcurrencyConfig {
	// 根据CPU核心数计算默认工作池大小
	cpuCount := runtime.NumCPU()
	defaultWorkerPoolSize := cpuCount * 4

	return ConcurrencyConfig{
		// 工作池配置
		WorkerPoolSize:    getAppInt("PERF_WORKER_POOL_SIZE", defaultWorkerPoolSize),
		WorkerQueueSize:   getAppInt("PERF_WORKER_QUEUE_SIZE", 1000),
		WorkerIdleTimeout: time.Duration(getAppInt("PERF_WORKER_IDLE_TIMEOUT_SECONDS", 60)) * time.Second,

		// 并发限制
		MaxConcurrentRequests: getAppInt("PERF_MAX_CONCURRENT_REQUESTS", 1000),
		MaxConcurrentQueries:  getAppInt("PERF_MAX_CONCURRENT_QUERIES", 100),
	}
}

// loadMonitoringConfig 加载监控配置
func loadMonitoringConfig() MonitoringConfig {
	return MonitoringConfig{
		// 指标配置
		EnableMetrics:    getAppBool("PERF_METRICS_ENABLED", true),
		MetricsInterval:  time.Duration(getAppInt("PERF_METRICS_INTERVAL_SECONDS", 60)) * time.Second,
		MetricsRetention: time.Duration(getAppInt("PERF_METRICS_RETENTION_HOURS", 24)) * time.Hour,

		// 慢查询配置
		SlowQueryEnabled: getAppBool("PERF_SLOW_QUERY_ENABLED", true),
		SlowQueryLimit:   getAppInt("PERF_SLOW_QUERY_LIMIT", 100),

		// Pprof配置
		PprofEnabled: getAppBool("PERF_PPROF_ENABLED", false),
		PprofPort:    getAppInt("PERF_PPROF_PORT", 6060),
	}
}

// GetPerformanceConfig 获取性能配置
func GetPerformanceConfig() PerformanceConfig {
	return GetConfig().Performance
}
