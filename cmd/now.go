package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	_ "time/tzdata"

	"gitlab.com/monokuro/era/localiser"
	"gitlab.com/monokuro/era/parser"

	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en_GB"

	"github.com/spf13/cobra"
)

// Format string flag to output the time with
var Format string

// Canonical time zone string flag to convert the time to before displaying
var TimeZone string

// Locale string to use when formatting date times
var Locale string

func init() {
	nowCmd.Flags().StringVarP(&Format, "formatter", "F", "", "Formatter to interpret and display the current datetime with")
	nowCmd.Flags().StringVarP(&TimeZone, "timezone", "t", "", "Time zone to set the time to")
	nowCmd.Flags().StringVarP(&Locale, "locale", "l", "", "Locale to use in formatting")
	rootCmd.AddCommand(nowCmd)
}

var nowCmd = &cobra.Command{
	Use:   "now",
	Short: "Get and manipulate the current time",
	Long:  "Convert, format and print the current time ",
	RunE: func(cmd *cobra.Command, args []string) error {
		now := time.Now()
		locale := en_GB.New()

		if TimeZone != "" {
			location, err := time.LoadLocation(TimeZone)
			if err != nil {
				return err
			}
			now = now.In(location)
		}

		if Locale != "" {
			parsedLocale, err := localiser.Parse(Locale)
			if err != nil {
				return err
			}
			locale = parsedLocale
		}

		parseStr := ""
		if len(args) > 0 {
			parseStr = args[0]
		}

		nowFormatted, err := FormatTime(now, locale, Format, parseStr)
		if err != nil {
			return err
		}
		fmt.Println(nowFormatted)

		return nil
	},
}

func FormatTime(dt time.Time, locale locales.Translator, formatter string, parseStr string) (string, error) {
	formattedTime := ""

	switch strings.ToLower(formatter) {
	case "unix", "timestamp", "ts":
		formattedTime = strconv.FormatInt(dt.Unix(), 10)
	case "rfc", "rfc3339":
		formattedTime = dt.Format(time.RFC3339)
	case "iso", "iso8601":
		formattedTime = dt.Format("2006-01-02T15:04:05.999Z07:00")
	case "go":
		if len(parseStr) == 0 {
			parseStr = "2006-01-02 15:04:05.999999999 -0700 MST"
		}
		formattedTime = parser.Go.Format(dt, locale, &parseStr)
	case "moment", "momentjs":
		if len(parseStr) == 0 {
			return formattedTime, fmt.Errorf("No format string provided")
		}
		formattedTime = parser.MomentJs.Format(dt, locale, &parseStr)
	case "luxon":
		if len(parseStr) == 0 {
			return formattedTime, fmt.Errorf("No format string provided")
		}
		formattedTime = parser.Luxon.Format(dt, locale, &parseStr)
	case "go:strftime", "go:strptime":
		if len(parseStr) == 0 {
			return formattedTime, fmt.Errorf("No format string provided")
		}
		formattedTime = parser.Strftime.Format(dt, locale, &parseStr)
	case "c", "strftime", "strptime":
		formattedTime = parser.CStr.Format(dt, locale, &parseStr)
	case "":
		formattedTime = dt.String()
	default:
		return formattedTime, fmt.Errorf("%q is not a supported formatter", formatter)
	}

	return formattedTime, nil
}
