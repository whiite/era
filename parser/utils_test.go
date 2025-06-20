// Package exclusively used to aid testing
package parser

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/go-playground/locales"
)

type CompareCtx struct {
	Dt        time.Time
	Locale    locales.Translator
	Formatter DateFormatter
	Format    string
}

func (ctx *CompareCtx) Error(got, want string) string {
	return fmt.Sprintf(`Formatter did not match format string: %q
location: %q
date:     %q | %d
locale:   %q

got:    %q
wanted: %q`,
		ctx.Format,
		ctx.Dt.Location().String(),
		ctx.Dt.String(),
		ctx.Dt.Unix(),
		ctx.Locale.Locale(),
		got,
		want)
}

func (ctx CompareCtx) Test(compare func(ctx CompareCtx, t *testing.T) string, t *testing.T) {
	got := ctx.Formatter.Format(ctx.Dt, ctx.Locale, &ctx.Format)
	want := compare(ctx, t)
	if got != want {
		t.Error(ctx.Error(got, want))
	}
}

func execDate(ctx CompareCtx, t *testing.T) string {
	cmd := exec.Command("date", "-r", fmt.Sprintf("%d", ctx.Dt.Unix()), fmt.Sprintf("+%s", ctx.Format))
	cmd.Env = append(cmd.Environ(), fmt.Sprintf("TZ=%s", ctx.Dt.Location()))
	cmd.Env = append(cmd.Environ(), fmt.Sprintf("LC_ALL=%s", ctx.Locale.Locale()))
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

func execLuxon(ctx CompareCtx, t *testing.T) string {
	localeCompat := strings.ReplaceAll(ctx.Locale.Locale(), "_", "-")
	evalStr := fmt.Sprintf(
		"(await import('npm:luxon')).DateTime.fromSeconds(%d).setLocale(%q).setZone(%q).toFormat(%q)",
		ctx.Dt.Unix(),
		localeCompat,
		ctx.Dt.Location(),
		ctx.Format,
	)
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

func execMoment(ctx CompareCtx, t *testing.T) string {
	evalStr := fmt.Sprintf("(await import('npm:moment-timezone')).default.tz(%d, %q).format(%q)", ctx.Dt.UnixMilli(), ctx.Dt.Location(), ctx.Format)
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
