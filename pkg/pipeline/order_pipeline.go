package pipeline

import (
	"github.com/Kevin-Jii/tower-go/model"
)

// OrderItemProcessor 订单明细处理管道
type OrderItemProcessor struct {
	pipeline *Pipeline[[]model.PurchaseOrderItem]
}

// NewOrderItemProcessor 创建订单明细处理器
func NewOrderItemProcessor() *OrderItemProcessor {
	return &OrderItemProcessor{
		pipeline: New[[]model.PurchaseOrderItem](),
	}
}

// AddValidation 添加验证阶段
func (p *OrderItemProcessor) AddValidation(validator func([]model.PurchaseOrderItem) error) *OrderItemProcessor {
	p.pipeline.Add(func(items []model.PurchaseOrderItem) ([]model.PurchaseOrderItem, error) {
		if err := validator(items); err != nil {
			return nil, err
		}
		return items, nil
	})
	return p
}

// AddTransform 添加转换阶段
func (p *OrderItemProcessor) AddTransform(transformer func([]model.PurchaseOrderItem) []model.PurchaseOrderItem) *OrderItemProcessor {
	p.pipeline.Add(func(items []model.PurchaseOrderItem) ([]model.PurchaseOrderItem, error) {
		return transformer(items), nil
	})
	return p
}

// Process 执行处理
func (p *OrderItemProcessor) Process(items []model.PurchaseOrderItem) ([]model.PurchaseOrderItem, error) {
	return p.pipeline.Execute(items)
}

// FilterBySupplier 按供应商过滤
func FilterBySupplier(supplierID uint) func([]model.PurchaseOrderItem) []model.PurchaseOrderItem {
	return func(items []model.PurchaseOrderItem) []model.PurchaseOrderItem {
		return Filter(items, func(item model.PurchaseOrderItem) bool {
			return item.SupplierID == supplierID
		})
	}
}

// CalculateAmounts 计算金额
func CalculateAmounts(items []model.PurchaseOrderItem) []model.PurchaseOrderItem {
	for i := range items {
		items[i].Amount = items[i].Quantity * items[i].UnitPrice
	}
	return items
}

// SumTotal 计算总金额
func SumTotal(items []model.PurchaseOrderItem) float64 {
	return Reduce(items, 0.0, func(total float64, item model.PurchaseOrderItem) float64 {
		return total + item.Amount
	})
}

// GroupBySupplier 按供应商分组
func GroupBySupplier(items []model.PurchaseOrderItem) map[uint][]model.PurchaseOrderItem {
	result := make(map[uint][]model.PurchaseOrderItem)
	for _, item := range items {
		result[item.SupplierID] = append(result[item.SupplierID], item)
	}
	return result
}
