package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	_ "time/tzdata"

	"gitlab.com/monokuro/era/localiser"
	"gitlab.com/monokuro/era/parser"

	"github.com/go-playground/locales/en_GB"
	"github.com/spf13/cobra"
)

// Format string flag to output the time with
var Parser string

func init() {
	parseCmd.Flags().StringVarP(&Format, "format", "f", "", "Format to display the datetime with")
	parseCmd.Flags().StringVarP(&Parser, "formatter", "F", "", "Formatter to interpret and display the supplied datetime with")
	parseCmd.Flags().StringVarP(&TimeZone, "timezone", "t", "", "Time zone to set the time to")
	parseCmd.Flags().StringVarP(&Locale, "locale", "l", "", "Locale to use in formatting")
	rootCmd.AddCommand(parseCmd)
}

var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "Parse a given time",
	Long:  "Parse a given time in order to manipulate; convert or output it in a different format",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		location := time.Now().Local().Location()
		locale := en_GB.New()

		if TimeZone != "" {
			loc, err := time.LoadLocation(TimeZone)
			if err != nil {
				return err
			}
			location = loc
		}

		if Locale != "" {
			parsedLocale, err := localiser.Parse(Locale)
			if err != nil {
				return err
			}
			locale = parsedLocale
		}

		var dt time.Time
		switch strings.ToLower(Parser) {
		case "unix", "timestamp", "ts":
			unixVal, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("Unable to parse %q as a unix timestamp", args[0])
			}
			dt = time.Unix(unixVal, 0).In(location)
		case "rfc", "rfc3339":
			time, err := time.Parse(time.RFC3339, args[0])
			if err != nil {
				return fmt.Errorf("Unable to parse %q as an RFC3339 string", args[0])
			}
			dt = time.In(location)
		case "iso", "iso8601":
			time, err := time.Parse("2006-01-02T15:04:05.999Z07:00", args[0])
			if err != nil {
				return fmt.Errorf("Unable to parse %q as an ISO8601 string", args[0])
			}
			dt = time.In(location)
		case "go", "":
			formatStr := "2006-01-02 15:04:05.999999999 -0700 MST"
			if len(args) > 1 {
				formatStr = args[1]
			}
			time, err := time.Parse(formatStr, args[0])
			if err != nil {
				return fmt.Errorf("Unable to parse %q as a Go format string", args[0])
			}
			dt = time.In(location)
		case "go:strftime", "go:strptime":
			if len(args) == 1 {
				return fmt.Errorf("Missing specified format argument")
			}
			time, err := parser.GoStrptime.Parse(args[0], args[1])
			if err != nil {
				return fmt.Errorf("Failed to parse %q via the strptime parser", args[0])
			}
			dt = time.In(location)
		case "c", "strftime", "strptime":
			if len(args) == 1 {
				return fmt.Errorf("Missing specified format argument")
			}
			time, err := parser.CStr.Parse(args[0], args[1])
			if err != nil {
				return fmt.Errorf("Failed to parse %q via the strptime parser", args[0])
			}
			dt = time.In(location)
		default:
			return fmt.Errorf("%q is not a supported parser", Parser)
		}

		formattedTime, err := FormatTime(dt, locale, Parser, Format)
		if err != nil {
			return err
		}
		fmt.Println(formattedTime)

		return nil
	},
}
