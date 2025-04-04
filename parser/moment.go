package parser

import (
	"fmt"
	"gitlab.com/monokuro/era/dateutils"
	"strconv"
	"time"

	"github.com/go-playground/locales"
)

// Handler for parsing and formatting using the `momentjs` tokens and rules
var MomentJs DateHandlerString

func init() {
	mapExpanded := expandTokenMap(&tokenMapMoment)
	MomentJs = DateHandlerString{
		escapeChars: []rune{'[', ']'},
		tokenDef:    tokenMapMoment,
		tokenGraph:  createTokenGraph(&mapExpanded),
	}
}

var tokenMapMoment = TokenMap{
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
	"N": {
		Desc: "Era name abbreviated - 'BC', 'AD'",
		expand: func(dt time.Time, locale locales.Translator) string {
			if dt.Year() < 0 {
				return "BC"
			}
			return "AD"
		},
		aliases: []string{"NN", "NNN", "NNNNN"},
	},
	"NNNN": {
		Desc: "Era name in full - 'Before Christ', 'Anno Domini'",
		expand: func(dt time.Time, locale locales.Translator) string {
			if dt.Year() < 0 {
				return "Before Christ"
			}
			return "Anno Domini"
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
	"gg": {
		Desc: "Year of week where the last Sunday of the current week is used, truncated and zero padded to the last two digits - '97', '07'",
		expand: func(dt time.Time, locale locales.Translator) string {
			weekEnd := dateutils.WeekEnd(time.Sunday, dt)
			return fmt.Sprintf("%02d", weekEnd.Year()%100)

		},
	},
	"gggg": {
		Desc: "Year of week where the last Sunday of the current week is used - '1997', '2007'",
		expand: func(dt time.Time, locale locales.Translator) string {
			weekEnd := dateutils.WeekEnd(time.Sunday, dt)
			return strconv.Itoa(weekEnd.Year())

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
	"w": {
		Desc: "Week of year where the first Sunday before January 1st is considered week one (1-53)",
		expand: func(dt time.Time, locale locales.Translator) string {
			midnight := dateutils.DayStart(dt)
			jan1st := dateutils.YearStart(midnight)
			weeksInYear := dateutils.YearEnd(dt).Sub(jan1st).Hours() / 24 / 7

			firstSunday := dateutils.WeekStart(time.Sunday, jan1st)
			weekDiff := midnight.Sub(firstSunday).Hours() / 24 / 7
			return strconv.Itoa(int(weekDiff)%int(weeksInYear) + 1)
		},
	},
	"wo": {
		Desc: "Week of year suffixed where the first Sunday before January 1st is considered week one (1st-53rd)",
		expand: func(dt time.Time, locale locales.Translator) string {
			midnight := dateutils.DayStart(dt)
			jan1st := dateutils.YearStart(midnight)
			weeksInYear := dateutils.YearEnd(dt).Sub(jan1st).Hours() / 24 / 7

			firstSunday := dateutils.WeekStart(time.Sunday, jan1st)
			weekDiff := midnight.Sub(firstSunday).Hours() / 24 / 7
			return numberSuffixed(int(weekDiff)%int(weeksInYear) + 1)
		},
	},
	"ww": {
		Desc: "Week of year where the first Sunday before January 1st is considered week one padded to two digits (01-53)",
		expand: func(dt time.Time, locale locales.Translator) string {
			midnight := dateutils.DayStart(dt)
			jan1st := dateutils.YearStart(midnight)
			weeksInYear := dateutils.YearEnd(dt).Sub(jan1st).Hours() / 24 / 7

			firstSunday := dateutils.WeekStart(time.Sunday, jan1st)
			weekDiff := midnight.Sub(firstSunday).Hours() / 24 / 7
			return fmt.Sprintf("%02d", int(weekDiff)%int(weeksInYear)+1)
		},
	},
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
		Desc:    "Year number - '1999', '2007'",
		expand:  func(dt time.Time, locale locales.Translator) string { return strconv.Itoa(dt.Year()) },
		aliases: []string{"y"},
	},
	"YY": {
		Desc: "Year number truncated to last two digits - '99', '07'",
		expand: func(dt time.Time, locale locales.Translator) string {
			return fmt.Sprintf("%02d", dt.Year()%100)
		},
		aliases: []string{"GG"},
	},
	"YYYY": {
		Desc:    "Year number - '1999', '2007'",
		expand:  func(dt time.Time, locale locales.Translator) string { return strconv.Itoa(dt.Year()) },
		aliases: []string{"GGGG"},
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
	"S": {
		Desc: "Fractional seconds to one digit (0-9)",
		expand: func(dt time.Time, locale locales.Translator) string {
			return strconv.Itoa(dt.Nanosecond() / 100_000_000)
		},
	},
	"SS": {
		Desc: "Fractional seconds to two digits (00-99)",
		expand: func(dt time.Time, locale locales.Translator) string {
			return fmt.Sprintf("%02d", (dt.Nanosecond() / 10_000_000))
		},
	},
	"SSS": {
		Desc: "Fractional seconds to three digits (000-999)",
		expand: func(dt time.Time, locale locales.Translator) string {
			return fmt.Sprintf("%03d", (dt.Nanosecond() / 1_000_000))
		},
	},
	"SSSS": {
		Desc: "Fractional seconds to four digits (0000-9999)",
		expand: func(dt time.Time, locale locales.Translator) string {
			return fmt.Sprintf("%04d", (dt.Nanosecond() / 100_000))
		},
	},
	"SSSSS": {
		Desc: "Fractional seconds to five digits (00000-99999)",
		expand: func(dt time.Time, locale locales.Translator) string {
			return fmt.Sprintf("%05d", (dt.Nanosecond() / 10_000))
		},
	},
	"SSSSSS": {
		Desc: "Fractional seconds to six digits (000000-999999)",
		expand: func(dt time.Time, locale locales.Translator) string {
			return fmt.Sprintf("%06d", (dt.Nanosecond() / 1_000))
		},
	},
	"SSSSSSS": {
		Desc: "Fractional seconds to seven digits (0000000-9999999)",
		expand: func(dt time.Time, locale locales.Translator) string {
			return fmt.Sprintf("%07d", (dt.Nanosecond() / 100))
		},
	},
	"SSSSSSSS": {
		Desc: "Fractional seconds to eight digits (00000000-99999999)",
		expand: func(dt time.Time, locale locales.Translator) string {
			return fmt.Sprintf("%08d", dt.Nanosecond()/10)
		},
	},
	"SSSSSSSSS": {
		Desc: "Fractional seconds to eight digits (00000000-99999999)",
		expand: func(dt time.Time, locale locales.Translator) string {
			return fmt.Sprintf("%09d", dt.Nanosecond())
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
}
