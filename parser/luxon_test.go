package parser

import (
	"testing"
	"time"

	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en_GB"
	// "github.com/go-playground/locales/fr_FR"
)

func TestTokensLuxon(t *testing.T) {
	tokens := Luxon.TokenMap()
	// NOTE: tokens known to be problematic but fixing is difficult
	excludeList := map[string]bool{"DDDD": true, "ttt": true, "tttt": true, "ZZZZ": true, "TTT": true}
	testDates := []string{"2024-01-07", "1997-01-04", "1989-12-31", "2007-01-01"}
	testLocales := []locales.Translator{en_GB.New()}

	for _, locale := range testLocales {
		for _, datestr := range testDates {
			dt, _ := time.Parse(time.DateOnly, datestr)
			for _, loc := range []string{"America/Los_Angeles", "Europe/London", "Europe/Paris"} {
				loc, err := time.LoadLocation(loc)
				if err != nil {
					t.Errorf("Location %q is unsupported: %s", loc, err)
					break
				}

				for token := range tokens {
					if excludeList[token] {
						continue
					}
					t.Run("", func(t *testing.T) {
						t.Parallel()
						CompareCtx{
							Dt:        dt.In(loc),
							Locale:    locale,
							Formatter: &Luxon,
							Format:    token,
						}.Test(execLuxon, t)
					})
				}
			}
		}
	}
}

func TestFormatStringsLuxon(t *testing.T) {
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
				"'output: D'DDanGeROUS",
				"'HH:mm' HH:mm",
				"h:mmd/L/yyyy",
			} {
				t.Run("", func(t *testing.T) {
					t.Parallel()
					CompareCtx{
						Dt:        dt.In(loc),
						Locale:    en_GB.New(),
						Formatter: &Luxon,
						Format:    formatstr,
					}.Test(execLuxon, t)
				})
			}
		}
	}
}
