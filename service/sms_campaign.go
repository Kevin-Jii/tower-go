package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/Kevin-Jii/tower-go/config"
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
	"github.com/Kevin-Jii/tower-go/pkg/aliyunsms"
	"github.com/Kevin-Jii/tower-go/pkg/apicode"
)

const smsTimezone = "Asia/Shanghai"

type SmsCampaignService struct {
	campaignModule    *module.SmsCampaignModule
	memberModule      *module.MemberModule
	tagModule         *module.MemberTagModule
	storeModule       *module.StoreModule
	storeSmsConfigSvc *StoreSmsConfigService
	smsClient         *aliyunsms.Client // 全局兑底客户端；多门店时以 storeSmsConfigSvc.ResolveEffectiveContext 为准。
}

func NewSmsCampaignService(campaignModule *module.SmsCampaignModule, memberModule *module.MemberModule, tagModules ...*module.MemberTagModule) *SmsCampaignService {
	cfg := config.GetAliyunSMSConfig()
	client, err := aliyunsms.NewClient(aliyunsms.Config{AccessKeyID: cfg.AccessKeyID, AccessKeySecret: cfg.AccessKeySecret, RegionID: cfg.RegionID, SignName: cfg.SignName, Enabled: cfg.Enabled})
	if err != nil {
		fmt.Printf("[SmsCampaign] 初始化阿里云短信客户端失败: %v\n", err)
	}
	var tags *module.MemberTagModule
	if len(tagModules) > 0 {
		tags = tagModules[0]
	}
	return &SmsCampaignService{campaignModule: campaignModule, memberModule: memberModule, tagModule: tags, storeModule: nil, storeSmsConfigSvc: nil, smsClient: client}
}

// SetStoreModule 注入门店模块以读取每个门店的默认短信签名。
func (s *SmsCampaignService) SetStoreModule(m *module.StoreModule) { s.storeModule = m }

// SetStoreSmsConfigService 注入门店独立 SMS 配置服务。
func (s *SmsCampaignService) SetStoreSmsConfigService(svc *StoreSmsConfigService) {
	s.storeSmsConfigSvc = svc
}

func (s *SmsCampaignService) GetStoreSignName(storeID uint) string {
	if storeID == 0 || s.storeModule == nil {
		return ""
	}
	store, err := s.storeModule.GetByID(storeID)
	if err != nil || store == nil {
		return ""
	}
	return strings.TrimSpace(store.SmsSignName)
}

func (s *SmsCampaignService) GetConfig(storeID uint, fallbackSign string) model.SmsServiceConfigResp {
	start, end := smsSendWindow()
	region, configured, enabled, accessKeyID, signName := "", false, false, "", strings.TrimSpace(fallbackSign)
	if s.storeSmsConfigSvc != nil {
		ctx, _ := s.storeSmsConfigSvc.ResolveEffectiveContext(storeID)
		if ctx != nil {
			configured = ctx.AccessKeyID != "" && ctx.AccessKeySecret != ""
			enabled = configured && ctx.Enabled
			region = ctx.RegionID
			accessKeyID = ctx.AccessKeyID
			if signName == "" {
				signName = ctx.SignName
			}
			if ctx.WindowConfigured {
				start = ctx.SendWindowStart
				end = ctx.SendWindowEnd
			}
		}
	}
	if region == "" {
		region = config.GetAliyunSMSConfig().RegionID
	}
	if !configured {
		cfg := config.GetAliyunSMSConfig()
		if region == "" {
			region = cfg.RegionID
		}
		configured = cfg.AccessKeyID != "" && cfg.AccessKeySecret != ""
		if accessKeyID == "" {
			accessKeyID = cfg.AccessKeyID
		}
		if signName == "" {
			signName = cfg.SignName
		}
	}
	_ = accessKeyID
	return model.SmsServiceConfigResp{Enabled: enabled, DefaultSign: signName, Region: region, Configured: configured, HelpURL: "https://help.aliyun.com/zh/sms/getting-started/sms-skill-guide", TemplateHint: "请填写审核通过的阿里云 TemplateCode；变量使用 JSON。排期按中国标准时间执行。", MaxBatchPhones: aliyunsms.MaxPhonesPerRequest, Timezone: smsTimezone, SendWindowStart: start, SendWindowEnd: end, SendWindowEndExclusive: true, QualificationHint: "门店短信资质 ID（阿里云账号级）不会透传；调用 CreateSmsTemplate / SendSms 时阿里云会以你账号下已审核通过的资质为依据。"}
}

