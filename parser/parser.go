package parser

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type FormatToken[T any] struct {
	// Description of what the token represents
	Desc string
	// Equivalent string for token given a `time.Time`
	expand  func(dt time.Time) string
	aliases []T
}

type DateFormatterPrefix struct {
	Prefix   rune
	TokenMap map[rune]FormatToken[rune]
}

func (formatter DateFormatterPrefix) Parse(dt time.Time, str *string) string {
	var formattedDate strings.Builder
	interpretToken := false

	for _, char := range *str {
		if char == formatter.Prefix {
			if interpretToken {
				formattedDate.WriteRune(formatter.Prefix)
			}
			interpretToken = true
			continue
		}

		if token, hasToken := formatter.TokenMap[char]; interpretToken && hasToken {
			formattedDate.WriteString(token.expand(dt))
			interpretToken = false
			continue
		}

		formattedDate.WriteRune(char)
		interpretToken = false
	}

	return formattedDate.String()
}

type DateFormatterNoPrefix struct {
	escapeChars []rune
	TokenMap    map[string]FormatToken[string]
}

func (formatter DateFormatterNoPrefix) Parse(dt time.Time, str *string) string {
	var formattedDate strings.Builder
	var tokens strings.Builder

	tokenNodeRoot := createTokenGraph(&formatter.TokenMap)
	tokenNode := tokenNodeRoot
	escapeSupport := len(formatter.escapeChars) > 1
	escapeMode := false

	for _, char := range *str {
		if escapeSupport && char == formatter.escapeChars[0] {
			escapeMode = true
		}
		if escapeSupport && char == formatter.escapeChars[1] && escapeMode {
			escapeMode = false
		}

		if node, hasToken := tokenNode.children[char]; hasToken && !escapeMode {
			tokens.WriteRune(char)
			tokenNode = node
			continue
		}

		if formatFunc := tokenNode.value.expand; formatFunc != nil {
			formattedDate.WriteString(formatFunc(dt))
		} else {
			formattedDate.WriteString(tokens.String())
		}
		tokens.Reset()
		formattedDate.WriteRune(char)

		tokenNode = tokenNodeRoot
		if node, hasToken := tokenNodeRoot.children[char]; hasToken && !escapeMode {
			tokenNode = node
		}
	}

	if formatFunc := tokenNode.value.expand; formatFunc != nil {
		formattedDate.WriteString(formatFunc(dt))
	} else {
		formattedDate.WriteString(tokens.String())
	}

	return formattedDate.String()
}

func numberSuffixed(num int) string {
	numstr := strconv.Itoa(num)
	switch numstr {
	case "1", "21", "31":
		return numstr + "st"
	case "2", "22":
		return numstr + "nd"
	case "3", "23":
		return numstr + "rd"
	default:
		return numstr + "th"
	}
}

