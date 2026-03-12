package config

import (
	"os"
	"runtime"
	"testing"
	"time"
)

func TestLoadPerformanceConfig(t *testing.T) {
	// Test with default values
	perfConfig := loadPerformanceConfig()

	// Verify database config defaults
	cpuCount := runtime.NumCPU()
	expectedMaxOpenConns := cpuCount * 2
	expectedMaxIdleConns := cpuCount

	if perfConfig.Database.MaxOpenConns != expectedMaxOpenConns {
		t.Errorf("Expected MaxOpenConns to be %d, got %d", expectedMaxOpenConns, perfConfig.Database.MaxOpenConns)
	}

	if perfConfig.Database.MaxIdleConns != expectedMaxIdleConns {
		t.Errorf("Expected MaxIdleConns to be %d, got %d", expectedMaxIdleConns, perfConfig.Database.MaxIdleConns)
	}

	if perfConfig.Database.SlowQueryThreshold != 100*time.Millisecond {
		t.Errorf("Expected SlowQueryThreshold to be 100ms, got %v", perfConfig.Database.SlowQueryThreshold)
	}

	if perfConfig.Database.DefaultBatchSize != 100 {
		t.Errorf("Expected DefaultBatchSize to be 100, got %d", perfConfig.Database.DefaultBatchSize)
	}

	// Verify cache config defaults
	if !perfConfig.Cache.LocalCacheEnabled {
		t.Error("Expected LocalCacheEnabled to be true")
	}

	if perfConfig.Cache.LocalCacheSize != 10000 {
		t.Errorf("Expected LocalCacheSize to be 10000, got %d", perfConfig.Cache.LocalCacheSize)
	}

	if perfConfig.Cache.RedisPipelineSize != 100 {
		t.Errorf("Expected RedisPipelineSize to be 100, got %d", perfConfig.Cache.RedisPipelineSize)
	}

	// Verify concurrency config defaults
	expectedWorkerPoolSize := cpuCount * 4
	if perfConfig.Concurrency.WorkerPoolSize != expectedWorkerPoolSize {
		t.Errorf("Expected WorkerPoolSize to be %d, got %d", expectedWorkerPoolSize, perfConfig.Concurrency.WorkerPoolSize)
	}

	if perfConfig.Concurrency.MaxConcurrentRequests != 1000 {
		t.Errorf("Expected MaxConcurrentRequests to be 1000, got %d", perfConfig.Concurrency.MaxConcurrentRequests)
	}

	// Verify monitoring config defaults
	if !perfConfig.Monitoring.EnableMetrics {
		t.Error("Expected EnableMetrics to be true")
	}

	if perfConfig.Monitoring.SlowQueryLimit != 100 {
		t.Errorf("Expected SlowQueryLimit to be 100, got %d", perfConfig.Monitoring.SlowQueryLimit)
	}
}

func TestLoadPerformanceConfigFromEnv(t *testing.T) {
	// Set environment variables
	os.Setenv("PERF_DB_MAX_OPEN_CONNS", "50")
	os.Setenv("PERF_DB_SLOW_QUERY_MS", "200")
	os.Setenv("PERF_CACHE_LOCAL_SIZE", "5000")
	os.Setenv("PERF_WORKER_POOL_SIZE", "16")
	os.Setenv("PERF_METRICS_ENABLED", "false")

	defer func() {
		os.Unsetenv("PERF_DB_MAX_OPEN_CONNS")
		os.Unsetenv("PERF_DB_SLOW_QUERY_MS")
		os.Unsetenv("PERF_CACHE_LOCAL_SIZE")
		os.Unsetenv("PERF_WORKER_POOL_SIZE")
		os.Unsetenv("PERF_METRICS_ENABLED")
	}()

	perfConfig := loadPerformanceConfig()

	// Verify environment variables are loaded
	if perfConfig.Database.MaxOpenConns != 50 {
		t.Errorf("Expected MaxOpenConns to be 50, got %d", perfConfig.Database.MaxOpenConns)
	}

	if perfConfig.Database.SlowQueryThreshold != 200*time.Millisecond {
		t.Errorf("Expected SlowQueryThreshold to be 200ms, got %v", perfConfig.Database.SlowQueryThreshold)
	}

	if perfConfig.Cache.LocalCacheSize != 5000 {
		t.Errorf("Expected LocalCacheSize to be 5000, got %d", perfConfig.Cache.LocalCacheSize)
	}

	if perfConfig.Concurrency.WorkerPoolSize != 16 {
		t.Errorf("Expected WorkerPoolSize to be 16, got %d", perfConfig.Concurrency.WorkerPoolSize)
	}

	if perfConfig.Monitoring.EnableMetrics {
		t.Error("Expected EnableMetrics to be false")
	}
}

