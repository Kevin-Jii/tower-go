package middleware

import (
	"sync"
	"time"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/utils/database"
	"github.com/Kevin-Jii/tower-go/utils/logging"
	"go.uber.org/zap"
)

const (
	auditLogQueueSize  = 2048
	auditLogBatchSize  = 100
	auditLogFlushEvery = 500 * time.Millisecond
)

// auditLogWriter decouples request latency from audit log persistence.
// The queue is bounded so a database outage applies backpressure instead of
// allowing unbounded memory growth or silently dropping audit records.
type auditLogWriter struct {
	queue chan *model.AuditLog
	once  sync.Once
}

func newAuditLogWriter() *auditLogWriter {
	w := &auditLogWriter{
		queue: make(chan *model.AuditLog, auditLogQueueSize),
	}
	w.once.Do(func() { go w.run() })
	return w
}

func (w *auditLogWriter) enqueue(log *model.AuditLog) {
	if w == nil || log == nil || database.DB == nil {
		return
	}
	w.queue <- log
}

func (w *auditLogWriter) run() {
	ticker := time.NewTicker(auditLogFlushEvery)
	defer ticker.Stop()

	batch := make([]*model.AuditLog, 0, auditLogBatchSize)
	flush := func() {
		if len(batch) == 0 {
			return
		}
		db := database.DB
		if db == nil {
			return
		}
		if err := db.CreateInBatches(batch, auditLogBatchSize).Error; err != nil {
			logging.LogWarn("批量写入操作日志失败，将重试", zap.Error(err), zap.Int("count", len(batch)))
			return
		}
		batch = batch[:0]
	}

	for {
		select {
		case log := <-w.queue:
			batch = append(batch, log)
			if len(batch) >= auditLogBatchSize {
				flush()
			}
		case <-ticker.C:
			flush()
		}
	}
}
