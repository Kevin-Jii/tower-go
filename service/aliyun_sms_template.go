package service

import (
	"errors"
	"strings"

	"github.com/Kevin-Jii/tower-go/config"
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
	"github.com/Kevin-Jii/tower-go/pkg/aliyunsms"
)

type AliyunSmsTemplateService struct {
	module         *module.AliyunSmsTemplateModule
	client         *aliyunsms.Client
	storeSmsConfig *StoreSmsConfigService
}

func NewAliyunSmsTemplateService(m *module.AliyunSmsTemplateModule) *AliyunSmsTemplateService {
	cfg := config.GetAliyunSMSConfig()
	client, err := aliyunsms.NewClient(aliyunsms.Config{
		AccessKeyID:     cfg.AccessKeyID,
		AccessKeySecret: cfg.AccessKeySecret,
		RegionID:        cfg.RegionID,
		SignName:        cfg.SignName,
		Enabled:         cfg.Enabled,
	})
	if err != nil {
		client = nil
	}
	return &AliyunSmsTemplateService{module: m, client: client, storeSmsConfig: nil}
}

// SetStoreSmsConfigService 注入门店独立 SMS 配置服务，用于按门店客户端调用模板 API。
func (s *AliyunSmsTemplateService) SetStoreSmsConfigService(svc *StoreSmsConfigService) {
	s.storeSmsConfig = svc
}

func (s *AliyunSmsTemplateService) resolveClient(storeID uint) (*aliyunsms.Client, error) {
	if s.storeSmsConfig != nil {
		ctx, err := s.storeSmsConfig.ResolveEffectiveContext(storeID)
		if err == nil && ctx != nil && ctx.Client != nil {
			return ctx.Client, nil
		}
	}
	if s.client != nil && s.client.Enabled() {
		return s.client, nil
	}
	return nil, errors.New("未配置可用的阿里云短信凭证：请在会员推广 → 基础设置 中为本门店配置")
}

func (s *AliyunSmsTemplateService) List(keyword, auditStatus string) ([]model.AliyunSmsTemplate, error) {
	return s.module.List(keyword, auditStatus)
}

func (s *AliyunSmsTemplateService) Create(req *model.CreateAliyunSmsTemplateReq, createdBy uint) (*model.AliyunSmsTemplate, error) {
	storeID := req.OwnerStoreID
	client, err := s.resolveClient(storeID)
	if err != nil {
		return nil, err
	}
	name := strings.TrimSpace(req.Name)
	content := strings.TrimSpace(req.Content)
	relatedSign := strings.TrimSpace(req.RelatedSign)
	remark := strings.TrimSpace(req.Remark)
	if name == "" || content == "" || req.TemplateType == nil {
		return nil, errors.New("请填写模板名称、内容与类型")
	}
	if existing, _ := s.module.GetByCode(""); existing != nil && existing.Name == name {
		return nil, errors.New("已存在同名模板")
	}
	templateCode, err := client.CreateTemplate(name, content, relatedSign, remark, *req.TemplateType)
	if err != nil {
		return nil, err
	}
	row := &model.AliyunSmsTemplate{
		TemplateCode:    templateCode,
		Name:            name,
		Content:         content,
		TemplateType:    *req.TemplateType,
		RelatedSign:     relatedSign,
		Remark:          remark,
		AuditStatus:     model.SmsTemplateAuditPending,
		SourceCreatedBy: createdBy,
	}
	if err := s.module.Upsert(row); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *AliyunSmsTemplateService) Refresh(templateCode string) (*model.AliyunSmsTemplate, error) {
	client, err := s.resolveClient(0)
	if err != nil {
		return nil, err
	}
	row, err := s.module.GetByCode(templateCode)
	if err != nil {
		return nil, err
	}
	snap, err := client.GetTemplate(templateCode)
	if err != nil {
		return nil, err
	}
	row.AuditStatus = mapAliyunAuditStatus(snap.TemplateStatus)
	row.AuditReason = snap.Reason
	if snap.TemplateContent != "" {
		row.Content = snap.TemplateContent
	}
	if err := s.module.Upsert(row); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *AliyunSmsTemplateService) Delete(templateCode string) error {
	client, err := s.resolveClient(0)
	if err != nil {
		return err
	}
	if err := client.DeleteTemplate(templateCode); err != nil {
		return err
	}
	return s.module.DeleteByCode(templateCode)
}

func (s *AliyunSmsTemplateService) ListApproved() ([]model.AliyunSmsTemplate, error) {
	return s.module.List("", model.SmsTemplateAuditApproved)
}

// mapAliyunAuditStatus converts the Aliyun TemplateStatus string to our local enum.
func mapAliyunAuditStatus(raw string) string {
	switch raw {
	case "0", "AUDIT_STATE_INIT":
		return model.SmsTemplateAuditPending
	case "1", "AUDIT_STATE_PASS":
		return model.SmsTemplateAuditApproved
	case "2", "AUDIT_STATE_NOT_PASS":
		return model.SmsTemplateAuditRejected
	case "10", "AUDIT_STATE_CANCEL":
		return model.SmsTemplateAuditCancelled
	default:
		return model.SmsTemplateAuditUnknown
	}
}
