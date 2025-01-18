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
	cmd.Env = append(cmd.Environ(), fmt.Sprintf("TZ=%s", dt.Location().String()))
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
	want := execDate(ctx.format, ctx.dt, t)
	if got != want {
		t.Errorf(`Formatter: '%s' did not match format string: '%s'
location: '%s'
date:     '%s' | %d

got:    '%s'
wanted: '%s'`, ctx.formatter, ctx.format, ctx.dt.Location().String(), ctx.dt.String(), ctx.dt.Unix(), got, want)
	}
}

func TestTokensStrptime(t *testing.T) {
	tokens := parser.Strptime.TokenMapExpanded()

	for _, datestr := range []string{"2024-01-07", "1997-01-04", "1989-12-31"} {
		dt, _ := time.Parse(time.DateOnly, datestr)
		for _, loc := range []string{"America/Los_Angeles", "Europe/London", "Europe/Paris"} {
			loc, err := time.LoadLocation(loc)
			if err != nil {
				t.Errorf("Location '%s' is unsupported: %s", loc, err)
				break
			}

			for token := range tokens {
				format := fmt.Sprintf("%%%c", token)
				compareFormat(compareCtx{
					dt:        dt.In(loc),
					locale:    en_GB.New(),
					formatter: "strptime",
					format:    format,
				}, t)
			}
		}
	}

}

func TestFormatStringsStrptime(t *testing.T) {
	for _, datestr := range []string{"2024-01-07", "1997-01-04", "1989-12-31"} {
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
				compareFormat(compareCtx{
					dt:        dt.In(loc),
					locale:    en_GB.New(),
					formatter: "strptime",
					format:    formatstr,
				}, t)
			}
		}
	}

}

// func TestScenario(t *testing.T) {
// 	dt, _ := time.Parse(time.DateOnly, "1997-01-04")
//
// 	loc, err := time.LoadLocation("America/Los_Angeles")
// 	// loc, err := time.LoadLocation("Europe/London")
// 	if err != nil {
// 		t.Errorf("Location '%s' is unsupported: %s", loc, err)
// 		return
// 	}
//
// 	compareFormat(compareCtx{
// 		dt:        dt.In(loc),
// 		locale:    en_GB.New(),
// 		formatter: "strptime",
// 		format:    "%W",
// 	}, t)
//
// }
