package helpers

import (
	"time"
)

const DMYLayout = "02-01-2006"

func StartOfDay(t time.Time) time.Time {
	year, month, day := t.In(time.Local).Date()
	dayStartTime := time.Date(year, month, day, 0, 0, 0, 0, time.Local)

	return dayStartTime
}

func EndOfDay(t time.Time) time.Time {
	year, month, day := t.In(time.Local).Date()
	dayEndTime := time.Date(year, month, day, 23, 59, 59, 0, time.Local)

	return dayEndTime
}

func StartOfMonth(t time.Time) time.Time {
	year, month, _ := t.In(time.Local).Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
}

func EndOfMonth(t time.Time) time.Time {
	year, month, _ := t.In(time.Local).Date()
	nextMonth := month + 1
	if nextMonth > 12 {
		nextMonth = 1
		year++
	}
	startOfNextMonth := time.Date(year, nextMonth, 1, 0, 0, 0, 0, time.Local)
	return startOfNextMonth.Add(-time.Nanosecond)
}

func OneMonth() string {
	currentTime := time.Now()
	subscriptionEndDate := currentTime.AddDate(0, 1, 0)

	return subscriptionEndDate.Format("02 January 2006")
}
