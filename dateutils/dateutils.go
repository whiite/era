package dateutils

import "time"

// Equivalent to midnight of January 1st for the year of the provided date time
// taking into account the location
func YearStart(t time.Time) time.Time {
	return time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
}

// Equivalent to 1 nanosecond before midnight of January 1st for the following
// year of the provided date time taking into account the location
func YearEnd(t time.Time) time.Time {
	return time.Date(t.Year(), 12, 31, 23, 59, 59, 999999999, t.Location())
}

// Equivalent to midnight of the provided date time taking into account the location
func DayStart(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// Returns midnight of the next desired weekday
//
// If the desired weekday is the current day then the current day will be returned
func NextWeekday(day time.Weekday, t time.Time) time.Time {
	daysUntilDay := (7 - t.Weekday() + day) % 7
	return time.Date(t.Year(), t.Month(), t.Day()+int(daysUntilDay), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
}

// Returns midnight of the next desired weekday
//
// If the desired weekday is the current day then the current day will be returned
func PreviousWeekday(day time.Weekday, t time.Time) time.Time {
	daysUntilDay := (t.Weekday() + day) % 7
	return time.Date(t.Year(), t.Month(), t.Day()-int(daysUntilDay), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
}

// Quarter of the provided year in the range of 1-4
func YearQuarter(t time.Time) int {
	daysInYear := float64(YearEnd(t).YearDay())
	quarterZeroed := (float64(t.YearDay()) / daysInYear) * 4
	if quarterZeroed == 4 {
		quarterZeroed = 3
	}
	return int(quarterZeroed + 1)
}
