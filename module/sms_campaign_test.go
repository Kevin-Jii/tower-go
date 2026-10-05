package module

import (
	"strings"
	"testing"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func newSMSCampaignDryRunDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       "gorm:gorm@tcp(localhost:9910)/gorm?charset=utf8&parseTime=True&loc=Local",
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		DryRun:                 true,
		DisableAutomaticPing:   true,
		SkipDefaultTransaction: true,
	})
	require.NoError(t, err)
	return db
}

func TestSmsCampaignScopeRestrictsStoreBoundQueries(t *testing.T) {
	db := newSMSCampaignDryRunDB(t)
	m := NewSmsCampaignModule(db)

	scopedSQL := db.ToSQL(func(_ *gorm.DB) *gorm.DB {
		return m.scoped(7, false).Where("id = ?", 42).Find(&[]model.SmsCampaign{})
	})
	require.Contains(t, scopedSQL, "owner_store_id")
	require.Contains(t, scopedSQL, "= 7")

	hqSQL := db.ToSQL(func(_ *gorm.DB) *gorm.DB {
		return m.scoped(0, true).Where("id = ?", 42).Find(&[]model.SmsCampaign{})
	})
	require.NotContains(t, hqSQL, "owner_store_id")
}

func TestSmsCampaignClaimUsesSingleConditionalUpdate(t *testing.T) {
	db := newSMSCampaignDryRunDB(t)
	var sql string
	var vars []interface{}
	require.NoError(t, db.Callback().Update().After("gorm:update").Register("test:capture_sms_campaign_claim", func(tx *gorm.DB) {
		sql = tx.Statement.SQL.String()
		vars = append([]interface{}(nil), tx.Statement.Vars...)
	}))

	claimed, err := NewSmsCampaignModule(db).Claim(99)
	require.NoError(t, err)
	require.False(t, claimed, "dry-run statements do not report an affected row")
	require.Contains(t, strings.ToLower(sql), "update")
	require.Contains(t, sql, "id = ?")
	require.Contains(t, sql, "status IN (?,?)")
	require.Contains(t, vars, uint(99))
	require.Contains(t, vars, model.SmsCampaignStatusDraft)
	require.Contains(t, vars, model.SmsCampaignStatusScheduled)
	require.Contains(t, vars, model.SmsCampaignStatusSending)
}

func TestClaimFailedSendRecordIsScopedAndAtomic(t *testing.T) {
	db := newSMSCampaignDryRunDB(t)
	var sql string
	var vars []interface{}
	require.NoError(t, db.Callback().Update().After("gorm:update").Register("test:capture_sms_retry_claim", func(tx *gorm.DB) {
		sql = tx.Statement.SQL.String()
		vars = append([]interface{}(nil), tx.Statement.Vars...)
	}))

	claimed, err := NewSmsCampaignModule(db).ClaimFailedSendRecord(12, 34)
	require.NoError(t, err)
	require.False(t, claimed, "dry-run statements do not report an affected row")
	require.Contains(t, sql, "campaign_id = ? AND id = ? AND status = ?")
	require.Contains(t, vars, uint(12))
	require.Contains(t, vars, uint(34))
	require.Contains(t, vars, model.SmsSendRecordFailed)
	require.Contains(t, vars, model.SmsSendRecordRetrying)
}
