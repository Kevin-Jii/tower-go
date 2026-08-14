package module

import (
	"testing"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestUpsertByProductAndUnitPreservesFalseIsSaleable(t *testing.T) {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       "gorm:gorm@tcp(localhost:9910)/gorm?charset=utf8&parseTime=True&loc=Local",
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		DryRun:                 true,
		DisableAutomaticPing:   true,
		SkipDefaultTransaction: true,
	})
	require.NoError(t, err)

	var capturedSQL string
	var capturedVars []interface{}
	require.NoError(t, db.Callback().Create().After("gorm:create").Register("test:capture_product_unit_spec_upsert", func(tx *gorm.DB) {
		capturedSQL = tx.Statement.SQL.String()
		capturedVars = append([]interface{}(nil), tx.Statement.Vars...)
	}))

	unitSpecModule := NewProductUnitSpecModule(db)
	err = unitSpecModule.UpsertByProductAndUnit(&model.ProductUnitSpec{
		ProductID:    1,
		UnitCode:     "liter",
		UnitName:     "1L",
		FactorToBase: 1,
		CostPrice:    8.2,
		SalePrice:    20,
		IsSaleable:   false,
		IsEnabled:    true,
	})
	require.NoError(t, err)
	require.Contains(t, capturedSQL, "`is_saleable`")
	require.Contains(t, capturedVars, false)
}
