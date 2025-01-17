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
		t.Errorf(`Formatter: '%s' did not match for format '%s'
	got:    '%s'
	wanted: '%s'`, ctx.formatter, ctx.format, got, want)
	}
}

func TestTokensStrptime(t *testing.T) {
	tokens := parser.Strptime.TokenMapExpanded()

	for token := range tokens {
		format := fmt.Sprintf("%%%c", token)
		compareFormat(compareCtx{
			dt:        time.Now(),
			locale:    en_GB.New(),
			formatter: "strptime",
			format:    format,
		}, t)
	}

}
