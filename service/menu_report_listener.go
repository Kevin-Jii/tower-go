package service

import (
	"fmt"
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/utils"
	"github.com/Kevin-Jii/tower-go/utils/events"
	"github.com/Kevin-Jii/tower-go/utils/logging"
)

// MenuReportOrderCreatedEvent 报菜记录单创建事件
type MenuReportOrderCreatedEvent struct {
	Order     *model.MenuReportOrder
	StoreName string
	UserName  string
}

// Name 实现 Event 接口
func (e MenuReportOrderCreatedEvent) Name() string {
	return "menu_report_order.created"
}

// MenuReportEventListener 报菜事件监听器
type MenuReportEventListener struct {
	dingTalkSvc *DingTalkService
}

func NewMenuReportEventListener(dingTalkSvc *DingTalkService) *MenuReportEventListener {
	return &MenuReportEventListener{
		dingTalkSvc: dingTalkSvc,
	}
}

// OnMenuReportOrderCreated 处理报菜记录单创建事件
func (l *MenuReportEventListener) OnMenuReportOrderCreated(event events.Event) error {
	e, ok := event.(MenuReportOrderCreatedEvent)
	if !ok {
		return fmt.Errorf("invalid event type")
	}

	// 构建推送消息
	title := "📋 新报菜通知"
	content := l.buildNotificationContent(e)

	// 生成PNG图片
	imageData, err := utils.GenerateMenuReportImage(e.Order, e.StoreName, e.UserName)
	if err != nil {
		if logging.SugaredLogger != nil {
			logging.SugaredLogger.Warnw("Failed to generate menu report image, sending text only",
				"orderID", e.Order.ID,
				"error", err)
		}
		// 图片生成失败,仍然发送文本消息
		imageData = nil
	} else {
		if logging.SugaredLogger != nil {
			logging.SugaredLogger.Infow("Menu report image generated successfully",
				"orderID", e.Order.ID,
				"imageSize", len(imageData))
		}
	}

	// 广播到门店的所有机器人（带图片）
	if err := l.dingTalkSvc.BroadcastToStoreWithImage(e.Order.StoreID, "markdown", title, content, imageData); err != nil {
		if logging.SugaredLogger != nil {
			logging.SugaredLogger.Errorw("Failed to broadcast menu report order",
				"orderID", e.Order.ID,
				"storeID", e.Order.StoreID,
				"error", err)
		}
		return err
	}

	return nil
}

// buildNotificationContent 构建通知内容
func (l *MenuReportEventListener) buildNotificationContent(e MenuReportOrderCreatedEvent) string {
	createdAt := e.Order.CreatedAt.Format("2006-01-02 15:04:05")

	content := fmt.Sprintf(`## 📋 新报菜通知

**门店名称:** %s
**操作人员:** %s
**报菜时间:** %s

**报菜明细:**
`, e.StoreName, e.UserName, createdAt)

	for _, item := range e.Order.Items {
		if item.Dish != nil {
			content += fmt.Sprintf("- **%s**: 数量 %d", item.Dish.Name, item.Quantity)
			if item.Remark != "" {
				content += fmt.Sprintf(" (%s)", item.Remark)
			}
			content += "\n"
		}
	}

	if e.Order.Remark != "" {
		content += fmt.Sprintf("\n**备注:** %s\n", e.Order.Remark)
	}

	content += "\n---\n"
	content += fmt.Sprintf("*报菜记录单ID: %d*", e.Order.ID)

	return content
}

// RegisterMenuReportEventListeners 注册报菜事件监听器
func RegisterMenuReportEventListeners(eventBus *events.EventBus, listener *MenuReportEventListener) {
	eventBus.Subscribe("menu_report_order.created", listener.OnMenuReportOrderCreated)
}
