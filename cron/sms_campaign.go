package cron

import (
	"fmt"
	"time"

	"github.com/Kevin-Jii/tower-go/service"
	"github.com/robfig/cron/v3"
)

func StartSmsCampaignScheduler(smsService *service.SmsCampaignService) (*cron.Cron, error) {
	if smsService == nil {
		return nil, fmt.Errorf("短信推广服务未初始化")
	}
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return nil, fmt.Errorf("加载短信推广时区失败: %w", err)
	}
	c := cron.New(cron.WithSeconds(), cron.WithLocation(location))
	if _, err := c.AddFunc("0 */1 * * * *", func() {
		if err := smsService.ProcessDueScheduled(time.Now()); err != nil {
			fmt.Printf("[SmsCampaign] 定时任务失败: %v\n", err)
		}
	}); err != nil {
		return nil, fmt.Errorf("添加短信推广定时任务失败: %w", err)
	}
	c.Start()
	fmt.Println("[SmsCampaign] 短信推广定时任务已启动 (每分钟检查排期)")
	return c, nil
}
