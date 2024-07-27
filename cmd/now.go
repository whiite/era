package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var Format string

func init() {
	nowCmd.Flags().StringVarP(&Format, "format", "f", "", "Format to display the datetime with")
	rootCmd.AddCommand(nowCmd)
}

var nowCmd = &cobra.Command{
	Use:   "now",
	Short: "Get and manipulate the current time",
	Long:  "Convert, format and print the current time ",
	Run: func(cmd *cobra.Command, args []string) {
		now := time.Now()
		switch Format {
		case "unix", "timestamp", "ts":
			fmt.Println(now.Unix())
		case "go", "":
			fmt.Println(now)
		default:
			fmt.Println(now.Format(Format))
		}
	},
}
