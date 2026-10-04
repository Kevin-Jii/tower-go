package module

import (
	"testing"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestAliyunSmsTemplateScopeRestrictsStoreBoundQueries(t *testing.T) {
	db := newSMSCampaignDryRunDB(t)
	m := NewAliyunSmsTemplateModule(db)

	storeSQL := db.ToSQL(func(_ *gorm.DB) *gorm.DB {
		return m.scoped(7, false).Where("template_code = ?", "SMS_1").Find(&model.AliyunSmsTemplate{})
	})
	require.Contains(t, storeSQL, "owner_store_id")
	require.Contains(t, storeSQL, "= 7")

	hqSQL := db.ToSQL(func(_ *gorm.DB) *gorm.DB {
		return m.scoped(0, true).Where("template_code = ?", "SMS_1").Find(&model.AliyunSmsTemplate{})
	})
	require.NotContains(t, hqSQL, "owner_store_id")
}
