package businessdate

import "time"

const (
	BusinessDayStartHour = 16
	BusinessDayEndHour   = 5
)

// Date returns the business date for store accounting.
// Business hours are 16:00 through 05:00 the next day, so entries made
// after midnight and before 05:00 belong to the previous calendar date.
func Date(t time.Time) time.Time {
	if t.Hour() < BusinessDayEndHour {
		t = t.AddDate(0, 0, -1)
	}
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func DateString(t time.Time) string {
	return Date(t).Format("2006-01-02")
}
