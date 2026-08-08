package service

import (
	"context"
	"reflect"
	"testing"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
)

type fakeDailyTurnoverRepository struct {
	stores    []module.DailyTurnoverStoreRow
	summaries []module.DailyTurnoverSummaryRow
	channels  []module.DailyTurnoverChannelRow
	admins    []module.DailyTurnoverAdminRow
}

func (f *fakeDailyTurnoverRepository) ListActiveStores(context.Context) ([]module.DailyTurnoverStoreRow, error) {
	return f.stores, nil
}

func (f *fakeDailyTurnoverRepository) ListSummaries(context.Context, string) ([]module.DailyTurnoverSummaryRow, error) {
	return f.summaries, nil
}

func (f *fakeDailyTurnoverRepository) ListChannels(context.Context, string) ([]module.DailyTurnoverChannelRow, error) {
	return f.channels, nil
}

func (f *fakeDailyTurnoverRepository) ListAdministrators(context.Context, []uint) ([]module.DailyTurnoverAdminRow, error) {
	return f.admins, nil
}

func TestDailyTurnoverServiceGroupsStoreData(t *testing.T) {
	repository := &fakeDailyTurnoverRepository{
		stores: []module.DailyTurnoverStoreRow{
			{StoreID: 1, StoreName: "一店"},
			{StoreID: 2, StoreName: "二店"},
		},
		summaries: []module.DailyTurnoverSummaryRow{
			{StoreID: 1, TotalAmount: 1264.14, OrderCount: 15},
		},
		channels: []module.DailyTurnoverChannelRow{
			{StoreID: 1, Channel: "wechat_mini", Amount: 424, OrderCount: 3},
			{StoreID: 99, Channel: "offline", Amount: 100, OrderCount: 1},
		},
		admins: []module.DailyTurnoverAdminRow{
			{StoreID: 1, UserID: 12, OpenID: "openid-1"},
			{StoreID: 2, UserID: 23, OpenID: "openid-2"},
		},
	}
	service := NewDailyTurnoverService(repository, nil)

	got, err := service.List(context.Background(), "2026-08-07")
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	want := []model.DailyTurnoverReport{
		{
			StoreID:      1,
			StoreName:    "一店",
			BusinessDate: "2026-08-07",
			TotalAmount:  1264.14,
			OrderCount:   15,
			Channels: []model.DailyTurnoverChannel{
				{Channel: "wechat_mini", ChannelName: "wechat_mini", Amount: 424, OrderCount: 3},
			},
			Admins: []model.DailyTurnoverAdmin{
				{UserID: 12, OpenID: "openid-1"},
			},
		},
		{
			StoreID:      2,
			StoreName:    "二店",
			BusinessDate: "2026-08-07",
			Channels:     []model.DailyTurnoverChannel{},
			Admins: []model.DailyTurnoverAdmin{
				{UserID: 23, OpenID: "openid-2"},
			},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("List() = %#v, want %#v", got, want)
	}
}
