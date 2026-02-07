package database

import (
	"database/sql"
	"fmt"
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

// InitDB 初始化数据库连接
func InitDB(cfg config.DatabaseConfig) error {
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
			SingularTable: true,
			TablePrefix:   "tower_",
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

	// 配置连接池
	sqlDB.SetMaxOpenConns(100)           // 最大打开连接数
	sqlDB.SetMaxIdleConns(10)            // 最大空闲连接数
	sqlDB.SetConnMaxLifetime(time.Hour)  // 连接最大生存时间

	logging.LogInfo("Database connected successfully",
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
		zap.String("database", cfg.Database),
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

// CreateOptimizedIndexes 创建优化后的数据库索引
func CreateOptimizedIndexes(db *gorm.DB) error {
	// 这里的实现需要根据实际需求添加
	// 可以使用 db.Exec() 执行 CREATE INDEX 语句
	return nil
}