// resolveClient 返回活动归属门店的独立 aliyunsms.Client（创建时不依赖全局兑底）。
func (s *SmsCampaignService) resolveClient(storeID uint) (*aliyunsms.Client, error) {
	if s.storeSmsConfigSvc != nil {
		ctx, err := s.storeSmsConfigSvc.ResolveEffectiveContext(storeID)
		if err != nil {
			return nil, err
		}
		if ctx.Client != nil {
			return ctx.Client, nil
		}
	}
	if s.smsClient != nil && s.smsClient.Enabled() {
		return s.smsClient, nil
	}
	return nil, errors.New("未配置可用的阿里云短信凭证：请在会员推广 → 基础设置 中为本门店配置，或在 .env 中保留 ALIYUN_SMS_* 全局兑底")
}

// resolveWindow 返回当前活动应使用的发送窗口（分钟）。
func (s *SmsCampaignService) resolveWindow(storeID uint) (startMinute, endMinute int, windowConfigured bool) {
	if s.storeSmsConfigSvc != nil {
		ctx, _ := s.storeSmsConfigSvc.ResolveEffectiveContext(storeID)
		if ctx != nil && ctx.WindowConfigured {
			if s, e := parseClock(ctx.SendWindowStart); e == nil {
				startMinute = s
			}
			if e, e2 := parseClock(ctx.SendWindowEnd); e2 == nil {
				endMinute = e
			}
			windowConfigured = true
			return
		}
	}
	start, end := smsSendWindow()
	if s, e := parseClock(start); e == nil {
		startMinute = s
	}
	if e, e2 := parseClock(end); e2 == nil {
		endMinute = e
	}
	return
}

// resolveSign returns the SMS signature to use for a campaign's group execution in this order:
// 1. Segment-level sign_name (per-template override)
// 2. Campaign-level sign_name (legacy default)
// 3. Owner store's sms_sign_name (per-store default)
// 4. Global env default ALIYUN_SMS_SIGN_NAME
func (s *SmsCampaignService) resolveSign(row *model.SmsCampaign, groupSign string) string {
	if v := strings.TrimSpace(groupSign); v != "" {
		return v
	}
	if v := strings.TrimSpace(row.SignName); v != "" {
		return v
	}
	if row.OwnerStoreID > 0 && s.storeSmsConfigSvc != nil {
		ctx, _ := s.storeSmsConfigSvc.ResolveEffectiveContext(row.OwnerStoreID)
		if ctx != nil && strings.TrimSpace(ctx.SignName) != "" {
			return strings.TrimSpace(ctx.SignName)
		}
	}
	if row.OwnerStoreID > 0 {
		if store, err := s.storeModule.GetByID(row.OwnerStoreID); err == nil && store != nil && strings.TrimSpace(store.SmsSignName) != "" {
			return strings.TrimSpace(store.SmsSignName)
		}
	}
	if s.smsClient != nil {
		return strings.TrimSpace(s.smsClient.DefaultSignName())
	}
	return ""
}

func (s *SmsCampaignService) List(storeID uint, allStores bool) ([]*model.SmsCampaign, error) {
	return s.campaignModule.List(storeID, allStores)
}
func (s *SmsCampaignService) GetByID(id, storeID uint, allStores bool) (*model.SmsCampaign, error) {
	return s.campaignModule.GetByID(id, storeID, allStores)
}
func (s *SmsCampaignService) ListSendRecords(campaignID, storeID uint, allStores bool) ([]*model.SmsSendRecord, error) {
	if _, err := s.campaignModule.GetByID(campaignID, storeID, allStores); err != nil {
		return nil, err
	}
	return s.campaignModule.ListSendRecords(campaignID, 500)
}

