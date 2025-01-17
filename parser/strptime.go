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
	tokenMap: map[rune]FormatToken[rune]{
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
			Desc:   "Date and time for the current locale (hardcoded to UK format currently)",
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
			Desc: "Day of month space padded to two characters ( 1-31)",
			expand: func(dt time.Time, locale locales.Translator) string {
				fmtstr := "%d"
				if dt.Day() < 10 {
					fmtstr = " %d"
				}
				return fmt.Sprintf(fmtstr, dt.Day())
			},
		},
		'F': {
			Desc:   "Date in year-month-day format equivalent to '%Y-%m-%d' - '2024-01-04', '1997-10-31'",
			expand: func(dt time.Time, locale locales.Translator) string { return fmt.Sprintf("% 2d", dt.Day()) },
		},
		'g': {
			Desc: "ISO week year shortened to the last two digits (00-99) ",
			expand: func(dt time.Time, locale locales.Translator) string {
				year, _ := dt.ISOWeek()
				return fmt.Sprintf("%02d", year%100)
			},
		},
		'G': {
			Desc: "ISO week year - '1999', '2007'",
			expand: func(dt time.Time, locale locales.Translator) string {
				year, _ := dt.ISOWeek()
				return fmt.Sprintf("%d", year)
			},
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
			Desc: "Hour in 12 hour format zero padded to two digits (01-12)",
			expand: func(dt time.Time, locale locales.Translator) string {
				hour := dt.Hour() % 12
				if hour == 0 {
					hour = 12
				}
				return fmt.Sprintf("%02d", hour)
			},
		},
		'j': {
			Desc:   "Day of year zero padded to three digits (001-366)",
			expand: func(dt time.Time, locale locales.Translator) string { return fmt.Sprintf("%03d", dt.YearDay()) },
		},
		'k': {
			Desc: "Hour in 24 hour format space padded to two digits ( 0-23)",
			expand: func(dt time.Time, locale locales.Translator) string {
				hour := dt.Hour()
				if hour > 9 {
					return fmt.Sprintf("%2d", hour)
				}
				return fmt.Sprintf(" %d", hour)
			},
		},
		'l': {
			Desc: "Hour in 12 hour format space padded to two digits ( 0-12)",
			expand: func(dt time.Time, locale locales.Translator) string {
				hour := dt.Hour() % 12
				if hour == 0 {
					hour = 12
				}
				if hour > 9 {
					return fmt.Sprintf("%2d", hour)
				}
				return fmt.Sprintf(" %d", hour)
			},
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
		'N': {
			Desc:   "Nanoseconds zero padded to nine digits (000000000-999999999)",
			expand: func(dt time.Time, locale locales.Translator) string { return strconv.Itoa(dt.Nanosecond()) },
		},
		'p': {
			Desc: "The locale's equivalent of AM or PM (hardcoded to English am/pm)",
			expand: func(dt time.Time, locale locales.Translator) string {
				// TODO: locale
				if dt.Hour() < 12 {
					return "am"
				}
				return "pm"
			},
		},
		'r': {
			Desc: "12 hour time represented as hours, minutes, seconds and am/pm equivalent to \"%I:%M:%S %p\" (hardcoded to English am/pm) - '11:24:52 pm', '04:09:20 am'",
			expand: func(dt time.Time, locale locales.Translator) string {
				hour := dt.Hour() % 12
				if hour == 0 {
					hour = 12
				}

				ampm := "pm"
				if dt.Hour() < 12 {
					ampm = "am"
				}
				return fmt.Sprintf("%02d:%02d:%02d %s", hour, dt.Minute(), dt.Second(), ampm)
			},
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
		'u': {
			Desc: "Day of week where Monday = 1 and Sunday = 7 (1-7)",
			expand: func(dt time.Time, locale locales.Translator) string {
				return strconv.Itoa((int(dt.Weekday())+6)%7 + 1)
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
		'V': {
			Desc: "ISO8601 week number of the year zero padded to two digits - (01-53)",
			expand: func(dt time.Time, locale locales.Translator) string {
				_, week := dt.ISOWeek()
				return fmt.Sprintf("%02d", week)
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
		'x': {
			Desc: "Locale date format - '04/12/1999', '11/02/2007'",
			expand: func(dt time.Time, locale locales.Translator) string {
				return locale.FmtDateShort(dt)
			},
		},
		'X': {
			Desc: "Locale time including seconds - '03:57:22', '18:08:01'",
			expand: func(dt time.Time, locale locales.Translator) string {
				return locale.FmtTimeMedium(dt)
			},
		},
		'y': {
			Desc:   "The year within the century zero padded to two digits (00-99)",
			expand: func(dt time.Time, locale locales.Translator) string { return fmt.Sprintf("%02d", dt.Year()%100) },
		},
		'Y': {
			Desc:   "Year number - '1999', '2007'",
			expand: func(dt time.Time, locale locales.Translator) string { return strconv.Itoa(dt.Year()) },
		},
		'z': {
			Desc: "Time zone offset in +hhmm format - '-0400', '+0530'",
			expand: func(dt time.Time, locale locales.Translator) string {
				_, offsetSeconds := dt.Zone()
				offsetMinutes := offsetSeconds / 60
				offsetHours := offsetMinutes / 60
				sign := '+'
				if offsetHours < 0 {
					sign = '-'
				}
				return fmt.Sprintf("%c%02d%02d", sign, offsetHours, offsetMinutes%60)

			},
		},
		'Z': {
			Desc: "Abbreviated time zone offset - 'GMT', 'CEST', '+0530'",
			expand: func(dt time.Time, locale locales.Translator) string {
				offsetName, _ := dt.Zone()
				return offsetName
			},
		},
	},
}
