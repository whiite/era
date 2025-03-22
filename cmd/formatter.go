package cmd

import (
	"fmt"
	"strings"
	_ "time/tzdata"

	"gitlab.com/monokuro/era/parser"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(formatterCmd)
}

var formatterCmd = &cobra.Command{
	Use:   "formatter",
	Short: "List all available formatters",
	Long:  "List all formatters available to use with any command that supports a formatter argument",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: create string from map below
		var output strings.Builder

		for name, meta := range formatterMap {
			output.WriteString(fmt.Sprintf("%s\n", name))
			if len(meta.alias) > 0 {
				output.WriteString("  aliases:")
				for _, alias := range meta.alias {
					output.WriteString(fmt.Sprintf(" %s", alias))
				}
				output.WriteString("\n")
			}
		}

		fmt.Print(output.String())
	},
}

type formatterDesc struct {
	formatter parser.DateFormatter
	alias     []string
}

var formatterMap = map[string]formatterDesc{
	"unix": {
		alias: []string{"timestamp", "ts"},
	},
	"rfc": {
		alias: []string{"rfc3339"},
	},
	"iso": {
		alias: []string{"iso8601"},
	},
	"go": {formatter: &parser.Go},
	"moment": {
		formatter: &parser.MomentJs,
		alias:     []string{"momentjs"},
	},
	"luxon": {formatter: &parser.Luxon},
	"strftime": {
		formatter: &parser.CStr,
		alias:     []string{"c", "strptime"},
	},
	"go:strftime": {
		formatter: &parser.GoStrptime,
		alias:     []string{"go:strptime"},
	},
}