// Formats date strings via the same system as `strptime`
var Strptime = DateFormatterPrefix{
	Prefix: '%',
	TokenMap: map[rune]FormatToken[rune]{
		'%': {
			Desc:   "'%' character literal",
			expand: func(dt time.Time) string { return "%" },
		},
		'A': {
			Desc:   "Weekday name - 'Monday', 'Tuesday'",
			expand: func(dt time.Time) string { return dt.Weekday().String() },
		},
		'a': {
			Desc:   "Weekday name truncated to three characters - 'Mon', 'Tue'",
			expand: func(dt time.Time) string { return dt.Weekday().String()[:3] },
		},
		'B': {
			Desc:   "Month name - 'January', 'February'",
			expand: func(dt time.Time) string { return dt.Month().String() },
		},
		'b': {
			Desc:   "Month month name truncated to three characters - 'Jan', 'Feb'",
			expand: func(dt time.Time) string { return dt.Month().String()[:3] },
		},
		'h': {
			Desc:   "Month name truncated to three characters - 'Jan', 'Feb'",
			expand: func(dt time.Time) string { return dt.Month().String()[:3] },
		},
		'c': {
			Desc:   "Date and time for the current locale (different to strptime and hardcoded to UK format currently)",
			expand: func(dt time.Time) string { return dt.Format("Mon _2 Jan 15:04:05 2006") },
		},
		'C': {
			Desc:   "The century number (0–99)",
			expand: func(dt time.Time) string { return strconv.Itoa(dt.Year() / 100) },
		},
		'd': {
			Desc:   "Day of month zero padded to two digits (01-31)",
			expand: func(dt time.Time) string { return fmt.Sprintf("%02d", dt.Day()) },
		},
		'e': {
			Desc:   "Day of month (1-31)",
			expand: func(dt time.Time) string { return strconv.Itoa(dt.Day()) },
		},
		'D': {
			Desc: "American style date (month first) equivalent to '%m/%d/%y' where the year is truncated to the last two digits - '01/31/97', '02/28/01'",
			expand: func(dt time.Time) string {
				return fmt.Sprintf("%02d/%02d/%02d", dt.Month(), dt.Day(), dt.Year()%100)
			},
		},
		'H': {
			Desc:   "Hour in 24 hour format zero padded to two digits (00-23)",
			expand: func(dt time.Time) string { return fmt.Sprintf("%02d", dt.Hour()) },
		},
		'I': {
			Desc:   "Hour in 12 hour format zero padded to two digits (01-12)",
			expand: func(dt time.Time) string { return fmt.Sprintf("%02d", dt.Hour()%12) },
		},
		'j': {
			Desc:   "Day of year zero padded to three digits (001-366)",
			expand: func(dt time.Time) string { return fmt.Sprintf("%03d", dt.YearDay()) },
		},
		'm': {
			Desc:   "Month number zero padded to two digits (01-12)",
			expand: func(dt time.Time) string { return fmt.Sprintf("%02d", dt.Month()) },
		},
		'M': {
			Desc:   "Minutes zero padded to two digits (00-59)",
			expand: func(dt time.Time) string { return fmt.Sprintf("%02d", dt.Minute()) },
		},
		'n': {
			Desc:   "Newline whitespace - '\\n'",
			expand: func(dt time.Time) string { return "\n" },
		},
		'R': {
			Desc:   "Time represented as hours and minutes equivalent to %H:%M - '12:24', '04:09'",
			expand: func(dt time.Time) string { return fmt.Sprintf("%02d:%02d", dt.Hour(), dt.Minute()) },
		},
		'S': {
			Desc:   "Seconds zero padded to two digits (00-60; 60 may occur for for leap seconds)",
			expand: func(dt time.Time) string { return fmt.Sprintf("%02d", dt.Second()) },
		},
		't': {
			Desc:   "Tab whitespace - '\\t'",
			expand: func(dt time.Time) string { return "\t" },
		},
		'T': {
			Desc: "Time represented as hours, minutes and seconds equivalent to %H:%M:%S - '12:34:03', '04:09:59'",
			expand: func(dt time.Time) string {
				return fmt.Sprintf("%02d:%02d:%02d", dt.Hour(), dt.Minute(), dt.Second())
			},
		},
		// TODO: implement same rules as strptime - first week is the first Sunday of January
		'U': {
			Desc: "ISO8601 week number of the year (this is different to the typical strptime token currently)",
			expand: func(dt time.Time) string {
				_, weekNumber := dt.ISOWeek()
				return fmt.Sprintf("%02d", weekNumber)
			},
		},
		'w': {
			Desc:   "Day of week number (0-6) where Sunday is 0 and Saturday is 6",
			expand: func(dt time.Time) string { return strconv.Itoa(int(dt.Weekday())) },
		},
		// TODO: implement same rules as strptime - first week is the first Monday of January
		'W': {
			Desc: "ISO8601 week number of the year (this is different to the typical strptime token currently)",
			expand: func(dt time.Time) string {
				_, week := dt.ISOWeek()
				return fmt.Sprintf("%02d", week)
			},
		},
		// NOTE: haven't found a locale aware version of these:
		// 'x': func(dt time.Time) string { return dt.Format("01/02/06") },
		// 'X': func(dt time.Time) string { return dt.Format(time.TimeOnly) },
		'y': {
			Desc:   "The year within the century zero padded to two digits (00-99)",
			expand: func(dt time.Time) string { return fmt.Sprintf("%02d", dt.Year()%100) },
		},
		'Y': {
			Desc:   "Year number - '1999', '2007'",
			expand: func(dt time.Time) string { return strconv.Itoa(dt.Year()) },
		},
	},
}

