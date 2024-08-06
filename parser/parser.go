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
	tokenMap map[string]func(dt time.Time) string
}

func (formatter DateFormatterNoPrefix) Parse(dt time.Time, str *string) string {
	var formattedDate strings.Builder

	tokenNodeRoot := createTokenGraph(&formatter.tokenMap)
	tokenNode := tokenNodeRoot

	for _, char := range *str {
		if node, hasToken := tokenNode.children[char]; hasToken {
			tokenNode = node
			continue
		}

		if formatFunc := tokenNode.value; formatFunc != nil {
			formattedDate.WriteString(formatFunc(dt))
		}
		formattedDate.WriteRune(char)

		tokenNode = tokenNodeRoot
		if node, hasToken := tokenNodeRoot.children[char]; hasToken {
			tokenNode = node
		}
	}

	if formatFunc := tokenNode.value; formatFunc != nil {
		formattedDate.WriteString(formatFunc(dt))
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
		'Y': func(dt time.Time) string { return strconv.Itoa(dt.Year()) },
		'm': func(dt time.Time) string { return strconv.Itoa(int(dt.Month())) },
		'd': func(dt time.Time) string { return strconv.Itoa(dt.Day()) },
		'j': func(dt time.Time) string { return strconv.Itoa(dt.YearDay()) },
		'u': func(dt time.Time) string { return strconv.Itoa(int(dt.Weekday())) },
		'H': func(dt time.Time) string { return strconv.Itoa(dt.Hour()) },
		'I': func(dt time.Time) string { return strconv.Itoa(dt.Hour() % 12) },
		'M': func(dt time.Time) string { return strconv.Itoa(dt.Minute()) },
		'S': func(dt time.Time) string { return strconv.Itoa(dt.Second()) },
		'n': func(dt time.Time) string { return "\n" },
		'R': func(dt time.Time) string { return strconv.Itoa(dt.Hour()) + ":" + strconv.Itoa(dt.Minute()) },
		'T': func(dt time.Time) string {
			return strconv.Itoa(dt.Hour()) + ":" + strconv.Itoa(dt.Minute()) + ":" + strconv.Itoa(dt.Second())
		},
		'U': func(dt time.Time) string {
			_, weekNumber := dt.ISOWeek()
			return strconv.Itoa(weekNumber)
		},
	},
}

// Formats date strings via the same system as `strptime`
var MomentJs = DateFormatterNoPrefix{
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
		"h":    func(dt time.Time) string { return strconv.Itoa(dt.Hour()) },
		"hh": func(dt time.Time) string {
			hour := strconv.Itoa(dt.Hour())
			if len(hour) < 2 {
				hour = "0" + hour
			}
			return hour
		},
	},
}
