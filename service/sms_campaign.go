package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Kevin-Jii/tower-go/config"
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
	"github.com/Kevin-Jii/tower-go/pkg/aliyunsms"
)

type SmsCampaignService struct {
	campaignModule *module.SmsCampaignModule
	memberModule   *module.MemberModule
	smsClient      *aliyunsms.Client
}

func NewSmsCampaignService(campaignModule *module.SmsCampaignModule, memberModule *module.MemberModule) *SmsCampaignService {
	smsCfg := config.GetAliyunSMSConfig()
	client, err := aliyunsms.NewClient(aliyunsms.Config{
		AccessKeyID:     smsCfg.AccessKeyID,
		AccessKeySecret: smsCfg.AccessKeySecret,
		RegionID:        smsCfg.RegionID,
		SignName:        smsCfg.SignName,
		Enabled:         smsCfg.Enabled,
	})
	if err != nil {
		fmt.Printf("[SmsCampaign] 初始化阿里云短信客户端失败: %v\n", err)
	}
	return &SmsCampaignService{
		campaignModule: campaignModule,
		memberModule:   memberModule,
		smsClient:      client,
	}
}

func (s *SmsCampaignService) GetConfig() model.SmsServiceConfigResp {
	cfg := config.GetAliyunSMSConfig()
	return model.SmsServiceConfigResp{
		Enabled:        s.smsClient != nil && s.smsClient.Enabled(),
		DefaultSign:    cfg.SignName,
		Region:         cfg.RegionID,
		Configured:     cfg.AccessKeyID != "" && cfg.AccessKeySecret != "",
		HelpURL:        "https://help.aliyun.com/zh/sms/getting-started/sms-skill-guide",
		TemplateHint:   "请在阿里云短信控制台申请营销/通知类模板，填写审核通过的 TemplateCode；变量以 JSON 形式传入，如 {\"activity\":\"双十一\"}",
		MaxBatchPhones: aliyunsms.MaxPhonesPerRequest,
	}
}

func (s *SmsCampaignService) List() ([]*model.SmsCampaign, error) {
	return s.campaignModule.List()
}

func (s *SmsCampaignService) GetByID(id uint) (*model.SmsCampaign, error) {
	return s.campaignModule.GetByID(id)
}

func (s *SmsCampaignService) ListSendRecords(campaignID uint) ([]*model.SmsSendRecord, error) {
	return s.campaignModule.ListSendRecords(campaignID, 500)
}

func (s *SmsCampaignService) Create(req *model.CreateSmsCampaignReq, createdBy uint) (*model.SmsCampaign, error) {
	row := &model.SmsCampaign{
		Name:            req.Name,
		CampaignType:    req.CampaignType,
		SignName:        strings.TrimSpace(req.SignName),
		TemplateCode:    strings.TrimSpace(req.TemplateCode),
		TemplateParam:   strings.TrimSpace(req.TemplateParam),
		PersonalizeName: req.PersonalizeName != nil && *req.PersonalizeName,
		TargetType:      req.TargetType,
		StoreIDs:        req.StoreIDs,
		CustomPhones:    normalizePhoneList(req.CustomPhones),
		ScheduledAt:     req.ScheduledAt,
		Status:          model.SmsCampaignStatusDraft,
		CreatedBy:       createdBy,
	}
	if err := s.validateCampaign(row); err != nil {
		return nil, err
	}
	if row.ScheduledAt != nil && row.ScheduledAt.After(time.Now()) {
		row.Status = model.SmsCampaignStatusScheduled
	}
	if err := s.campaignModule.Create(row); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *SmsCampaignService) Update(id uint, req *model.UpdateSmsCampaignReq) error {
	row, err := s.campaignModule.GetByID(id)
	if err != nil {
		return err
	}
	if row.Status == model.SmsCampaignStatusSending || row.Status == model.SmsCampaignStatusSent {
		return errors.New("已发送或发送中的活动不可编辑")
	}
	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.CampaignType != nil {
		updates["campaign_type"] = *req.CampaignType
	}
	if req.SignName != nil {
		updates["sign_name"] = strings.TrimSpace(*req.SignName)
	}
	if req.TemplateCode != nil {
		updates["template_code"] = strings.TrimSpace(*req.TemplateCode)
	}
	if req.TemplateParam != nil {
		updates["template_param"] = strings.TrimSpace(*req.TemplateParam)
	}
	if req.PersonalizeName != nil {
		updates["personalize_name"] = *req.PersonalizeName
	}
	if req.TargetType != nil {
		updates["target_type"] = *req.TargetType
	}
	if req.StoreIDs != nil {
		updates["store_ids"] = *req.StoreIDs
	}
	if req.CustomPhones != nil {
		updates["custom_phones"] = normalizePhoneList(*req.CustomPhones)
	}
	if req.ScheduledAt != nil {
		updates["scheduled_at"] = req.ScheduledAt
		if req.ScheduledAt.After(time.Now()) {
			updates["status"] = model.SmsCampaignStatusScheduled
		} else {
			updates["status"] = model.SmsCampaignStatusDraft
		}
	}
	if len(updates) == 0 {
		return nil
	}
	merged := *row
	for k, v := range updates {
		switch k {
		case "name":
			merged.Name = v.(string)
		case "campaign_type":
			merged.CampaignType = v.(string)
		case "sign_name":
			merged.SignName = v.(string)
		case "template_code":
			merged.TemplateCode = v.(string)
		case "template_param":
			merged.TemplateParam = v.(string)
		case "personalize_name":
			merged.PersonalizeName = v.(bool)
		case "target_type":
			merged.TargetType = v.(string)
		case "store_ids":
			merged.StoreIDs = v.(model.UintList)
		case "custom_phones":
			merged.CustomPhones = v.(model.StringList)
		}
	}
	if err := s.validateCampaign(&merged); err != nil {
		return err
	}
	return s.campaignModule.Update(id, updates)
}

