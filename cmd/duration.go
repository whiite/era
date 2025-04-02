package cmd

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// TODO: features:
// - Support fractional values? (e.g. 1.5h)
// - Support addition and subtraction (e.g. 1s + 2s => 3000)

var OutputDur string
var Separator string
var Round bool

func init() {
	durationCmd.Flags().StringVarP(&OutputDur, "output", "o", "ms", "Output units to display the duration as")
	durationCmd.Flags().StringVarP(&Separator, "separator", "_", "", "Output units displayed with a visual separator with '_' by default (e.g. 1_000_000)")
	durationCmd.Flags().BoolVarP(&Round, "int", "i", false, "Output duration as an integer rounded down")
	sepFlag := durationCmd.Flags().Lookup("separator")
	sepFlag.NoOptDefVal = "_"
	rootCmd.AddCommand(durationCmd)
}

var durationCmd = &cobra.Command{
	Use:     "duration",
	Aliases: []string{"dur"},
	Short:   "Parse and convert durations",
	Long:    "Parse and convert human readable durations into different units and formats",
	RunE: func(cmd *cobra.Command, args []string) error {
		outUnit, err := durationUnit(OutputDur)
		if err != nil {
			return err
		}

		res, err := parseDur(strings.TrimSpace(args[0]))
		if err != nil {
			return err
		}

		unitRes := float64(res) / float64(outUnit)
		if Round {
			unitRes = math.Floor(unitRes)
		}

		if Separator != "" {
			fmt.Println(formatFloatWithSeparator(unitRes, Separator))
		} else {
			fmt.Println(unitRes)
		}
		return nil
	},
}

// Parses a provided duration string into a provided output following the `time`
// package definition for time: 1 = 1 nanosecond
func parseDur(durStr string) (int, error) {
	valStack := []rune{}
	unitStack := []rune{}

	parse := func(valStack, unitStack []rune) (int, error) {
		durVal, err := strconv.Atoi(string(valStack))
		if err != nil {
			return 0, fmt.Errorf("Invalid duration value %q", string(valStack))
		}

		durUnit, err := durationUnit(string(unitStack))
		if err != nil {
			return 0, err
		}

		return durVal * durUnit, nil
	}

	total := 0
	for _, char := range durStr {
		switch {
		case '0' <= char && char <= '9':
			if len(unitStack) > 0 {
				val, err := parse(valStack, unitStack)
				if err != nil {
					return 0, err
				}

				total += val
				valStack = []rune{}
				unitStack = []rune{}
			}

			valStack = append(valStack, char)
		case 'A' <= char && char <= 'Z':
		case 'a' <= char && char <= 'z':
			unitStack = append(unitStack, char)
		// TODO: add default units in case of ' '
		case char == ' ' || char == '_':
			continue
		default:
			return 0, fmt.Errorf("Invalid character: %q", char)
		}
	}

	if len(valStack) > 0 || len(unitStack) > 0 {
		val, err := parse(valStack, unitStack)
		if err != nil {
			return 0, err
		}

		total += val
	}

	return total, nil
}

func durationUnit(durUnit string) (int, error) {
	switch durUnit {
	case "d", "day", "days":
		return int(time.Hour * 24), nil
	case "h", "hour", "hours":
		return int(time.Hour), nil
	case "m", "minute", "minutes":
		return int(time.Minute), nil
	case "s", "second", "seconds":
		return int(time.Second), nil
	case "ms", "millisecond", "milliseconds":
		return int(time.Millisecond), nil
	case "ns", "nanosecond", "nanoseconds":
		return int(time.Nanosecond), nil
	}

	return 0, fmt.Errorf("Invalid unit: %q", durUnit)
}

// Formats an int value into one with character separators for visual clarity
func formatIntWithSeparator(val int, separator string) string {
	numStr := strconv.Itoa(val)
	output := numStr
	for idx := len(numStr) - 1; idx > 0; idx -= 1 {
		if (idx+1)%3 == 0 {
			output = output[:idx] + separator + output[idx:]
		}
	}

	return output
}

// Formats an int value into one with character separators for visual clarity
func formatFloatWithSeparator(val float64, separator string) string {
	numStr := strconv.FormatFloat(val, 'f', -1, 64)
	output := numStr
	pointOffset := strings.IndexRune(numStr, '.')
	if pointOffset == -1 {
		pointOffset = len(numStr)
	}
	for idx := pointOffset - 1; idx > 0; idx -= 1 {
		if (pointOffset-idx)%3 == 0 {
			output = output[:idx] + separator + output[idx:]
		}
	}

	return output
}
