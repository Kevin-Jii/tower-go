package cron

import (
	"context"
	"fmt"
	"time"

	"github.com/Kevin-Jii/tower-go/service"
	"github.com/robfig/cron/v3"
)

func StartGalleryUploadCleanup(galleryService *service.GalleryService) (*cron.Cron, error) {
	if galleryService == nil {
		return nil, nil
	}
	cleanup := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()
		if err := galleryService.CleanupExpiredMultipartUploads(ctx); err != nil {
			fmt.Printf("[GalleryUploadCleanup] 清理失败: %v\n", err)
		}
	}

	c := cron.New(cron.WithSeconds())
	if _, err := c.AddFunc("0 15 * * * *", cleanup); err != nil {
		return nil, fmt.Errorf("添加图库上传会话清理任务失败: %w", err)
	}
	cleanup()
	c.Start()
	fmt.Println("[GalleryUploadCleanup] 图库上传会话清理任务已启动 (每小时执行)")
	return c, nil
}
