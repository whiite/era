package cmd

import (
	"fmt"
	"strings"
	_ "time/tzdata"

	"gitlab.com/monokuro/era/parser"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(parserCmd)
}

var parserCmd = &cobra.Command{
	Use:   "parser",
	Short: "List all available parsers",
	Long:  "List all parsers available to use with any command that supports a parser argument",
	Run: func(cmd *cobra.Command, args []string) {
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

type parserDesc struct {
	formatter parser.DateFormatter
	alias     []string
}

// luxon and moment are not supported currently
var parserMap = map[string]parserDesc{
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
	"strftime": {
		formatter: &parser.CStr,
		alias:     []string{"c", "strptime"},
	},
	"go:strftime": {
		formatter: &parser.GoStrptime,
		alias:     []string{"go:strptime"},
	},
}
