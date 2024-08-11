package cmd

import (
	"fmt"
	"monokuro/era/parser"

	"github.com/spf13/cobra"
)

func init() {
	tokensCmd.Flags().StringVarP(&Parser, "parser", "p", "", "Parser to interpret supplied time wFormat to display the datetime withith")
	rootCmd.AddCommand(tokensCmd)
}

var tokensCmd = &cobra.Command{
	Use:   "tokens",
	Short: "Token info for parsers",
	Long:  "Get and query tokens for specified parsers to see what they map to",
	RunE: func(cmd *cobra.Command, args []string) error {
		switch Parser {
		case "moment", "momentjs":
			if len(args) == 0 {
				return fmt.Errorf("Showing all tokens not supported yet. Please specify at least one token")
			}
			token, hasToken := parser.MomentJs.TokenMap[args[0]]
			if !hasToken {
				return fmt.Errorf("Token '%s' is not supported for this parser", args[0])
			}
			fmt.Printf("'%s': %s\n", args[0], token.Desc)
		case "strptime":
			if len(args) == 0 {
				return fmt.Errorf("Showing all tokens not supported yet. Please specify at least one token")
			}
			tokenChar := []rune(args[0])[0]
			token, hasToken := parser.Strptime.TokenMap[tokenChar]
			if !hasToken {
				return fmt.Errorf("Token '%c' is not supported for this parser", tokenChar)
			}
			fmt.Printf("'%c%c': %s\n", parser.Strptime.Prefix, tokenChar, token.Desc)
		default:
			return fmt.Errorf("No parser specified")
		}
		return nil
	},
}
