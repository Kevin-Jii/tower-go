package service

import (
	"errors"
	"strings"

	"github.com/Kevin-Jii/tower-go/config"
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
	"github.com/Kevin-Jii/tower-go/pkg/aliyunsms"
	"github.com/Kevin-Jii/tower-go/utils/encryption"
)

type StoreSmsConfigService struct {
	module      *module.StoreSmsConfigModule
	storeModule *module.StoreModule
}

func NewStoreSmsConfigService(m *module.StoreSmsConfigModule, storeModule *module.StoreModule) *StoreSmsConfigService {
	return &StoreSmsConfigService{module: m, storeModule: storeModule}
}

func (s *StoreSmsConfigService) GetByStoreID(storeID uint) (*model.StoreSmsConfigResp, error) {
	row, err := s.module.GetByStoreID(storeID)
	if err != nil {
		return nil, err
	}
	return s.toResp(row, false), nil
}

func (s *StoreSmsConfigService) Upsert(storeID uint, req *model.UpsertStoreSmsConfigReq, userID uint) (*model.StoreSmsConfigResp, error) {
	if storeID == 0 {
		return nil, errors.New("缺少门店 ID")
	}
	if _, err := s.storeModule.GetByID(storeID); err != nil {
		return nil, errors.New("门店不存在")
	}
	row, _ := s.module.GetByStoreID(storeID)
	if row == nil {
		row = &model.StoreSmsConfig{StoreID: storeID}
	}
	row.AccessKeyID = strings.TrimSpace(req.AccessKeyID)
	if req.RegionID != "" {
		row.RegionID = strings.TrimSpace(req.RegionID)
	} else if row.RegionID == "" {
		row.RegionID = "cn-hangzhou"
	}
	row.SignName = strings.TrimSpace(req.SignName)
	row.Enabled = req.Enabled
	row.SendWindowStart = strings.TrimSpace(req.SendWindowStart)
	row.SendWindowEnd = strings.TrimSpace(req.SendWindowEnd)
	row.ConfiguredBy = &userID

	// Secret 留空 = 不变更
	if secret := strings.TrimSpace(req.AccessKeySecret); secret != "" {
		cipherText, err := encryption.Encrypt(secret)
		if err != nil {
			return nil, errors.New("加密 AccessKey Secret 失败")
		}
		row.AccessKeySecretCipher = cipherText
	} else if row.ID == 0 {
		// 新建且未提供 Secret → 报错
		return nil, errors.New("请填写 AccessKey Secret")
	}
	if err := s.module.Upsert(row); err != nil {
		return nil, err
	}
	return s.toResp(row, false), nil
}

// TestConnection 用提供的密钥（不写库）发起一次阿里云模板查询，验证可用性。
func (s *StoreSmsConfigService) TestConnection(req *model.StoreSmsConfigTestReq) (string, error) {
	if strings.TrimSpace(req.AccessKeyID) == "" || strings.TrimSpace(req.AccessKeySecret) == "" {
		return "", errors.New("请填写 AccessKey ID 与 Secret")
	}
	region := strings.TrimSpace(req.RegionID)
	if region == "" {
		region = "cn-hangzhou"
	}
	client, err := aliyunsms.NewClient(aliyunsms.Config{
		AccessKeyID:     req.AccessKeyID,
		AccessKeySecret: req.AccessKeySecret,
		RegionID:        region,
		Enabled:         true,
	})
	if err != nil {
		return "", err
	}
	// 用 QuerySmsTemplateList 查询任意模板（拉取一页即可），错误包含阿里云侧原因。
	snap, err := client.GetTemplate("")
	if err == nil && snap != nil {
		// GetTemplate 找不到模板会返回 TemplateCode 空，不视为失败。
		return "OK：AccessKey 已通过阿里云认证。", nil
	}
	return "连通性测试：" + err.Error(), nil
}

func (s *StoreSmsConfigService) toResp(row *model.StoreSmsConfig, includePlainSecret bool) *model.StoreSmsConfigResp {
	resp := &model.StoreSmsConfigResp{
		ID:              row.ID,
		StoreID:         row.StoreID,
		AccessKeyID:     row.AccessKeyID,
		RegionID:        row.RegionID,
		SignName:        row.SignName,
		Enabled:         row.Enabled,
		SendWindowStart: row.SendWindowStart,
		SendWindowEnd:   row.SendWindowEnd,
		LastTestedAt:    row.LastTestedAt,
		LastTestMessage: row.LastTestMessage,
		CreatedAt:       row.CreatedAt,
		UpdatedAt:       row.UpdatedAt,
	}
	if includePlainSecret {
		plain, _ := encryption.Decrypt(row.AccessKeySecretCipher)
		resp.AccessKeySecret = plain
	} else {
		plain, err := encryption.Decrypt(row.AccessKeySecretCipher)
		if err != nil || plain == "" {
			resp.AccessKeySecret = ""
		} else {
			resp.AccessKeySecret = encryption.MaskSecret(plain)
		}
	}
	return resp
}

// EffectiveSmsContext 解析一个门店可用的 SMS 客户端与默认签名/窗口。
// 优先使用门店独立配置，否则回退到 .env 的 ALIYUN_SMS_* 全局值（多租户未配置门店兜底）。
type EffectiveSmsContext struct {
	Client           *aliyunsms.Client
	AccessKeyID      string
	AccessKeySecret  string
	RegionID         string
	SignName         string
	Enabled          bool
	SendWindowStart  string
	SendWindowEnd    string
	WindowConfigured bool
}

// ResolveEffectiveContext 返回 storeID 对应的有效 SMS 配置。
func (s *StoreSmsConfigService) ResolveEffectiveContext(storeID uint) (*EffectiveSmsContext, error) {
	if row, err := s.module.GetByStoreID(storeID); err == nil && row != nil {
		secret, _ := encryption.Decrypt(row.AccessKeySecretCipher)
		if row.Enabled && strings.TrimSpace(secret) != "" && strings.TrimSpace(row.AccessKeyID) != "" {
			client, err := aliyunsms.NewClient(aliyunsms.Config{
				AccessKeyID:     row.AccessKeyID,
				AccessKeySecret: secret,
				RegionID:        row.RegionID,
				SignName:        row.SignName,
				Enabled:         true,
			})
			if err != nil {
				return nil, err
			}
			return &EffectiveSmsContext{
				Client:           client,
				AccessKeyID:      row.AccessKeyID,
				AccessKeySecret:  secret,
				RegionID:         row.RegionID,
				SignName:         row.SignName,
				Enabled:          true,
				SendWindowStart:  row.SendWindowStart,
				SendWindowEnd:    row.SendWindowEnd,
				WindowConfigured: row.SendWindowStart != "" && row.SendWindowEnd != "",
			}, nil
		}
	}
	// 回退：全局 env
	cfg := config.GetAliyunSMSConfig()
	return &EffectiveSmsContext{
		Client:          nil,
		AccessKeyID:     cfg.AccessKeyID,
		AccessKeySecret: cfg.AccessKeySecret,
		RegionID:        cfg.RegionID,
		SignName:        cfg.SignName,
		Enabled:         cfg.Enabled && cfg.AccessKeyID != "" && cfg.AccessKeySecret != "",
	}, nil
}