func (s *SmsCampaignService) Delete(id uint) error {
	row, err := s.campaignModule.GetByID(id)
	if err != nil {
		return err
	}
	if row.Status == model.SmsCampaignStatusSending {
		return errors.New("发送中的活动不可删除")
	}
	return s.campaignModule.Delete(id)
}

func (s *SmsCampaignService) Cancel(id uint) error {
	row, err := s.campaignModule.GetByID(id)
	if err != nil {
		return err
	}
	if row.Status != model.SmsCampaignStatusScheduled && row.Status != model.SmsCampaignStatusDraft {
		return errors.New("仅草稿或已排期的活动可取消")
	}
	return s.campaignModule.Update(id, map[string]interface{}{
		"status":       model.SmsCampaignStatusCancelled,
		"scheduled_at": nil,
	})
}

func (s *SmsCampaignService) SendNow(id uint, isAdmin bool, scopedStoreID uint) error {
	row, err := s.campaignModule.GetByID(id)
	if err != nil {
		return err
	}
	if row.Status == model.SmsCampaignStatusSending {
		return errors.New("活动正在发送中")
	}
	if row.Status == model.SmsCampaignStatusSent {
		return errors.New("活动已发送完成")
	}
	if row.Status == model.SmsCampaignStatusCancelled {
		return errors.New("活动已取消")
	}
	return s.executeCampaign(row, isAdmin, scopedStoreID)
}

func (s *SmsCampaignService) ProcessDueScheduled(now time.Time) error {
	rows, err := s.campaignModule.ListDueScheduled(now, 20)
	if err != nil {
		return err
	}
	for _, row := range rows {
		if err := s.executeCampaign(row, true, 0); err != nil {
			fmt.Printf("[SmsCampaign] 定时发送失败 id=%d: %v\n", row.ID, err)
		}
	}
	return nil
}

func (s *SmsCampaignService) executeCampaign(row *model.SmsCampaign, isAdmin bool, scopedStoreID uint) error {
	if s.smsClient == nil || !s.smsClient.Enabled() {
		return errors.New("阿里云短信未配置，请在环境变量中设置 ALIYUN_SMS_*")
	}
	recipients, err := s.resolveRecipients(row, isAdmin, scopedStoreID)
	if err != nil {
		return err
	}
	if len(recipients) == 0 {
		return errors.New("没有可发送的手机号")
	}

	_ = s.campaignModule.Update(row.ID, map[string]interface{}{
		"status":      model.SmsCampaignStatusSending,
		"total_count": len(recipients),
		"last_error":  "",
	})

	signName := row.SignName
	if signName == "" {
		signName = s.smsClient.DefaultSignName()
	}

	success := 0
	fail := 0
	var lastErr string

	if row.PersonalizeName {
		for _, r := range recipients {
			param, err := mergeTemplateParam(row.TemplateParam, map[string]string{"name": displayMemberName(r.Name)})
			if err != nil {
				fail++
				lastErr = err.Error()
				s.saveRecord(row.ID, r, "", model.SmsSendRecordFailed, lastErr)
				continue
			}
			bizID, err := s.smsClient.Send([]string{r.Phone}, signName, row.TemplateCode, param)
			if err != nil {
				fail++
				lastErr = err.Error()
				s.saveRecord(row.ID, r, "", model.SmsSendRecordFailed, err.Error())
				continue
			}
			success++
			s.saveRecord(row.ID, r, bizID, model.SmsSendRecordSuccess, "")
		}
	} else {
		param := row.TemplateParam
		if param == "" {
			param = "{}"
		}
		phones := make([]string, len(recipients))
		for i, r := range recipients {
			phones[i] = r.Phone
		}
		for i := 0; i < len(phones); i += aliyunsms.MaxPhonesPerRequest {
			end := i + aliyunsms.MaxPhonesPerRequest
			if end > len(phones) {
				end = len(phones)
			}
			chunk := phones[i:end]
			bizID, err := s.smsClient.Send(chunk, signName, row.TemplateCode, param)
			if err != nil {
				fail += len(chunk)
				lastErr = err.Error()
				for _, phone := range chunk {
					r := findRecipient(recipients, phone)
					s.saveRecord(row.ID, r, "", model.SmsSendRecordFailed, err.Error())
				}
				continue
			}
			success += len(chunk)
			for _, phone := range chunk {
				r := findRecipient(recipients, phone)
				s.saveRecord(row.ID, r, bizID, model.SmsSendRecordSuccess, "")
			}
		}
	}

	now := time.Now()
	status := model.SmsCampaignStatusSent
	if fail > 0 && success == 0 {
		status = model.SmsCampaignStatusFailed
	}
	return s.campaignModule.Update(row.ID, map[string]interface{}{
		"status":        status,
		"success_count": success,
		"fail_count":    fail,
		"last_error":    truncateErr(lastErr),
		"sent_at":       &now,
	})
}

