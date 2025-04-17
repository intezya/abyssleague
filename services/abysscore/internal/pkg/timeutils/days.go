package timeutils

import "time"

func IsDayBeforeToday(day time.Time) bool {
	now := time.Now()

	dayDate := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
	nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	return dayDate.AddDate(0, 0, 1).Equal(nowDate)
}
