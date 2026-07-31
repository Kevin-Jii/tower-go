package service

import (
	"testing"
	"time"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/utils/businessdate"
)

func TestIsTakeoutChannelValue(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{name: "known mall value", value: "mall", want: true},
		{name: "known group buy value", value: "group_buy", want: true},
		{name: "known Chinese label", value: "美团外卖", want: true},
		{name: "embedded known platform", value: "store-meituan", want: true},
		{name: "small is not mall", value: "small", want: false},
		{name: "small store is not mall", value: "small_store", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isTakeoutChannelValue(tt.value); got != tt.want {
				t.Fatalf("isTakeoutChannelValue(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestStoreAccountEditWindow_CurrentBusinessDayOnly(t *testing.T) {
	svc := &StoreAccountService{}
	now := time.Now()
	previousBusinessDay := businessdate.Date(now).AddDate(0, 0, -1).Add(16 * time.Hour)

	tests := []struct {
		name    string
		account *model.StoreAccount
		want    bool
	}{
		{
			name:    "current business day is editable",
			account: &model.StoreAccount{CreatedAt: now},
			want:    true,
		},
		{
			name:    "previous business day is not editable",
			account: &model.StoreAccount{CreatedAt: previousBusinessDay},
			want:    false,
		},
		{
			name:    "zero created time is not editable",
			account: &model.StoreAccount{},
			want:    false,
		},
		{
			name:    "nil account is not editable",
			account: nil,
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := svc.CanUpdateAccount(tt.account, &model.UpdateStoreAccountReq{}); got != tt.want {
				t.Fatalf("CanUpdateAccount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStoreAccountCanBindConsumables_UnboundOnly(t *testing.T) {
	svc := &StoreAccountService{}
	now := time.Now()

	tests := []struct {
		name    string
		account *model.StoreAccount
		want    bool
	}{
		{
			name:    "unbound account can bind consumables",
			account: &model.StoreAccount{CreatedAt: now},
			want:    true,
		},
		{
			name: "bound account cannot bind consumables again",
			account: &model.StoreAccount{
				CreatedAt:   now,
				Consumables: []model.StoreAccountConsumable{{ID: 1, AccountID: 1, ProductName: "纸杯", Quantity: 1}},
			},
			want: false,
		},
		{
			name:    "canceled account cannot bind consumables",
			account: &model.StoreAccount{CreatedAt: now, IsCanceled: true},
			want:    false,
		},
		{
			name:    "zero created time cannot bind consumables",
			account: &model.StoreAccount{},
			want:    false,
		},
		{
			name:    "nil account cannot bind consumables",
			account: nil,
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := svc.CanBindConsumables(tt.account); got != tt.want {
				t.Fatalf("CanBindConsumables() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStoreAccountCanCancel_UncanceledAccount(t *testing.T) {
	svc := &StoreAccountService{}
	now := time.Now()
	previousBusinessDay := businessdate.Date(now).AddDate(0, 0, -1).Add(16 * time.Hour)

	tests := []struct {
		name    string
		account *model.StoreAccount
		want    bool
	}{
		{
			name:    "current business day account can be canceled",
			account: &model.StoreAccount{CreatedAt: now},
			want:    true,
		},
		{
			name:    "previous business day account can be canceled",
			account: &model.StoreAccount{CreatedAt: previousBusinessDay},
			want:    true,
		},
		{
			name:    "account without created time can be canceled",
			account: &model.StoreAccount{},
			want:    true,
		},
		{
			name:    "already canceled account cannot be canceled again",
			account: &model.StoreAccount{CreatedAt: now, IsCanceled: true},
			want:    false,
		},
		{
			name:    "nil account cannot be canceled",
			account: nil,
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := svc.CanCancelAccount(tt.account); got != tt.want {
				t.Fatalf("CanCancelAccount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildCancelRestoreOrderUsesStableAccountOrderNo(t *testing.T) {
	svc := &StoreAccountService{}
	account := &model.StoreAccount{
		ID:        466,
		AccountNo: "JZ202607290001",
		StoreID:   2,
		Items: []model.StoreAccountItem{
			{
				ProductID:   123,
				ProductName: "测试商品",
				Quantity:    2,
				Unit:        "瓶",
			},
			{
				ProductID:   model.StoreAccountItemCustomProductID,
				ProductName: "自定义商品",
				Quantity:    1,
				Unit:        "项",
			},
		},
	}

	first := svc.buildCancelRestoreOrder(account, 1003)
	second := svc.buildCancelRestoreOrder(account, 1003)
	if first == nil || second == nil {
		t.Fatal("buildCancelRestoreOrder() returned nil")
	}
	if first.OrderNo != "RKZF0000000466" {
		t.Fatalf("OrderNo = %q, want %q", first.OrderNo, "RKZF0000000466")
	}
	if second.OrderNo != first.OrderNo {
		t.Fatalf("OrderNo is not stable: first=%q second=%q", first.OrderNo, second.OrderNo)
	}
	if other := storeAccountCancelInventoryOrderNo(467); other == first.OrderNo {
		t.Fatalf("different accounts share OrderNo %q", other)
	}
	if len(first.Items) != 1 || first.Items[0].ProductID != 123 {
		t.Fatalf("restore items = %#v, want only system product 123", first.Items)
	}
}

func TestBuildStoreAccountItemsRecalculatesCustomItemAmount(t *testing.T) {
	svc := &StoreAccountService{}
	items, totalAmount, itemCostAmount, _, err := svc.buildStoreAccountItems([]model.CreateStoreAccountItemReq{
		{
			ProductID:   model.StoreAccountItemCustomProductID,
			ProductName: "临时服务费",
			Quantity:    3,
			Unit:        "次",
			Price:       12.5,
			Amount:      1,
		},
	})
	if err != nil {
		t.Fatalf("buildStoreAccountItems() error = %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("len(items) = %d, want 1", len(items))
	}
	if items[0].Amount != 37.5 {
		t.Fatalf("custom item amount = %.2f, want recalculated amount 37.50", items[0].Amount)
	}
	if totalAmount != 37.5 {
		t.Fatalf("totalAmount = %.2f, want 37.50", totalAmount)
	}
	if itemCostAmount != 0 {
		t.Fatalf("itemCostAmount = %.2f, want 0", itemCostAmount)
	}
}

func TestPaymentStatusOnlyUpdateRejectsItemReplacement(t *testing.T) {
	paid := model.StoreAccountPaymentPaid
	svc := &StoreAccountService{}
	account := &model.StoreAccount{PaymentStatus: model.StoreAccountPaymentUnpaid}
	req := &model.UpdateStoreAccountReq{
		PaymentStatus: &paid,
		Items: []model.CreateStoreAccountItemReq{
			{ProductID: 1, Quantity: 1, Unit: "瓶"},
		},
	}
	if svc.canApplyPaymentStatusOnlyUpdate(account, req) {
		t.Fatal("item replacement must not be treated as a payment-status-only update")
	}
}

func TestBuildAccountItemAdjustmentOrdersUsesStockDelta(t *testing.T) {
	svc := &StoreAccountService{}
	account := &model.StoreAccount{
		ID:        88,
		AccountNo: "JZ202607300088",
		StoreID:   3,
		Items: []model.StoreAccountItem{
			{ProductID: 11, ProductName: "商品A", Quantity: 2, Unit: "瓶"},
			{ProductID: 22, ProductName: "商品B", Quantity: 4, Unit: "瓶"},
			{ProductID: 0, ProductName: "自定义", Quantity: 9, Unit: "项"},
		},
	}
	newItems := []model.StoreAccountItem{
		{ProductID: 11, ProductName: "商品A", Quantity: 5, Unit: "瓶"},
		{ProductID: 22, ProductName: "商品B", Quantity: 1, Unit: "瓶"},
		{ProductID: 0, ProductName: "自定义", Quantity: 1, Unit: "项"},
	}

	inOrder, outOrder, err := svc.buildAccountItemAdjustmentOrders(account, newItems, 1003)
	if err != nil {
		t.Fatalf("buildAccountItemAdjustmentOrders() error = %v", err)
	}
	if inOrder == nil || len(inOrder.Items) != 1 {
		t.Fatalf("inOrder = %#v, want one stock return item", inOrder)
	}
	if inOrder.Items[0].ProductID != 22 || inOrder.Items[0].Quantity != 3 {
		t.Fatalf("inOrder item = %#v, want product 22 quantity 3", inOrder.Items[0])
	}
	if outOrder == nil || len(outOrder.Items) != 1 {
		t.Fatalf("outOrder = %#v, want one stock deduction item", outOrder)
	}
	if outOrder.Items[0].ProductID != 11 || outOrder.Items[0].Quantity != 3 {
		t.Fatalf("outOrder item = %#v, want product 11 quantity 3", outOrder.Items[0])
	}
}
