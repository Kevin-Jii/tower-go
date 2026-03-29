package xpyun

import (
	"github.com/Kevin-Jii/tower-go/pkg/xpyun/model"
	"github.com/Kevin-Jii/tower-go/pkg/xpyun/service"
	"github.com/Kevin-Jii/tower-go/pkg/xpyun/util"
)

// Config 芯烨云配置
type Config struct {
	User    string
	UserKey string
	BaseURL string
}

// Client 芯烨云客户端
type Client struct {
	config *Config
}

// NewClient 创建芯烨云客户端
func NewClient(user, userKey string) *Client {
	return NewClientWithBaseURL(user, userKey, "https://open.xpyun.net/api/openapi")
}

// NewClientWithBaseURL 创建芯烨云客户端（带自定义API地址）
func NewClientWithBaseURL(user, userKey, baseURL string) *Client {
	return &Client{
		config: &Config{
			User:    user,
			UserKey: userKey,
			BaseURL: baseURL,
		},
	}
}

// AddPrinter 添加打印机
func (c *Client) AddPrinter(sn, name string) *model.XPYunResp {
	request := model.AddPrinterRequest{}
	request.User = c.config.User
	request.UserKey = c.config.UserKey
	request.Timestamp = util.GetMillisecond()
	request.GenerateSign()

	requestItem := model.AddPrinterRequestItem{}
	requestItem.Sn = sn
	if name != "" {
		requestItem.Name = name
	}

	request.Items = []*model.AddPrinterRequestItem{&requestItem}
	return service.XpYunAddPrintersWithURL(&request, c.config.BaseURL)
}

// DelPrinter 删除打印机
func (c *Client) DelPrinter(snList []string) *model.XPYunResp {
	request := model.DelPrinterRequest{}
	request.User = c.config.User
	request.UserKey = c.config.UserKey
	request.Timestamp = util.GetMillisecond()
	request.GenerateSign()
	request.SnList = snList
	return service.XpYunDelPrintersWithURL(&request, c.config.BaseURL)
}

// UpdatePrinter 更新打印机信息
func (c *Client) UpdatePrinter(sn, name string) *model.XPYunResp {
	request := model.UpdPrinterRequest{}
	request.User = c.config.User
	request.UserKey = c.config.UserKey
	request.Timestamp = util.GetMillisecond()
	request.GenerateSign()
	request.Sn = sn
	request.Name = name
	return service.XpYunUpdatePrinterWithURL(&request, c.config.BaseURL)
}

// PrintReceipt 打印小票
// content: 使用标签格式化的打印内容，如 "<C>标题<BR><BR>内容"
func (c *Client) PrintReceipt(sn, content string, copies int) *model.XPYunResp {
	request := model.PrintRequest{}
	request.User = c.config.User
	request.UserKey = c.config.UserKey
	request.Timestamp = util.GetMillisecond()
	request.GenerateSign()
	request.Sn = sn
	request.Content = content
	request.Copies = copies
	request.Mode = 1 // 不检查打印机是否在线
	request.Voice = 2 // 来单播放模式
	return service.XpYunPrintWithURL(&request, c.config.BaseURL)
}

// PrintLabel 打印标签
func (c *Client) PrintLabel(sn, content string, copies int) *model.XPYunResp {
	request := model.PrintLabelRequest{}
	request.User = c.config.User
	request.UserKey = c.config.UserKey
	request.Timestamp = util.GetMillisecond()
	request.GenerateSign()
	request.Sn = sn
	request.Content = content
	request.Copies = copies
	request.Mode = 1
	return service.XpYunPrintLabelWithURL(&request, c.config.BaseURL)
}

// QueryPrinterStatus 查询打印机状态
// 返回: 0-离线, 1-在线正常, 2-在线异常
func (c *Client) QueryPrinterStatus(sn string) *model.XPYunResp {
	request := model.PrinterRequest{}
	request.User = c.config.User
	request.UserKey = c.config.UserKey
	request.Timestamp = util.GetMillisecond()
	request.GenerateSign()
	request.Sn = sn
	return service.XpYunQueryPrinterStatusWithURL(&request, c.config.BaseURL)
}

// QueryOrderState 查询订单状态
func (c *Client) QueryOrderState(orderId string) *model.XPYunResp {
	request := model.QueryOrderStateRequest{}
	request.User = c.config.User
	request.UserKey = c.config.UserKey
	request.Timestamp = util.GetMillisecond()
	request.GenerateSign()
	request.OrderId = orderId
	return service.XpYunQueryOrderStateWithURL(&request, c.config.BaseURL)
}

// ClearPrintQueue 清空打印队列
func (c *Client) ClearPrintQueue(sn string) *model.XPYunResp {
	request := model.ClearPrintOrderRequest{}
	request.User = c.config.User
	request.UserKey = c.config.UserKey
	request.Timestamp = util.GetMillisecond()
	request.GenerateSign()
	request.Sn = sn
	return service.XpYunDelPrinterQueueWithURL(&request, c.config.BaseURL)
}

// PlayVoice 金额播报
func (c *Client) PlayVoice(sn string, payType, payMode int, money float64) *model.XPYunResp {
	request := model.VoiceRequest{}
	request.User = c.config.User
	request.UserKey = c.config.UserKey
	request.Timestamp = util.GetMillisecond()
	request.GenerateSign()
	request.Sn = sn
	request.PayType = payType
	request.PayMode = payMode
	request.Money = money
	return service.XpYunPlayVoiceWithURL(&request, c.config.BaseURL)
}

// GetPrinterInfo 获取打印机信息
func (c *Client) GetPrinterInfo(sn string) *model.XPYunResp {
	request := model.PrinterRequest{}
	request.User = c.config.User
	request.UserKey = c.config.UserKey
	request.Timestamp = util.GetMillisecond()
	request.GenerateSign()
	request.Sn = sn
	return service.XpYunPrinterInfoWithURL(&request, c.config.BaseURL)
}