func (s *SmsCampaignService) RetryFailedRecord(campaignID, recordID, storeID uint, allStores bool) error {
	row, err := s.campaignModule.GetByID(campaignID, storeID, allStores)
	if err != nil {
		return err
	}
	record, err := s.campaignModule.GetSendRecord(campaignID, recordID)
	if err != nil {
		return errors.New("发送记录不存在")
	}
	if record.Status != model.SmsSendRecordFailed {
		return errors.New("只有发送失败的记录可以重新发送")
	}
	if !s.inSendWindow(row.OwnerStoreID, time.Now()) {
		start, end := s.windowFor(row.OwnerStoreID)
		return fmt.Errorf("当前不在短信发送时段（中国时间 %s–%s，结束时间不含）", start, end)
	}
	client, err := s.resolveClient(row.OwnerStoreID)
	if err != nil {
		return err
	}

	recipient, err := s.retryRecipient(row, record)
	if err != nil {
		return err
	}
	sign := s.resolveSign(row, recipient.Config.SignName)
	if sign == "" {
		return errors.New("短信签名为空，请在活动分组、所属门店或基础设置中至少配置一个")
	}
	param := recipient.Config.TemplateParam
	if strings.TrimSpace(param) == "" {
		param = "{}"
	}
	if recipient.Config.PersonalizeName {
		param, err = mergeTemplateParam(param, map[string]string{"name": displayMemberName(recipient.Name)})
		if err != nil {
			return err
		}
	}

	claimed, err := s.campaignModule.ClaimFailedSendRecord(campaignID, recordID)
	if err != nil {
		return err
	}
	if !claimed {
		return errors.New("该失败记录正在重发或已重发，请刷新发送记录")
	}
	bizID, sendErr := client.Send([]string{record.Phone}, sign, recipient.Config.TemplateCode, param)
	if sendErr != nil {
		message := truncateErr(sendErr.Error())
		if finishErr := s.campaignModule.FinishSendRecord(campaignID, recordID, model.SmsSendRecordFailed, "", message); finishErr != nil {
			return fmt.Errorf("短信重发失败：%s；保存失败日志时发生错误：%v", message, finishErr)
		}
		return fmt.Errorf("短信重发失败：%s", message)
	}
	if err := s.campaignModule.MarkRetrySuccess(campaignID, recordID, bizID); err != nil {
		return fmt.Errorf("短信已提交至阿里云，但更新发送记录失败：%w", err)
	}
	return nil
}

func (s *SmsCampaignService) retryRecipient(row *model.SmsCampaign, record *model.SmsSendRecord) (smsRecipient, error) {
	config := smsTemplate{SignName: row.SignName, TemplateCode: row.TemplateCode, TemplateParam: row.TemplateParam, PersonalizeName: row.PersonalizeName}
	if record.SegmentID != nil {
		found := false
		for i := range row.Segments {
			segment := &row.Segments[i]
			if segment.ID != *record.SegmentID {
				continue
			}
			id := segment.ID
			config = smsTemplate{SegmentID: &id, SignName: segment.SignName, TemplateCode: segment.TemplateCode, TemplateParam: segment.TemplateParam, PersonalizeName: segment.PersonalizeName}
			found = true
			break
		}
		if !found {
			return smsRecipient{}, errors.New("原发送分组不存在，无法重新发送")
		}
	}
	if strings.TrimSpace(config.TemplateCode) == "" {
		return smsRecipient{}, errors.New("短信模板 CODE 为空，无法重新发送")
	}
	recipient := smsRecipient{Phone: record.Phone, MemberID: record.MemberID, Config: config}
	if config.PersonalizeName && record.MemberID != nil {
		member, err := s.memberModule.GetMember(*record.MemberID, row.OwnerStoreID, row.OwnerStoreID == 0)
		if err != nil {
			return smsRecipient{}, errors.New("无法读取会员姓名，无法重新发送")
		}
		recipient.Name = member.Name
	}
	return recipient, nil
}

func (s *SmsCampaignService) Create(req *model.CreateSmsCampaignReq, createdBy, effectiveStoreID uint, hqUnbound bool) (*model.SmsCampaign, error) {
	owner := effectiveStoreID
	if hqUnbound && owner == 0 {
		owner = req.OwnerStoreID
	}
	row := &model.SmsCampaign{OwnerStoreID: owner, Name: strings.TrimSpace(req.Name), CampaignType: req.CampaignType, SignName: strings.TrimSpace(req.SignName), TemplateCode: strings.TrimSpace(req.TemplateCode), TemplateParam: strings.TrimSpace(req.TemplateParam), PersonalizeName: req.PersonalizeName != nil && *req.PersonalizeName, TargetType: req.TargetType, StoreIDs: req.StoreIDs, CustomPhones: normalizePhoneList(req.CustomPhones), Status: model.SmsCampaignStatusDraft, CreatedBy: createdBy}
	row.Segments = buildSegments(0, req.Segments)
	if req.ScheduledAt != nil {
		normalized, err := s.normalizeAndValidateSchedule(*req.ScheduledAt, time.Now())
		if err != nil {
			return nil, err
		}
		row.ScheduledAt = &normalized
		row.Status = model.SmsCampaignStatusScheduled
	}
	if err := s.validateCampaign(row, hqUnbound && owner == 0); err != nil {
		return nil, err
	}
	if err := s.campaignModule.Create(row); err != nil {
		return nil, err
	}
	return s.campaignModule.GetUnscopedByID(row.ID)
}

