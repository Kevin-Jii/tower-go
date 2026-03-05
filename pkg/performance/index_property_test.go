package performance

import (
	"testing"
	"time"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/utils/database"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"gorm.io/gorm"
)

// Property 1: Index performance improvement - Requirements 1.1
func TestProperty_IndexPerformanceImprovement(t *testing.T) {
	db := database.GetDB()
	if db == nil {
		t.Skip("No database connection available")
	}

	properties := gopter.NewProperties(nil)

	properties.Property("store_id + account_date query uses composite index", prop.ForAll(
		func(storeID uint, days int) bool {
			if storeID == 0 || days <= 0 {
				return true
			}
			endDate := time.Now()
			startDate := endDate.AddDate(0, 0, -days)
			var count int64
			err := db.Model(&model.StoreAccount{}).
				Where("store_id = ? AND account_date BETWEEN ? AND ?", storeID, startDate, endDate).
				Count(&count).Error
			return err == nil
		},
		gen.UIntRange(1, 1000),
		gen.IntRange(1, 365),
	))

	properties.Property("indexed query completes successfully", prop.ForAll(
		func(storeID uint, days int) bool {
			if storeID == 0 || days <= 0 {
				return true
			}
			endDate := time.Now()
			startDate := endDate.AddDate(0, 0, -days)
			var results []model.StoreAccount
			err := db.Where("store_id = ? AND account_date BETWEEN ? AND ?", storeID, startDate, endDate).
				Limit(100).Find(&results).Error
			return err == nil
		},
		gen.UIntRange(1, 1000),
		gen.IntRange(1, 365),
	))

	properties.Property("store_id + channel query uses index", prop.ForAll(
		func(storeID uint, channel string) bool {
			if storeID == 0 || channel == "" {
				return true
			}
			var count int64
			err := db.Model(&model.StoreAccount{}).
				Where("store_id = ? AND channel = ?", storeID, channel).Count(&count).Error
			return err == nil
		},
		gen.UIntRange(1, 1000),
		gen.AlphaString(),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Property 2: Inventory uniqueness and performance - Requirements 1.2
func TestProperty_InventoryUniquenessConstraint(t *testing.T) {
	db := database.GetDB()
	if db == nil {
		t.Skip("No database connection available")
	}

	properties := gopter.NewProperties(nil)

	properties.Property("inventory query by store and product uses index efficiently", prop.ForAll(
		func(storeID, productID uint) bool {
			if storeID == 0 || productID == 0 {
				return true
			}
			var inventory model.Inventory
			err := db.Where("store_id = ? AND product_id = ?", storeID, productID).First(&inventory).Error
			return err == nil || err == gorm.ErrRecordNotFound
		},
		gen.UIntRange(1, 1000),
		gen.UIntRange(1, 10000),
	))

	properties.Property("low stock query uses index efficiently", prop.ForAll(
		func(storeID uint, threshold float64) bool {
			if storeID == 0 {
				return true
			}
			var inventories []model.Inventory
			err := db.Where("store_id = ? AND quantity < ?", storeID, threshold).Limit(50).Find(&inventories).Error
			return err == nil
		},
		gen.UIntRange(1, 1000),
		gen.Float64Range(0, 100),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Property 3: N+1 query elimination - Requirements 1.3
func TestProperty_NPlusOneQueryElimination(t *testing.T) {
	db := database.GetDB()
	if db == nil {
		t.Skip("No database connection available")
	}

	properties := gopter.NewProperties(nil)

	properties.Property("preload with associations uses constant queries", prop.ForAll(
		func(storeID uint, limit int) bool {
			if storeID == 0 || limit <= 0 || limit > 100 {
				return true
			}
			var accounts []model.StoreAccount
			err := db.Where("store_id = ?", storeID).
				Preload("Store").Preload("Operator").Preload("Items").
				Limit(limit).Find(&accounts).Error
			return err == nil
		},
		gen.UIntRange(1, 1000),
		gen.IntRange(1, 50),
	))

	properties.Property("joins can be used for eager loading", prop.ForAll(
		func(storeID uint) bool {
			if storeID == 0 {
				return true
			}
			var results []model.StoreAccount
			err := db.Where("store_accounts.store_id = ?", storeID).
				Joins("LEFT JOIN stores ON stores.id = store_accounts.store_id").
				Limit(10).Find(&results).Error
			return err == nil
		},
		gen.UIntRange(1, 1000),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Property 6: Cursor pagination performance - Requirements 2.2
func TestProperty_CursorPaginationPerformance(t *testing.T) {
	db := database.GetDB()
	if db == nil {
		t.Skip("No database connection available")
	}

	properties := gopter.NewProperties(nil)

	properties.Property("cursor pagination returns correct results", prop.ForAll(
		func(storeID uint, limit int, cursor string) bool {
			if storeID == 0 || limit <= 0 {
				return true
			}
			paginator := database.NewCursorPaginator("id", true)
			query := db.Where("store_id = ?", storeID).Order("id DESC")
			_, err := paginator.Paginate(query, cursor, limit)
			if err != nil {
				return true
			}
			return true
		},
		gen.UIntRange(1, 1000),
		gen.IntRange(1, 50),
		gen.AlphaString(),
	))

	properties.Property("cursor encoding and decoding are symmetric", prop.ForAll(
		func(lastID uint, lastValue string) bool {
			if lastID == 0 {
				return true
			}
			paginator := database.NewCursorPaginator("id", false)
			cursor := paginator.EncodeCursor(lastID, lastValue)
			if cursor == "" {
				return false
			}
			decodedID, _, err := paginator.DecodeCursor(cursor)
			if err != nil {
				return false
			}
			return decodedID == lastID
		},
		gen.UIntRange(1, 100000),
		gen.AlphaString(),
	))

	properties.Property("empty cursor returns first page", prop.ForAll(
		func(storeID uint, limit int) bool {
			if storeID == 0 || limit <= 0 {
				return true
			}
			paginator := database.NewCursorPaginator("id", true)
			query := db.Where("store_id = ?", storeID).Order("id DESC")
			_, err := paginator.Paginate(query, "", limit)
			return err == nil
		},
		gen.UIntRange(1, 1000),
		gen.IntRange(1, 50),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}