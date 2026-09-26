package service

import (
	"testing"
	"time"

	"github.com/Kevin-Jii/tower-go/model"
)

func TestChooseTemplateUsesFirstMatchingSegmentThenDefault(t *testing.T) {
	segments := []model.SmsCampaignSegment{
		{ID: 10, Position: 0, TagIDs: model.UintList{1, 2}, TemplateCode: "FIRST"},
		{ID: 11, Position: 1, TagIDs: model.UintList{2}, TemplateCode: "SECOND"},
		{ID: 12, Position: 2, IsDefault: true, TemplateCode: "DEFAULT"},
	}
	legacy := smsTemplate{TemplateCode: "LEGACY"}
	got := chooseTemplate(smsRecipient{}, map[uint]struct{}{2: {}}, segments, legacy)
	if got.TemplateCode != "FIRST" || got.SegmentID == nil || *got.SegmentID != 10 {
		t.Fatalf("matching template = %#v, want first segment", got)
	}
	got = chooseTemplate(smsRecipient{}, nil, segments, legacy)
	if got.TemplateCode != "DEFAULT" || got.SegmentID == nil || *got.SegmentID != 12 {
		t.Fatalf("unmatched template = %#v, want default segment", got)
	}
	got = chooseTemplate(smsRecipient{}, nil, nil, legacy)
	if got.TemplateCode != "LEGACY" || got.SegmentID != nil {
		t.Fatalf("legacy template = %#v", got)
	}
}

func TestGroupRecipientsKeepsDeterministicInputOrder(t *testing.T) {
	in := []smsRecipient{
		{Phone: "1", Config: smsTemplate{TemplateCode: "A"}},
		{Phone: "2", Config: smsTemplate{TemplateCode: "B"}},
		{Phone: "3", Config: smsTemplate{TemplateCode: "A"}},
	}
	groups := groupRecipients(in)
	if len(groups) != 2 || groups[0].Config.TemplateCode != "A" || groups[1].Config.TemplateCode != "B" {
		t.Fatalf("groups = %#v", groups)
	}
	if len(groups[0].Recipients) != 2 || groups[0].Recipients[0].Phone != "1" || groups[0].Recipients[1].Phone != "3" {
		t.Fatalf("group A recipients = %#v", groups[0].Recipients)
	}
}

func TestScheduleNormalizesOffsetAndHonorsChinaWindow(t *testing.T) {
	t.Setenv("ALIYUN_SMS_SEND_WINDOW_START", "08:00")
	t.Setenv("ALIYUN_SMS_SEND_WINDOW_END", "22:00")
	svc := &SmsCampaignService{}
	loc, err := time.LoadLocation(smsTimezone)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2030, 1, 1, 7, 0, 0, 0, loc)
	// The +09:00 input is 08:30 in China and must be stored as one UTC instant.
	input := time.Date(2030, 1, 1, 9, 30, 0, 0, time.FixedZone("JST", 9*60*60))
	got, err := svc.normalizeAndValidateSchedule(input, now)
	if err != nil {
		t.Fatalf("normalize schedule: %v", err)
	}
	want := time.Date(2030, 1, 1, 0, 30, 0, 0, time.UTC)
	if !got.Equal(want) || got.Location() != time.UTC {
		t.Fatalf("normalized = %s (%s), want %s UTC", got, got.Location(), want)
	}
	endExclusive := time.Date(2030, 1, 1, 22, 0, 0, 0, loc)
	if _, err := svc.normalizeAndValidateSchedule(endExclusive, now); err == nil {
		t.Fatal("22:00 China time should be excluded")
	}
}
