package service

import (
	"fmt"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/utils"
)

// MenuReportNotifier 报菜记录钉钉通知器
type MenuReportNotifier struct {
	dingTalkService *DingTalkService
}

// NewMenuReportNotifier 创建报菜记录通知器
func NewMenuReportNotifier(dingTalkService *DingTalkService) *MenuReportNotifier {
	return &MenuReportNotifier{
		dingTalkService: dingTalkService,
	}
}

// Update 实现 Observer 接口
func (n *MenuReportNotifier) Update(event utils.EventType, data interface{}) {
	if event == utils.EventMenuReportCreated {
		if order, ok := data.(*model.MenuReportOrder); ok {
			n.sendMenuReportNotification(order)
		}
	}
}

// sendMenuReportNotification 发送报菜记录通知
func (n *MenuReportNotifier) sendMenuReportNotification(order *model.MenuReportOrder) {
	// 异步发送，不影响主流程
	go func() {
		// 构建消息内容
		title, content := n.buildMenuReportMessage(order)

		// 使用 DingTalkService 的 BroadcastToStore 方法推送到门店机器人
		err := n.dingTalkService.BroadcastToStore(order.StoreID, "markdown", title, content)
		if err != nil {
			fmt.Printf("钉钉通知发送失败: %v\n", err)
		}
	}()
}

// buildMenuReportMessage 构造报菜消息
func (n *MenuReportNotifier) buildMenuReportMessage(order *model.MenuReportOrder) (string, string) {
	var itemsText string
	for _, item := range order.Items {
		if item != nil && item.Dish != nil {
			itemsText += fmt.Sprintf(
				"- **%s**: %d 份\n",
				item.Dish.Name,
				item.Quantity,
			)
		}
	}

	title := "📋 新报菜通知"

	// 消息内容
	message := fmt.Sprintf(
		"### 📋 新报菜通知\n\n"+
			"**门店**: %s\n\n"+
			"**操作员**: %s\n\n",
		order.Store.Name,
	)

	if itemsText != "" {
		message += fmt.Sprintf("**报菜品项**:\n%s\n", itemsText)
	}

	if order.Remark != "" {
		message += fmt.Sprintf("### 📝 备注\n%s\n\n", order.Remark)
	}

	message += fmt.Sprintf("**时间**: %s", order.CreatedAt.Format("2006-01-02 15:04:05"))

	return title, message
}

// Register 注册到事件总线
func (n *MenuReportNotifier) Register() {
	utils.GlobalEventBus.Register(utils.EventMenuReportCreated, n)
}
