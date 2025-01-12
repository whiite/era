package cmd

import (
	"fmt"
	"monokuro/era/parser"
	"strconv"
	"strings"
	"time"
	_ "time/tzdata"

	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en_GB"

	"github.com/spf13/cobra"
)

// Format string flag to output the time with
var Format string

// Canonical time zone string flag to convert the time to before displaying
var TimeZone string

func init() {
	nowCmd.Flags().StringVarP(&Format, "formatter", "F", "", "Formatter to interpret and display the current datetime with")
	nowCmd.Flags().StringVarP(&TimeZone, "timezone", "t", "", "Time zone to set the time to")
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

func FormatTime(dt time.Time, locale locales.Translator, format string, parseStr string) (string, error) {
	formattedTime := ""

	switch strings.ToLower(format) {
	case "unix", "timestamp", "ts":
		formattedTime = strconv.FormatInt(dt.Unix(), 10)
	case "rfc", "rfc3339":
		formattedTime = dt.Format(time.RFC3339)
	case "iso", "iso8601":
		formattedTime = dt.Format("2006-01-02T15:04:05.999Z07:00")
	case "go", "":
		formattedTime = dt.String()
	case "moment", "momentjs":
		if len(parseStr) == 0 {
			return formattedTime, fmt.Errorf("No format string provided")
		}
		formattedTime = parser.MomentJs.Parse(dt, locale, &parseStr)
	case "luxon":
		if len(parseStr) == 0 {
			return formattedTime, fmt.Errorf("No format string provided")
		}
		formattedTime = parser.Luxon.Parse(dt, locale, &parseStr)
	case "strptime":
		if len(parseStr) == 0 {
			return formattedTime, fmt.Errorf("No format string provided")
		}
		formattedTime = parser.Strptime.Parse(dt, locale, &parseStr)
	default:
		return formattedTime, fmt.Errorf("'%s' is not a supported format", parseStr)
	}

	return formattedTime, nil
}
