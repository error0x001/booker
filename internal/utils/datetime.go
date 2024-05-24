package utils

import (
	"fmt"
	"time"
)

func DaysBetween(from, to time.Time) []time.Time {
	if from.After(to) {
		return nil
	}

	from = ToDay(from)
	to = ToDay(to)

	numDays := int(to.Sub(from).Hours()/24) + 1
	days := make([]time.Time, 0, numDays)

	for d := from; !d.After(to); d = d.AddDate(0, 0, 1) {
		days = append(days, d)
	}

	return days
}

func ToDay(timestamp time.Time) time.Time {
	return time.Date(timestamp.Year(), timestamp.Month(), timestamp.Day(), 0, 0, 0, 0, time.UTC)
}

func Date(year, month, day int) time.Time {
	// validate the month and day input
	if month < 1 || month > 12 {
		panic(fmt.Sprintf("invalid month: %d", month))
	}
	if day < 1 || day > 31 {
		panic(fmt.Sprintf("invalid day: %d", day))
	}
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}
