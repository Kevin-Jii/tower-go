package database

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/Kevin-Jii/tower-go/config"
	"github.com/Kevin-Jii/tower-go/utils/logging"
	"go.uber.org/zap"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var db *gorm.DB

// DB 公开的数据库实例引用（兼容旧代码）
var DB *gorm.DB

var (
	// dbConfig 保存数据库配置用于重连
	dbConfig config.DatabaseConfig
	// dbMutex 保护数据库实例的并发访问
	dbMutex sync.RWMutex
	// healthCheckTicker 健康检查定时器
	healthCheckTicker *time.Ticker
	// healthCheckStop 停止健康检查的通道
	healthCheckStop chan struct{}
	// connectionStats 连接池统计信息
	connectionStats *ConnectionPoolStats
	// statsLock 保护统计信息的并发访问
	statsLock sync.RWMutex
)

// ConnectionPoolStats 连接池统计信息
type ConnectionPoolStats struct {
	// 基础统计
	MaxOpenConnections int           // 最大打开连接数
	OpenConnections    int           // 当前打开连接数
	InUse              int           // 正在使用的连接数
	Idle               int           // 空闲连接数
	WaitCount          int64         // 等待连接的总次数
	WaitDuration       time.Duration // 等待连接的总时间
	MaxIdleClosed      int64         // 因超过最大空闲时间而关闭的连接数
	MaxLifetimeClosed  int64         // 因超过最大生命周期而关闭的连接数

	// 健康检查统计
	LastHealthCheck   time.Time // 最后一次健康检查时间
	HealthCheckCount  int64     // 健康检查总次数
	HealthCheckFailed int64     // 健康检查失败次数
	LastHealthStatus  bool      // 最后一次健康检查状态

	// 重连统计
	ReconnectCount     int64     // 重连次数
	LastReconnectTime  time.Time // 最后一次重连时间
	LastReconnectError string    // 最后一次重连错误
}

// InitDB 初始化数据库连接
func InitDB(cfg config.DatabaseConfig) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	// 保存配置用于重连
	dbConfig = cfg

	// 初始化统计信息
	if connectionStats == nil {
		connectionStats = &ConnectionPoolStats{
			LastHealthCheck:  time.Now(),
			LastHealthStatus: false,
		}
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: false,
		},
	})

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// 同步更新 DB 引用
	DB = db

	// 获取底层 sql.DB 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// 使用性能配置优化连接池
	if err := configureConnectionPool(sqlDB); err != nil {
		return fmt.Errorf("failed to configure connection pool: %w", err)
	}

	// 启动健康检查
	StartHealthCheck()

	logging.LogInfo("Database connected successfully",
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
		zap.String("database", cfg.Database),
	)

	return nil
}

// configureConnectionPool 配置数据库连接池
// 根据CPU核心数和性能配置动态优化连接池参数
func configureConnectionPool(sqlDB *sql.DB) error {
	// 获取性能配置
	perfConfig := config.GetPerformanceConfig()
	dbConfig := perfConfig.Database

	// 配置最大打开连接数
	// 默认值基于CPU核心数计算: CPU核心数 * 2
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)

	// 配置最大空闲连接数
	// 默认值基于CPU核心数计算: CPU核心数
	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)

	// 配置连接最大生命周期
	// 防止连接长时间存活导致的问题(如数据库重启、网络问题等)
	sqlDB.SetConnMaxLifetime(dbConfig.ConnMaxLifetime)

	// 配置连接最大空闲时间
	// 自动回收空闲超时的连接,释放资源
	sqlDB.SetConnMaxIdleTime(dbConfig.ConnMaxIdleTime)

	logging.LogInfo("Connection pool configured",
		zap.Int("maxOpenConns", dbConfig.MaxOpenConns),
		zap.Int("maxIdleConns", dbConfig.MaxIdleConns),
		zap.Duration("connMaxLifetime", dbConfig.ConnMaxLifetime),
		zap.Duration("connMaxIdleTime", dbConfig.ConnMaxIdleTime),
	)

	return nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return db
}

