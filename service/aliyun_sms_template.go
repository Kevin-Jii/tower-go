package service

import (
	"errors"
	"strings"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/pkg/aliyunsms"
)

// AliyunSmsClient is the narrow provider contract needed by template management.
// Keeping the generated Aliyun SDK behind this interface makes the business service testable.
type AliyunSmsClient interface {
	CreateTemplate(name, content, relatedSign, remark string, templateType int32) (string, error)
	GetTemplate(templateCode string) (*aliyunsms.TemplateSnapshot, error)
	DeleteTemplate(templateCode string) error
	ListSignatures() ([]aliyunsms.SignatureSnapshot, error)
}

type AliyunSmsClientResolver interface {
	Resolve(storeID uint) (AliyunSmsClient, error)
}

type effectiveAliyunSmsClientResolver struct {
	storeConfig *StoreSmsConfigService
	fallback    *aliyunsms.Client
}

func NewAliyunSmsClientResolver(storeConfig *StoreSmsConfigService, fallback *aliyunsms.Client) AliyunSmsClientResolver {
	return &effectiveAliyunSmsClientResolver{storeConfig: storeConfig, fallback: fallback}
}

func (r *effectiveAliyunSmsClientResolver) Resolve(storeID uint) (AliyunSmsClient, error) {
	if r.storeConfig != nil && storeID > 0 {
		ctx, err := r.storeConfig.ResolveEffectiveContext(storeID)
		if err != nil {
			return nil, err
		}
		if ctx != nil && ctx.Client != nil {
			return ctx.Client, nil
		}
	}
	if r.fallback != nil && r.fallback.Enabled() {
		return r.fallback, nil
	}
	return nil, errors.New("未配置可用的阿里云短信凭证：请在会员推广 → 基础设置 中为本门店配置")
}

type AliyunSmsTemplateRepository interface {
	List(storeID uint, allStores bool, keyword, auditStatus string) ([]model.AliyunSmsTemplate, error)
	GetByCode(code string, storeID uint, allStores bool) (*model.AliyunSmsTemplate, error)
	ExistsByName(name string, storeID uint) (bool, error)
	Upsert(row *model.AliyunSmsTemplate) error
	Delete(row *model.AliyunSmsTemplate) error
}

type AliyunSmsTemplateService struct {
	repository AliyunSmsTemplateRepository
	resolver   AliyunSmsClientResolver
}

func NewAliyunSmsTemplateService(repository AliyunSmsTemplateRepository, resolver AliyunSmsClientResolver) *AliyunSmsTemplateService {
	return &AliyunSmsTemplateService{repository: repository, resolver: resolver}
}

func (s *AliyunSmsTemplateService) resolveClient(storeID uint) (AliyunSmsClient, error) {
	if s.resolver == nil {
		return nil, errors.New("阿里云短信客户端解析器未配置")
	}
	return s.resolver.Resolve(storeID)
}

func (s *AliyunSmsTemplateService) List(storeID uint, allStores bool, keyword, auditStatus string) ([]model.AliyunSmsTemplate, error) {
	return s.repository.List(storeID, allStores, keyword, auditStatus)
}

func (s *AliyunSmsTemplateService) ListApproved(storeID uint, allStores bool) ([]model.AliyunSmsTemplate, error) {
	return s.repository.List(storeID, allStores, "", model.SmsTemplateAuditApproved)
}

func (s *AliyunSmsTemplateService) ListSignatures(storeID uint, approvedOnly bool) ([]aliyunsms.SignatureSnapshot, error) {
	client, err := s.resolveClient(storeID)
	if err != nil {
		return nil, err
	}
	rows, err := client.ListSignatures()
	if err != nil {
		return nil, err
	}
	if !approvedOnly {
		return rows, nil
	}
	approved := make([]aliyunsms.SignatureSnapshot, 0, len(rows))
	for _, row := range rows {
		if row.AuditStatus == "AUDIT_STATE_PASS" || row.AuditStatus == "1" {
			approved = append(approved, row)
		}
	}
	return approved, nil
}

func (s *AliyunSmsTemplateService) Create(req *model.CreateAliyunSmsTemplateReq, ownerStoreID, createdBy uint) (*model.AliyunSmsTemplate, error) {
	client, err := s.resolveClient(ownerStoreID)
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
	exists, err := s.repository.ExistsByName(name, ownerStoreID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("当前门店已存在同名模板")
	}
	templateCode, err := client.CreateTemplate(name, content, relatedSign, remark, *req.TemplateType)
	if err != nil {
		return nil, err
	}
	row := &model.AliyunSmsTemplate{
		OwnerStoreID:    ownerStoreID,
		TemplateCode:    templateCode,
		Name:            name,
		Content:         content,
		TemplateType:    *req.TemplateType,
		RelatedSign:     relatedSign,
		Remark:          remark,
		AuditStatus:     model.SmsTemplateAuditPending,
		SourceCreatedBy: createdBy,
	}
	if err := s.repository.Upsert(row); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *AliyunSmsTemplateService) Refresh(templateCode string, storeID uint, allStores bool) (*model.AliyunSmsTemplate, error) {
	row, err := s.repository.GetByCode(templateCode, storeID, allStores)
	if err != nil {
		return nil, err
	}
	client, err := s.resolveClient(row.OwnerStoreID)
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
	if err := s.repository.Upsert(row); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *AliyunSmsTemplateService) Delete(templateCode string, storeID uint, allStores bool) error {
	row, err := s.repository.GetByCode(templateCode, storeID, allStores)
	if err != nil {
		return err
	}
	client, err := s.resolveClient(row.OwnerStoreID)
	if err != nil {
		return err
	}
	if err := client.DeleteTemplate(templateCode); err != nil {
		return err
	}
	return s.repository.Delete(row)
}

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
