package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// Format string flag to output the time with
var Format string

// Canonical time zone string flag to convert the time to before displaying
var TimeZone string

func init() {
	nowCmd.Flags().StringVarP(&Format, "format", "f", "", "Format to display the datetime with")
	nowCmd.Flags().StringVarP(&TimeZone, "timezone", "t", "", "Time zone to set the time to")
	rootCmd.AddCommand(nowCmd)
}

var nowCmd = &cobra.Command{
	Use:   "now",
	Short: "Get and manipulate the current time",
	Long:  "Convert, format and print the current time ",
	RunE: func(cmd *cobra.Command, args []string) error {
		now := time.Now()

		if TimeZone != "" {
			location, err := time.LoadLocation(TimeZone)
			if err != nil {
				return err
			}
			now = now.In(location)
		}

		switch strings.ToLower(Format) {
		case "unix", "timestamp", "ts":
			fmt.Println(now.Unix())
		case "rfc", "rfc3339":
			fmt.Println(now.Format(time.RFC3339))
		case "iso", "iso8601":
			fmt.Println(now.Format("2006-01-02T15:04:05.999Z07:00"))
		case "go", "":
			fmt.Println(now)
		default:
			fmt.Println(parseFormatLinux(now, &Format, &dateFormatter))
		}

		return nil
	},
}

type DateFormatterPrefix struct {
	prefix   rune
	tokenMap map[rune]func(dt time.Time) string
}

// Formats date strings using the unix `date` CLI tokens and prefix
var dateFormatter = DateFormatterPrefix{
	prefix: '%',
	tokenMap: map[rune]func(dt time.Time) string{
		'Y': func(dt time.Time) string { return strconv.Itoa(dt.Year()) },
		'm': func(dt time.Time) string { return strconv.Itoa(int(dt.Month())) },
		'B': func(dt time.Time) string { return dt.Month().String() },
		'b': func(dt time.Time) string { return dt.Month().String()[:3] },
		'd': func(dt time.Time) string { return strconv.Itoa(dt.Day()) },
		'j': func(dt time.Time) string { return strconv.Itoa(dt.YearDay()) },
		'u': func(dt time.Time) string { return strconv.Itoa(int(dt.Weekday())) },
		'A': func(dt time.Time) string { return dt.Weekday().String() },
		'a': func(dt time.Time) string { return dt.Weekday().String()[:3] },
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

func parseFormatLinux(date time.Time, formatString *string, formatter *DateFormatterPrefix) string {
	var formattedDate strings.Builder
	interpretToken := false

	for _, char := range *formatString {
		if char == formatter.prefix {
			if interpretToken {
				formattedDate.WriteRune(formatter.prefix)
			}
			interpretToken = true
			continue
		}

		formatFunc := formatter.tokenMap[char]
		if interpretToken && formatFunc != nil {
			formattedDate.WriteString(formatFunc(date))
			interpretToken = false
			continue
		}

		formattedDate.WriteRune(char)
		interpretToken = false
	}

	return formattedDate.String()
}