// Formats date strings via the same system as `strptime`
var MomentJs = DateFormatterNoPrefix{
	escapeChars: []rune{'[', ']'},
	TokenMap: map[string]FormatToken[string]{
		"M": {
			Desc:   "Month number (1-12)",
			expand: func(dt time.Time) string { return strconv.Itoa(int(dt.Month())) },
		},
		"Mo": {
			Desc:   "Month number suffixed - '1st', '13th', '22nd'",
			expand: func(dt time.Time) string { return numberSuffixed(int(dt.Month())) },
		},
		"MM": {
			Desc:   "Month number zero padded to two digits - (01-12)",
			expand: func(dt time.Time) string { return fmt.Sprintf("%02d", dt.Month()) },
		},
		"MMM": {
			Desc:   "Month name truncated to three characters - 'Jan', 'Feb'",
			expand: func(dt time.Time) string { return dt.Month().String()[:3] },
		},
		"MMMM": {
			Desc:   "Month name - 'January', 'February'",
			expand: func(dt time.Time) string { return dt.Month().String() },
		},
		"Q": {
			Desc:   "Quarter of year (1-4)",
			expand: func(dt time.Time) string { return strconv.Itoa(time.Now().YearDay() % 4) },
		},
		"Qo": {
			Desc:   "Quarter of year suffixed - '1st', '2nd', '3rd', '4th'",
			expand: func(dt time.Time) string { return numberSuffixed(time.Now().YearDay() % 4) },
		},
		"D": {
			Desc:   "Day of month (1-31)",
			expand: func(dt time.Time) string { return strconv.Itoa(dt.Day()) },
		},
		"Do": {
			Desc:   "Day of month suffixed (1st-31st)",
			expand: func(dt time.Time) string { return numberSuffixed(dt.Day()) },
		},
		"DD": {
			Desc:   "Day of month zero padded to two digits (01-31)",
			expand: func(dt time.Time) string { return fmt.Sprintf("%02d", dt.Day()) },
		},
		"DDD": {
			Desc:   "Day of year (1-366)",
			expand: func(dt time.Time) string { return strconv.Itoa(dt.YearDay()) },
		},
		"DDDo": {
			Desc:   "Day of year suffixed (1st-366th)",
			expand: func(dt time.Time) string { return numberSuffixed(dt.YearDay()) },
		},
		"DDDD": {
			Desc:   "Day of year zero padded to three digits (001-366)",
			expand: func(dt time.Time) string { return fmt.Sprintf("%03d", dt.YearDay()) },
		},
		"d": {
			Desc:   "Day of week where Sunday = 0 and Saturday = 6 (0-6)",
			expand: func(dt time.Time) string { return strconv.Itoa(int(dt.Weekday())) },
		},
		"do": {
			Desc:   "Day of week suffixed where Sunday = 0th and Saturday = 6th (0th-6th)",
			expand: func(dt time.Time) string { return numberSuffixed(int(dt.Weekday())) },
		},
		"dd": {
			Desc:   "Day of week name truncated to two characters - 'Su', 'Mo'",
			expand: func(dt time.Time) string { return dt.Weekday().String()[:2] },
		},
		"ddd": {
			Desc:   "Day of week name truncated to three characters - 'Sun', 'Mon'",
			expand: func(dt time.Time) string { return dt.Weekday().String()[:3] },
		},
		"dddd": {
			Desc:   "Day of week name - 'Sunday', 'Monday'",
			expand: func(dt time.Time) string { return dt.Weekday().String() },
		},
		"H": {
			Desc:   "Hour in 24 hour format (0-23)",
			expand: func(dt time.Time) string { return strconv.Itoa(dt.Hour()) },
		},
		"HH": {
			Desc:   "Hour in 24 hour format zero padded to two digits (00-23)",
			expand: func(dt time.Time) string { return fmt.Sprintf("%02d", dt.Hour()) },
		},
		"h": {
			Desc:   "Hour in 12 hour format (1-12)",
			expand: func(dt time.Time) string { return strconv.Itoa(dt.Hour() % 12) },
		},
		"hh": {
			Desc:   "Hour in 12 hour format zero padded to two digits (01-12)",
			expand: func(dt time.Time) string { return fmt.Sprintf("%02d", dt.Hour()%12) },
		},
		"k": {
			Desc:   "Hour in 24 hour format starting from 1 (1-24)",
			expand: func(dt time.Time) string { return strconv.Itoa(dt.Hour()) },
		},
		"kk": {
			Desc:   "Hour in 24 hour format starting from 1 zero padded to two digits (01-24)",
			expand: func(dt time.Time) string { return fmt.Sprintf("%02d", dt.Hour()) },
		},
		"Y": {
			Desc:   "Year number - '1999', '2007'",
			expand: func(dt time.Time) string { return strconv.Itoa(dt.Year()) },
		},
		"YY": {
			Desc:   "Year number truncated to last two digits - '99', '07'",
			expand: func(dt time.Time) string { return strconv.Itoa(dt.Year() % 100) },
		},
		"YYYY": {
			Desc:   "Year number - '1999', '2007'",
			expand: func(dt time.Time) string { return strconv.Itoa(dt.Year()) },
		},
		"YYYYYY": {
			Desc:   "Year number zeo padded to 6 digits - '001999', '002007'",
			expand: func(dt time.Time) string { return fmt.Sprintf("%06d", dt.Year()) },
		},
		"m": {
			Desc:   "Minutes (0-59)",
			expand: func(dt time.Time) string { return strconv.Itoa(dt.Minute()) },
		},
		"mm": {
			Desc:   "Minutes zero padded to two digits (00-59)",
			expand: func(dt time.Time) string { return fmt.Sprintf("%02d", dt.Minute()) },
		},
		"s": {
			Desc:   "Seconds (0-59)",
			expand: func(dt time.Time) string { return strconv.Itoa(dt.Second()) },
		},
		"ss": {
			Desc:   "Seconds zero padded to two digits (00-59)",
			expand: func(dt time.Time) string { return fmt.Sprintf("%02d", dt.Second()) },
		},
		"X": {
			Desc:   "Unix timestamp in seconds",
			expand: func(dt time.Time) string { return strconv.Itoa(int(dt.Unix())) },
		},
		"x": {
			Desc:   "Unix timestamp in milliseconds",
			expand: func(dt time.Time) string { return strconv.Itoa(int(dt.UnixMilli())) },
		},
	},
}

