package parser

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-playground/locales"
)

// Formats date strings via the same system as `strptime`
var Strptime = DateFormatterPrefix{
	Prefix: '%',
	TokenMap: map[rune]FormatToken[rune]{
		'%': {
			Desc:   "'%' character literal",
			expand: func(dt time.Time, locale locales.Translator) string { return "%" },
		},
		'A': {
			Desc:   "Weekday name - 'Monday', 'Tuesday'",
			expand: func(dt time.Time, locale locales.Translator) string { return dt.Weekday().String() },
		},
		'a': {
			Desc:   "Weekday name truncated to three characters - 'Mon', 'Tue'",
			expand: func(dt time.Time, locale locales.Translator) string { return dt.Weekday().String()[:3] },
		},
		'B': {
			Desc:   "Month name - 'January', 'February'",
			expand: func(dt time.Time, locale locales.Translator) string { return dt.Month().String() },
		},
		'b': {
			Desc:    "Month month name truncated to three characters - 'Jan', 'Feb'",
			expand:  func(dt time.Time, locale locales.Translator) string { return dt.Month().String()[:3] },
			aliases: []rune{'h'},
		},
		'c': {
			Desc:   "Date and time for the current locale (different to strptime and hardcoded to UK format currently)",
			expand: func(dt time.Time, locale locales.Translator) string { return dt.Format("Mon _2 Jan 15:04:05 2006") },
		},
		'C': {
			Desc:   "The century number (0–99)",
			expand: func(dt time.Time, locale locales.Translator) string { return strconv.Itoa(dt.Year() / 100) },
		},
		'd': {
			Desc:   "Day of month zero padded to two digits (01-31)",
			expand: func(dt time.Time, locale locales.Translator) string { return fmt.Sprintf("%02d", dt.Day()) },
		},
		'e': {
			Desc:   "Day of month (1-31)",
			expand: func(dt time.Time, locale locales.Translator) string { return strconv.Itoa(dt.Day()) },
		},
		'D': {
			Desc: "American style date (month first) equivalent to '%m/%d/%y' where the year is truncated to the last two digits - '01/31/97', '02/28/01'",
			expand: func(dt time.Time, locale locales.Translator) string {
				return fmt.Sprintf("%02d/%02d/%02d", dt.Month(), dt.Day(), dt.Year()%100)
			},
		},
		'H': {
			Desc:   "Hour in 24 hour format zero padded to two digits (00-23)",
			expand: func(dt time.Time, locale locales.Translator) string { return fmt.Sprintf("%02d", dt.Hour()) },
		},
		'I': {
			Desc:   "Hour in 12 hour format zero padded to two digits (01-12)",
			expand: func(dt time.Time, locale locales.Translator) string { return fmt.Sprintf("%02d", dt.Hour()%12) },
		},
		'j': {
			Desc:   "Day of year zero padded to three digits (001-366)",
			expand: func(dt time.Time, locale locales.Translator) string { return fmt.Sprintf("%03d", dt.YearDay()) },
		},
		'm': {
			Desc:   "Month number zero padded to two digits (01-12)",
			expand: func(dt time.Time, locale locales.Translator) string { return fmt.Sprintf("%02d", dt.Month()) },
		},
		'M': {
			Desc:   "Minutes zero padded to two digits (00-59)",
			expand: func(dt time.Time, locale locales.Translator) string { return fmt.Sprintf("%02d", dt.Minute()) },
		},
		'n': {
			Desc:   "Newline whitespace - '\\n'",
			expand: func(dt time.Time, locale locales.Translator) string { return "\n" },
		},
		'R': {
			Desc: "Time represented as hours and minutes equivalent to %H:%M - '12:24', '04:09'",
			expand: func(dt time.Time, locale locales.Translator) string {
				return fmt.Sprintf("%02d:%02d", dt.Hour(), dt.Minute())
			},
		},
		'S': {
			Desc:   "Seconds zero padded to two digits (00-60; 60 may occur for for leap seconds)",
			expand: func(dt time.Time, locale locales.Translator) string { return fmt.Sprintf("%02d", dt.Second()) },
		},
		't': {
			Desc:   "Tab whitespace - '\\t'",
			expand: func(dt time.Time, locale locales.Translator) string { return "\t" },
		},
		'T': {
			Desc: "Time represented as hours, minutes and seconds equivalent to %H:%M:%S - '12:34:03', '04:09:59'",
			expand: func(dt time.Time, locale locales.Translator) string {
				return fmt.Sprintf("%02d:%02d:%02d", dt.Hour(), dt.Minute(), dt.Second())
			},
		},
		// TODO: implement same rules as strptime - first week is the first Sunday of January
		'U': {
			Desc: "ISO8601 week number of the year (this is different to the typical strptime token currently)",
			expand: func(dt time.Time, locale locales.Translator) string {
				_, weekNumber := dt.ISOWeek()
				return fmt.Sprintf("%02d", weekNumber)
			},
		},
		'w': {
			Desc:   "Day of week number (0-6) where Sunday is 0 and Saturday is 6",
			expand: func(dt time.Time, locale locales.Translator) string { return strconv.Itoa(int(dt.Weekday())) },
		},
		// TODO: implement same rules as strptime - first week is the first Monday of January
		'W': {
			Desc: "ISO8601 week number of the year (this is different to the typical strptime token currently)",
			expand: func(dt time.Time, locale locales.Translator) string {
				_, week := dt.ISOWeek()
				return fmt.Sprintf("%02d", week)
			},
		},
		// NOTE: haven't found a locale aware version of these:
		// 'x': func(dt time.Time, locale locales.Translator) string { return dt.Format("01/02/06") },
		// 'X': func(dt time.Time, locale locales.Translator) string { return dt.Format(time.TimeOnly) },
		'y': {
			Desc:   "The year within the century zero padded to two digits (00-99)",
			expand: func(dt time.Time, locale locales.Translator) string { return fmt.Sprintf("%02d", dt.Year()%100) },
		},
		'Y': {
			Desc:   "Year number - '1999', '2007'",
			expand: func(dt time.Time, locale locales.Translator) string { return strconv.Itoa(dt.Year()) },
		},
	},
}
