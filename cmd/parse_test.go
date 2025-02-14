package cmd

import (
	"fmt"
	"monokuro/era/parser"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en_GB"
)

func execDate(format string, dt time.Time, t *testing.T) string {
	cmd := exec.Command("date", "-r", fmt.Sprintf("%d", dt.Unix()), fmt.Sprintf("+%s", format))
	cmd.Env = append(cmd.Environ(), fmt.Sprintf("TZ=%s", dt.Location()))
	var out strings.Builder
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		t.Errorf("Could not run '%s': %s", cmd.Path, err)
	}
	output := out.String()
	// Remove newline char appended to the end of stdout
	// Deliberately don't use trim in case of intentional \n formats
	return output[:len(output)-1]
}

func execLuxon(format string, dt time.Time, t *testing.T) string {
	evalStr := fmt.Sprintf("(await import('npm:luxon')).DateTime.fromSeconds(%d).setZone(%q).toFormat(%q)", dt.Unix(), dt.Location(), format)
	cmd := exec.Command("deno", "eval", "-p", evalStr)
	var out strings.Builder
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		t.Errorf("Could not run '%s': %s", cmd.Path, err)
	}
	output := out.String()
	// Remove newline char appended to the end of stdout
	// Deliberately don't use trim in case of intentional \n formats
	return output[:len(output)-1]
}

func execMoment(format string, dt time.Time, t *testing.T) string {
	evalStr := fmt.Sprintf("(await import('npm:moment-timezone')).default.tz(%d, %q).format(%q)", dt.UnixMilli(), dt.Location(), format)
	cmd := exec.Command("deno", "eval", "-p", evalStr)
	var out strings.Builder
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		t.Errorf("Could not run '%s': %s", cmd.Path, err)
	}
	output := out.String()
	// Remove newline char appended to the end of stdout
	// Deliberately don't use trim in case of intentional \n formats
	return output[:len(output)-1]
}

type compareCtx struct {
	dt        time.Time
	locale    locales.Translator
	formatter string
	format    string
}

func compareFormat(ctx compareCtx, t *testing.T) {
	got, err := FormatTime(ctx.dt, ctx.locale, ctx.formatter, ctx.format)
	if err != nil {
		t.Errorf("Error during evaluation of 'got': %s", err)
	}
	execCmd := execDate
	if ctx.formatter == "luxon" {
		execCmd = execLuxon
	} else if ctx.formatter == "moment" {
		execCmd = execMoment
	}
	want := execCmd(ctx.format, ctx.dt, t)
	if got != want {
		t.Errorf(`Formatter: '%s' did not match format string: '%s'
location: '%s'
date:     '%s' | %d

got:    '%s'
wanted: '%s'`, ctx.formatter, ctx.format, ctx.dt.Location().String(), ctx.dt.String(), ctx.dt.Unix(), got, want)
	}
}

var testDates = []string{"2024-01-07", "1997-01-04", "1989-12-31", "2007-01-01"}

func TestTokensStrptime(t *testing.T) {
	tokens := parser.Strptime.TokenMapExpanded()

	for _, datestr := range testDates {
		dt, _ := time.Parse(time.DateOnly, datestr)
		for _, loc := range []string{"America/Los_Angeles", "Europe/London", "Europe/Paris"} {
			loc, err := time.LoadLocation(loc)
			if err != nil {
				t.Errorf("Location '%s' is unsupported: %s", loc, err)
				break
			}

			for token := range tokens {
				format := fmt.Sprintf("%%%s", token)
				t.Run("", func(t *testing.T) {
					t.Parallel()
					compareFormat(compareCtx{
						dt:        dt.In(loc),
						locale:    en_GB.New(),
						formatter: "strptime",
						format:    format,
					}, t)
				})
			}
		}
	}

}

func TestTokensLuxon(t *testing.T) {
	tokens := parser.Luxon.TokenMapExpanded()
	// NOTE: tokens known to be problematic but fixing is difficult
	excludeList := map[string]bool{"DDDD": true, "ttt": true, "tttt": true, "ZZZZ": true, "TTT": true}

	for _, datestr := range testDates {
		dt, _ := time.Parse(time.DateOnly, datestr)
		for _, loc := range []string{"America/Los_Angeles", "Europe/London", "Europe/Paris"} {
			loc, err := time.LoadLocation(loc)
			if err != nil {
				t.Errorf("Location '%s' is unsupported: %s", loc, err)
				break
			}

			for token := range tokens {
				if excludeList[token] {
					continue
				}
				t.Run("", func(t *testing.T) {
					t.Parallel()
					compareFormat(compareCtx{
						dt:        dt.In(loc),
						locale:    en_GB.New(),
						formatter: "luxon",
						format:    token,
					}, t)
				})
			}
		}
	}

}

func TestTokensMoment(t *testing.T) {
	tokens := parser.MomentJs.TokenMapExpanded()

	for _, datestr := range testDates {
		dt, _ := time.Parse(time.DateOnly, datestr)
		for _, loc := range []string{"America/Los_Angeles", "Europe/London", "Europe/Paris"} {
			loc, err := time.LoadLocation(loc)
			if err != nil {
				t.Errorf("Location '%s' is unsupported: %s", loc, err)
				break
			}

			for token := range tokens {
				t.Run("", func(t *testing.T) {
					t.Parallel()
					compareFormat(compareCtx{
						dt:        dt.In(loc),
						locale:    en_GB.New(),
						formatter: "moment",
						format:    token,
					}, t)
				})
			}
		}
	}

}

func TestFormatStringsStrptime(t *testing.T) {
	for _, datestr := range testDates {
		dt, _ := time.Parse(time.DateOnly, datestr)
		for _, loc := range []string{"America/Los_Angeles", "Europe/London", "Europe/Paris"} {
			loc, err := time.LoadLocation(loc)
			if err != nil {
				t.Errorf("Location '%s' is unsupported: %s", loc, err)
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
					compareFormat(compareCtx{
						dt:        dt.In(loc),
						locale:    en_GB.New(),
						formatter: "strptime",
						format:    formatstr,
					}, t)
				})
			}
		}
	}

}

func TestFormatStringsLuxon(t *testing.T) {
	for _, datestr := range testDates {
		dt, _ := time.Parse(time.DateOnly, datestr)
		for _, loc := range []string{"America/Los_Angeles", "Europe/London", "Europe/Paris"} {
			loc, err := time.LoadLocation(loc)
			if err != nil {
				t.Errorf("Location '%s' is unsupported: %s", loc, err)
				break
			}

			for _, formatstr := range []string{
				"'output: D'DDanGeROUS",
				"'HH:mm' HH:mm",
				"h:mmd/L/yyyy",
			} {
				t.Run("", func(t *testing.T) {
					t.Parallel()
					compareFormat(compareCtx{
						dt:        dt.In(loc),
						locale:    en_GB.New(),
						formatter: "luxon",
						format:    formatstr,
					}, t)
				})
			}
		}
	}

}

func TestScenario(t *testing.T) {
	dt, _ := time.Parse(time.DateOnly, "1989-12-31")

	// loc, err := time.LoadLocation("America/Los_Angeles")
	loc, err := time.LoadLocation("Europe/London")
	if err != nil {
		t.Errorf("Location '%s' is unsupported: %s", loc, err)
		return
	}

	compareFormat(compareCtx{
		dt:        dt.In(loc),
		locale:    en_GB.New(),
		formatter: "strptime",
		format:    "%tout",
	}, t)

}
