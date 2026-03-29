package cron

import (
	"fmt"
	"time"

	"github.com/Kevin-Jii/tower-go/service"
	"github.com/robfig/cron/v3"
)

// PrinterSyncJob 打印机状态同步任务
type PrinterSyncJob struct {
	printerService *service.PrinterService
}

// NewPrinterSyncJob 创建打印机同步任务
func NewPrinterSyncJob(printerService *service.PrinterService) *PrinterSyncJob {
	return &PrinterSyncJob{
		printerService: printerService,
	}
}

// Run 执行同步任务
func (j *PrinterSyncJob) Run() {
	start := time.Now()
	fmt.Printf("[PrinterSync] 开始同步打印机状态...\n")

	if err := j.printerService.SyncAllPrinterStatus(); err != nil {
		fmt.Printf("[PrinterSync] 同步失败: %v\n", err)
		return
	}

	fmt.Printf("[PrinterSync] 同步完成，耗时: %v\n", time.Since(start))
}

// StartPrinterSync 启动打印机状态同步定时任务
// 启动时立即同步一次，之后每30分钟同步一次
func StartPrinterSync(printerService *service.PrinterService) (*cron.Cron, error) {
	c := cron.New(cron.WithSeconds())

	job := NewPrinterSyncJob(printerService)

	// 启动时立即执行一次
	fmt.Println("[PrinterSync] 启动时执行首次同步...")
	job.Run()

	// 每30分钟同步一次 (0 */30 * * * *)
	_, err := c.AddFunc("0 */30 * * * *", job.Run)
	if err != nil {
		return nil, fmt.Errorf("添加打印机同步任务失败: %w", err)
	}

	c.Start()
	fmt.Println("[PrinterSync] 打印机状态同步定时任务已启动 (每30分钟执行)")

	return c, nil
}