// GetDBConn 获取底层 sql.DB
func GetDBConn() *sql.DB {
	if db == nil {
		return nil
	}
	sqlDB, _ := db.DB()
	return sqlDB
}

// Ping 检查数据库连接
func Ping() error {
	sqlDB := GetDBConn()
	if sqlDB == nil {
		return fmt.Errorf("database not initialized")
	}
	return sqlDB.Ping()
}

// CloseDB 关闭数据库连接
func CloseDB() error {
	sqlDB := GetDBConn()
	if sqlDB != nil {
		return sqlDB.Close()
	}
	return nil
}

// HealthCheck 数据库健康检查
func HealthCheck() error {
	if db == nil {
		return fmt.Errorf("database not initialized")
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// WithTransaction 事务处理
func WithTransaction(txFunc func(*gorm.DB) error) error {
	return db.Transaction(txFunc)
}

// SetMaxOpenConns 设置最大打开连接数
func SetMaxOpenConns(n int) {
	if sqlDB, err := db.DB(); err == nil {
		sqlDB.SetMaxOpenConns(n)
	}
}

// SetMaxIdleConns 设置最大空闲连接数
func SetMaxIdleConns(n int) {
	if sqlDB, err := db.DB(); err == nil {
		sqlDB.SetMaxIdleConns(n)
	}
}

// GetConnectionPoolStats 获取连接池统计信息
func GetConnectionPoolStats() (*ConnectionPoolStats, error) {
	sqlDB := GetDBConn()
	if sqlDB == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	stats := sqlDB.Stats()

	statsLock.Lock()
	defer statsLock.Unlock()

	// 更新基础统计信息
	connectionStats.MaxOpenConnections = stats.MaxOpenConnections
	connectionStats.OpenConnections = stats.OpenConnections
	connectionStats.InUse = stats.InUse
	connectionStats.Idle = stats.Idle
	connectionStats.WaitCount = stats.WaitCount
	connectionStats.WaitDuration = stats.WaitDuration
	connectionStats.MaxIdleClosed = stats.MaxIdleClosed
	connectionStats.MaxLifetimeClosed = stats.MaxLifetimeClosed

	// 返回统计信息的副本
	statsCopy := *connectionStats
	return &statsCopy, nil
}

// StartHealthCheck 启动连接池健康检查
func StartHealthCheck() {
	// 如果已经在运行，先停止
	StopHealthCheck()

	// 创建新的定时器和停止通道
	healthCheckTicker = time.NewTicker(30 * time.Second) // 每30秒检查一次
	healthCheckStop = make(chan struct{})

	// 启动健康检查goroutine
	go func() {
		logging.LogInfo("Connection pool health check started")

		for {
			select {
			case <-healthCheckTicker.C:
				performHealthCheck()
			case <-healthCheckStop:
				logging.LogInfo("Connection pool health check stopped")
				return
			}
		}
	}()
}

// StopHealthCheck 停止连接池健康检查
func StopHealthCheck() {
	if healthCheckTicker != nil {
		healthCheckTicker.Stop()
		healthCheckTicker = nil
	}

	if healthCheckStop != nil {
		close(healthCheckStop)
		healthCheckStop = nil
	}
}

// performHealthCheck 执行健康检查
func performHealthCheck() {
	// 确保统计信息已初始化
	statsLock.Lock()
	if connectionStats == nil {
		connectionStats = &ConnectionPoolStats{
			LastHealthCheck:  time.Now(),
			LastHealthStatus: false,
		}
	}
	connectionStats.HealthCheckCount++
	connectionStats.LastHealthCheck = time.Now()
	statsLock.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := PingWithContext(ctx)

	statsLock.Lock()
	if err != nil {
		connectionStats.HealthCheckFailed++
		connectionStats.LastHealthStatus = false
		statsLock.Unlock()

		logging.LogError("Database health check failed",
			zap.Error(err),
			zap.Int64("failedCount", connectionStats.HealthCheckFailed),
		)

		// 尝试重连
		if err := reconnect(); err != nil {
			logging.LogError("Database reconnection failed", zap.Error(err))
		}
	} else {
		connectionStats.LastHealthStatus = true
		statsLock.Unlock()

		// 记录连接池统计信息
		if stats, err := GetConnectionPoolStats(); err == nil {
			logging.LogDebug("Connection pool stats",
				zap.Int("openConnections", stats.OpenConnections),
				zap.Int("inUse", stats.InUse),
				zap.Int("idle", stats.Idle),
				zap.Int64("waitCount", stats.WaitCount),
			)
		}
	}
}

// reconnect 重新连接数据库
func reconnect() error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	logging.LogWarn("Attempting to reconnect to database")

	statsLock.Lock()
	if connectionStats == nil {
		connectionStats = &ConnectionPoolStats{
			LastHealthCheck:  time.Now(),
			LastHealthStatus: false,
		}
	}
	connectionStats.ReconnectCount++
	connectionStats.LastReconnectTime = time.Now()
	statsLock.Unlock()

	// 关闭旧连接
	if db != nil {
		if sqlDB, err := db.DB(); err == nil {
			sqlDB.Close()
		}
	}

	// 创建新连接
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Database,
	)

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: false,
		},
	})

	if err != nil {
		statsLock.Lock()
		connectionStats.LastReconnectError = err.Error()
		statsLock.Unlock()

		return fmt.Errorf("failed to reconnect to database: %w", err)
	}

	// 同步更新 DB 引用
	DB = db

	// 获取底层 sql.DB 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		statsLock.Lock()
		connectionStats.LastReconnectError = err.Error()
		statsLock.Unlock()

		return fmt.Errorf("failed to get sql.DB after reconnect: %w", err)
	}

	// 重新配置连接池
	if err := configureConnectionPool(sqlDB); err != nil {
		statsLock.Lock()
		connectionStats.LastReconnectError = err.Error()
		statsLock.Unlock()

		return fmt.Errorf("failed to configure connection pool after reconnect: %w", err)
	}

	// 验证连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := PingWithContext(ctx); err != nil {
		statsLock.Lock()
		connectionStats.LastReconnectError = err.Error()
		statsLock.Unlock()

		return fmt.Errorf("failed to ping database after reconnect: %w", err)
	}

	statsLock.Lock()
	connectionStats.LastReconnectError = ""
	connectionStats.LastHealthStatus = true
	statsLock.Unlock()

	logging.LogInfo("Database reconnected successfully",
		zap.Int64("reconnectCount", connectionStats.ReconnectCount),
	)

	return nil
}

// PingWithContext 带上下文的数据库连接检查
func PingWithContext(ctx context.Context) error {
	sqlDB := GetDBConn()
	if sqlDB == nil {
		return fmt.Errorf("database not initialized")
	}
	return sqlDB.PingContext(ctx)
}

// IsHealthy 检查数据库连接是否健康
func IsHealthy() bool {
	statsLock.RLock()
	defer statsLock.RUnlock()
	if connectionStats == nil {
		return false
	}
	return connectionStats.LastHealthStatus
}

// GetRawConnectionPoolStats 获取原始连接池统计信息（sql.DBStats）
func GetRawConnectionPoolStats() (*sql.DBStats, error) {
	sqlDB := GetDBConn()
	if sqlDB == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	stats := sqlDB.Stats()
	return &stats, nil
}

// CreateOptimizedIndexes 创建优化后的数据库索引
func CreateOptimizedIndexes(db *gorm.DB) error {
	// 这里的实现需要根据实际需求添加
	// 可以使用 db.Exec() 执行 CREATE INDEX 语句
	return nil
}
