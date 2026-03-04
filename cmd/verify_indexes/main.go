package main

import (
	"fmt"
	"log"

	"github.com/Kevin-Jii/tower-go/bootstrap"
	"github.com/Kevin-Jii/tower-go/config"
	"github.com/Kevin-Jii/tower-go/utils/database"
)

type IndexInfo struct {
	TableName  string `gorm:"column:TABLE_NAME"`
	IndexName  string `gorm:"column:INDEX_NAME"`
	ColumnName string `gorm:"column:COLUMN_NAME"`
	NonUnique  int    `gorm:"column:NON_UNIQUE"`
	SeqInIndex int    `gorm:"column:SEQ_IN_INDEX"`
}

func main() {
	// 初始化配置
	config.InitConfig()

	// 初始化数据库连接
	bootstrap.InitDatabase()

	// 获取数据库实例
	db := database.GetDB()
	if db == nil {
		log.Fatal("数据库连接失败")
	}

	fmt.Println("==============================================")
	fmt.Println("数据库索引验证工具")
	fmt.Println("==============================================")
	fmt.Println()

	// 需要验证的表
	tables := []string{
		"store_accounts",
		"store_account_items",
		"inventories",
		"inventory_orders",
		"inventory_order_items",
		"supplier_products",
		"users",
		"stores",
	}

	// 期望的索引
	expectedIndexes := map[string][]string{
		"store_accounts": {
			"idx_store_account_date_range",
			"idx_store_account_store_channel",
			"idx_store_account_date_channel",
			"idx_store_account_all",
		},
		"store_account_items": {
			"idx_account_items_product_time",
		},
		"inventories": {
			"idx_inventory_unique",
			"idx_inventory_store_qty",
		},
		"inventory_orders": {
			"idx_inv_order_store_type_date",
			"idx_inv_order_type_date",
		},
		"inventory_order_items": {
			"idx_order_items_product_qty",
		},
		"supplier_products": {
			"idx_supplier_prod_name",
			"idx_supplier_prod_category",
		},
		"users": {
			"idx_users_store_name",
		},
		"stores": {
			"idx_stores_name",
		},
	}

	totalExpected := 0
	totalFound := 0

	for _, tableName := range tables {
		fmt.Printf("📊 表: %s\n", tableName)

		// 查询表的所有索引
		var indexes []IndexInfo
		query := `
			SELECT 
				TABLE_NAME,
				INDEX_NAME,
				COLUMN_NAME,
				NON_UNIQUE,
				SEQ_IN_INDEX
			FROM information_schema.statistics
			WHERE table_schema = DATABASE()
			AND table_name = ?
			ORDER BY INDEX_NAME, SEQ_IN_INDEX
		`
		db.Raw(query, tableName).Scan(&indexes)

		// 按索引名分组
		indexMap := make(map[string][]string)
		for _, idx := range indexes {
			indexMap[idx.IndexName] = append(indexMap[idx.IndexName], idx.ColumnName)
		}

		// 检查期望的索引
		expected := expectedIndexes[tableName]
		for _, expectedIdx := range expected {
			totalExpected++
			if columns, exists := indexMap[expectedIdx]; exists {
				fmt.Printf("  ✅ %s (%v)\n", expectedIdx, columns)
				totalFound++
			} else {
				fmt.Printf("  ❌ %s (缺失)\n", expectedIdx)
			}
		}

		fmt.Println()
	}

	// 输出统计
	fmt.Println("==============================================")
	fmt.Printf("索引验证结果: %d/%d\n", totalFound, totalExpected)
	if totalFound == totalExpected {
		fmt.Println("✅ 所有索引创建成功！")
	} else {
		fmt.Printf("⚠️  缺失 %d 个索引\n", totalExpected-totalFound)
	}
	fmt.Println("==============================================")
}
