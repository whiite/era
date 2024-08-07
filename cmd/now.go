package cmd

import (
	"fmt"
	"monokuro/era/parser"
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
		case "moment", "momentjs":
			if len(args) == 0 {
				return fmt.Errorf("No format string provided")
			}
			fmt.Println(parser.MomentJs.Parse(now, &args[0]))
		case "strptime":
			if len(args) == 0 {
				return fmt.Errorf("No format string provided")
			}
			fmt.Println(parser.Strptime.Parse(now, &args[0]))
		default:
			return fmt.Errorf("'%s' is not a supported format", Format)
		}

		return nil
	},
}
