package service

import (
	"strings"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
)

type AuditLogService struct {
	module *module.AuditLogModule
}

func NewAuditLogService(module *module.AuditLogModule) *AuditLogService {
	return &AuditLogService{module: module}
}

func (s *AuditLogService) Create(log *model.AuditLog) error {
	if s == nil || s.module == nil || log == nil {
		return nil
	}
	normalizeAuditLog(log)
	return s.module.Create(log)
}

func (s *AuditLogService) List(req model.AuditLogListReq, allowedStoreID uint, crossStore bool) ([]*model.AuditLog, int64, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 10
	}
	return s.module.List(req, allowedStoreID, crossStore)
}

func (s *AuditLogService) Get(id uint, allowedStoreID uint, crossStore bool) (*model.AuditLog, error) {
	log, err := s.module.GetByID(id)
	if err != nil {
		return nil, err
	}
	if !crossStore && log.StoreID != allowedStoreID {
		return nil, ErrForbiddenAuditLog
	}
	return log, nil
}

var ErrForbiddenAuditLog = &auditLogError{"无权查看该操作日志"}

type auditLogError struct {
	msg string
}

func (e *auditLogError) Error() string {
	return e.msg
}

func normalizeAuditLog(log *model.AuditLog) {
	log.Status = strings.TrimSpace(log.Status)
	if log.Status == "" {
		log.Status = model.AuditStatusSuccess
	}
	log.Module = strings.TrimSpace(log.Module)
	if log.Module == "" {
		log.Module = "system"
	}
	log.ModuleName = strings.TrimSpace(log.ModuleName)
	if log.ModuleName == "" {
		log.ModuleName = "系统"
	}
	log.Action = strings.TrimSpace(log.Action)
	if log.Action == "" {
		log.Action = model.AuditActionOther
	}
	log.ActionName = strings.TrimSpace(log.ActionName)
	if log.ActionName == "" {
		log.ActionName = actionName(log.Action)
	}
	log.ResourceType = strings.TrimSpace(log.ResourceType)
	if log.ResourceType == "" {
		log.ResourceType = log.Module
	}
}

func actionName(action string) string {
	switch action {
	case model.AuditActionLogin:
		return "登录"
	case model.AuditActionCreate:
		return "新增"
	case model.AuditActionUpdate:
		return "修改"
	case model.AuditActionDelete:
		return "删除"
	case model.AuditActionQuery:
		return "查询"
	default:
		return "操作"
	}
}
