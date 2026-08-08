package controller

import (
	"testing"
	"time"
)

func TestResolveDailyTurnoverBusinessDate(t *testing.T) {
	now := time.Date(2026, time.August, 8, 10, 30, 0, 0, chinaStandardTime)

	got, err := resolveDailyTurnoverBusinessDate("", now)
	if err != nil {
		t.Fatalf("default date error = %v", err)
	}
	if got != "2026-08-07" {
		t.Fatalf("default date = %q, want %q", got, "2026-08-07")
	}

	got, err = resolveDailyTurnoverBusinessDate(" 2026-07-31 ", now)
	if err != nil {
		t.Fatalf("explicit date error = %v", err)
	}
	if got != "2026-07-31" {
		t.Fatalf("explicit date = %q, want %q", got, "2026-07-31")
	}

	if _, err := resolveDailyTurnoverBusinessDate("2026/07/31", now); err == nil {
		t.Fatal("invalid date unexpectedly passed validation")
	}
}

func TestResolveDailyTurnoverBusinessDateBeforeBusinessDayEnd(t *testing.T) {
	now := time.Date(2026, time.August, 8, 3, 0, 0, 0, chinaStandardTime)
	got, err := resolveDailyTurnoverBusinessDate("", now)
	if err != nil {
		t.Fatalf("default date error = %v", err)
	}
	if got != "2026-08-06" {
		t.Fatalf("default date = %q, want previous finalized business date %q", got, "2026-08-06")
	}
}
