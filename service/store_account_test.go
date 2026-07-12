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

	tests := []struct {
		name    string
		account *model.StoreAccount
		want    bool
	}{
		{
			name:    "unbound account can bind consumables",
			account: &model.StoreAccount{},
			want:    true,
		},
		{
			name: "bound account cannot bind consumables again",
			account: &model.StoreAccount{
				Consumables: []model.StoreAccountConsumable{{ID: 1, AccountID: 1, ProductName: "纸杯", Quantity: 1}},
			},
			want: false,
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
