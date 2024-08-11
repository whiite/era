package cmd

import (
	"fmt"
	"monokuro/era/parser"
	"strings"

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
				args = make([]string, len(parser.MomentJs.TokenMap))
				idx := 0
				for token := range parser.MomentJs.TokenMap {
					args[idx] = token
					idx += 1
				}
			}

			var validOutput strings.Builder
			var invalidOutput strings.Builder
			for _, arg := range args {
				if token, hasToken := parser.MomentJs.TokenMap[arg]; hasToken {
					validOutput.WriteString(fmt.Sprintf("'%s': %s\n", arg, token.Desc))
					continue
				}
				invalidOutput.WriteString(fmt.Sprintf("Token '%s' is not supported for this parser", arg))
			}

			if validOutput.Len() == 0 {
				return fmt.Errorf(invalidOutput.String())
			}

			fmt.Print(validOutput.String(), invalidOutput.String())
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