func (s *SmsCampaignService) Update(id uint, req *model.UpdateSmsCampaignReq, storeID uint, allStores bool) error {
	row, err := s.campaignModule.GetByID(id, storeID, allStores)
	if err != nil {
		return err
	}
	if row.Status == model.SmsCampaignStatusSending || row.Status == model.SmsCampaignStatusSent {
		return errors.New("已发送或发送中的活动不可编辑")
	}
	updates := map[string]interface{}{}
	if req.Name != nil {
		row.Name = strings.TrimSpace(*req.Name)
		updates["name"] = row.Name
	}
	if req.CampaignType != nil {
		row.CampaignType = *req.CampaignType
		updates["campaign_type"] = row.CampaignType
	}
	if req.SignName != nil {
		row.SignName = strings.TrimSpace(*req.SignName)
		updates["sign_name"] = row.SignName
	}
	if req.TemplateCode != nil {
		row.TemplateCode = strings.TrimSpace(*req.TemplateCode)
		updates["template_code"] = row.TemplateCode
	}
	if req.TemplateParam != nil {
		row.TemplateParam = strings.TrimSpace(*req.TemplateParam)
		updates["template_param"] = row.TemplateParam
	}
	if req.PersonalizeName != nil {
		row.PersonalizeName = *req.PersonalizeName
		updates["personalize_name"] = row.PersonalizeName
	}
	if req.TargetType != nil {
		row.TargetType = *req.TargetType
		updates["target_type"] = row.TargetType
	}
	if req.StoreIDs != nil {
		row.StoreIDs = *req.StoreIDs
		updates["store_ids"] = row.StoreIDs
	}
	if req.CustomPhones != nil {
		row.CustomPhones = normalizePhoneList(*req.CustomPhones)
		updates["custom_phones"] = row.CustomPhones
	}
	if req.ScheduledAt != nil {
		normalized, e := s.normalizeAndValidateSchedule(*req.ScheduledAt, time.Now())
		if e != nil {
			return e
		}
		row.ScheduledAt = &normalized
		row.Status = model.SmsCampaignStatusScheduled
		updates["scheduled_at"] = normalized
		updates["status"] = row.Status
	} else if req.ClearScheduledAt {
		row.ScheduledAt = nil
		row.Status = model.SmsCampaignStatusDraft
		updates["scheduled_at"] = nil
		updates["status"] = row.Status
	}
	var segments *[]model.SmsCampaignSegment
	if req.Segments != nil {
		built := buildSegments(id, *req.Segments)
		row.Segments = built
		segments = &built
	}
	if err := s.validateCampaign(row, allStores && row.OwnerStoreID == 0); err != nil {
		return err
	}
	updated, err := s.campaignModule.UpdateEditable(id, updates, segments)
	if err != nil {
		return err
	}
	if !updated {
		return errors.New("活动已开始发送，不可编辑")
	}
	return nil
}

func (s *SmsCampaignService) Delete(id, storeID uint, allStores bool) error {
	row, err := s.campaignModule.GetByID(id, storeID, allStores)
	if err != nil {
		return err
	}
	if row.Status == model.SmsCampaignStatusSending {
		return errors.New("发送中的活动不可删除")
	}
	deleted, err := s.campaignModule.DeleteEditable(id)
	if err != nil {
		return err
	}
	if !deleted {
		return errors.New("活动已开始发送，不可删除")
	}
	return nil
}
func (s *SmsCampaignService) Cancel(id, storeID uint, allStores bool) error {
	row, err := s.campaignModule.GetByID(id, storeID, allStores)
	if err != nil {
		return err
	}
	if row.Status != model.SmsCampaignStatusScheduled && row.Status != model.SmsCampaignStatusDraft {
		return errors.New("仅草稿或已排期的活动可取消")
	}
	cancelled, err := s.campaignModule.CancelEditable(id)
	if err != nil {
		return err
	}
	if !cancelled {
		return errors.New("活动已开始发送，不可取消")
	}
	return nil
}

