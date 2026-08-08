package service

import (
	"context"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
)

type dailyTurnoverRepository interface {
	ListActiveStores(ctx context.Context) ([]module.DailyTurnoverStoreRow, error)
	ListSummaries(ctx context.Context, businessDate string) ([]module.DailyTurnoverSummaryRow, error)
	ListChannels(ctx context.Context, businessDate string) ([]module.DailyTurnoverChannelRow, error)
	ListAdministrators(ctx context.Context, storeIDs []uint) ([]module.DailyTurnoverAdminRow, error)
}

type DailyTurnoverService struct {
	repository dailyTurnoverRepository
	dictModule *module.DictModule
}

func NewDailyTurnoverService(repository dailyTurnoverRepository, dictModule *module.DictModule) *DailyTurnoverService {
	return &DailyTurnoverService{
		repository: repository,
		dictModule: dictModule,
	}
}

func (s *DailyTurnoverService) List(ctx context.Context, businessDate string) ([]model.DailyTurnoverReport, error) {
	stores, err := s.repository.ListActiveStores(ctx)
	if err != nil {
		return nil, err
	}

	reports := make([]model.DailyTurnoverReport, len(stores))
	storeIndexes := make(map[uint]int, len(stores))
	storeIDs := make([]uint, len(stores))
	for index, store := range stores {
		reports[index] = model.DailyTurnoverReport{
			StoreID:      store.StoreID,
			StoreName:    store.StoreName,
			BusinessDate: businessDate,
			Channels:     make([]model.DailyTurnoverChannel, 0),
			Admins:       make([]model.DailyTurnoverAdmin, 0),
		}
		storeIndexes[store.StoreID] = index
		storeIDs[index] = store.StoreID
	}

	summaries, err := s.repository.ListSummaries(ctx, businessDate)
	if err != nil {
		return nil, err
	}
	for _, summary := range summaries {
		if index, ok := storeIndexes[summary.StoreID]; ok {
			reports[index].TotalAmount = summary.TotalAmount
			reports[index].OrderCount = summary.OrderCount
		}
	}

	channelNames, err := s.salesChannelNames()
	if err != nil {
		return nil, err
	}
	channels, err := s.repository.ListChannels(ctx, businessDate)
	if err != nil {
		return nil, err
	}
	for _, channel := range channels {
		index, ok := storeIndexes[channel.StoreID]
		if !ok {
			continue
		}
		channelName := channelNames[channel.Channel]
		if channelName == "" {
			channelName = channel.Channel
		}
		reports[index].Channels = append(reports[index].Channels, model.DailyTurnoverChannel{
			Channel:     channel.Channel,
			ChannelName: channelName,
			Amount:      channel.Amount,
			OrderCount:  channel.OrderCount,
		})
	}

	admins, err := s.repository.ListAdministrators(ctx, storeIDs)
	if err != nil {
		return nil, err
	}
	for _, admin := range admins {
		if index, ok := storeIndexes[admin.StoreID]; ok {
			reports[index].Admins = append(reports[index].Admins, model.DailyTurnoverAdmin{
				UserID: admin.UserID,
				OpenID: admin.OpenID,
			})
		}
	}

	return reports, nil
}

func (s *DailyTurnoverService) salesChannelNames() (map[string]string, error) {
	channelNames := make(map[string]string)
	if s.dictModule == nil {
		return channelNames, nil
	}
	items, err := s.dictModule.ListDataByTypeCode("sales_channel")
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		if item != nil {
			channelNames[item.Value] = item.Label
		}
	}
	return channelNames, nil
}
