package parser

import (
	"time"

	"github.com/go-playground/locales"
)

var Go = DateHandlerWrapper{
	format: func(dt time.Time, locale locales.Translator, formatStr string) string {
		return dt.Format(formatStr)
	},
	parse: func(input, format string) (time.Time, error) {
		return time.Parse(format, input)
	},
	tokenDef: map[string]FormatToken[string]{
		"1": {
			Desc: "Month number (1-12)",
		},
		"01": {
			Desc: "Month number zero padded to two characters (01-12)",
		},
		"2": {
			Desc: "Day of month (1-31)",
		},
		"_2": {
			Desc: "Day of month space padded to two characters ( 1-31)",
		},
		"02": {
			Desc: "Day of month zero padded to two characters (01-31)",
		},
		"__2": {
			Desc: "Day of year space padded to three characters (  1-366)",
		},
		"002": {
			Desc: "Day of year zero padded to three characters (001-366)",
		},
		"3": {
			Desc: "Hour in 12 hour format (0-12)",
		},
		"03": {
			Desc: "Hour in 12 hour format zero padded to two characters (00-12)",
		},
		"15": {
			Desc: "Hour in 24 hour format (00-23)",
		},
		"4": {
			Desc: "Minute (0-59)",
		},
		"04": {
			Desc: "Minute zero padded to two characters (00-59)",
		},
		"5": {
			Desc: "Second (0-59)",
		},
		"05": {
			Desc: "Second zero padded to two characters (00-59)",
		},
		"06": {
			Desc: "Year number to two characters (00-99)",
		},
		"2006": {
			Desc: "Year number to two characters (0000-9999)",
		},
		"-0700": {
			Desc: "Time zone offset as '±hhmm'",
		},
		"-07:00": {
			Desc: "Time zone offset as '±hh:mm'",
		},
		"-07": {
			Desc: "Time zone offset as '±hh'",
		},
		"-070000": {
			Desc: "Time zone offset as '±hhmmss'",
		},
		"-07:00:00": {
			Desc: "Time zone offset as '±hh:mm:ss'",
		},
		"Z0700": {
			Desc: "ISO 8601 time zone offset name or formatted '±hhmm'",
		},
		"Z07:00": {
			Desc: "ISO 8601 time zone offset name or formatted '±hh:mm'",
		},
		"Z07": {
			Desc: "ISO 8601 time zone offset name or formatted '±hh'",
		},
		"Z070000": {
			Desc: "ISO 8601 time zone offset name or formatted '±hhmmss'",
		},
		"Z07:00:00": {
			Desc: "ISO 8601 time zone offset name or formatted '±hh:mm:ss'",
		},
		"January": {
			Desc: "Month name",
		},
		"Jan": {
			Desc: "Month name shortened to three characters",
		},
		"Monday": {
			Desc: "Day of week name",
		},
		"Mon": {
			Desc: "Day of week name shortened to three characters",
		},
		"PM": {Desc: "AM/PM label"},
	},
}
