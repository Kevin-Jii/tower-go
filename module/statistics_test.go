package module

import (
	"testing"

	"github.com/Kevin-Jii/tower-go/model"
)

func TestCalculateBusinessOverviewNetProfitUsesStoreAccountScope(t *testing.T) {
	tests := []struct {
		name           string
		stats          *model.BusinessOverviewStats
		itemCostAmount float64
		want           float64
	}{
		{
			name: "canceled-only period does not turn negative from standalone store expense",
			stats: &model.BusinessOverviewStats{
				StoreExpenseAmount: 120,
			},
			want: 0,
		},
		{
			name: "active account deductions follow store account profit scope",
			stats: &model.BusinessOverviewStats{
				SalesAmount:        200,
				OtherExpenseAmount: 10,
				ErrandFeeAmount:    5,
				ConsumableAmount:   12,
				GiftWineCostAmount: 8,
				RoundAmount:        2,
				StoreExpenseAmount: 100,
			},
			itemCostAmount: 40,
			want:           123,
		},
		{
			name: "nil stats",
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateBusinessOverviewNetProfit(tt.stats, tt.itemCostAmount); got != tt.want {
				t.Fatalf("calculateBusinessOverviewNetProfit() = %.2f, want %.2f", got, tt.want)
			}
		})
	}
}
