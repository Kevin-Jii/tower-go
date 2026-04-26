package service

import (
	"strings"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
)

// convertToBaseQuantity 将任意单位数量换算到商品基础单位数量。
// 优先使用 product_unit_specs 配置，找不到时兼容旧逻辑（箱->瓶）。
func convertToBaseQuantity(
	unitSpecModule *module.ProductUnitSpecModule,
	product *model.SupplierProduct,
	productID uint,
	quantity float64,
	inputUnit string,
) (float64, string) {
	baseUnit := strings.TrimSpace(inputUnit)
	if product != nil && strings.TrimSpace(product.Unit) != "" {
		baseUnit = strings.TrimSpace(product.Unit)
	}
	if baseUnit == "" {
		baseUnit = "件"
	}
	if quantity <= 0 {
		return quantity, baseUnit
	}

	normalizedInput := strings.TrimSpace(inputUnit)
	if normalizedInput == "" {
		normalizedInput = baseUnit
	}

	if unitSpecModule != nil && productID > 0 {
		if spec, err := unitSpecModule.GetByProductAndUnit(productID, normalizedInput); err == nil && spec != nil && spec.FactorToBase > 0 {
			return quantity * spec.FactorToBase, baseUnit
		}
	}

	// 兼容旧数据：未配置规格表时，沿用“箱->每箱瓶数”的换算方式
	if product != nil && strings.Contains(normalizedInput, "箱") && product.BottlesPerCase > 0 {
		return quantity * float64(product.BottlesPerCase), baseUnit
	}

	return quantity, baseUnit
}
