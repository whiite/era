package cmd

import (
	"fmt"
	"gitlab.com/monokuro/era/parser"

	"github.com/spf13/cobra"
)

func init() {
	tokensCmd.Flags().StringVarP(&Parser, "formatter", "F", "", "Formatter to display supported tokens")
	rootCmd.AddCommand(tokensCmd)
}

var tokensCmd = &cobra.Command{
	Use:   "tokens",
	Short: "Token info for parsers",
	Long:  "Get and query tokens for specified parsers to see what they map to",
	RunE: func(cmd *cobra.Command, args []string) error {
		var selectedParser parser.DateDescriptor
		switch Parser {
		case "moment", "momentjs":
			selectedParser = &parser.MomentJs
		case "luxon":
			selectedParser = &parser.Luxon
		case "c", "strftime", "strptime":
			selectedParser = &parser.CStr
		case "go:strptime", "go:strftime":
			selectedParser = &parser.GoStrptime
		case "go":
			selectedParser = &parser.Go
		case "":
			return fmt.Errorf("No parser specified")
		default:
			return fmt.Errorf("Parser %q is not supported", Parser)
		}
		fmt.Print(selectedParser.TokenDesc())
		return nil
	},
}
