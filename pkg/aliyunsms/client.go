package aliyunsms

import (
	"fmt"
	"strconv"
	"strings"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v4/client"
)

const MaxPhonesPerRequest = 100

// Client 阿里云短信客户端
type Client struct {
	inner    *dysmsapi.Client
	signName string
	enabled  bool
}

// Config 客户端配置
type Config struct {
	AccessKeyID     string
	AccessKeySecret string
	RegionID        string
	SignName        string
	Enabled         bool
}

func NewClient(cfg Config) (*Client, error) {
	region := cfg.RegionID
	if region == "" {
		region = "cn-hangzhou"
	}
	c := &Client{
		signName: strings.TrimSpace(cfg.SignName),
		enabled:  cfg.Enabled && cfg.AccessKeyID != "" && cfg.AccessKeySecret != "",
	}
	if !c.enabled {
		return c, nil
	}
	apiCfg := &openapi.Config{
		AccessKeyId:     &cfg.AccessKeyID,
		AccessKeySecret: &cfg.AccessKeySecret,
		RegionId:        &region,
	}
	apiCfg.Endpoint = teaString("dysmsapi.aliyuncs.com")
	inner, err := dysmsapi.NewClient(apiCfg)
	if err != nil {
		return nil, fmt.Errorf("init dysms client: %w", err)
	}
	c.inner = inner
	return c, nil
}

func (c *Client) Enabled() bool {
	return c != nil && c.enabled && c.inner != nil
}

func (c *Client) DefaultSignName() string {
	if c == nil {
		return ""
	}
	return c.signName
}

// Send 向单个或多个手机号发送相同模板参数（最多 100 个号码，逗号分隔）
func (c *Client) Send(phones []string, signName, templateCode, templateParam string) (bizID string, err error) {
	if !c.Enabled() {
		return "", fmt.Errorf("阿里云短信未配置或未启用，请设置 ALIYUN_SMS_ACCESS_KEY_ID、ALIYUN_SMS_ACCESS_KEY_SECRET 与 ALIYUN_SMS_ENABLED=true")
	}
	if len(phones) == 0 {
		return "", fmt.Errorf("手机号列表为空")
	}
	if len(phones) > MaxPhonesPerRequest {
		return "", fmt.Errorf("单次最多发送 %d 个号码", MaxPhonesPerRequest)
	}
	sign := strings.TrimSpace(signName)
	if sign == "" {
		sign = c.signName
	}
	if sign == "" {
		return "", fmt.Errorf("短信签名为空，请配置 ALIYUN_SMS_SIGN_NAME 或在活动中填写签名")
	}
	if strings.TrimSpace(templateCode) == "" {
		return "", fmt.Errorf("模板 CODE 不能为空")
	}

	phoneStr := strings.Join(normalizePhones(phones), ",")
	req := &dysmsapi.SendSmsRequest{
		PhoneNumbers:  teaString(phoneStr),
		SignName:      teaString(sign),
		TemplateCode:  teaString(templateCode),
		TemplateParam: teaString(templateParam),
	}
	resp, err := c.inner.SendSms(req)
	if err != nil {
		return "", err
	}
	if resp == nil || resp.Body == nil {
		return "", fmt.Errorf("短信接口无响应")
	}
	code := teaStringValue(resp.Body.Code)
	if code != "OK" {
		msg := teaStringValue(resp.Body.Message)
		return "", fmt.Errorf("短信发送失败: %s (%s)", msg, code)
	}
	return teaStringValue(resp.Body.BizId), nil
}

// TemplateSnapshot describes an Aliyun SMS template as returned by GetSmsTemplate.
type TemplateSnapshot struct {
	TemplateCode    string
	TemplateName    string
	TemplateContent string
	TemplateType    int32
	TemplateStatus  string // 0=pending, 1=approved, 2=rejected, 10=cancelled
	Reason          string
	CreateDate      string
}

