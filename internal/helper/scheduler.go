package helper

import "time"

func GetNextTime(hour, min int) time.Time {
	now := time.Now()
	location := now.Location()
	next := time.Date(now.Year(), now.Month(), now.Day(), hour, min, 0, 0, location)
	if now.After(next) {
		next = next.Add(24 * time.Hour)
	}
	return next
}
