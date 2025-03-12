package parser

import (
	"fmt"
	"gitlab.com/monokuro/era/dateutils"
	"strconv"
	"time"

	"github.com/go-playground/locales"
)

var Luxon DateFormatterString

func init() {
	mapExpanded := expandTokenMap(&tokenMapLuxon)
	Luxon = DateFormatterString{
		escapeChars: []rune{'\''},
		tokenDef:    tokenMapLuxon,
		tokenGraph:  createTokenGraph(&mapExpanded),
	}
}

// TODO: missing tokens:
// - "ZZZZZ" - full offset name e.g. Eastern Standard Time
// - "TTTT" - localised 24 hour time with full time zone name
// - "f", "ff", "fff", "ffff" - localised date and time
// - "F", "FF", "FFF", "FFFF" - localised date and time with seconds

// Formats date strings via the same system as `strftime`
var tokenMapLuxon = TokenMap{
	"a": {
		Desc: "Meridiem - 'AM'",
		expand: func(dt time.Time, locale locales.Translator) string {
			// TODO: locale
			if dt.Hour() < 12 {
				return "AM"
			}
			return "PM"
		},
	},
	"c": {
		Desc: "Day of week where Monday = 1 and Sunday = 7 (1-7)",
		expand: func(dt time.Time, locale locales.Translator) string {
			return strconv.Itoa((int(dt.Weekday())+6)%7 + 1)
		},
		aliases: []string{"E"},
	},
	"ccc": {
		Desc: "Day of week name truncated to three characters - 'Sun', 'Mon'",
		expand: func(dt time.Time, locale locales.Translator) string {
			return locale.WeekdayWide(dt.Weekday())[:3]
		},
		aliases: []string{"EEE"},
	},
	"cccc": {
		Desc: "Day of week name - 'Sunday', 'Monday'",
		expand: func(dt time.Time, locale locales.Translator) string {
			return locale.WeekdayWide(dt.Weekday())
		},
		aliases: []string{"EEEE"},
	},
	"ccccc": {
		Desc: "Day of week name truncated to one character - 'S', 'M'",
		expand: func(dt time.Time, locale locales.Translator) string {
			return locale.WeekdayNarrow(dt.Weekday())
		},
		aliases: []string{"EEEEE"},
	},
	"d": {
		Desc:   "Day of month (1-31)",
		expand: func(dt time.Time, locale locales.Translator) string { return strconv.Itoa(dt.Day()) },
	},
	"dd": {
		Desc:   "Day of month zero padded to two digits (01-31)",
		expand: func(dt time.Time, locale locales.Translator) string { return fmt.Sprintf("%02d", dt.Day()) },
	},
	"D": {
		Desc: "Localised numerical date - '08/11/24'",
		expand: func(dt time.Time, locale locales.Translator) string {
			return locale.FmtDateShort(dt)
		},
	},
	"DD": {
		Desc: "Localised date with abbreviated month name - 'Nov 8, 2024'",
		expand: func(dt time.Time, locale locales.Translator) string {
			return locale.FmtDateMedium(dt)
		},
	},
	"DDD": {
		Desc: "Localised date with month name - 'November 8, 2024'",
		expand: func(dt time.Time, locale locales.Translator) string {
			return locale.FmtDateLong(dt)
		},
	},
	"DDDD": {
		Desc: "Localised date with weekday and month name - 'Friday, November 8, 2024'",
		expand: func(dt time.Time, locale locales.Translator) string {
			return locale.FmtDateFull(dt)
		},
	},
	"G": {
		Desc: "Era name abbreviated - 'BC', 'AD'",
		expand: func(dt time.Time, locale locales.Translator) string {
			// TODO: locale
			if dt.Year() < 0 {
				return "BC"
			}
			return "AD"
		},
	},
	"GG": {
		Desc: "Era name abbreviated - 'Before Christ', 'Anno Domini'",
		expand: func(dt time.Time, locale locales.Translator) string {
			// TODO: locale
			if dt.Year() < 0 {
				return "Before Christ"
			}
			return "Anno Domini"
		},
	},
	"GGGGG": {
		Desc: "Era name abbreviated to one character - 'B', 'A'",
		expand: func(dt time.Time, locale locales.Translator) string {
			// TODO: locale
			if dt.Year() < 0 {
				return "B"
			}
			return "A"
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
	"kk": {
		Desc: "ISO week year shortened to the last two digits - '99, '07'",
		expand: func(dt time.Time, locale locales.Translator) string {
			year, _ := dt.ISOWeek()
			return fmt.Sprintf("%02d", year%100)
		},
	},
	"kkkk": {
		Desc: "ISO week year zero padded to four digits - '1999', '2007'",
		expand: func(dt time.Time, locale locales.Translator) string {
			year, _ := dt.ISOWeek()
			return fmt.Sprintf("%04d", year)
		},
	},
	"L": {
		Desc:    "Month number (1-12)",
		expand:  func(dt time.Time, locale locales.Translator) string { return strconv.Itoa(int(dt.Month())) },
		aliases: []string{"M"},
	},
	"LL": {
		Desc:    "Month number zero padded to two digits - (01-12)",
		expand:  func(dt time.Time, locale locales.Translator) string { return fmt.Sprintf("%02d", dt.Month()) },
		aliases: []string{"MM"},
	},
	"LLL": {
		Desc: "Month name truncated to three characters - 'Jan', 'Feb'",
		expand: func(dt time.Time, locale locales.Translator) string {
			return locale.MonthWide(dt.Month())[:3]
		},
		aliases: []string{"MMM"},
	},
	"LLLL": {
		Desc: "Month name - 'January', 'February'",
		expand: func(dt time.Time, locale locales.Translator) string {
			return locale.MonthWide(dt.Month())
		},
		aliases: []string{"MMMM"},
	},
	"LLLLL": {
		Desc: "Month name truncated to one character - 'J', 'F'",
		expand: func(dt time.Time, locale locales.Translator) string {
			return locale.MonthNarrow(dt.Month())
		},
		aliases: []string{"MMMMM"},
	},
	"m": {
		Desc:   "Minutes (0-59)",
		expand: func(dt time.Time, locale locales.Translator) string { return strconv.Itoa(dt.Minute()) },
	},
	"mm": {
		Desc:   "Minutes zero padded to two digits (00-59)",
		expand: func(dt time.Time, locale locales.Translator) string { return fmt.Sprintf("%02d", dt.Minute()) },
	},
	"n": {
		Desc: "Week of year where the week containing January 1st is considered week one (1-53)",
		expand: func(dt time.Time, locale locales.Translator) string {
			midnight := dateutils.DayStart(dt)
			jan1st := dateutils.YearStart(midnight)
			weeksInYear := dateutils.YearEnd(dt).Sub(jan1st).Hours() / 24 / 7

			lastSunday := dateutils.WeekEnd(time.Sunday, jan1st)
			weekDiff := midnight.Sub(lastSunday).Hours() / 24 / 7
			return strconv.Itoa(int(weekDiff)%int(weeksInYear) + 1)
		},
	},
	"nn": {
		Desc: "Week of year where the week containing January 1st is considered week one, zero padded to two digits (01-53)",
		expand: func(dt time.Time, locale locales.Translator) string {
			midnight := dateutils.DayStart(dt)
			jan1st := dateutils.YearStart(midnight)
			weeksInYear := dateutils.YearEnd(dt).Sub(jan1st).Hours() / 24 / 7

			lastSunday := dateutils.WeekEnd(time.Sunday, jan1st)
			weekDiff := midnight.Sub(lastSunday).Hours() / 24 / 7
			return fmt.Sprintf("%02d", int(weekDiff)%int(weeksInYear)+1)
		},
	},
	"o": {
		Desc:   "Ordinal day of year (1-366)",
		expand: func(dt time.Time, locale locales.Translator) string { return strconv.Itoa(dt.YearDay()) },
	},
	"ooo": {
		Desc:   "Ordinal day of year zero padded to three digits (001-366)",
		expand: func(dt time.Time, locale locales.Translator) string { return fmt.Sprintf("%03d", dt.YearDay()) },
	},
	"q": {
		Desc: "Quarter of year (1-4)",
		expand: func(dt time.Time, locale locales.Translator) string {
			return strconv.Itoa(dateutils.YearQuarter(dt))
		},
	},
	"qq": {
		Desc: "Quarter of year zero padded to two digits (01-04)",
		expand: func(dt time.Time, locale locales.Translator) string {
			return fmt.Sprintf("%02d", dateutils.YearQuarter(dt))
		},
	},
	"s": {
		Desc:   "Seconds (0-59)",
		expand: func(dt time.Time, locale locales.Translator) string { return strconv.Itoa(dt.Second()) },
	},
	"ss": {
		Desc:   "Seconds zero padded to two digits (00-59)",
		expand: func(dt time.Time, locale locales.Translator) string { return fmt.Sprintf("%02d", dt.Second()) },
	},
	"S": {
		Desc:   "Milliseconds (0-999)",
		expand: func(dt time.Time, locale locales.Translator) string { return strconv.Itoa(dt.Nanosecond() / 1_000_000) },
	},
	"SSS": {
		Desc: "Milliseconds zero padded to three digits (000-999)",
		expand: func(dt time.Time, locale locales.Translator) string {
			return fmt.Sprintf("%03d", dt.Nanosecond()/1_000_000)
		},
		aliases: []string{"u"},
	},
	"t": {
		Desc: "Localised time - '9:07 AM'",
		expand: func(dt time.Time, locale locales.Translator) string {
			return locale.FmtTimeShort(dt)
		},
	},
	"tt": {
		Desc: "Localised time with seconds - '9:07:53 AM'",
		expand: func(dt time.Time, locale locales.Translator) string {
			return locale.FmtTimeMedium(dt)
		},
	},
	"ttt": {
		Desc: "Localised time with seconds and abbreviated offset name - '9:07:53 AM EDT'",
		expand: func(dt time.Time, locale locales.Translator) string {
			return locale.FmtTimeLong(dt)
		},
	},
	"tttt": {
		Desc: "Localised time with seconds and offset name - '9:07:53 AM Eastern Daylight Time'",
		expand: func(dt time.Time, locale locales.Translator) string {
			return locale.FmtTimeFull(dt)
		},
	},
	"T": {
		Desc: "Localised 24 hour time - '13:07'",
		expand: func(dt time.Time, locale locales.Translator) string {
			return fmt.Sprintf("%02d:%02d", dt.Hour(), dt.Minute())
		},
	},
	"TT": {
		Desc: "Localised 24 hour time with seconds - '13:07'",
		expand: func(dt time.Time, locale locales.Translator) string {
			return fmt.Sprintf("%02d:%02d:%02d", dt.Hour(), dt.Minute(), dt.Second())
		},
	},
	"TTT": {
		Desc: "Localised 24 hour time with seconds and abbreviated offset - '13:07 CST'",
		expand: func(dt time.Time, locale locales.Translator) string {
			offsetName, _ := dt.Zone()
			return fmt.Sprintf("%d:%02d:%02d %s", dt.Hour(), dt.Minute(), dt.Second(), offsetName)
		},
	},
	"W": {
		Desc: "ISO week (1-53)",
		expand: func(dt time.Time, locale locales.Translator) string {
			_, week := dt.ISOWeek()
			return strconv.Itoa(week)
		},
		aliases: []string{"n"},
	},
	"WW": {
		Desc: "ISO week zero padded to two digits (01-53)",
		expand: func(dt time.Time, locale locales.Translator) string {
			_, week := dt.ISOWeek()
			return fmt.Sprintf("%02d", week)
		},
		aliases: []string{"nn"},
	},
	"uu": {
		Desc: "Fractional seconds zero padded to two digits (00-99)",
		expand: func(dt time.Time, locale locales.Translator) string {
			return fmt.Sprintf("%02d", dt.Nanosecond()/10_000_000)
		},
	},
	"uuu": {
		Desc: "Fractional seconds between 0 and 9 (0-9)",
		expand: func(dt time.Time, locale locales.Translator) string {
			return fmt.Sprintf("%d", dt.Nanosecond()/100_000_000)
		},
	},
	"X": {
		Desc:   "Unix timestamp in seconds",
		expand: func(dt time.Time, locale locales.Translator) string { return strconv.Itoa(int(dt.Unix())) },
	},
	"x": {
		Desc:   "Unix timestamp in milliseconds",
		expand: func(dt time.Time, locale locales.Translator) string { return strconv.Itoa(int(dt.UnixMilli())) },
	},
	"y": {
		Desc:   "Year number - '1999', '2007'",
		expand: func(dt time.Time, locale locales.Translator) string { return strconv.Itoa(dt.Year()) },
	},
	"yy": {
		Desc:    "Year number truncated to last two digits - '99', '07'",
		expand:  func(dt time.Time, locale locales.Translator) string { return fmt.Sprintf("%02d", dt.Year()%100) },
		aliases: []string{"ii"},
	},
	"yyyy": {
		Desc:    "Year number zero padded to four digits - '1999', '0007'",
		expand:  func(dt time.Time, locale locales.Translator) string { return fmt.Sprintf("%04d", dt.Year()) },
		aliases: []string{"iiii"},
	},
	"z": {
		Desc: "IANA canonical time zone string - 'Europe/London'",
		expand: func(dt time.Time, locale locales.Translator) string {
			return dt.Location().String()
		},
	},
	"Z": {
		Desc: "Time zone offset shortened to one digit - '+5', '-3'",
		expand: func(dt time.Time, locale locales.Translator) string {
			_, offsetSeconds := dt.Zone()
			offsetHours := offsetSeconds / (60 * 60)
			return fmt.Sprintf("%+d", offsetHours)
		},
	},
	"ZZ": {
		Desc: "Time zone offset - '+05:30', '-03:00'",
		expand: func(dt time.Time, locale locales.Translator) string {
			_, offsetSeconds := dt.Zone()
			offsetMinutes := offsetSeconds / 60
			offsetHours := offsetMinutes / 60
			return fmt.Sprintf("%+03d:%02d", offsetHours, offsetMinutes%60)
		},
	},
	"ZZZ": {
		Desc: "Time zone offset formatted without the dividing ':' - '+0530', '-0300'",
		expand: func(dt time.Time, locale locales.Translator) string {
			_, offsetSeconds := dt.Zone()
			offsetMinutes := offsetSeconds / 60
			offsetHours := offsetMinutes / 60
			return fmt.Sprintf("%+03d%02d", offsetHours, offsetMinutes%60)
		},
	},
	"ZZZZ": {
		Desc: "Abbreviated time zone offset - 'GMT', 'CEST', '+0530'",
		expand: func(dt time.Time, locale locales.Translator) string {
			offsetName, _ := dt.Zone()
			return offsetName
		},
	},
}
