package cmd

import (
	"fmt"
	"monokuro/era/parser"

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
		var selectedParser parser.DateFormatter
		switch Parser {
		case "moment", "momentjs":
			selectedParser = parser.MomentJs
		case "luxon":
			selectedParser = parser.Luxon
		case "strptime":
			selectedParser = parser.Strptime
		default:
			return fmt.Errorf("No parser specified")
		}
		fmt.Print(selectedParser.TokenDesc())
		return nil
	},
}
