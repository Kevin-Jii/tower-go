package service

import (
	"testing"
	"time"

	"github.com/Kevin-Jii/tower-go/model"
)

func TestBuildPreOrderReminderTargets(t *testing.T) {
	now := time.Date(2026, 8, 8, 9, 30, 0, 0, preOrderLocation)
	targets := BuildPreOrderReminderTargets(now, PreOrderReminderSlot0930)
	if len(targets) != 2 {
		t.Fatalf("target count = %d, want 2", len(targets))
	}
	if got := targets[0].Start.Format("2006-01-02"); got != "2026-08-09" {
		t.Fatalf("previous-day target date = %s", got)
	}
	if targets[0].ReminderKey != model.PreOrderReminderPreviousDay0930 || targets[0].RelativeDate != "明天" {
		t.Fatalf("unexpected previous-day target: %#v", targets[0])
	}
	if got := targets[1].Start.Format("2006-01-02"); got != "2026-08-08" {
		t.Fatalf("due-day target date = %s", got)
	}
	if targets[1].ReminderKey != model.PreOrderReminderDueDay0930 || targets[1].RelativeDate != "今天" {
		t.Fatalf("unexpected due-day target: %#v", targets[1])
	}

	afternoon := BuildPreOrderReminderTargets(now, PreOrderReminderSlot1600)
	if afternoon[0].ReminderKey != model.PreOrderReminderPreviousDay1600 || afternoon[1].ReminderKey != model.PreOrderReminderDueDay1600 {
		t.Fatalf("unexpected afternoon keys: %#v", afternoon)
	}
}

func TestPreOrderStatusAllowsReminder(t *testing.T) {
	for _, status := range []int8{model.PreOrderStatusPending, model.PreOrderStatusPrepared} {
		if !PreOrderStatusAllowsReminder(status) {
			t.Fatalf("status %d should allow reminder", status)
		}
	}
	for _, status := range []int8{model.PreOrderStatusDelivered, model.PreOrderStatusCancelled} {
		if PreOrderStatusAllowsReminder(status) {
			t.Fatalf("status %d should not allow reminder", status)
		}
	}
}

func TestPreOrderStatusTransitionAllowed(t *testing.T) {
	if !preOrderStatusTransitionAllowed(model.PreOrderStatusPending, model.PreOrderStatusPrepared) {
		t.Fatal("pending should transition to prepared")
	}
	if !preOrderStatusTransitionAllowed(model.PreOrderStatusPrepared, model.PreOrderStatusDelivered) {
		t.Fatal("prepared should transition to delivered")
	}
	if preOrderStatusTransitionAllowed(model.PreOrderStatusDelivered, model.PreOrderStatusPending) {
		t.Fatal("delivered must be terminal")
	}
}

func TestSamePreOrderDate(t *testing.T) {
	morning := time.Date(2026, 8, 8, 9, 30, 0, 0, preOrderLocation)
	afternoon := time.Date(2026, 8, 8, 16, 0, 0, 0, preOrderLocation)
	tomorrow := time.Date(2026, 8, 9, 9, 30, 0, 0, preOrderLocation)
	if !samePreOrderDate(morning, afternoon) {
		t.Fatal("times on the same local date should match")
	}
	if samePreOrderDate(morning, tomorrow) {
		t.Fatal("different local dates should not match")
	}
}
