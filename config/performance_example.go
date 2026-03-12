package config

import (
	"fmt"
)

// ExampleUsage demonstrates how to use the performance configuration
func ExampleUsage() {
	// Initialize configuration (loads from environment variables or uses defaults)
	InitConfig()

	// Get the performance configuration
	perfConfig := GetPerformanceConfig()

	// Access database performance settings
	fmt.Printf("Database Max Open Connections: %d\n", perfConfig.Database.MaxOpenConns)
	fmt.Printf("Database Max Idle Connections: %d\n", perfConfig.Database.MaxIdleConns)
	fmt.Printf("Slow Query Threshold: %v\n", perfConfig.Database.SlowQueryThreshold)
	fmt.Printf("Default Batch Size: %d\n", perfConfig.Database.DefaultBatchSize)

	// Access cache performance settings
	fmt.Printf("Local Cache Enabled: %v\n", perfConfig.Cache.LocalCacheEnabled)
	fmt.Printf("Local Cache Size: %d\n", perfConfig.Cache.LocalCacheSize)
	fmt.Printf("Redis Pipeline Size: %d\n", perfConfig.Cache.RedisPipelineSize)

	// Access concurrency settings
	fmt.Printf("Worker Pool Size: %d\n", perfConfig.Concurrency.WorkerPoolSize)
	fmt.Printf("Max Concurrent Requests: %d\n", perfConfig.Concurrency.MaxConcurrentRequests)

	// Access monitoring settings
	fmt.Printf("Metrics Enabled: %v\n", perfConfig.Monitoring.EnableMetrics)
	fmt.Printf("Slow Query Enabled: %v\n", perfConfig.Monitoring.SlowQueryEnabled)
	fmt.Printf("Pprof Enabled: %v\n", perfConfig.Monitoring.PprofEnabled)
}

// ExampleDatabaseConfig demonstrates database performance configuration usage
func ExampleDatabaseConfig() {
	perfConfig := GetPerformanceConfig()
	dbConfig := perfConfig.Database

	// Use in database connection setup
	fmt.Printf("Setting up database with:\n")
	fmt.Printf("  Max Open Connections: %d\n", dbConfig.MaxOpenConns)
	fmt.Printf("  Max Idle Connections: %d\n", dbConfig.MaxIdleConns)
	fmt.Printf("  Connection Max Lifetime: %v\n", dbConfig.ConnMaxLifetime)
	fmt.Printf("  Connection Max Idle Time: %v\n", dbConfig.ConnMaxIdleTime)

	// Use in query optimization
	if dbConfig.EnableQueryLog {
		fmt.Println("Query logging is enabled")
	}

	if dbConfig.EnableExplain {
		fmt.Println("Query EXPLAIN is enabled")
	}

	// Use in batch processing
	fmt.Printf("Using batch size: %d (max: %d)\n", dbConfig.DefaultBatchSize, dbConfig.MaxBatchSize)
}

// ExampleCacheConfig demonstrates cache performance configuration usage
func ExampleCacheConfig() {
	perfConfig := GetPerformanceConfig()
	cacheConfig := perfConfig.Cache

	// Use in cache setup
	if cacheConfig.LocalCacheEnabled {
		fmt.Printf("Setting up local cache with size: %d, TTL: %v\n",
			cacheConfig.LocalCacheSize, cacheConfig.LocalCacheTTL)
	}

	if cacheConfig.RedisEnabled {
		fmt.Printf("Setting up Redis cache with pool size: %d, pipeline size: %d\n",
			cacheConfig.RedisPoolSize, cacheConfig.RedisPipelineSize)
	}

	// Use cache strategies
	if cacheConfig.EnableCacheLock {
		fmt.Println("Cache lock is enabled to prevent cache stampede")
	}

	if cacheConfig.EnableCacheWarmup {
		fmt.Println("Cache warmup is enabled")
	}
}

// ExampleConcurrencyConfig demonstrates concurrency configuration usage
func ExampleConcurrencyConfig() {
	perfConfig := GetPerformanceConfig()
	concurrencyConfig := perfConfig.Concurrency

	// Use in worker pool setup
	fmt.Printf("Setting up worker pool with:\n")
	fmt.Printf("  Pool Size: %d workers\n", concurrencyConfig.WorkerPoolSize)
	fmt.Printf("  Queue Size: %d tasks\n", concurrencyConfig.WorkerQueueSize)
	fmt.Printf("  Idle Timeout: %v\n", concurrencyConfig.WorkerIdleTimeout)

	// Use in concurrency limiting
	fmt.Printf("Max concurrent requests: %d\n", concurrencyConfig.MaxConcurrentRequests)
	fmt.Printf("Max concurrent queries: %d\n", concurrencyConfig.MaxConcurrentQueries)
}

// ExampleMonitoringConfig demonstrates monitoring configuration usage
func ExampleMonitoringConfig() {
	perfConfig := GetPerformanceConfig()
	monitoringConfig := perfConfig.Monitoring

	// Use in metrics collection
	if monitoringConfig.EnableMetrics {
		fmt.Printf("Metrics collection enabled:\n")
		fmt.Printf("  Interval: %v\n", monitoringConfig.MetricsInterval)
		fmt.Printf("  Retention: %v\n", monitoringConfig.MetricsRetention)
	}

	// Use in slow query logging
	if monitoringConfig.SlowQueryEnabled {
		fmt.Printf("Slow query logging enabled, keeping last %d queries\n",
			monitoringConfig.SlowQueryLimit)
	}

	// Use in pprof setup
	if monitoringConfig.PprofEnabled {
		fmt.Printf("Pprof enabled on port: %d\n", monitoringConfig.PprofPort)
	}
}