func (s *SmsCampaignService) SendNow(id, storeID uint, allStores bool) error {
	row, err := s.campaignModule.GetByID(id, storeID, allStores)
	if err != nil {
		return err
	}
	if err := s.ensureSendable(row); err != nil {
		return apicode.Newf(apicode.InvalidParameter, "%s", err.Error())
	}
	if !s.inSendWindow(row.OwnerStoreID, time.Now()) {
		start, end := s.windowFor(row.OwnerStoreID)
		return apicode.Newf(apicode.InvalidParameter, "当前不在短信发送时段（中国时间 %s–%s，结束时间不含）", start, end)
	}
	err = s.executeCampaign(row)
	if err == nil {
		return nil
	}
	if _, recognized := apicode.Resolve(err); recognized {
		return err
	}
	return apicode.Newf(apicode.InvalidParameter, "%s", err.Error())
}

func (s *SmsCampaignService) ProcessDueScheduled(now time.Time) error {
	rows, err := s.campaignModule.ListDueScheduled(now, 20)
	if err != nil {
		return err
	}
	for _, row := range rows {
		if !s.inSendWindow(row.OwnerStoreID, now) {
			continue
		}
		if err := s.executeCampaign(row); err != nil {
			fmt.Printf("[SmsCampaign] 定时发送失败 id=%d owner_store_id=%d: %v\n", row.ID, row.OwnerStoreID, err)
		}
	}
	return nil
}

func (s *SmsCampaignService) windowFor(storeID uint) (string, string) {
	start, end, custom := s.resolveWindow(storeID)
	if !custom {
		return smsSendWindow()
	}
	return formatClock(start), formatClock(end)
}

func formatClock(total int) string {
	h := total / 60
	m := total % 60
	return fmt.Sprintf("%02d:%02d", h, m)
}

func (s *SmsCampaignService) ensureSendable(row *model.SmsCampaign) error {
	client, err := s.resolveClient(row.OwnerStoreID)
	if err != nil {
		return err
	}
	if !client.Enabled() {
		return errors.New("阿里云短信未配置，请在会员推广 → 基础设置 中配置门店 SMS 凭证")
	}
	switch row.Status {
	case model.SmsCampaignStatusDraft, model.SmsCampaignStatusScheduled:
		return nil
	case model.SmsCampaignStatusSending:
		return errors.New("活动正在发送中")
	case model.SmsCampaignStatusSent:
		return errors.New("活动已发送完成")
	case model.SmsCampaignStatusCancelled:
		return errors.New("活动已取消")
	case model.SmsCampaignStatusFailed:
		if message := strings.TrimSpace(row.LastError); message != "" {
			return fmt.Errorf("活动上次发送失败：%s；请在发送记录中对失败号码重新发送", message)
		}
		return errors.New("活动上次发送失败，请在发送记录中对失败号码重新发送")
	default:
		return errors.New("活动状态不可发送")
	}
}