// Formats date strings via the same system as `strptime`
var Luxon = DateFormatterNoPrefix{
	escapeChars: []rune{'[', ']'},
	TokenMap: map[string]FormatToken[string]{
		"a": {
			Desc: "Meridiem - 'AM'",
			expand: func(dt time.Time) string {
				if dt.Hour() < 12 {
					return "AM"
				}
				return "PM"
			},
			aliases: []string{"L"},
		},
		"M": {
			Desc:    "Month number (1-12)",
			expand:  func(dt time.Time) string { return strconv.Itoa(int(dt.Month())) },
			aliases: []string{"L"},
		},
		"MM": {
			Desc:    "Month number zero padded to two digits - (01-12)",
			expand:  func(dt time.Time) string { return fmt.Sprintf("%02d", dt.Month()) },
			aliases: []string{"LL"},
		},
		"MMM": {
			Desc:    "Month name truncated to three characters - 'Jan', 'Feb'",
			expand:  func(dt time.Time) string { return dt.Month().String()[:3] },
			aliases: []string{"LLL"},
		},
		"MMMM": {
			Desc:    "Month name - 'January', 'February'",
			expand:  func(dt time.Time) string { return dt.Month().String() },
			aliases: []string{"LLLL"},
		},
		"MMMMM": {
			Desc:    "Month name truncated to one character - 'J', 'F'",
			expand:  func(dt time.Time) string { return dt.Month().String()[:1] },
			aliases: []string{"LLLLL"},
		},
		"q": {
			Desc:   "Quarter of year (1-4)",
			expand: func(dt time.Time) string { return strconv.Itoa(time.Now().YearDay() % 4) },
		},
		"qq": {
			Desc:   "Quarter of year zero padded to two digits (01-04)",
			expand: func(dt time.Time) string { return fmt.Sprintf("%02d", time.Now().YearDay()%4) },
		},
		"d": {
			Desc:   "Day of month (1-31)",
			expand: func(dt time.Time) string { return strconv.Itoa(dt.Day()) },
		},
		"dd": {
			Desc:   "Day of month zero padded to two digits (01-31)",
			expand: func(dt time.Time) string { return fmt.Sprintf("%02d", dt.Day()) },
		},
		"o": {
			Desc:   "Ordinal day of year (1-366)",
			expand: func(dt time.Time) string { return strconv.Itoa(dt.YearDay()) },
		},
		"ooo": {
			Desc:   "Ordinal day of year zero padded to three digits (001-366)",
			expand: func(dt time.Time) string { return fmt.Sprintf("%03d", dt.YearDay()) },
		},
		"c": {
			Desc:    "Day of week where Monday = 1 and Sunday = 7 (1-7)",
			expand:  func(dt time.Time) string { return strconv.Itoa((7+int(dt.Weekday()))%7 + 1) },
			aliases: []string{"E"},
		},
		"ccc": {
			Desc:    "Day of week name truncated to three characters - 'Sun', 'Mon'",
			expand:  func(dt time.Time) string { return dt.Weekday().String()[:3] },
			aliases: []string{"EEE"},
		},
		"cccc": {
			Desc:    "Day of week name - 'Sunday', 'Monday'",
			expand:  func(dt time.Time) string { return dt.Weekday().String() },
			aliases: []string{"EEEE"},
		},
		"ccccc": {
			Desc:    "Day of week name truncated to one character - 'S', 'M'",
			expand:  func(dt time.Time) string { return dt.Weekday().String()[:1] },
			aliases: []string{"EEEEE"},
		},
		"H": {
			Desc:   "Hour in 24 hour format (0-23)",
			expand: func(dt time.Time) string { return strconv.Itoa(dt.Hour()) },
		},
		"HH": {
			Desc:   "Hour in 24 hour format zero padded to two digits (00-23)",
			expand: func(dt time.Time) string { return fmt.Sprintf("%02d", dt.Hour()) },
		},
		"h": {
			Desc:   "Hour in 12 hour format (1-12)",
			expand: func(dt time.Time) string { return strconv.Itoa(dt.Hour() % 12) },
		},
		"hh": {
			Desc:   "Hour in 12 hour format zero padded to two digits (01-12)",
			expand: func(dt time.Time) string { return fmt.Sprintf("%02d", dt.Hour()%12) },
		},
		"kK": {
			Desc: "ISO week year shortened to the last two digits - '99, '07'",
			expand: func(dt time.Time) string {
				year, _ := dt.ISOWeek()
				return fmt.Sprintf("%02d", year%100)
			},
		},
		"kkkk": {
			Desc: "ISO week year zero padded to four digits - '1999', '2007'",
			expand: func(dt time.Time) string {
				year, _ := dt.ISOWeek()
				return fmt.Sprintf("%04d", year)
			},
		},
		"y": {
			Desc:   "Year number - '1999', '2007'",
			expand: func(dt time.Time) string { return strconv.Itoa(dt.Year()) },
		},
		"yy": {
			Desc:   "Year number truncated to last two digits - '99', '07'",
			expand: func(dt time.Time) string { return strconv.Itoa(dt.Year() % 100) },
		},
		"yyyy": {
			Desc:   "Year number zero padded to four digits - '1999', '0007'",
			expand: func(dt time.Time) string { return fmt.Sprintf("%04d", dt.Year()) },
		},
		"m": {
			Desc:   "Minutes (0-59)",
			expand: func(dt time.Time) string { return strconv.Itoa(dt.Minute()) },
		},
		"mm": {
			Desc:   "Minutes zero padded to two digits (00-59)",
			expand: func(dt time.Time) string { return fmt.Sprintf("%02d", dt.Minute()) },
		},
		"s": {
			Desc:   "Seconds (0-59)",
			expand: func(dt time.Time) string { return strconv.Itoa(dt.Second()) },
		},
		"ss": {
			Desc:   "Seconds zero padded to two digits (00-59)",
			expand: func(dt time.Time) string { return fmt.Sprintf("%02d", dt.Second()) },
		},
		"S": {
			Desc:   "Milliseconds (0-999)",
			expand: func(dt time.Time) string { return strconv.Itoa(dt.Nanosecond() / 1_000_000) },
		},
		"SSS": {
			Desc:    "Milliseconds zero padded to three digits (000-999)",
			expand:  func(dt time.Time) string { return fmt.Sprintf("%03d", dt.Nanosecond()/1_000_000) },
			aliases: []string{"u"},
		},
		"W": {
			Desc: "ISO week (1-53)",
			expand: func(dt time.Time) string {
				_, week := dt.ISOWeek()
				return strconv.Itoa(week)
			},
		},
		"WW": {
			Desc: "ISO week zero padded to two digits (01-53)",
			expand: func(dt time.Time) string {
				_, week := dt.ISOWeek()
				return fmt.Sprintf("%02d", week)
			},
		},
		"X": {
			Desc:   "Unix timestamp in seconds",
			expand: func(dt time.Time) string { return strconv.Itoa(int(dt.Unix())) },
		},
		"x": {
			Desc:   "Unix timestamp in milliseconds",
			expand: func(dt time.Time) string { return strconv.Itoa(int(dt.UnixMilli())) },
		},
		"z": {
			Desc:   "IANA canonical time zone string - 'Europe/London'",
			expand: func(dt time.Time) string { return dt.Location().String() },
		},
		"Z": {
			Desc: "Time zone offset shortened to one digit - '+5', '-3'",
			expand: func(dt time.Time) string {
				_, offsetSeconds := dt.Zone()
				offsetHours := offsetSeconds / (60 * 60)
				if offsetHours < 0 {
					return fmt.Sprintf("-%d", offsetHours)
				}
				return fmt.Sprintf("+%d", offsetHours)
			},
		},
		"ZZ": {
			Desc: "Time zone offset - '+05:30', '-03:00'",
			expand: func(dt time.Time) string {
				_, offsetSeconds := dt.Zone()
				offsetMinutes := offsetSeconds / 60
				offsetHours := offsetMinutes / 60
				sign := '+'
				if offsetHours < 0 {
					sign = '-'
				}
				return fmt.Sprintf("%c%02d:%02d", sign, offsetHours, offsetMinutes%60)
			},
		},
		"ZZZ": {
			Desc: "Time zone offset formatted without the dividing ':' - '+0530', '-0300'",
			expand: func(dt time.Time) string {
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
		"ZZZZ": {
			Desc: "Abbreviated time zone offset - 'GMT', 'CEST', '+0530'",
			expand: func(dt time.Time) string {
				offsetName, _ := dt.Zone()
				return offsetName
			},
		},
	},
}
