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
	Order        *model.MenuReportOrder
	StoreName    string
	UserName     string
	StorePhone   string
	StoreAddress string
	BotID        uint // 指定发送的机器人ID
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
	imageData, err := utils.GenerateMenuReportImage(e.Order, e.StoreName, e.UserName, e.StorePhone, e.StoreAddress)
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

	// 如果生成了图片，保存到文件系统（即使没有机器人也保存）
	if imageData != nil {
		imageURL, err := utils.SaveImageFile("menu_report.png", imageData)
		if err != nil {
			if logging.SugaredLogger != nil {
				logging.SugaredLogger.Warnw("Failed to save image to file system",
					"orderID", e.Order.ID,
					"error", err)
			}
		} else {
			if logging.SugaredLogger != nil {
				logging.SugaredLogger.Infow("Image saved to file system",
					"orderID", e.Order.ID,
					"imageURL", imageURL,
					"imageSize", len(imageData))
			}
		}
	}

	// 如果没有指定机器人ID，跳过发送通知
	if e.BotID == 0 {
		if logging.SugaredLogger != nil {
			logging.SugaredLogger.Infow("No bot specified, skipping notification",
				"orderID", e.Order.ID,
				"storeID", e.Order.StoreID)
		}
		return nil
	}

	// 获取指定的机器人
	bot, err := l.dingTalkSvc.GetBot(e.BotID)
	if err != nil {
		if logging.SugaredLogger != nil {
			logging.SugaredLogger.Errorw("Failed to get bot",
				"orderID", e.Order.ID,
				"botID", e.BotID,
				"error", err)
		}
		return err
	}

	// 发送到指定机器人
	if bot.BotType == "stream" {
		// Stream 模式：通过钉钉服务端 API 发送
		if imageData != nil {
			err = l.dingTalkSvc.SendStreamImageText(bot, title, content, imageData)
		} else {
			err = l.dingTalkSvc.SendStreamMarkdown(bot, title, content)
		}
	} else {
		// Webhook 模式：直接 HTTP POST（不支持直接显示图片，但可以发送图片链接）
		contentWithImage := content
		if imageData != nil {
			// 获取图片 URL（已经在前面保存过了）
			imageURL, err := utils.SaveImageFile("menu_report.png", imageData)
			if err == nil {
				// 在内容末尾添加图片链接
				contentWithImage = fmt.Sprintf("%s\n\n**📷 查看报菜图片:**\n[点击查看](%s)", content, imageURL)
				if logging.SugaredLogger != nil {
					logging.SugaredLogger.Infow("Added image link to webhook message",
						"botID", bot.ID,
						"imageURL", imageURL)
				}
			}
		}
		err = l.dingTalkSvc.SendMarkdownToBot(bot, title, contentWithImage)
	}

	if err != nil {
		if logging.SugaredLogger != nil {
			logging.SugaredLogger.Errorw("Failed to send menu report notification",
				"orderID", e.Order.ID,
				"botID", bot.ID,
				"botType", bot.BotType,
				"error", err)
		}
		return err
	}

	if logging.SugaredLogger != nil {
		logging.SugaredLogger.Infow("Menu report notification sent successfully",
			"orderID", e.Order.ID,
			"botID", bot.ID,
			"botName", bot.Name)
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
