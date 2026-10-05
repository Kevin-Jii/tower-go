package service

import (
	"errors"
	"testing"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/pkg/aliyunsms"
	"github.com/stretchr/testify/require"
)

type fakeAliyunSmsClient struct {
	signatures  []aliyunsms.SignatureSnapshot
	templates   []aliyunsms.TemplateSnapshot
	template    *aliyunsms.TemplateSnapshot
	deletedCode string
}

func (f *fakeAliyunSmsClient) CreateTemplate(string, string, string, string, int32) (string, error) {
	return "SMS_001", nil
}
func (f *fakeAliyunSmsClient) GetTemplate(string) (*aliyunsms.TemplateSnapshot, error) {
	if f.template == nil {
		return nil, errors.New("missing template")
	}
	return f.template, nil
}
func (f *fakeAliyunSmsClient) DeleteTemplate(code string) error {
	f.deletedCode = code
	return nil
}
func (f *fakeAliyunSmsClient) ListTemplates() ([]aliyunsms.TemplateSnapshot, error) {
	return f.templates, nil
}
func (f *fakeAliyunSmsClient) ListSignatures() ([]aliyunsms.SignatureSnapshot, error) {
	return f.signatures, nil
}

type fakeAliyunResolver struct {
	clients      map[uint]AliyunSmsClient
	resolvedWith []uint
}

func (f *fakeAliyunResolver) Resolve(storeID uint) (AliyunSmsClient, error) {
	f.resolvedWith = append(f.resolvedWith, storeID)
	client := f.clients[storeID]
	if client == nil {
		return nil, errors.New("client not found")
	}
	return client, nil
}

type fakeAliyunTemplateRepository struct {
	row          *model.AliyunSmsTemplate
	getStoreID   uint
	getAllStores bool
	upserted     *model.AliyunSmsTemplate
	synced       []*model.AliyunSmsTemplate
	deleted      *model.AliyunSmsTemplate
	exists       bool
}

func (f *fakeAliyunTemplateRepository) List(uint, bool, string, string) ([]model.AliyunSmsTemplate, error) {
	return nil, nil
}
func (f *fakeAliyunTemplateRepository) GetByCode(_ string, storeID uint, allStores bool) (*model.AliyunSmsTemplate, error) {
	f.getStoreID, f.getAllStores = storeID, allStores
	if f.row == nil {
		return nil, errors.New("not found")
	}
	copy := *f.row
	return &copy, nil
}
func (f *fakeAliyunTemplateRepository) ExistsByName(string, uint) (bool, error) {
	return f.exists, nil
}
func (f *fakeAliyunTemplateRepository) Upsert(row *model.AliyunSmsTemplate) error {
	f.upserted = row
	return nil
}
func (f *fakeAliyunTemplateRepository) SyncFromAliyun(row *model.AliyunSmsTemplate) error {
	f.synced = append(f.synced, row)
	return nil
}
func (f *fakeAliyunTemplateRepository) Delete(row *model.AliyunSmsTemplate) error {
	f.deleted = row
	return nil
}

func TestAliyunSmsTemplateListSynchronizesAliyunAccount(t *testing.T) {
	client := &fakeAliyunSmsClient{templates: []aliyunsms.TemplateSnapshot{{
		TemplateCode: "SMS_REMOTE", TemplateName: "remote", TemplateContent: "content",
		TemplateType: 2, TemplateStatus: "AUDIT_STATE_PASS",
	}}}
	resolver := &fakeAliyunResolver{clients: map[uint]AliyunSmsClient{7: client}}
	repo := &fakeAliyunTemplateRepository{}
	svc := NewAliyunSmsTemplateService(repo, resolver)

	_, err := svc.List(7, false, "", "")
	require.NoError(t, err)
	require.Len(t, repo.synced, 1)
	require.Equal(t, uint(7), repo.synced[0].OwnerStoreID)
	require.Equal(t, "SMS_REMOTE", repo.synced[0].TemplateCode)
	require.Equal(t, model.SmsTemplateAuditApproved, repo.synced[0].AuditStatus)
}

func TestAliyunSmsTemplateListSignaturesFiltersApproved(t *testing.T) {
	client := &fakeAliyunSmsClient{signatures: []aliyunsms.SignatureSnapshot{
		{SignName: "approved", AuditStatus: "AUDIT_STATE_PASS"},
		{SignName: "legacy-approved", AuditStatus: "1"},
		{SignName: "pending", AuditStatus: "AUDIT_STATE_INIT"},
	}}
	resolver := &fakeAliyunResolver{clients: map[uint]AliyunSmsClient{7: client}}
	svc := NewAliyunSmsTemplateService(&fakeAliyunTemplateRepository{}, resolver)

	rows, err := svc.ListSignatures(7, true)
	require.NoError(t, err)
	require.Equal(t, []string{"approved", "legacy-approved"}, []string{rows[0].SignName, rows[1].SignName})
	require.Equal(t, []uint{7}, resolver.resolvedWith)
}

func TestAliyunSmsTemplateRefreshUsesPersistedOwnerClient(t *testing.T) {
	repo := &fakeAliyunTemplateRepository{row: &model.AliyunSmsTemplate{
		OwnerStoreID: 9, TemplateCode: "SMS_009", AuditStatus: model.SmsTemplateAuditPending,
	}}
	client := &fakeAliyunSmsClient{template: &aliyunsms.TemplateSnapshot{
		TemplateCode: "SMS_009", TemplateStatus: "1", TemplateContent: "updated",
	}}
	resolver := &fakeAliyunResolver{clients: map[uint]AliyunSmsClient{9: client}}
	svc := NewAliyunSmsTemplateService(repo, resolver)

	row, err := svc.Refresh("SMS_009", 9, false)
	require.NoError(t, err)
	require.Equal(t, uint(9), repo.getStoreID)
	require.False(t, repo.getAllStores)
	require.Equal(t, []uint{9}, resolver.resolvedWith)
	require.Equal(t, model.SmsTemplateAuditApproved, row.AuditStatus)
	require.Equal(t, "updated", repo.upserted.Content)
}

func TestAliyunSmsTemplateDeleteUsesPersistedOwnerClient(t *testing.T) {
	repo := &fakeAliyunTemplateRepository{row: &model.AliyunSmsTemplate{OwnerStoreID: 12, TemplateCode: "SMS_012"}}
	client := &fakeAliyunSmsClient{}
	resolver := &fakeAliyunResolver{clients: map[uint]AliyunSmsClient{12: client}}
	svc := NewAliyunSmsTemplateService(repo, resolver)

	require.NoError(t, svc.Delete("SMS_012", 12, false))
	require.Equal(t, []uint{12}, resolver.resolvedWith)
	require.Equal(t, "SMS_012", client.deletedCode)
	require.Equal(t, "SMS_012", repo.deleted.TemplateCode)
}
