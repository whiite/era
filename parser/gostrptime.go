package parser

import (
	"fmt"
	"strconv"
	"time"

	"gitlab.com/monokuro/era/dateutils"

	"github.com/go-playground/locales"
)

// Go port of the C time functions `strptime` and `strftime`
//
// Most formatting tokens are supported for English locales and are near one to one with the
// C equivalents with some small differences when it comes to spacing between characters
//
// Parsing is very much a work in progress only supporting a very limited set.
//
// Due to these incompatibilities it's recommend to use the CStr parser if possible
// for the time being
var GoStrptime DateHandlerPrefix

func init() {
	mapExpanded := expandTokenMap(&tokenMapStrftime)
	GoStrptime = DateHandlerPrefix{
		Prefix:     '%',
		tokenDef:   tokenMapStrftime,
		tokenGraph: createTokenGraph(&mapExpanded),
	}
}

// Formats date strings via the same system as `strftime`
var tokenMapStrftime = TokenMap{
	"%": {
		Desc:   "'%' character literal",
		expand: func(dt time.Time, locale locales.Translator) string { return "%" },
	},
	"A": {
		Desc: "Weekday name - 'Monday', 'Tuesday'",
		expand: func(dt time.Time, locale locales.Translator) string {
			return locale.WeekdayWide(dt.Weekday())
		},
	},
	"a": {
		Desc: "Weekday name truncated to three characters - 'Mon', 'Tue'",
		expand: func(dt time.Time, locale locales.Translator) string {
			return locale.WeekdayAbbreviated(dt.Weekday())
		},
	},
	"B": {
		Desc:   "Month name - 'January', 'February'",
		expand: func(dt time.Time, locale locales.Translator) string { return dt.Month().String() },
	},
	"b": {
		Desc:    "Month month name truncated to three characters - 'Jan', 'Feb'",
		expand:  func(dt time.Time, locale locales.Translator) string { return dt.Month().String()[:3] },
		aliases: []string{"h"},
	},
	"c": {
		Desc:   "Date and time for the current locale (hardcoded to UK format currently)",
		expand: func(dt time.Time, locale locales.Translator) string { return dt.Format("Mon _2 Jan 15:04:05 2006") },
	},
	"C": {
		Desc:   "The century number (0–99)",
		expand: func(dt time.Time, locale locales.Translator) string { return strconv.Itoa(dt.Year() / 100) },
	},
	"d": {
		Desc:   "Day of month zero padded to two digits (01-31)",
		expand: func(dt time.Time, locale locales.Translator) string { return fmt.Sprintf("%02d", dt.Day()) },
		parse: func(dt *time.Time, str *string) (int, error) {
			date, err := strconv.Atoi((*str)[:2])
			if err != nil {
				return 0, fmt.Errorf("Unable to parse date")
			}
			*dt = time.Date(dt.Year(), dt.Month(), date, dt.Hour(), dt.Minute(), dt.Second(), dt.Nanosecond(), dt.Location())
			return 1, nil
		},
	},
	"D": {
		Desc: "American style date (month first) equivalent to '%m/%d/%y' where the year is truncated to the last two digits - '01/31/97', '02/28/01'",
		expand: func(dt time.Time, locale locales.Translator) string {
			return fmt.Sprintf("%02d/%02d/%02d", dt.Month(), dt.Day(), dt.Year()%100)
		},
	},
	"e": {
		Desc: "Day of month space padded to two characters ( 1-31)",
		expand: func(dt time.Time, locale locales.Translator) string {
			fmtstr := " %d"
			if dt.Day() > 9 {
				fmtstr = fmtstr[1:]
			}
			return fmt.Sprintf(fmtstr, dt.Day())
		},
		parse: func(dt *time.Time, str *string) (int, error) {
			// TODO:
			date, err := strconv.Atoi((*str)[:2])
			if err != nil {
				return 0, fmt.Errorf("Unable to parse date")
			}
			*dt = time.Date(dt.Year(), dt.Month(), date, dt.Hour(), dt.Minute(), dt.Second(), dt.Nanosecond(), dt.Location())
			return 1, nil
		},
	},
	"F": {
		Desc: "Date in year-month-day format equivalent to '%Y-%m-%d' - '2024-01-04', '1997-10-31'",
		expand: func(dt time.Time, locale locales.Translator) string {
			return fmt.Sprintf("%d-%02d-%02d", dt.Year(), dt.Month(), dt.Day())
		},
	},
	"g": {
		Desc: "ISO week year shortened to the last two digits (00-99) ",
		expand: func(dt time.Time, locale locales.Translator) string {
			year, _ := dt.ISOWeek()
			return fmt.Sprintf("%02d", year%100)
		},
	},
	"G": {
		Desc: "ISO week year - '1999', '2007'",
		expand: func(dt time.Time, locale locales.Translator) string {
			year, _ := dt.ISOWeek()
			return fmt.Sprintf("%d", year)
		},
	},
	"H": {
		Desc:   "Hour in 24 hour format zero padded to two digits (00-23)",
		expand: func(dt time.Time, locale locales.Translator) string { return fmt.Sprintf("%02d", dt.Hour()) },
	},
	"I": {
		Desc: "Hour in 12 hour format zero padded to two digits (01-12)",
		expand: func(dt time.Time, locale locales.Translator) string {
			hour := dt.Hour() % 12
			if hour == 0 {
				hour = 12
			}
			return fmt.Sprintf("%02d", hour)
		},
	},
	"j": {
		Desc:   "Day of year zero padded to three digits (001-366)",
		expand: func(dt time.Time, locale locales.Translator) string { return fmt.Sprintf("%03d", dt.YearDay()) },
	},
	"k": {
		Desc: "Hour in 24 hour format space padded to two digits ( 0-23)",
		expand: func(dt time.Time, locale locales.Translator) string {
			hour := dt.Hour()
			if hour > 9 {
				return fmt.Sprintf("%2d", hour)
			}
			return fmt.Sprintf(" %d", hour)
		},
	},
	"l": {
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
	"m": {
		Desc:   "Month number zero padded to two digits (01-12)",
		expand: func(dt time.Time, locale locales.Translator) string { return fmt.Sprintf("%02d", dt.Month()) },
		parse: func(dt *time.Time, str *string) (int, error) {
			month, err := strconv.Atoi((*str)[:2])
			if err != nil {
				return 0, fmt.Errorf("Unable to parse month")
			}
			*dt = time.Date(dt.Year(), time.Month(month), dt.Day(), dt.Hour(), dt.Minute(), dt.Second(), dt.Nanosecond(), dt.Location())
			return 1, nil
		},
	},
	"M": {
		Desc:   "Minutes zero padded to two digits (00-59)",
		expand: func(dt time.Time, locale locales.Translator) string { return fmt.Sprintf("%02d", dt.Minute()) },
	},
	"n": {
		Desc:   "Newline whitespace - '\\n'",
		expand: func(dt time.Time, locale locales.Translator) string { return "\n" },
	},
	"p": {
		Desc: "The locale's equivalent of AM or PM (hardcoded to English am/pm)",
		expand: func(dt time.Time, locale locales.Translator) string {
			// TODO: locale
			if dt.Hour() < 12 {
				return "am"
			}
			return "pm"
		},
	},
	"r": {
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
	"R": {
		Desc: "Time represented as hours and minutes equivalent to %H:%M - '12:24', '04:09'",
		expand: func(dt time.Time, locale locales.Translator) string {
			return fmt.Sprintf("%02d:%02d", dt.Hour(), dt.Minute())
		},
	},
	"s": {
		Desc:   "Seconds since the unix epoch 1970-01-01 00:00:00 +0000 (UTC)",
		expand: func(dt time.Time, locale locales.Translator) string { return fmt.Sprintf("%d", dt.Unix()) },
	},
	"S": {
		Desc:   "Seconds zero padded to two digits (00-60; 60 may occur for for leap seconds)",
		expand: func(dt time.Time, locale locales.Translator) string { return fmt.Sprintf("%02d", dt.Second()) },
	},
	"t": {
		Desc:   "Tab whitespace - '\\t'",
		expand: func(dt time.Time, locale locales.Translator) string { return "\t" },
	},
	"T": {
		Desc: "Time represented as hours, minutes and seconds equivalent to %H:%M:%S - '12:34:03', '04:09:59'",
		expand: func(dt time.Time, locale locales.Translator) string {
			return fmt.Sprintf("%02d:%02d:%02d", dt.Hour(), dt.Minute(), dt.Second())
		},
	},
	"u": {
		Desc: "Day of week where Monday = 1 and Sunday = 7 (1-7)",
		expand: func(dt time.Time, locale locales.Translator) string {
			return strconv.Itoa((int(dt.Weekday())+6)%7 + 1)
		},
	},
	"U": {
		Desc: "Week number of the year where the first Sunday of January is considered week 1 - (00-53)",
		expand: func(dt time.Time, locale locales.Translator) string {
			midnight := dateutils.DayStart(dt)
			jan1st := dateutils.YearStart(midnight)

			firstSundayJan := dateutils.NextWeekday(time.Sunday, jan1st)
			weekDiff := int((midnight.Sub(firstSundayJan).Hours() / 24 / 7) + 1)
			return fmt.Sprintf("%02d", weekDiff)
		},
	},
	"v": {
		Desc: "Date with space padded day; truncated month name and year equivalent to %d-%b-%Y - ' 4-Jan-1997'",
		expand: func(dt time.Time, locale locales.Translator) string {
			day := dt.Day()
			formatStr := " %d-%s-%d"
			if day > 9 {
				formatStr = formatStr[1:]
			}
			return fmt.Sprintf(formatStr, dt.Day(), locale.MonthWide(dt.Month())[:3], dt.Year())
		},
	},
	"V": {
		Desc: "ISO8601 week number of the year zero padded to two digits - (01-53)",
		expand: func(dt time.Time, locale locales.Translator) string {
			_, week := dt.ISOWeek()
			return fmt.Sprintf("%02d", week)
		},
	},
	"w": {
		Desc:   "Day of week number (0-6) where Sunday is 0 and Saturday is 6",
		expand: func(dt time.Time, locale locales.Translator) string { return strconv.Itoa(int(dt.Weekday())) },
	},
	"W": {
		Desc: "Week number of the year where the first Monday of January is considered week 1 - (00-53)",
		expand: func(dt time.Time, locale locales.Translator) string {
			midnight := dateutils.DayStart(dt)
			jan1st := dateutils.YearStart(midnight)

			firstSundayJan := dateutils.NextWeekday(time.Monday, jan1st)
			weekDiff := int((midnight.Sub(firstSundayJan).Hours() / 24 / 7) + 1)
			return fmt.Sprintf("%02d", weekDiff)
		},
	},
	"x": {
		Desc: "Locale date format - '04/12/1999', '11/02/2007'",
		expand: func(dt time.Time, locale locales.Translator) string {
			return locale.FmtDateShort(dt)
		},
	},
	"X": {
		Desc: "Locale time including seconds - '03:57:22', '18:08:01'",
		expand: func(dt time.Time, locale locales.Translator) string {
			return locale.FmtTimeMedium(dt)
		},
	},
	"y": {
		Desc:   "The year within the century zero padded to two digits (00-99)",
		expand: func(dt time.Time, locale locales.Translator) string { return fmt.Sprintf("%02d", dt.Year()%100) },
		parse: func(dt *time.Time, str *string) (int, error) {
			year, err := strconv.Atoi((*str)[:2])
			if err != nil {
				return 0, fmt.Errorf("Unable to parse year")
			}
			if year > time.Now().Year()%100 {
				year += 1900
			} else {
				year += 2000
			}
			*dt = time.Date(year, dt.Month(), dt.Day(), dt.Hour(), dt.Minute(), dt.Second(), dt.Nanosecond(), dt.Location())
			return 1, nil
		},
	},
	"Y": {
		Desc:   "Year number - '1999', '2007'",
		expand: func(dt time.Time, locale locales.Translator) string { return strconv.Itoa(dt.Year()) },
	},
	"z": {
		Desc: "Time zone offset in +hhmm format - '-0400', '+0530'",
		expand: func(dt time.Time, locale locales.Translator) string {
			_, offsetSeconds := dt.Zone()
			offsetMinutes := offsetSeconds / 60
			offsetHours := offsetMinutes / 60
			return fmt.Sprintf("%+03d%02d", offsetHours, offsetMinutes%60)
		},
	},
	"Z": {
		Desc: "Abbreviated time zone offset - 'GMT', 'CEST', '+0530'",
		expand: func(dt time.Time, locale locales.Translator) string {
			offsetName, _ := dt.Zone()
			return offsetName
		},
	},
	"Ec": {
		Desc:   "Alternative representation for date and time for the current locale (hardcoded to UK format currently)",
		expand: func(dt time.Time, locale locales.Translator) string { return dt.Format("Mon _2 Jan 15:04:05 2006") },
	},
	"EC": {
		Desc: "Base year/period - (0-99)",
		expand: func(dt time.Time, locale locales.Translator) string {
			return strconv.Itoa(dt.Year() / 100)
		},
	},
	"Ex": {
		Desc: "Short date format in the specified locale",
		expand: func(dt time.Time, locale locales.Translator) string {
			return locale.FmtDateShort(dt)
		},
	},
	"EX": {
		Desc: "Time format in the specified locale",
		expand: func(dt time.Time, locale locales.Translator) string {
			return locale.FmtTimeMedium(dt)
		},
	},
	"Ey": {
		Desc: "Year to two digits - '97', '07'",
		expand: func(dt time.Time, locale locales.Translator) string {
			return fmt.Sprintf("%02d", dt.Year()%100)
		},
	},
	"EY": {
		Desc: "Alternative year number - '1997', '2007'",
		expand: func(dt time.Time, locale locales.Translator) string {
			return strconv.Itoa(dt.Year())
		},
	},
	"Od": {
		Desc: "Day of the month using the locale's alternative numeric symbols, zero padded - (00-31)",
		expand: func(dt time.Time, locale locales.Translator) string {
			return fmt.Sprintf("%02d", dt.Day())
		},
	},
	"Oe": {
		Desc: "Day of the month using the locale's alternative numeric symbols, space padded - ( 0-31)",
		expand: func(dt time.Time, locale locales.Translator) string {
			fmtstr := " %d"
			if dt.Day() > 9 {
				fmtstr = fmtstr[1:]
			}
			return fmt.Sprintf(fmtstr, dt.Day())
		},
	},
	"OH": {
		Desc: "Hour in 24 hour format using the locale's alternative numeric symbols, zero padded - (00-23)",
		expand: func(dt time.Time, locale locales.Translator) string {
			return fmt.Sprintf("%02d", dt.Hour())
		},
	},
	"OI": {
		Desc: "Hour in 12 hour format using the locale's alternative numeric symbols, zero padded - (00-12)",
		expand: func(dt time.Time, locale locales.Translator) string {
			hour := dt.Hour() % 12
			if hour == 0 {
				hour = 12
			}
			return fmt.Sprintf("%02d", hour)
		},
	},
	"Om": {
		Desc: "Month number using the locale's alternative numeric symbols, zero padded (01-12)",
		expand: func(dt time.Time, locale locales.Translator) string {
			return fmt.Sprintf("%02d", dt.Month())
		},
	},
	"OM": {
		Desc:   "Minutes using the locale's alternative numeric symbols, zero padded (00-59)",
		expand: func(dt time.Time, locale locales.Translator) string { return fmt.Sprintf("%02d", dt.Minute()) },
	},
	"OS": {
		Desc:   "Seconds using the locale's alternative numeric symbols, zero padded (00-60; 60 may occur for for leap seconds)",
		expand: func(dt time.Time, locale locales.Translator) string { return fmt.Sprintf("%02d", dt.Second()) },
	},
	"OU": {
		Desc: "Week number of the year using the locale's alternative numeric symbols where the first Sunday of January is considered week 1 - (00-53)",
		expand: func(dt time.Time, locale locales.Translator) string {
			midnight := dateutils.DayStart(dt)
			jan1st := dateutils.YearStart(midnight)

			firstSundayJan := dateutils.NextWeekday(time.Sunday, jan1st)
			weekDiff := int((midnight.Sub(firstSundayJan).Hours() / 24 / 7) + 1)
			return fmt.Sprintf("%02d", weekDiff)
		},
	},
	"Ow": {
		Desc:   "Day of week number (0-6) using the locale's alternative numeric symbols where Sunday is 0 and Saturday is 6",
		expand: func(dt time.Time, locale locales.Translator) string { return strconv.Itoa(int(dt.Weekday())) },
	},
	"OW": {
		Desc: "Week number of the year using the locale's alternative numeric symbols where the first Monday of January is considered week 1 - (00-53)",
		expand: func(dt time.Time, locale locales.Translator) string {
			midnight := dateutils.DayStart(dt)
			jan1st := dateutils.YearStart(midnight)

			firstSundayJan := dateutils.NextWeekday(time.Monday, jan1st)
			weekDiff := int((midnight.Sub(firstSundayJan).Hours() / 24 / 7) + 1)
			return fmt.Sprintf("%02d", weekDiff)
		},
	},
	"Oy": {
		Desc: "Year number offset from the century using the locale's alternative numeric symbols - '99', '07'",
		expand: func(dt time.Time, locale locales.Translator) string {
			return fmt.Sprintf("%02d", dt.Year()%100)
		},
	},
}
