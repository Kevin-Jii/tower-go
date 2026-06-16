package businessdate

import (
	"testing"
	"time"
)

func TestDate(t *testing.T) {
	loc := time.FixedZone("CST", 8*60*60)
	tests := []struct {
		name string
		at   time.Time
		want string
	}{
		{
			name: "before five belongs to previous day",
			at:   time.Date(2026, 6, 16, 0, 1, 0, 0, loc),
			want: "2026-06-15",
		},
		{
			name: "five belongs to current day",
			at:   time.Date(2026, 6, 16, 5, 0, 0, 0, loc),
			want: "2026-06-16",
		},
		{
			name: "afternoon belongs to current day",
			at:   time.Date(2026, 6, 15, 16, 0, 0, 0, loc),
			want: "2026-06-15",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DateString(tt.at); got != tt.want {
				t.Fatalf("DateString() = %s, want %s", got, tt.want)
			}
		})
	}
}
