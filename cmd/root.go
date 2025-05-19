package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var NoColor bool

func init() {
	rootCmd.PersistentFlags().BoolVar(&NoColor, "no-color", false, "Disable coloured output")
}

var rootCmd = &cobra.Command{
	Use:     "era",
	Version: "0.4.3",
	Short:   "Simple utility for working with time and dates",
	Long:    `A simple and intuitive tool for working with and manipulating time and dates`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