// CreateTemplate submits a new template to Aliyun and returns its TemplateCode.
func (c *Client) CreateTemplate(name, content, relatedSign, remark string, templateType int32) (string, error) {
	if !c.Enabled() {
		return "", fmt.Errorf("阿里云短信未配置或未启用，请设置 ALIYUN_SMS_ACCESS_KEY_ID、ALIYUN_SMS_ACCESS_KEY_SECRET 与 ALIYUN_SMS_ENABLED=true")
	}
	req := &dysmsapi.CreateSmsTemplateRequest{
		TemplateName:    teaString(name),
		TemplateContent: teaString(content),
		TemplateType:    teaInt32(templateType),
	}
	if relatedSign = strings.TrimSpace(relatedSign); relatedSign != "" {
		req.RelatedSignName = teaString(relatedSign)
	}
	if remark = strings.TrimSpace(remark); remark != "" {
		req.Remark = teaString(remark)
	}
	resp, err := c.inner.CreateSmsTemplate(req)
	if err != nil {
		return "", fmt.Errorf("提交模板失败: %w", err)
	}
	if resp == nil || resp.Body == nil {
		return "", fmt.Errorf("创建模板接口无响应")
	}
	if code := teaStringValue(resp.Body.Code); code != "OK" {
		return "", fmt.Errorf("创建模板失败: %s (%s)", teaStringValue(resp.Body.Message), code)
	}
	return teaStringValue(resp.Body.TemplateCode), nil
}

// GetTemplate queries Aliyun for the current audit status of a template.
func (c *Client) GetTemplate(templateCode string) (*TemplateSnapshot, error) {
	if !c.Enabled() {
		return nil, fmt.Errorf("阿里云短信未配置或未启用")
	}
	resp, err := c.inner.GetSmsTemplate(&dysmsapi.GetSmsTemplateRequest{TemplateCode: teaString(templateCode)})
	if err != nil {
		return nil, fmt.Errorf("查询模板失败: %w", err)
	}
	if resp == nil || resp.Body == nil {
		return nil, fmt.Errorf("查询模板接口无响应")
	}
	if code := teaStringValue(resp.Body.Code); code != "OK" {
		return nil, fmt.Errorf("查询模板失败: %s (%s)", teaStringValue(resp.Body.Message), code)
	}
	snap := &TemplateSnapshot{
		TemplateCode:    teaStringValue(resp.Body.TemplateCode),
		TemplateName:    teaStringValue(resp.Body.TemplateName),
		TemplateContent: teaStringValue(resp.Body.TemplateContent),
		TemplateStatus:  teaStringValue(resp.Body.TemplateStatus),
		CreateDate:      teaStringValue(resp.Body.CreateDate),
	}
	if resp.Body.TemplateType != nil {
		if v, err := strconv.Atoi(*resp.Body.TemplateType); err == nil {
			snap.TemplateType = int32(v)
		}
	}
	if resp.Body.AuditInfo != nil && resp.Body.AuditInfo.RejectInfo != nil {
		snap.Reason = teaStringValue(resp.Body.AuditInfo.RejectInfo)
	}
	if snap.TemplateCode == "" {
		snap.TemplateCode = templateCode
	}
	return snap, nil
}

// DeleteTemplate removes a template on Aliyun. Approved templates are typically not deletable
// until the audit period allows; caller should surface the error if it occurs.
func (c *Client) DeleteTemplate(templateCode string) error {
	if !c.Enabled() {
		return fmt.Errorf("阿里云短信未配置或未启用")
	}
	resp, err := c.inner.DeleteSmsTemplate(&dysmsapi.DeleteSmsTemplateRequest{TemplateCode: teaString(templateCode)})
	if err != nil {
		return fmt.Errorf("删除模板失败: %w", err)
	}
	if resp == nil || resp.Body == nil {
		return fmt.Errorf("删除模板接口无响应")
	}
	if code := teaStringValue(resp.Body.Code); code != "OK" {
		return fmt.Errorf("删除模板失败: %s (%s)", teaStringValue(resp.Body.Message), code)
	}
	return nil
}

func teaInt32(v int32) *int32 { return &v }

func normalizePhones(phones []string) []string {
	out := make([]string, 0, len(phones))
	seen := make(map[string]struct{}, len(phones))
	for _, p := range phones {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if _, ok := seen[p]; ok {
			continue
		}
		seen[p] = struct{}{}
		out = append(out, p)
	}
	return out
}

func teaString(s string) *string {
	return &s
}

func teaStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
