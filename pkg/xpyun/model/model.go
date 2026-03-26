package model

import (
	"strconv"
	"github.com/Kevin-Jii/tower-go/pkg/xpyun/util"
)

// RestRequest 基础请求结构
type RestRequest struct {
	User      string `json:"user"`
	UserKey   string `json:"-"`
	Timestamp int64  `json:"timestamp"`
	Sign      string `json:"sign"`
	Debug     int    `json:"debug,omitempty"`
}

func (r *RestRequest) GenerateSign() {
	r.Sign = util.Sign(r.User + r.UserKey + strconv.FormatInt(r.Timestamp, 10))
}

// AddPrinterRequestItem 添加打印机项
type AddPrinterRequestItem struct {
	Sn   string `json:"sn"`
	Name string `json:"name,omitempty"`
}

// AddPrinterRequest 添加打印机请求
type AddPrinterRequest struct {
	RestRequest `json:",inline"`
	Items       []*AddPrinterRequestItem `json:"items"`
}

// DelPrinterRequest 删除打印机请求
type DelPrinterRequest struct {
	RestRequest `json:",inline"`
	SnList      []string `json:"snlist"`
}

// SetVoiceTypeRequest 设置语音类型请求
type SetVoiceTypeRequest struct {
	RestRequest `json:",inline"`
	Sn          string `json:"sn"`
	VoiceType   int    `json:"voiceType"`
	VolumeLevel int    `json:"volumeLevel,omitempty"`
}

// UpdPrinterRequest 更新打印机请求
type UpdPrinterRequest struct {
	RestRequest `json:",inline"`
	Sn          string `json:"sn"`
	Name        string `json:"name,omitempty"`
}

// PrinterRequest 通用打印机请求
type PrinterRequest struct {
	RestRequest `json:",inline"`
	Sn          string `json:"sn"`
}

// ClearPrintOrderRequest 清空打印队列请求
type ClearPrintOrderRequest struct {
	RestRequest `json:",inline"`
	Sn          string `json:"sn"`
	OrderId     string `json:"orderId,omitempty"`
}

// QueryOrderStateRequest 查询订单状态请求
type QueryOrderStateRequest struct {
	RestRequest `json:",inline"`
	OrderId     string `json:"orderId"`
}

// QueryOrderStatisRequest 查询订单统计请求
type QueryOrderStatisRequest struct {
	RestRequest `json:",inline"`
	Sn          string `json:"sn"`
	Date        string `json:"date"`
}

// VoiceRequest 语音播报请求
type VoiceRequest struct {
	RestRequest `json:",inline"`
	Sn          string  `json:"sn"`
	PayType     int     `json:"payType"`
	PayMode     int     `json:"payMode"`
	Money       float64 `json:"money"`
}

// VoicePlayMsgRequest 语音播报消息请求
type VoicePlayMsgRequest struct {
	RestRequest   `json:",inline"`
	Sn            string `json:"sn"`
	Content       string `json:"content"`
	VoiceTime     int    `json:"voiceTime,omitempty"`
	VoiceInterval int    `json:"voiceInterval,omitempty"`
}

// UploadLogoRequest 上传LOGO请求
type UploadLogoRequest struct {
	RestRequest `json:",inline"`
	Sn          string `json:"sn"`
	Content     string `json:"content"`
	LabelMode   int    `json:"labelMode,omitempty"`
	ImageSize   int    `json:"imageSize,omitempty"`
}