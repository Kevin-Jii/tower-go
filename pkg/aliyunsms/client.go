package aliyunsms

import (
	"fmt"
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
