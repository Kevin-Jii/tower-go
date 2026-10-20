package model

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestThirdPartyAccountRespDoesNotExposeCredentials(t *testing.T) {
	row := &ThirdPartyAccount{
		Password:       "password-secret",
		ApplicationKey: "application-secret",
		LastToken:      "access-token-secret",
	}
	body, err := json.Marshal(NewThirdPartyAccountResp(row))
	if err != nil {
		t.Fatalf("序列化第三方账号响应失败: %v", err)
	}
	text := string(body)
	for _, secret := range []string{row.Password, row.ApplicationKey, row.LastToken} {
		if strings.Contains(text, secret) {
			t.Fatalf("第三方账号响应泄漏敏感信息: %s", text)
		}
	}
}

func TestDingTalkBotRespDoesNotExposeCredentials(t *testing.T) {
	bot := &DingTalkBot{
		Webhook:      "https://example.com/webhook-secret",
		Secret:       "sign-secret",
		ClientID:     "client-id-secret",
		ClientSecret: "client-secret",
		AgentID:      "agent-id-secret",
	}
	body, err := json.Marshal(NewDingTalkBotResp(bot))
	if err != nil {
		t.Fatalf("序列化钉钉机器人响应失败: %v", err)
	}
	text := string(body)
	for _, secret := range []string{bot.Webhook, bot.Secret, bot.ClientID, bot.ClientSecret, bot.AgentID} {
		if strings.Contains(text, secret) {
			t.Fatalf("钉钉机器人响应泄漏敏感信息: %s", text)
		}
	}
}

func TestThirdPartyOrderDoesNotExposeRawJSON(t *testing.T) {
	row := &ThirdPartyOrder{RawJSON: `{"private":"raw-secret"}`}
	body, err := json.Marshal(row)
	if err != nil {
		t.Fatalf("序列化第三方订单响应失败: %v", err)
	}
	if strings.Contains(string(body), "raw-secret") || strings.Contains(string(body), "raw_json") {
		t.Fatalf("第三方订单响应泄漏原始报文: %s", body)
	}
}
