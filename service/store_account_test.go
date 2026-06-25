package service

import (
	"testing"
	"time"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/utils/businessdate"
)

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
			if got := svc.CanBindConsumables(tt.account); got != tt.want {
				t.Fatalf("CanBindConsumables() = %v, want %v", got, tt.want)
			}
		})
	}
}
