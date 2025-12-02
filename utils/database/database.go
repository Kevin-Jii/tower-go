package database

import (
	"fmt"
	"log"
	"time"
	"github.com/Kevin-Jii/tower-go/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(cfg config.DatabaseConfig) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
	if err != nil {
		return err
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxOpenConns(100)                 // 最大连接数
	sqlDB.SetMaxIdleConns(20)                  // 最大空闲连接数
	sqlDB.SetConnMaxLifetime(time.Hour)        // 连接最大生命周期
	sqlDB.SetConnMaxIdleTime(10 * time.Minute) // 空闲连接最大生命周期

	DB = db
	return nil
}

// indexExists 检查索引是否存在
func indexExists(db *gorm.DB, tableName, indexName string) bool {
	var count int64
	query := `SELECT COUNT(*) FROM information_schema.statistics 
	          WHERE table_schema = DATABASE() 
	          AND table_name = ? 
	          AND index_name = ?`
	db.Raw(query, tableName, indexName).Scan(&count)
	return count > 0
}

// getAllExistingIndexes 一次性获取所有已存在的索引（性能优化）
func getAllExistingIndexes(db *gorm.DB) map[string]map[string]bool {
	type IndexInfo struct {
		TableName string
		IndexName string
	}
	
	var indexes []IndexInfo
	query := `SELECT table_name as TableName, index_name as IndexName 
	          FROM information_schema.statistics 
	          WHERE table_schema = DATABASE()`
	db.Raw(query).Scan(&indexes)
	
	// 构建 map[tableName]map[indexName]bool 结构
	result := make(map[string]map[string]bool)
	for _, idx := range indexes {
		if result[idx.TableName] == nil {
			result[idx.TableName] = make(map[string]bool)
		}
		result[idx.TableName][idx.IndexName] = true
	}
	
	return result
}

// CreateOptimizedIndexes 创建优化索引（兼容 MySQL 5.x）
func CreateOptimizedIndexes(db *gorm.DB) error {
	log.Println("开始创建/检查优化索引...")

	// 一次性获取所有已存在的索引（性能优化：避免多次查询）
	existingIndexes := getAllExistingIndexes(db)

	indexes := []struct{
		table     string
		indexName string
		sql       string
		desc      string
	}{
		// 用户表索引
		{
			table:     "users",
			indexName: "idx_users_store_status",
			sql:       "CREATE INDEX idx_users_store_status ON users(store_id, status)",
			desc:      "用户-门店-状态复合索引",
		},
		{
			table:     "users",
			indexName: "idx_users_role_id",
			sql:       "CREATE INDEX idx_users_role_id ON users(role_id)",
			desc:      "用户-角色索引",
		},
		{
			table:     "users",
			indexName: "idx_users_phone_prefix",
			sql:       "CREATE INDEX idx_users_phone_prefix ON users(phone)",
			desc:      "手机号索引",
		},
		{
			table:     "users",
			indexName: "idx_users_username",
			sql:       "CREATE INDEX idx_users_username ON users(username)",
			desc:      "用户名索引",
		},

		// 菜单表索引
		{
			table:     "menus",
			indexName: "idx_menus_parent_sort_status",
			sql:       "CREATE INDEX idx_menus_parent_sort_status ON menus(parent_id, sort, status)",
			desc:      "菜单-父级-排序-状态复合索引",
		},
		{
			table:     "menus",
			indexName: "idx_menus_status",
			sql:       "CREATE INDEX idx_menus_status ON menus(status)",
			desc:      "菜单-状态索引",
		},

		// 角色菜单关联表索引
		{
			table:     "role_menus",
			indexName: "idx_role_menus_role",
			sql:       "CREATE INDEX idx_role_menus_role ON role_menus(role_id)",
			desc:      "角色菜单-角色索引",
		},
		{
			table:     "role_menus",
			indexName: "idx_role_menus_menu",
			sql:       "CREATE INDEX idx_role_menus_menu ON role_menus(menu_id)",
			desc:      "角色菜单-菜单索引",
		},

		// 门店角色菜单关联表索引
		{
			table:     "store_role_menus",
			indexName: "idx_store_role_menus_store_role",
			sql:       "CREATE INDEX idx_store_role_menus_store_role ON store_role_menus(store_id, role_id)",
			desc:      "门店角色菜单-门店角色复合索引",
		},
		{
			table:     "store_role_menus",
			indexName: "idx_store_role_menus_menu",
			sql:       "CREATE INDEX idx_store_role_menus_menu ON store_role_menus(menu_id)",
			desc:      "门店角色菜单-菜单索引",
		},

		// 门店表索引
		{
			table:     "stores",
			indexName: "idx_stores_status",
			sql:       "CREATE INDEX idx_stores_status ON stores(status)",
			desc:      "门店-状态索引",
		},

		// 供应商表索引
		{
			table:     "suppliers",
			indexName: "idx_suppliers_status",
			sql:       "CREATE INDEX idx_suppliers_status ON suppliers(status)",
			desc:      "供应商-状态索引",
		},

		// 供应商商品表索引
		{
			table:     "supplier_products",
			indexName: "idx_supplier_products_supplier",
			sql:       "CREATE INDEX idx_supplier_products_supplier ON supplier_products(supplier_id)",
			desc:      "供应商商品-供应商索引",
		},
		{
			table:     "supplier_products",
			indexName: "idx_supplier_products_category",
			sql:       "CREATE INDEX idx_supplier_products_category ON supplier_products(category_id)",
			desc:      "供应商商品-分类索引",
		},

		// 采购订单表索引
		{
			table:     "purchase_orders",
			indexName: "idx_purchase_orders_store_status",
			sql:       "CREATE INDEX idx_purchase_orders_store_status ON purchase_orders(store_id, status)",
			desc:      "采购订单-门店-状态复合索引",
		},
		{
			table:     "purchase_orders",
			indexName: "idx_purchase_orders_order_date",
			sql:       "CREATE INDEX idx_purchase_orders_order_date ON purchase_orders(order_date)",
			desc:      "采购订单-日期索引",
		},

		// 采购订单明细表索引
		{
			table:     "purchase_order_items",
			indexName: "idx_purchase_order_items_order",
			sql:       "CREATE INDEX idx_purchase_order_items_order ON purchase_order_items(order_id)",
			desc:      "采购订单明细-订单索引",
		},
		{
			table:     "purchase_order_items",
			indexName: "idx_purchase_order_items_supplier",
			sql:       "CREATE INDEX idx_purchase_order_items_supplier ON purchase_order_items(supplier_id)",
			desc:      "采购订单明细-供应商索引",
		},
	}

	createdCount := 0
	skippedCount := 0
	
	for _, idx := range indexes {
		// 使用预加载的索引信息检查（避免多次查询数据库）
		if existingIndexes[idx.table] != nil && existingIndexes[idx.table][idx.indexName] {
			skippedCount++
			continue
		}

		// 索引不存在，创建它
		if err := db.Exec(idx.sql).Error; err != nil {
			log.Printf("⚠️  索引创建失败 [%s]: %v", idx.desc, err)
			// 继续创建其他索引，不中断
		} else {
			log.Printf("✅ 索引创建成功 [%s]", idx.desc)
			createdCount++
		}
	}

	log.Printf("索引创建/检查完成 (新建: %d, 跳过: %d)", createdCount, skippedCount)
	return nil
}
