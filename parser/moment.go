package parser

import (
	"fmt"
	"monokuro/era/dateutils"
	"strconv"
	"time"

	"github.com/go-playground/locales"
)

// Formats time according to the momentjs tokens
var MomentJs = DateFormatterString{
	escapeChars: []rune{'[', ']'},
	tokenMap: map[string]FormatToken[string]{
		"a": {
			Desc: "Meridiem - 'am'",
			expand: func(dt time.Time, locale locales.Translator) string {
				// TODO: locale
				if dt.Hour() < 12 {
					return "am"
				}
				return "pm"
			},
		},
		"A": {
			Desc: "Meridiem capitalised - 'AM'",
			expand: func(dt time.Time, locale locales.Translator) string {
				// TODO: locale
				if dt.Hour() < 12 {
					return "AM"
				}
				return "PM"
			},
		},
		"M": {
			Desc:   "Month number (1-12)",
			expand: func(dt time.Time, locale locales.Translator) string { return strconv.Itoa(int(dt.Month())) },
		},
		"Mo": {
			Desc:   "Month number suffixed - '1st', '13th', '22nd'",
			expand: func(dt time.Time, locale locales.Translator) string { return numberSuffixed(int(dt.Month())) },
		},
		"MM": {
			Desc:   "Month number zero padded to two digits - (01-12)",
			expand: func(dt time.Time, locale locales.Translator) string { return fmt.Sprintf("%02d", dt.Month()) },
		},
		"MMM": {
			Desc: "Month name truncated to three characters - 'Jan', 'Feb'",
			expand: func(dt time.Time, locale locales.Translator) string {
				return locale.MonthWide(dt.Month())[:3]
			},
		},
		"MMMM": {
			Desc: "Month name - 'January', 'February'",
			expand: func(dt time.Time, locale locales.Translator) string {
				return locale.MonthWide(dt.Month())
			},
		},
		"Q": {
			Desc:   "Quarter of year (1-4)",
			expand: func(dt time.Time, locale locales.Translator) string { return strconv.Itoa(dateutils.YearQuarter(dt)) },
		},
		"Qo": {
			Desc:   "Quarter of year suffixed - '1st', '2nd', '3rd', '4th'",
			expand: func(dt time.Time, locale locales.Translator) string { return numberSuffixed(dateutils.YearQuarter(dt)) },
		},
		"D": {
			Desc:   "Day of month (1-31)",
			expand: func(dt time.Time, locale locales.Translator) string { return strconv.Itoa(dt.Day()) },
		},
		"Do": {
			Desc:   "Day of month suffixed (1st-31st)",
			expand: func(dt time.Time, locale locales.Translator) string { return numberSuffixed(dt.Day()) },
		},
		"DD": {
			Desc:   "Day of month zero padded to two digits (01-31)",
			expand: func(dt time.Time, locale locales.Translator) string { return fmt.Sprintf("%02d", dt.Day()) },
		},
		"DDD": {
			Desc:   "Day of year (1-366)",
			expand: func(dt time.Time, locale locales.Translator) string { return strconv.Itoa(dt.YearDay()) },
		},
		"DDDo": {
			Desc:   "Day of year suffixed (1st-366th)",
			expand: func(dt time.Time, locale locales.Translator) string { return numberSuffixed(dt.YearDay()) },
		},
		"DDDD": {
			Desc:   "Day of year zero padded to three digits (001-366)",
			expand: func(dt time.Time, locale locales.Translator) string { return fmt.Sprintf("%03d", dt.YearDay()) },
		},
		"d": {
			Desc:   "Day of week where Sunday = 0 and Saturday = 6 (0-6)",
			expand: func(dt time.Time, locale locales.Translator) string { return strconv.Itoa(int(dt.Weekday())) },
		},
		"do": {
			Desc:   "Day of week suffixed where Sunday = 0th and Saturday = 6th (0th-6th)",
			expand: func(dt time.Time, locale locales.Translator) string { return numberSuffixed(int(dt.Weekday())) },
		},
		"dd": {
			Desc: "Day of week name truncated to two characters - 'Su', 'Mo'",
			expand: func(dt time.Time, locale locales.Translator) string {
				return locale.WeekdayShort(dt.Weekday())
			},
		},
		"ddd": {
			Desc: "Day of week name truncated to three characters - 'Sun', 'Mon'",
			expand: func(dt time.Time, locale locales.Translator) string {
				return locale.WeekdayWide(dt.Weekday())[:3]
			},
		},
		"dddd": {
			Desc: "Day of week name - 'Sunday', 'Monday'",
			expand: func(dt time.Time, locale locales.Translator) string {
				return locale.WeekdayWide(dt.Weekday())
			},
		},
		"e": {
			Desc: "Day of week where Monday = 1 and Sunday = 0 - (0-6)",
			expand: func(dt time.Time, locale locales.Translator) string {
				return strconv.Itoa(int(dt.Weekday()))
			},
		},
		"E": {
			Desc: "Day of week where Monday = 1 and Sunday = 7 - (1-7)",
			expand: func(dt time.Time, locale locales.Translator) string {
				return strconv.Itoa((int(dt.Weekday())+6)%7 + 1)
			},
		},
		"H": {
			Desc:   "Hour in 24 hour format (0-23)",
			expand: func(dt time.Time, locale locales.Translator) string { return strconv.Itoa(dt.Hour()) },
		},
		"HH": {
			Desc:   "Hour in 24 hour format zero padded to two digits (00-23)",
			expand: func(dt time.Time, locale locales.Translator) string { return fmt.Sprintf("%02d", dt.Hour()) },
		},
		"h": {
			Desc: "Hour in 12 hour format (1-12)",
			expand: func(dt time.Time, locale locales.Translator) string {
				hour := dt.Hour() % 12
				if hour == 0 {
					hour = 12
				}
				return strconv.Itoa(hour)
			},
		},
		"hh": {
			Desc: "Hour in 12 hour format zero padded to two digits (01-12)",
			expand: func(dt time.Time, locale locales.Translator) string {
				hour := dt.Hour() % 12
				if hour == 0 {
					hour = 12
				}
				return fmt.Sprintf("%02d", hour)
			},
		},
		"k": {
			Desc: "Hour in 24 hour format starting from 1 (1-24)",
			expand: func(dt time.Time, locale locales.Translator) string {
				hour := dt.Hour()
				if hour == 0 {
					hour = 24
				}
				return strconv.Itoa(hour)
			},
		},
		"kk": {
			Desc: "Hour in 24 hour format starting from 1 zero padded to two digits (01-24)",
			expand: func(dt time.Time, locale locales.Translator) string {
				hour := dt.Hour()
				if hour == 0 {
					hour = 24
				}
				return fmt.Sprintf("%02d", hour)
			},
		},
		// TODO: "w", "wo", "ww" - what do these represent?
		"W": {
			Desc: "ISO week of year (1-53)",
			expand: func(dt time.Time, locale locales.Translator) string {
				_, week := dt.ISOWeek()
				return strconv.Itoa(week)
			},
		},
		"Wo": {
			Desc: "ISO week of year zero padded (01-53)",
			expand: func(dt time.Time, locale locales.Translator) string {
				_, week := dt.ISOWeek()
				return numberSuffixed(week)
			},
		},
		"WW": {
			Desc: "ISO week of year zero padded (01-53)",
			expand: func(dt time.Time, locale locales.Translator) string {
				_, week := dt.ISOWeek()
				return fmt.Sprintf("%02d", week)
			},
		},

		"Y": {
			Desc:   "Year number - '1999', '2007'",
			expand: func(dt time.Time, locale locales.Translator) string { return strconv.Itoa(dt.Year()) },
		},
		"YY": {
			Desc:   "Year number truncated to last two digits - '99', '07'",
			expand: func(dt time.Time, locale locales.Translator) string { return strconv.Itoa(dt.Year() % 100) },
		},
		"YYYY": {
			Desc:   "Year number - '1999', '2007'",
			expand: func(dt time.Time, locale locales.Translator) string { return strconv.Itoa(dt.Year()) },
		},
		"YYYYYY": {
			Desc:   "Year number zeo padded to 6 digits - '001999', '002007'",
			expand: func(dt time.Time, locale locales.Translator) string { return fmt.Sprintf("%+07d", dt.Year()) },
		},
		"m": {
			Desc:   "Minutes (0-59)",
			expand: func(dt time.Time, locale locales.Translator) string { return strconv.Itoa(dt.Minute()) },
		},
		"mm": {
			Desc:   "Minutes zero padded to two digits (00-59)",
			expand: func(dt time.Time, locale locales.Translator) string { return fmt.Sprintf("%02d", dt.Minute()) },
		},
		"s": {
			Desc:   "Seconds (0-59)",
			expand: func(dt time.Time, locale locales.Translator) string { return strconv.Itoa(dt.Second()) },
		},
		"ss": {
			Desc:   "Seconds zero padded to two digits (00-59)",
			expand: func(dt time.Time, locale locales.Translator) string { return fmt.Sprintf("%02d", dt.Second()) },
		},
		"X": {
			Desc:   "Unix timestamp in seconds",
			expand: func(dt time.Time, locale locales.Translator) string { return strconv.Itoa(int(dt.Unix())) },
		},
		"x": {
			Desc:   "Unix timestamp in milliseconds",
			expand: func(dt time.Time, locale locales.Translator) string { return strconv.Itoa(int(dt.UnixMilli())) },
		},
		"z": {
			Desc: "Abbreviated time zone offset - 'GMT', 'CEST', '+0530'",
			expand: func(dt time.Time, locale locales.Translator) string {
				offsetName, _ := dt.Zone()
				return offsetName
			},
			aliases: []string{"zz"},
		},
		"Z": {
			Desc: "Time zone offset - '+05:30', '-03:00'",
			expand: func(dt time.Time, locale locales.Translator) string {
				_, offsetSeconds := dt.Zone()
				offsetMinutes := offsetSeconds / 60
				offsetHours := offsetMinutes / 60
				return fmt.Sprintf("%+03d:%02d", offsetHours, offsetMinutes%60)
			},
		},
		"ZZ": {
			Desc: "Time zone offset formatted without the dividing ':' - '+0530', '-0300'",
			expand: func(dt time.Time, locale locales.Translator) string {
				_, offsetSeconds := dt.Zone()
				offsetMinutes := offsetSeconds / 60
				offsetHours := offsetMinutes / 60
				return fmt.Sprintf("%+03d%02d", offsetHours, offsetMinutes%60)
			},
		},
	},
}
