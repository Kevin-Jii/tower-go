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