func (s *SmsCampaignService) executeCampaign(row *model.SmsCampaign) error {
	if err := s.ensureSendable(row); err != nil {
		return err
	}
	recipients, err := s.resolveRecipients(row)
	if err != nil {
		return err
	}
	if len(recipients) == 0 {
		return errors.New("没有可发送的手机号")
	}
	claimed, err := s.campaignModule.Claim(row.ID)
	if err != nil {
		return err
	}
	if !claimed {
		return errors.New("活动已被其他发送任务处理")
	}
	if err := s.campaignModule.Update(row.ID, map[string]interface{}{"total_count": len(recipients)}, nil); err != nil {
		return err
	}
	client, err := s.resolveClient(row.OwnerStoreID)
	if err != nil {
		return err
	}
	success, fail, lastErr := 0, 0, ""
	groups := groupRecipients(recipients)
	for _, group := range groups {
		cfg := group.Config
		sign := s.resolveSign(row, cfg.SignName)
		if sign == "" {
			return errors.New("短信签名为空，请在活动分组、所属门店或 基础设置 中至少配置一个")
		}
		if cfg.PersonalizeName {
			for _, r := range group.Recipients {
				param, e := mergeTemplateParam(cfg.TemplateParam, map[string]string{"name": displayMemberName(r.Name)})
				if e == nil {
					var biz string
					biz, e = client.Send([]string{r.Phone}, sign, cfg.TemplateCode, param)
					if e == nil {
						success++
						s.saveRecord(row.ID, r, biz, model.SmsSendRecordSuccess, "")
						continue
					}
				}
				fail++
				lastErr = e.Error()
				s.saveRecord(row.ID, r, "", model.SmsSendRecordFailed, lastErr)
			}
			continue
		}
		param := cfg.TemplateParam
		if strings.TrimSpace(param) == "" {
			param = "{}"
		}
		for i := 0; i < len(group.Recipients); i += aliyunsms.MaxPhonesPerRequest {
			end := i + aliyunsms.MaxPhonesPerRequest
			if end > len(group.Recipients) {
				end = len(group.Recipients)
			}
			chunk := group.Recipients[i:end]
			phones := make([]string, len(chunk))
			for j := range chunk {
				phones[j] = chunk[j].Phone
			}
			biz, e := client.Send(phones, sign, cfg.TemplateCode, param)
			if e != nil {
				fail += len(chunk)
				lastErr = e.Error()
				for _, r := range chunk {
					s.saveRecord(row.ID, r, "", model.SmsSendRecordFailed, lastErr)
				}
			} else {
				success += len(chunk)
				for _, r := range chunk {
					s.saveRecord(row.ID, r, biz, model.SmsSendRecordSuccess, "")
				}
			}
		}
	}
	now := time.Now().UTC()
	status := model.SmsCampaignStatusSent
	if fail > 0 && success == 0 {
		status = model.SmsCampaignStatusFailed
	}
	if err := s.campaignModule.Update(row.ID, map[string]interface{}{"status": status, "success_count": success, "fail_count": fail, "last_error": truncateErr(lastErr), "sent_at": &now}, nil); err != nil {
		return err
	}
	if fail > 0 {
		message := strings.TrimSpace(lastErr)
		if message == "" {
			message = "阿里云未返回具体错误信息"
		}
		return apicode.Newf(apicode.InvalidParameter, "短信发送完成：成功 %d 条，失败 %d 条；阿里云错误：%s", success, fail, message)
	}
	return nil
}

type smsTemplate struct {
	SegmentID                             *uint
	SignName, TemplateCode, TemplateParam string
	PersonalizeName                       bool
}
type smsRecipient struct {
	Phone    string
	MemberID *uint
	Name     string
	Config   smsTemplate
}
type recipientGroup struct {
	Key        string
	Config     smsTemplate
	Recipients []smsRecipient
}

func (s *SmsCampaignService) resolveRecipients(row *model.SmsCampaign) ([]smsRecipient, error) {
	legacy := smsTemplate{SignName: row.SignName, TemplateCode: row.TemplateCode, TemplateParam: row.TemplateParam, PersonalizeName: row.PersonalizeName}
	if row.TargetType == model.SmsCampaignTargetCustom {
		out := make([]smsRecipient, 0, len(row.CustomPhones))
		for _, p := range normalizePhoneList(row.CustomPhones) {
			r := smsRecipient{Phone: p}
			r.Config = chooseTemplate(r, nil, row.Segments, legacy)
			out = append(out, r)
		}
		return out, nil
	}
	storeIDs := make([]uint, 0)
	if row.OwnerStoreID > 0 {
		storeIDs = []uint{row.OwnerStoreID}
	} else if row.TargetType == model.SmsCampaignTargetStores {
		storeIDs = []uint(row.StoreIDs)
	}
	members, err := s.memberModule.ListSMSRecipients(storeIDs, row.OwnerStoreID, row.OwnerStoreID == 0)
	if err != nil {
		return nil, err
	}
	ids := make([]uint, len(members))
	for i := range members {
		ids[i] = members[i].MemberID
	}
	tags := map[uint]map[uint]struct{}{}
	if s.tagModule != nil {
		tags, err = s.tagModule.MemberTagMap(ids, row.OwnerStoreID)
		if err != nil {
			return nil, err
		}
	}
	out := make([]smsRecipient, 0, len(members))
	seen := map[string]struct{}{}
	for _, m := range members {
		if _, ok := seen[m.Phone]; ok {
			continue
		}
		seen[m.Phone] = struct{}{}
		id := m.MemberID
		r := smsRecipient{Phone: m.Phone, MemberID: &id, Name: m.Name}
		r.Config = chooseTemplate(r, tags[id], row.Segments, legacy)
		out = append(out, r)
	}
	return out, nil
}

