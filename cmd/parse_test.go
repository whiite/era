package cmd

import (
	"fmt"
	"gitlab.com/monokuro/era/parser"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en_GB"
)

func execDate(ctx compareCtx, t *testing.T) string {
	cmd := exec.Command("date", "-r", fmt.Sprintf("%d", ctx.dt.Unix()), fmt.Sprintf("+%s", ctx.format))
	cmd.Env = append(cmd.Environ(), fmt.Sprintf("TZ=%s", ctx.dt.Location()))
	cmd.Env = append(cmd.Environ(), fmt.Sprintf("LC_ALL=%s", ctx.locale.Locale()))
	var out strings.Builder
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		t.Errorf("Could not run %q: %s", cmd.Path, err)
	}
	output := out.String()
	// Remove newline char appended to the end of stdout
	// Deliberately don't use trim in case of intentional \n formats
	return output[:len(output)-1]
}

func execLuxon(ctx compareCtx, t *testing.T) string {
	evalStr := fmt.Sprintf("(await import('npm:luxon')).DateTime.fromSeconds(%d).setZone(%q).toFormat(%q)", ctx.dt.Unix(), ctx.dt.Location(), ctx.format)
	cmd := exec.Command("deno", "eval", "-p", evalStr)
	var out strings.Builder
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		t.Errorf("Could not run %q: %s", cmd.Path, err)
	}
	output := out.String()
	// Remove newline char appended to the end of stdout
	// Deliberately don't use trim in case of intentional \n formats
	return output[:len(output)-1]
}

func execMoment(ctx compareCtx, t *testing.T) string {
	evalStr := fmt.Sprintf("(await import('npm:moment-timezone')).default.tz(%d, %q).format(%q)", ctx.dt.UnixMilli(), ctx.dt.Location(), ctx.format)
	cmd := exec.Command("deno", "eval", "-p", evalStr)
	var out strings.Builder
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		t.Errorf("Could not run %q: %s", cmd.Path, err)
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

func (ctx *compareCtx) Error(got, want string) string {
	return fmt.Sprintf(`Formatter: %q did not match format string: %q
location: %q
date:     %q | %d
locale:   %q

got:    %q
wanted: %q`,
		ctx.formatter,
		ctx.format,
		ctx.dt.Location().String(),
		ctx.dt.String(),
		ctx.dt.Unix(),
		ctx.locale.Locale(),
		got,
		want)
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
	want := execCmd(ctx, t)
	if got != want {
		t.Error(ctx.Error(got, want))
	}
}

var testDates = []string{"2024-01-07", "1997-01-04", "1989-12-31", "2007-01-01"}

func TestTokensStrftime(t *testing.T) {
	tokens := parser.Strftime.TokenMapExpanded()

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
					compareFormat(compareCtx{
						dt:        dt.In(loc),
						locale:    en_GB.New(),
						formatter: "go:strftime",
						format:    format,
					}, t)
				})
			}
		}
	}
}

func TestTokensLuxon(t *testing.T) {
	tokens := parser.Luxon.TokenMap()
	// NOTE: tokens known to be problematic but fixing is difficult
	excludeList := map[string]bool{"DDDD": true, "ttt": true, "tttt": true, "ZZZZ": true, "TTT": true}

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
	tokens := parser.MomentJs.TokenMap()

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

func TestFormatStringsStrftime(t *testing.T) {
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
					compareFormat(compareCtx{
						dt:        dt.In(loc),
						locale:    en_GB.New(),
						formatter: "go:strftime",
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
		t.Errorf("Location %q is unsupported: %s", loc, err)
		return
	}

	compareFormat(compareCtx{
		dt:        dt.In(loc),
		locale:    en_GB.New(),
		formatter: "strftime",
		format:    "%tout",
	}, t)

}
