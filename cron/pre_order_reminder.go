package cron

import (
	"fmt"
	"time"

	"github.com/Kevin-Jii/tower-go/service"
	"github.com/robfig/cron/v3"
)

func StartPreOrderReminders(preOrderService *service.PreOrderService) (*cron.Cron, error) {
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return nil, fmt.Errorf("加载预订单提醒时区失败: %w", err)
	}
	c := cron.New(cron.WithSeconds(), cron.WithLocation(location))
	if _, err := c.AddFunc("0 30 9 * * *", func() {
		if err := preOrderService.ProcessReminderSlot(time.Now(), service.PreOrderReminderSlot0930); err != nil {
			fmt.Printf("[PreOrderReminder] 09:30 提醒处理失败: %v\n", err)
		}
	}); err != nil {
		return nil, fmt.Errorf("添加预订单09:30提醒失败: %w", err)
	}
	if _, err := c.AddFunc("0 0 16 * * *", func() {
		if err := preOrderService.ProcessReminderSlot(time.Now(), service.PreOrderReminderSlot1600); err != nil {
			fmt.Printf("[PreOrderReminder] 16:00 提醒处理失败: %v\n", err)
		}
	}); err != nil {
		return nil, fmt.Errorf("添加预订单16:00提醒失败: %w", err)
	}
	c.Start()
	fmt.Println("[PreOrderReminder] 预订单提醒任务已启动 (09:30 / 16:00)")
	return c, nil
}