func chooseTemplate(_ smsRecipient, memberTags map[uint]struct{}, segments []model.SmsCampaignSegment, legacy smsTemplate) smsTemplate {
	var def *model.SmsCampaignSegment
	for i := range segments {
		seg := &segments[i]
		if seg.IsDefault {
			if def == nil {
				def = seg
			}
			continue
		}
		for _, tagID := range seg.TagIDs {
			if _, ok := memberTags[tagID]; ok {
				id := seg.ID
				return smsTemplate{SegmentID: &id, SignName: seg.SignName, TemplateCode: seg.TemplateCode, TemplateParam: seg.TemplateParam, PersonalizeName: seg.PersonalizeName}
			}
		}
	}
	if def != nil {
		id := def.ID
		return smsTemplate{SegmentID: &id, SignName: def.SignName, TemplateCode: def.TemplateCode, TemplateParam: def.TemplateParam, PersonalizeName: def.PersonalizeName}
	}
	return legacy
}

func groupRecipients(in []smsRecipient) []recipientGroup {
	byKey := map[string]int{}
	out := make([]recipientGroup, 0)
	for _, r := range in {
		segmentID := uint(0)
		if r.Config.SegmentID != nil {
			segmentID = *r.Config.SegmentID
		}
		key := fmt.Sprintf("%d|%s|%s|%s|%t", segmentID, r.Config.SignName, r.Config.TemplateCode, r.Config.TemplateParam, r.Config.PersonalizeName)
		i, ok := byKey[key]
		if !ok {
			i = len(out)
			byKey[key] = i
			out = append(out, recipientGroup{Key: key, Config: r.Config})
		}
		out[i].Recipients = append(out[i].Recipients, r)
	}
	return out
}

func (s *SmsCampaignService) saveRecord(campaignID uint, r smsRecipient, bizID, status, errMsg string) {
	rec := &model.SmsSendRecord{CampaignID: campaignID, SegmentID: r.Config.SegmentID, Phone: r.Phone, MemberID: r.MemberID, BizID: bizID, Status: status, ErrorMessage: truncateErr(errMsg)}
	_ = s.campaignModule.CreateSendRecords([]*model.SmsSendRecord{rec})
}

func (s *SmsCampaignService) validateCampaign(row *model.SmsCampaign, globalHQ bool) error {
	if strings.TrimSpace(row.Name) == "" {
		return errors.New("请填写活动名称")
	}
	if len(row.Segments) == 0 && strings.TrimSpace(row.TemplateCode) == "" {
		return errors.New("请填写阿里云短信模板 CODE")
	}
	if err := validateTemplateParam(row.TemplateParam); err != nil {
		return err
	}
	defaults := 0
	for i := range row.Segments {
		seg := &row.Segments[i]
		if strings.TrimSpace(seg.TemplateCode) == "" {
			return fmt.Errorf("第 %d 个分组缺少模板 CODE", i+1)
		}
		if err := validateTemplateParam(seg.TemplateParam); err != nil {
			return fmt.Errorf("第 %d 个分组: %w", i+1, err)
		}
		if seg.IsDefault {
			defaults++
			if len(seg.TagIDs) > 0 {
				return errors.New("默认分组不能指定标签")
			}
		} else if len(seg.TagIDs) == 0 {
			return errors.New("非默认分组必须选择标签")
		}
		if s.tagModule == nil {
			return errors.New("会员标签服务未初始化")
		}
		if !globalHQ && row.OwnerStoreID == 0 {
			return errors.New("请选择活动所属门店")
		}
		if err := s.tagModule.ValidateTagIDsAccess(seg.TagIDs, row.OwnerStoreID, globalHQ); err != nil {
			return err
		}
	}
	if defaults > 1 {
		return errors.New("只能配置一个默认分组")
	}
	if len(row.Segments) > 0 && defaults == 0 && strings.TrimSpace(row.TemplateCode) == "" {
		return errors.New("未配置活动级模板时，分组活动必须配置一个默认分组")
	}
	switch row.TargetType {
	case model.SmsCampaignTargetCustom:
		if len(normalizePhoneList(row.CustomPhones)) == 0 {
			return errors.New("自定义受众请至少填写一个手机号")
		}
	case model.SmsCampaignTargetStores:
		if len(row.StoreIDs) == 0 && row.OwnerStoreID == 0 {
			return errors.New("请选择至少一个门店")
		}
		if row.OwnerStoreID > 0 {
			for _, id := range row.StoreIDs {
				if id != row.OwnerStoreID {
					return errors.New("活动只能选择所属门店的会员")
				}
			}
		}
	case model.SmsCampaignTargetAllMembers:
	default:
		return errors.New("受众类型无效")
	}
	return nil
}

