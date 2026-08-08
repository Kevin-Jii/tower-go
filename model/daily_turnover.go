package model

// DailyTurnoverReport is the store-level payload consumed by the WeChat cloud service.
type DailyTurnoverReport struct {
	StoreID      uint                   `json:"store_id"`
	StoreName    string                 `json:"store_name"`
	BusinessDate string                 `json:"business_date"`
	TotalAmount  float64                `json:"total_amount"`
	OrderCount   int64                  `json:"order_count"`
	Channels     []DailyTurnoverChannel `json:"channels"`
	Admins       []DailyTurnoverAdmin   `json:"admins"`
}

type DailyTurnoverChannel struct {
	Channel     string  `json:"channel"`
	ChannelName string  `json:"channel_name"`
	Amount      float64 `json:"amount"`
	OrderCount  int64   `json:"order_count"`
}

type DailyTurnoverAdmin struct {
	UserID uint   `json:"user_id"`
	OpenID string `json:"openid"`
}