func TestGetPerformanceConfig(t *testing.T) {
	// Initialize config
	InitConfig()

	// Get performance config
	perfConfig := GetPerformanceConfig()

	// Verify it's not nil and has expected structure
	if perfConfig.Database.MaxOpenConns <= 0 {
		t.Error("Expected MaxOpenConns to be positive")
	}

	if perfConfig.Cache.LocalCacheSize <= 0 {
		t.Error("Expected LocalCacheSize to be positive")
	}

	if perfConfig.Concurrency.WorkerPoolSize <= 0 {
		t.Error("Expected WorkerPoolSize to be positive")
	}
}

func TestDatabasePerformanceConfigDefaults(t *testing.T) {
	config := loadDatabasePerformanceConfig()

	// Test connection pool defaults
	if config.MaxOpenConns <= 0 {
		t.Error("MaxOpenConns should be positive")
	}

	if config.MaxIdleConns <= 0 {
		t.Error("MaxIdleConns should be positive")
	}

	if config.ConnMaxLifetime <= 0 {
		t.Error("ConnMaxLifetime should be positive")
	}

	if config.ConnMaxIdleTime <= 0 {
		t.Error("ConnMaxIdleTime should be positive")
	}

	// Test query config defaults
	if config.SlowQueryThreshold <= 0 {
		t.Error("SlowQueryThreshold should be positive")
	}

	// Test batch config defaults
	if config.DefaultBatchSize <= 0 {
		t.Error("DefaultBatchSize should be positive")
	}

	if config.MaxBatchSize <= config.DefaultBatchSize {
		t.Error("MaxBatchSize should be greater than DefaultBatchSize")
	}
}

func TestCachePerformanceConfigDefaults(t *testing.T) {
	config := loadCachePerformanceConfig()

	// Test local cache defaults
	if config.LocalCacheSize <= 0 {
		t.Error("LocalCacheSize should be positive")
	}

	if config.LocalCacheTTL <= 0 {
		t.Error("LocalCacheTTL should be positive")
	}

	// Test Redis defaults
	if config.RedisPipelineSize <= 0 {
		t.Error("RedisPipelineSize should be positive")
	}

	if config.RedisPoolSize <= 0 {
		t.Error("RedisPoolSize should be positive")
	}
}

func TestConcurrencyConfigDefaults(t *testing.T) {
	config := loadConcurrencyConfig()

	// Test worker pool defaults
	if config.WorkerPoolSize <= 0 {
		t.Error("WorkerPoolSize should be positive")
	}

	if config.WorkerQueueSize <= 0 {
		t.Error("WorkerQueueSize should be positive")
	}

	if config.WorkerIdleTimeout <= 0 {
		t.Error("WorkerIdleTimeout should be positive")
	}

	// Test concurrency limits
	if config.MaxConcurrentRequests <= 0 {
		t.Error("MaxConcurrentRequests should be positive")
	}

	if config.MaxConcurrentQueries <= 0 {
		t.Error("MaxConcurrentQueries should be positive")
	}
}

func TestMonitoringConfigDefaults(t *testing.T) {
	config := loadMonitoringConfig()

	// Test metrics defaults
	if config.MetricsInterval <= 0 {
		t.Error("MetricsInterval should be positive")
	}

	if config.MetricsRetention <= 0 {
		t.Error("MetricsRetention should be positive")
	}

	// Test slow query defaults
	if config.SlowQueryLimit <= 0 {
		t.Error("SlowQueryLimit should be positive")
	}

	// Test pprof defaults
	if config.PprofPort <= 0 {
		t.Error("PprofPort should be positive")
	}
}
