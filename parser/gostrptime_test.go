package parser

import (
	"testing"
	"time"
)

type testCase struct {
	input  string
	format string
	want   time.Time
}

func TestParse(t *testing.T) {
	scenarios := []testCase{
		{input: "04/01/97", format: "%d/%m/%y", want: time.Date(1997, 1, 4, 0, 0, 0, 0, time.Local)},
		{input: " 4/01/97", format: "%e/%m/%y", want: time.Date(1997, 1, 4, 0, 0, 0, 0, time.Local)},
	}

	for _, testCase := range scenarios {
		t.Run("", func(t *testing.T) {
			t.Parallel()
			got, err := GoStrptime.Parse(testCase.input, testCase.format)
			if err != nil {
				t.Errorf("Failed to parse")
				t.Fail()
			}
			if got.Compare(testCase.want) != 0 {
				t.Errorf("Fail\nGot:  %s\nwant: %s", got, testCase.want)
				t.Fail()
			}
		})
	}
}
