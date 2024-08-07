package parser

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type DateFormatterPrefix struct {
	prefix   rune
	tokenMap map[rune]func(dt time.Time) string
}

func (formatter DateFormatterPrefix) Parse(dt time.Time, str *string) string {
	var formattedDate strings.Builder
	interpretToken := false

	for _, char := range *str {
		if char == formatter.prefix {
			if interpretToken {
				formattedDate.WriteRune(formatter.prefix)
			}
			interpretToken = true
			continue
		}

		formatFunc := formatter.tokenMap[char]
		if interpretToken && formatFunc != nil {
			formattedDate.WriteString(formatFunc(dt))
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
	tokenMap    map[string]func(dt time.Time) string
}

func (formatter DateFormatterNoPrefix) Parse(dt time.Time, str *string) string {
	var formattedDate strings.Builder
	var tokens strings.Builder

	tokenNodeRoot := createTokenGraph(&formatter.tokenMap)
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

		if formatFunc := tokenNode.value; formatFunc != nil {
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

	if formatFunc := tokenNode.value; formatFunc != nil {
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
	prefix: '%',
	tokenMap: map[rune]func(dt time.Time) string{
		'%': func(dt time.Time) string { return "%" },
		'A': func(dt time.Time) string { return dt.Weekday().String() },
		'a': func(dt time.Time) string { return dt.Weekday().String()[:3] },
		'B': func(dt time.Time) string { return dt.Month().String() },
		'b': func(dt time.Time) string { return dt.Month().String()[:3] },
		'h': func(dt time.Time) string { return dt.Month().String()[:3] },
		'c': func(dt time.Time) string { return dt.Format("Mon _2 Jan 15:04:05 2006") },
		'C': func(dt time.Time) string { return strconv.Itoa(dt.Year() / 100) },
		'd': func(dt time.Time) string { return fmt.Sprintf("%02d", dt.Day()) },
		'e': func(dt time.Time) string { return strconv.Itoa(dt.Day()) },
		'D': func(dt time.Time) string {
			return fmt.Sprintf("%02d/%02d/%02d", dt.Month(), dt.Day(), dt.Year()%100)
		},
		'H': func(dt time.Time) string { return strconv.Itoa(dt.Hour()) },
		'I': func(dt time.Time) string { return fmt.Sprintf("%02d", dt.Hour()%12) },
		'j': func(dt time.Time) string { return strconv.Itoa(dt.YearDay()) },
		'm': func(dt time.Time) string { return fmt.Sprintf("%02d", dt.Month()) },
		'M': func(dt time.Time) string { return fmt.Sprintf("%02d", dt.Minute()) },
		'n': func(dt time.Time) string { return "\n" },
		'S': func(dt time.Time) string { return fmt.Sprintf("%02d", dt.Second()) },
		'R': func(dt time.Time) string { return fmt.Sprintf("%02d:%02d", dt.Hour(), dt.Minute()) },
		't': func(dt time.Time) string { return "\t" },
		'T': func(dt time.Time) string {
			return fmt.Sprintf("%02d:%02d:%02d", dt.Hour(), dt.Minute(), dt.Second())
		},
		'U': func(dt time.Time) string {
			_, weekNumber := dt.ISOWeek()
			return fmt.Sprintf("%02d", weekNumber)
		},
		'w': func(dt time.Time) string { return strconv.Itoa(int(dt.Weekday())) },
		'W': func(dt time.Time) string {
			_, week := dt.ISOWeek()
			return fmt.Sprintf("%02d", week)
		},
		// NOTE: haven't found a locale aware version of these:
		// 'x': func(dt time.Time) string { return dt.Format("01/02/06") },
		// 'X': func(dt time.Time) string { return dt.Format(time.TimeOnly) },
		'y': func(dt time.Time) string { return strconv.Itoa(dt.Year() % 100) },
		'Y': func(dt time.Time) string { return strconv.Itoa(dt.Year()) },
	},
}

// Formats date strings via the same system as `strptime`
var MomentJs = DateFormatterNoPrefix{
	escapeChars: []rune{'[', ']'},
	tokenMap: map[string]func(dt time.Time) string{
		"M":    func(dt time.Time) string { return strconv.Itoa(int(dt.Month())) },
		"Mo":   func(dt time.Time) string { return numberSuffixed(int(dt.Month())) },
		"MM":   func(dt time.Time) string { return fmt.Sprintf("%02d", dt.Month()) },
		"MMM":  func(dt time.Time) string { return dt.Month().String()[:3] },
		"MMMM": func(dt time.Time) string { return dt.Month().String() },
		"Q":    func(dt time.Time) string { return strconv.Itoa(time.Now().YearDay() % 4) },
		"Qo":   func(dt time.Time) string { return numberSuffixed(time.Now().YearDay() % 4) },
		"D":    func(dt time.Time) string { return strconv.Itoa(dt.Day()) },
		"Do":   func(dt time.Time) string { return numberSuffixed(dt.Day()) },
		"DD":   func(dt time.Time) string { return fmt.Sprintf("%02d", dt.Day()) },
		"DDD":  func(dt time.Time) string { return fmt.Sprintf("%02d", dt.Day()) },
		"DDDo": func(dt time.Time) string { return numberSuffixed(dt.YearDay()) },
		"DDDD": func(dt time.Time) string { return fmt.Sprintf("%03d", dt.YearDay()) },
		"d":    func(dt time.Time) string { return strconv.Itoa(int(dt.Weekday())) },
		"do":   func(dt time.Time) string { return numberSuffixed(int(dt.Weekday())) },
		"dd":   func(dt time.Time) string { return dt.Weekday().String()[:2] },
		"ddd":  func(dt time.Time) string { return dt.Weekday().String()[:3] },
		"dddd": func(dt time.Time) string { return dt.Weekday().String() },
		"H":    func(dt time.Time) string { return strconv.Itoa(dt.Hour()) },
		"HH":   func(dt time.Time) string { return fmt.Sprintf("%02d", dt.Hour()) },
		"h":    func(dt time.Time) string { return strconv.Itoa(dt.Hour() % 12) },
		"hh":   func(dt time.Time) string { return fmt.Sprintf("%02d", dt.Hour()%12) },
		"k":    func(dt time.Time) string { return strconv.Itoa(dt.Hour()) },
		"kk":   func(dt time.Time) string { return fmt.Sprintf("%02d", dt.Hour()) },
		"Y":    func(dt time.Time) string { return strconv.Itoa(dt.Year()) },
		"YY":   func(dt time.Time) string { return strconv.Itoa(dt.Year() % 100) },
		"YYYY": func(dt time.Time) string { return strconv.Itoa(dt.Year()) },
		"m":    func(dt time.Time) string { return strconv.Itoa(dt.Minute()) },
		"mm":   func(dt time.Time) string { return fmt.Sprintf("%02d", dt.Minute()) },
		"s":    func(dt time.Time) string { return strconv.Itoa(dt.Second()) },
		"ss":   func(dt time.Time) string { return fmt.Sprintf("%02d", dt.Second()) },
		"X":    func(dt time.Time) string { return strconv.Itoa(int(dt.Unix())) },
		"x":    func(dt time.Time) string { return strconv.Itoa(int(dt.UnixMilli())) },
	},
}
