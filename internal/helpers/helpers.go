package helpers

import "time"

func IsToday(date time.Time) bool {
	now := time.Now()
	return date.Year() == now.Year() && date.YearDay() == now.YearDay()
}
