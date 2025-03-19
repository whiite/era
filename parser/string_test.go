package parser

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-playground/locales/en_GB"
)

func TestTokensStrftime(t *testing.T) {
	tokens := GoStrptime.TokenMapExpanded()

	var testDates = []string{"2024-01-07", "1997-01-04", "1989-12-31", "2007-01-01"}

	for _, datestr := range testDates {
		dt, _ := time.Parse(time.DateOnly, datestr)
		for _, loc := range []string{"America/Los_Angeles", "Europe/London", "Europe/Paris"} {
			loc, err := time.LoadLocation(loc)
			if err != nil {
				t.Errorf("Location %q is unsupported: %s", loc, err)
				break
			}

			for token := range tokens {
				format := fmt.Sprintf("%%%s", token)
				t.Run("", func(t *testing.T) {
					t.Parallel()
					CompareCtx{
						Dt:        dt.In(loc),
						Locale:    en_GB.New(),
						Formatter: &GoStrptime,
						Format:    format,
					}.Test(execDate, t)
				})
			}
		}
	}
}

func TestFormatStringsStrftime(t *testing.T) {
	var testDates = []string{"2024-01-07", "1997-01-04", "1989-12-31", "2007-01-01"}

	for _, datestr := range testDates {
		dt, _ := time.Parse(time.DateOnly, datestr)
		for _, loc := range []string{"America/Los_Angeles", "Europe/London", "Europe/Paris"} {
			loc, err := time.LoadLocation(loc)
			if err != nil {
				t.Errorf("Location %q is unsupported: %s", loc, err)
				break
			}

			for _, formatstr := range []string{
				"%toutput: %G%M%%%S%v",
				"%%%%%%%S",
				"% %V%% %t %t",
				"%n%n%n",
			} {
				t.Run("", func(t *testing.T) {
					t.Parallel()
					CompareCtx{
						Dt:        dt.In(loc),
						Locale:    en_GB.New(),
						Formatter: &GoStrptime,
						Format:    formatstr,
					}.Test(execDate, t)
				})
			}
		}
	}
}
