package parser

import (
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/locales"
)

type FormatToken[T any] struct {
	// Description of what the token represents
	Desc string
	// Equivalent string for token given a `time.Time`
	expand  func(dt time.Time, locale locales.Translator) string
	parse   func(dt *time.Time, str *string) (int, error)
	aliases []T
}

type TokenMap = map[string]FormatToken[string]

// Expands the inner token map to include aliases
func expandTokenMap(tmap *TokenMap) TokenMap {
	expandedTokenMap := map[string]FormatToken[string]{}
	for token, tokenDef := range *tmap {
		expandedTokenMap[token] = tokenDef
		for _, alias := range tokenDef.aliases {
			expandedTokenMap[alias] = tokenDef
		}
	}
	return expandedTokenMap
}

// Describes and contains a mapping of tokens
type DateDescriptor interface {
	// Outputs all supported tokens with their descriptions and any aliases
	TokenDesc() string
	// Mapping between any token and it's corresponding meta data
	TokenMap() TokenMap
}

// Allows formatting date times into other data types
type DateFormatter interface {
	// Formats a given time and locale according to the provided string
	Format(dt time.Time, locale locales.Translator, str *string) string
}

// Allows parsing data types into date times
type DateParser interface {
	// Parses a time string according to the provided format
	Parse(input, format string) (time.Time, error)
}

type DateHandler interface {
	*DateDescriptor
	*DateFormatter
	*DateParser
}

// Simplest date handler type wrapping functionality defined elsewhere that works with tokens
type DateHandlerTokenWrapper struct {
	format   func(dt time.Time, locale locales.Translator, formatStr string) string
	parse    func(input, format string) (time.Time, error)
	tokenDef TokenMap
	prefix   rune
}

func (formatter *DateHandlerTokenWrapper) TokenMap() TokenMap {
	return expandTokenMap(&formatter.tokenDef)
}

func (formatter *DateHandlerTokenWrapper) Format(dt time.Time, locale locales.Translator, str *string) string {
	return formatter.format(dt, locale, *str)
}

func (formatter *DateHandlerTokenWrapper) Parse(input, format string) (time.Time, error) {
	return formatter.parse(input, format)
}

func (formatter *DateHandlerTokenWrapper) TokenDesc() string {
	var output strings.Builder
	for _, tokenChar := range slices.Sorted(maps.Keys(formatter.tokenDef)) {
		tokenDef := formatter.tokenDef[tokenChar]
		output.WriteString(fmt.Sprintf("%c%s: %s\n", formatter.prefix, tokenChar, tokenDef.Desc))
		if len(tokenDef.aliases) > 0 {
			output.WriteString("  aliases:")
			for _, alias := range tokenDef.aliases {
				output.WriteString(fmt.Sprintf(" %s", alias))
			}
			output.WriteString("\n")
		}
	}
	return output.String()
}

func numberSuffixed(num int) string {
	numstr := strconv.Itoa(num)

	// Keep "th" suffix for 11, 12, 13 ending ints
	if twoDigit := num % 100; twoDigit >= 11 && twoDigit <= 13 {
		return numstr + "th"
	}

	switch num % 10 {
	case 1:
		return numstr + "st"
	case 2:
		return numstr + "nd"
	case 3:
		return numstr + "rd"
	default:
		return numstr + "th"
	}
}
