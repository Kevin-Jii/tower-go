package module

import (
	"strings"
	"time"

	"github.com/Kevin-Jii/tower-go/model"
	"gorm.io/gorm"
)

type AuditLogModule struct {
	db *gorm.DB
}

func NewAuditLogModule(db *gorm.DB) *AuditLogModule {
	return &AuditLogModule{db: db}
}

func (m *AuditLogModule) Create(log *model.AuditLog) error {
	if m == nil || m.db == nil || log == nil {
		return nil
	}
	return m.db.Create(log).Error
}

func (m *AuditLogModule) GetByID(id uint) (*model.AuditLog, error) {
	var log model.AuditLog
	if err := m.db.First(&log, id).Error; err != nil {
		return nil, err
	}
	return &log, nil
}

func (m *AuditLogModule) List(req model.AuditLogListReq, allowedStoreID uint, crossStore bool) ([]*model.AuditLog, int64, error) {
	var rows []*model.AuditLog
	var total int64

	q := m.db.Model(&model.AuditLog{})
	if !crossStore {
		q = q.Where("store_id = ?", allowedStoreID)
	} else if req.StoreID > 0 {
		q = q.Where("store_id = ?", req.StoreID)
	}
	if req.UserID > 0 {
		q = q.Where("user_id = ?", req.UserID)
	}
	if req.Module != "" {
		q = q.Where("module = ?", strings.TrimSpace(req.Module))
	}
	if req.Action != "" {
		q = q.Where("action = ?", strings.TrimSpace(req.Action))
	}
	if req.Status != "" {
		q = q.Where("status = ?", strings.TrimSpace(req.Status))
	}
	if req.Keyword != "" {
		kw := "%" + strings.TrimSpace(req.Keyword) + "%"
		q = q.Where(
			"username LIKE ? OR nickname LIKE ? OR phone LIKE ? OR module_name LIKE ? OR action_name LIKE ? OR resource_no LIKE ? OR resource_name LIKE ? OR path LIKE ?",
			kw, kw, kw, kw, kw, kw, kw, kw,
		)
	}
	if t, ok := parseAuditTime(req.StartTime); ok {
		q = q.Where("created_at >= ?", t)
	}
	if t, ok := parseAuditTime(req.EndTime); ok {
		q = q.Where("created_at <= ?", t)
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}
	offset := (req.Page - 1) * req.PageSize
	if err := q.Order("created_at DESC, id DESC").Offset(offset).Limit(req.PageSize).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func parseAuditTime(s string) (time.Time, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, false
	}
	layouts := []string{
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02",
	}
	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, s, time.Local); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}