func buildSegments(campaignID uint, reqs []model.SmsCampaignSegmentReq) []model.SmsCampaignSegment {
	out := make([]model.SmsCampaignSegment, 0, len(reqs))
	for i, req := range reqs {
		out = append(out, model.SmsCampaignSegment{CampaignID: campaignID, Position: i, TagIDs: model.UintList(uniqueIDs(req.TagIDs)), SignName: strings.TrimSpace(req.SignName), TemplateCode: strings.TrimSpace(req.TemplateCode), TemplateParam: strings.TrimSpace(req.TemplateParam), PersonalizeName: req.PersonalizeName, IsDefault: req.IsDefault})
	}
	return out
}
func uniqueIDs(ids []uint) []uint {
	seen := map[uint]struct{}{}
	out := make([]uint, 0, len(ids))
	for _, id := range ids {
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}
func validateTemplateParam(raw string) error {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	var tmp map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &tmp); err != nil {
		return fmt.Errorf("模板变量必须是合法 JSON 对象: %w", err)
	}
	if tmp == nil {
		return errors.New("模板变量必须是 JSON 对象，不能是 null")
	}
	return nil
}

func (s *SmsCampaignService) normalizeAndValidateSchedule(value, now time.Time) (time.Time, error) {
	loc, err := time.LoadLocation(smsTimezone)
	if err != nil {
		return time.Time{}, err
	}
	china := value.In(loc)
	if !china.After(now.In(loc)) {
		return time.Time{}, errors.New("计划发送时间必须晚于当前时间")
	}
	if !s.inSendWindow(0, china) {
		start, end := smsSendWindow()
		return time.Time{}, fmt.Errorf("计划发送时间必须在中国时间 %s–%s 之间（结束时间不含）", start, end)
	}
	return china.UTC(), nil
}
func (s *SmsCampaignService) inSendWindow(storeID uint, value time.Time) bool {
	loc, err := time.LoadLocation(smsTimezone)
	if err != nil {
		return false
	}
	start, end, custom := s.resolveWindow(storeID)
	if !custom {
		start, end = 8*60, 22*60
	}
	local := value.In(loc)
	minute := local.Hour()*60 + local.Minute()
	return minute >= start && minute < end
}
func smsSendWindow() (string, string) {
	start := strings.TrimSpace(os.Getenv("ALIYUN_SMS_SEND_WINDOW_START"))
	end := strings.TrimSpace(os.Getenv("ALIYUN_SMS_SEND_WINDOW_END"))
	if start == "" {
		start = "08:00"
	}
	if end == "" {
		end = "22:00"
	}
	startMinute, startErr := parseClock(start)
	endMinute, endErr := parseClock(end)
	if startErr != nil || endErr != nil || startMinute >= endMinute {
		return "08:00", "22:00"
	}
	return start, end
}

func parseClock(v string) (int, error) {
	var h, m int
	if _, err := fmt.Sscanf(strings.TrimSpace(v), "%d:%d", &h, &m); err != nil || h < 0 || h > 23 || m < 0 || m > 59 {
		return 0, errors.New("invalid clock")
	}
	return h*60 + m, nil
}

func normalizePhoneList(list model.StringList) model.StringList {
	out := make(model.StringList, 0, len(list))
	seen := map[string]struct{}{}
	for _, raw := range list {
		for _, part := range strings.FieldsFunc(raw, func(r rune) bool { return r == ',' || r == ';' || r == '\n' || r == '\r' || r == ' ' }) {
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
		if m == nil {
			m = map[string]interface{}{}
		}
	}
	for k, v := range extra {
		m[k] = v
	}
	data, err := json.Marshal(m)
	return string(data), err
}
func displayMemberName(name string) string {
	if name = strings.TrimSpace(name); name == "" {
		return "会员"
	}
	return name
}
func truncateErr(msg string) string {
	msg = strings.TrimSpace(msg)
	if len(msg) <= 500 {
		return msg
	}
	return msg[:500]
}
