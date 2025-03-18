package parser

import (
	"testing"
	"time"

	"github.com/go-playground/locales/en_GB"
)

func TestTokensMoment(t *testing.T) {
	tokens := MomentJs.TokenMap()
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
				t.Run("", func(t *testing.T) {
					t.Parallel()
					CompareCtx{
						Dt:        dt.In(loc),
						Locale:    en_GB.New(),
						Formatter: &MomentJs,
						Format:    token,
					}.Test(execMoment, t)
				})
			}
		}
	}
}
