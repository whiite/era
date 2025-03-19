package parser

/*
#cgo noescape strftime
#cgo nocallback strftime
#cgo noescape strptime
#cgo nocallback strptime

#include <stdio.h>
#include <stdlib.h>
#include <time.h>
#include <locale.h>
*/
import "C"

import (
	"fmt"
	"time"
	"unsafe"

	"github.com/go-playground/locales"
)

func tmToTime(tm *C.struct_tm, location *time.Location) time.Time {
	return time.Date(
		int(tm.tm_year),
		time.Month(tm.tm_mon+1),
		int(tm.tm_mday+1),
		int(tm.tm_hour),
		int(tm.tm_min),
		int(tm.tm_sec),
		0,
		location,
	)
}

// Converts a Go based time.Time to a C tm struct equivalent
//
// Careful with memory when using.
// struct_tm.tm_zone is a C string that should be freed when done
func timeToTm(dt *time.Time) C.struct_tm {
	zone, offset := dt.Zone()
	dst := 0
	if dt.IsDST() {
		dst = 1
	}
	return C.struct_tm{
		// Years since 1900
		tm_year: C.int(dt.Year() - 1900),
		// Zero indexed
		tm_mon:  C.int(dt.Month()) - 1,
		tm_mday: C.int(dt.Day()),
		tm_wday: C.int(dt.Weekday()),
		tm_hour: C.int(dt.Hour()),
		tm_min:  C.int(dt.Minute()),
		tm_sec:  C.int(dt.Second()),
		// Offset in minutes
		tm_gmtoff: C.long(offset / 60),
		tm_isdst:  C.int(dst),
		tm_yday:   C.int(dt.YearDay()),
		tm_zone:   C.CString(zone),
	}
}

// Formatter built around C FFI to strftime and strptime
//
// As this uses C FFI with `cgo`; not all OS's are supported or tested
var CStr = DateHandlerWrapper{
	format: func(dt time.Time, locale locales.Translator, formatStr string) string {
		format := C.CString(formatStr)
		result := C.CString("")
		defer C.free(unsafe.Pointer(format))
		defer C.free(unsafe.Pointer(result))

		tm := timeToTm(&dt)
		defer C.free(unsafe.Pointer(tm.tm_zone))

		localeStr := C.CString(fmt.Sprintf("%s.UTF-8", locale.Locale()))
		defer C.free(unsafe.Pointer(localeStr))
		C.setlocale(C.LC_TIME, localeStr)

		// %c => ~24 chars
		// '%', 'c' => 12 chars each
		// Nearest multiple of 8 = 16
		const MAX_CHAR_PRINT = 16
		bufferSize := len(formatStr) * MAX_CHAR_PRINT
		C.strftime(result, C.size_t(bufferSize*C.sizeof_char), format, &tm)

		return C.GoString(result)
	},
	parse: func(input, format string) (time.Time, error) {
		tm := C.struct_tm{
			// NOTE: this needs the full year not just the years since 1900
			// Unsure why in this particular case and not when formatting
			tm_year: C.int(time.Unix(0, 0).Year()),
		}

		cInput, cFormat := C.CString(input), C.CString(format)
		defer C.free(unsafe.Pointer(cInput))
		defer C.free(unsafe.Pointer(cFormat))

		C.strptime(cInput, cFormat, &tm)
		dt := tmToTime(&tm, time.Local)
		return dt, nil
	},
	prefix:   '%',
	tokenDef: tokenMapStrftime,
}
