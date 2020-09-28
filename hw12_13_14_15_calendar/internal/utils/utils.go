package utils

import "time"

func StartDay(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
}

func EndDay(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day()+1, 0, 0, 0, 0, date.Location())
}

func StartWeek(date time.Time) time.Time {
	// Opetarate rounded by midnight date.
	rDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	switch date.Weekday() {
	// Sunday is zero day of week.
	case time.Sunday:
		rDate = rDate.Add(-6 * time.Hour * 24)
	default:
		rDate = rDate.Add(-1 * time.Duration(rDate.Weekday()-1) * time.Hour * 24)
	}

	return rDate
}

func EndWeek(date time.Time) time.Time {
	startDate := StartWeek(date)

	return startDate.Add(7 * time.Hour * 24)
}

func StartMonth(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
}

func EndMonth(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month()+1, 1, 0, 0, 0, 0, date.Location())
}
