package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// Format string flag to output the time with
var Parser string

func init() {
	parseCmd.Flags().StringVarP(&Parser, "parser", "p", "", "Parser to interpret supplied time with")
	parseCmd.Flags().StringVarP(&TimeZone, "timezone", "t", "", "Time zone to set the time to")
	rootCmd.AddCommand(parseCmd)
}

var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "Parse a given time",
	Long:  "Parse a given time in order to maniuplate; convert or output it in a different format",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		location := time.Now().Local().Location()

		if TimeZone != "" {
			loc, err := time.LoadLocation(TimeZone)
			if err != nil {
				return err
			}
			location = loc
		}

		switch strings.ToLower(Parser) {
		case "unix", "timestamp", "ts", "":
			unixVal, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("Unable to parse '%s' as a unix timestamp", args[0])
			}
			res := time.Unix(int64(unixVal), 0).In(location)
			fmt.Println(res.String())
		default:
			return fmt.Errorf("'%s' is not a supported parser", Parser)
		}

		return nil
	},
}
