package parser

/*
#cgo noescape strftime
#cgo nocallback strftime
#cgo noescape strptime
#cgo nocallback strptime

#include <stdio.h>
#include <stdlib.h>
#include <time.h>
*/
import "C"

import (
	"time"
	"unsafe"
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
var CStr = DateFormatterWrapper{
	format: func(dt time.Time, formatStr string) string {
		format := C.CString(formatStr)
		result := C.CString("")
		defer C.free(unsafe.Pointer(result))

		tm := timeToTm(&dt)
		C.strftime(result, C.sizeof_char*1024, format, &tm)
		return C.GoString(result)
	},
	parse: func(input, format string) (time.Time, error) {
		tm := C.struct_tm{
			tm_year: C.int(time.Unix(0, 0).Year()),
		}

		cInput, cFormat := C.CString(input), C.CString(format)
		defer C.free(unsafe.Pointer(cInput))
		defer C.free(unsafe.Pointer(cFormat))

		C.strptime(cInput, cFormat, &tm)
		dt := tmToTime(&tm, time.Local)
		return dt, nil
	},
	tokenDef: Strftime.tokenDef,
}
