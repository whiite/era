package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// TODO: features:
// - Support and parse milti-rune units (e.g. "ms", "ns")
// - Choose output units (e.g. 1m --output seconds => 60)
// - Support printing with separators (e.g. 10000 => 10_000)
// - Support addition and subtraction (e.g. 1s + 2s => 3000)

func init() {
	rootCmd.AddCommand(durationCmd)
}

var durationCmd = &cobra.Command{
	Use:     "duration",
	Aliases: []string{"dur"},
	Short:   "Parse and convert durations",
	Long:    "Parse and convert durations into different output formats",
	RunE: func(cmd *cobra.Command, args []string) error {
		res, err := parseDur(strings.TrimSpace(args[0]))
		if err != nil {
			return err
		}
		fmt.Println(res)
		return nil
	},
}

// Parses a provided duration string into milliseconds
func parseDur(durStr string) (int, error) {
	durs := [][]rune{}
	durStack := []rune{}

	for _, char := range durStr {
		switch char {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			durStack = append(durStack, char)
		case 'd', 'h', 'm', 's':
			durStack = append(durStack, char)
			durs = append(durs, durStack)
			durStack = []rune{}
		}
	}

	total := 0
	for _, durStack := range durs {
		unit := durStack[len(durStack)-1]
		durVal, err := strconv.Atoi(string(durStack[:len(durStack)-1]))
		if err != nil {
			return 0, err
		}

		durMs, err := duration(string(unit))
		if err != nil {
			return 0, err
		}

		total += durVal * durMs
	}
	return total / int(time.Millisecond), nil
}

func duration(durUnit string) (int, error) {
	switch durUnit {
	case "d":
		return int(time.Hour * 24), nil
	case "h":
		return int(time.Hour), nil
	case "m":
		return int(time.Minute), nil
	case "s":
		return int(time.Second), nil
	}

	return 1, nil
}