type smsRecipient struct {
	Phone    string
	MemberID *uint
	Name     string
}

func (s *SmsCampaignService) resolveRecipients(row *model.SmsCampaign, isAdmin bool, scopedStoreID uint) ([]smsRecipient, error) {
	switch row.TargetType {
	case model.SmsCampaignTargetCustom:
		phones := normalizePhoneList(row.CustomPhones)
		out := make([]smsRecipient, 0, len(phones))
		for _, p := range phones {
			out = append(out, smsRecipient{Phone: p})
		}
		return out, nil
	case model.SmsCampaignTargetStores, model.SmsCampaignTargetAllMembers:
		storeIDs := []uint(row.StoreIDs)
		members, err := s.memberModule.ListSMSRecipients(storeIDs, scopedStoreID, isAdmin)
		if err != nil {
			return nil, err
		}
		out := make([]smsRecipient, 0, len(members))
		for _, m := range members {
			id := m.MemberID
			out = append(out, smsRecipient{Phone: m.Phone, MemberID: &id, Name: m.Name})
		}
		return out, nil
	default:
		return nil, fmt.Errorf("未知受众类型: %s", row.TargetType)
	}
}

func (s *SmsCampaignService) saveRecord(campaignID uint, r smsRecipient, bizID, status, errMsg string) {
	rec := &model.SmsSendRecord{
		CampaignID:   campaignID,
		Phone:        r.Phone,
		MemberID:     r.MemberID,
		BizID:        bizID,
		Status:       status,
		ErrorMessage: truncateErr(errMsg),
	}
	_ = s.campaignModule.CreateSendRecords([]*model.SmsSendRecord{rec})
}

func (s *SmsCampaignService) validateCampaign(row *model.SmsCampaign) error {
	if strings.TrimSpace(row.TemplateCode) == "" {
		return errors.New("请填写阿里云短信模板 CODE")
	}
	if row.TemplateParam != "" {
		var tmp map[string]interface{}
		if err := json.Unmarshal([]byte(row.TemplateParam), &tmp); err != nil {
			return fmt.Errorf("模板变量必须是合法 JSON: %w", err)
		}
	}
	switch row.TargetType {
	case model.SmsCampaignTargetCustom:
		if len(normalizePhoneList(row.CustomPhones)) == 0 {
			return errors.New("自定义受众请至少填写一个手机号")
		}
	case model.SmsCampaignTargetStores:
		if len(row.StoreIDs) == 0 {
			return errors.New("请选择至少一个门店")
		}
	case model.SmsCampaignTargetAllMembers:
	default:
		return errors.New("受众类型无效")
	}
	return nil
}

func normalizePhoneList(list model.StringList) model.StringList {
	if len(list) == 0 {
		return model.StringList{}
	}
	out := make(model.StringList, 0, len(list))
	seen := make(map[string]struct{})
	for _, raw := range list {
		for _, part := range strings.FieldsFunc(raw, func(r rune) bool {
			return r == ',' || r == ';' || r == '\n' || r == '\r' || r == ' '
		}) {
			p := strings.TrimSpace(part)
			if p == "" {
				continue
			}
			if _, ok := seen[p]; ok {
				continue
			}
			seen[p] = struct{}{}
			out = append(out, p)
		}
	}
	return out
}

func mergeTemplateParam(base string, extra map[string]string) (string, error) {
	m := map[string]interface{}{}
	if strings.TrimSpace(base) != "" {
		if err := json.Unmarshal([]byte(base), &m); err != nil {
			return "", err
		}
	}
	for k, v := range extra {
		m[k] = v
	}
	data, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func displayMemberName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "会员"
	}
	return name
}

func findRecipient(list []smsRecipient, phone string) smsRecipient {
	for _, r := range list {
		if r.Phone == phone {
			return r
		}
	}
	return smsRecipient{Phone: phone}
}

func truncateErr(msg string) string {
	msg = strings.TrimSpace(msg)
	if len(msg) <= 500 {
		return msg
	}
	return msg[:500]
}
