package service

import (
	"testing"
	"time"

	"github.com/Kevin-Jii/tower-go/model"
)

func TestBuildB2BSupplyOrderAccount(t *testing.T) {
	orderDate := time.Date(2026, time.August, 8, 0, 0, 0, 0, time.Local)
	order := &model.B2BSupplyOrder{
		OrderNo:       "B2B202608080001",
		StoreID:       3,
		OrderDate:     orderDate,
		TotalAmount:   168.5,
		ProfitAmount:  48.5,
		PaymentStatus: model.B2BPaymentPartial,
		Remark:        "送到后门",
		OperatorID:    17,
		Items: []model.B2BSupplyOrderItem{
			{
				ProductID:   11,
				ProductName: "测试商品",
				UnitName:    "整箱",
				Quantity:    2,
				SupplyPrice: 84.25,
				Amount:      168.5,
				Remark:      "轻放",
			},
		},
	}

	account := buildB2BSupplyOrderAccount(order)
	if account == nil {
		t.Fatal("buildB2BSupplyOrderAccount() returned nil")
	}
	if account.AccountNo != "JZ-"+order.OrderNo || account.OrderNo != order.OrderNo {
		t.Fatalf("account numbers = %q / %q", account.AccountNo, account.OrderNo)
	}
	if account.SourceType != model.StoreAccountSourceB2BSupplyOrder || account.SourceID != 0 {
		t.Fatalf("account source = %q/%d", account.SourceType, account.SourceID)
	}
	if account.PaymentStatus != model.StoreAccountPaymentUnpaid {
		t.Fatalf("partial B2B payment mapped to %d, want unpaid", account.PaymentStatus)
	}
	if account.TotalAmount != order.TotalAmount || account.NetIncomeAmount != order.ProfitAmount {
		t.Fatalf("account amounts = %.2f / %.2f", account.TotalAmount, account.NetIncomeAmount)
	}
	if account.StoreID != order.StoreID || account.OperatorID != order.OperatorID || !account.AccountDate.Equal(orderDate) {
		t.Fatalf("account ownership/date mismatch: %#v", account)
	}
	if account.Channel != "B2B供货" || account.Remark != order.Remark || account.ItemCount != 1 {
		t.Fatalf("account metadata mismatch: %#v", account)
	}
	if len(account.Items) != 1 {
		t.Fatalf("len(account.Items) = %d, want 1", len(account.Items))
	}
	item := account.Items[0]
	if item.ProductID != 11 || item.Spec != "整箱" || item.Unit != "整箱" || item.Quantity != 2 || item.Price != 84.25 || item.Amount != 168.5 {
		t.Fatalf("account item mismatch: %#v", item)
	}
}

func TestBuildB2BSupplyOrderAccountMapsPaidStatus(t *testing.T) {
	account := buildB2BSupplyOrderAccount(&model.B2BSupplyOrder{
		OrderNo:       "B2B202608080002",
		PaymentStatus: model.B2BPaymentPaid,
	})
	if account == nil || account.PaymentStatus != model.StoreAccountPaymentPaid {
		t.Fatalf("paid B2B status mapped to %#v", account)
	}
}

func TestBuildB2BSupplyOrderAccountHandlesNil(t *testing.T) {
	if account := buildB2BSupplyOrderAccount(nil); account != nil {
		t.Fatalf("nil order mapped to %#v", account)
	}
}
