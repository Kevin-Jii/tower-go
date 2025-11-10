package service

import (
	"fmt"
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/utils/events"
	"github.com/Kevin-Jii/tower-go/utils/logging"
)

// MenuReportCreatedEvent 报菜记录创建事件
type MenuReportCreatedEvent struct {
	Report    *model.MenuReport
	StoreName string
	DishName  string
	UserName  string
}

// Name 实现 Event 接口
func (e MenuReportCreatedEvent) Name() string {
	return "menu_report.created"
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

// OnMenuReportCreated 处理报菜创建事件
func (l *MenuReportEventListener) OnMenuReportCreated(event events.Event) error {
	e, ok := event.(MenuReportCreatedEvent)
	if !ok {
		return fmt.Errorf("invalid event type")
	}

	// 构建推送消息
	title := "📋 新报菜通知"
	content := l.buildNotificationContent(e)

	// 广播到门店的所有机器人
	if err := l.dingTalkSvc.BroadcastToStore(e.Report.StoreID, "markdown", title, content); err != nil {
		if logging.SugaredLogger != nil {
			logging.SugaredLogger.Errorw("Failed to broadcast menu report",
				"reportID", e.Report.ID,
				"storeID", e.Report.StoreID,
				"error", err)
		}
		return err
	}

	return nil
}

// buildNotificationContent 构建通知内容
func (l *MenuReportEventListener) buildNotificationContent(e MenuReportCreatedEvent) string {
	createdAt := e.Report.CreatedAt.Format("2006-01-02 15:04:05")

	content := fmt.Sprintf(`## 📋 新报菜通知

**菜品名称:** %s  
**报菜数量:** %d  
**门店名称:** %s  
**操作人员:** %s  
**报菜时间:** %s  
`, e.DishName, e.Report.Quantity, e.StoreName, e.UserName, createdAt)

	if e.Report.Remark != "" {
		content += fmt.Sprintf("**备注:** %s  \n", e.Report.Remark)
	}

	content += "\n---\n"
	content += fmt.Sprintf("*报菜记录ID: %d*", e.Report.ID)

	return content
}

// RegisterMenuReportEventListeners 注册报菜事件监听器
func RegisterMenuReportEventListeners(eventBus *events.EventBus, listener *MenuReportEventListener) {
	eventBus.Subscribe("menu_report.created", listener.OnMenuReportCreated)
}
