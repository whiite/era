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

type DateFormatter interface {
	TokenDesc() string
	Format(dt time.Time, locale locales.Translator, str *string) string
	TokenMap() TokenMap
}

type DateFormatterWrapper struct {
	format   func(dt time.Time, formatStr string) string
	tokenDef TokenMap
}

func (formatter *DateFormatterWrapper) TokenMap() TokenMap {
	return expandTokenMap(&formatter.tokenDef)
}

func (formatter DateFormatterWrapper) Format(dt time.Time, locale locales.Translator, str *string) string {
	return formatter.format(dt, *str)
}

func (formatter DateFormatterWrapper) TokenDesc() string {
	var output strings.Builder
	for _, tokenChar := range slices.Sorted(maps.Keys(formatter.tokenDef)) {
		tokenDef := formatter.tokenDef[tokenChar]
		output.WriteString(fmt.Sprintf("%s: %s\n", tokenChar, tokenDef.Desc))
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

type DateFormatterPrefix struct {
	Prefix     rune
	tokenDef   TokenMap
	tokenGraph *TokenGraphNode[FormatToken[string]]
}

func (formatter *DateFormatterPrefix) Parse(input *string, format *string) (time.Time, error) {
	dt := time.Unix(0, 0).UTC()

	tokenNode := formatter.tokenGraph
	interpretMode := false
	offsetIdx := 0
	var startIdx int

	for idx, char := range *format {
		if char == formatter.Prefix && !interpretMode {
			interpretMode = true
			startIdx = idx + offsetIdx
			continue
		}
		if !interpretMode {
			continue
		}

		if node, hasToken := tokenNode.children[char]; hasToken {
			tokenNode = node
			continue
		}

		if parseFunc := tokenNode.value.parse; parseFunc != nil {
			slice := (*input)[startIdx:]
			offset, err := parseFunc(&dt, &slice)
			if err != nil {
				return dt, err
			}
			offsetIdx += offset - 1
		}

		tokenNode = formatter.tokenGraph
		interpretMode = char == formatter.Prefix
	}

	if parseFunc := tokenNode.value.parse; parseFunc != nil {
		slice := (*input)[startIdx:]
		_, err := parseFunc(&dt, &slice)
		if err != nil {
			return dt, err
		}
	}

	return dt, nil
}

func (formatter *DateFormatterPrefix) TokenMap() TokenMap {
	return expandTokenMap(&formatter.tokenDef)
}

func (formatter DateFormatterPrefix) Format(dt time.Time, locale locales.Translator, str *string) string {
	var formattedDate strings.Builder
	var tokens strings.Builder

	tokenNode := formatter.tokenGraph
	interpretMode := false

	for _, char := range *str {
		if char == formatter.Prefix && !interpretMode {
			interpretMode = true
			continue
		}
		if !interpretMode {
			formattedDate.WriteRune(char)
			continue
		}

		if node, hasToken := tokenNode.children[char]; hasToken {
			tokens.WriteRune(char)
			tokenNode = node
			continue
		}

		if formatFunc := tokenNode.value.expand; formatFunc != nil {
			formattedDate.WriteString(formatFunc(dt, locale))
		} else {
			formattedDate.WriteString(tokens.String())
		}

		tokens.Reset()
		interpretMode = false
		tokenNode = formatter.tokenGraph

		if char == formatter.Prefix {
			interpretMode = true
			continue
		}

		formattedDate.WriteRune(char)
	}

	if formatFunc := tokenNode.value.expand; formatFunc != nil {
		formattedDate.WriteString(formatFunc(dt, locale))
	} else {
		formattedDate.WriteString(tokens.String())
	}

	return formattedDate.String()
}

// Expands the inner token map to include aliases
func (formatter DateFormatterPrefix) TokenMapExpanded() map[string]FormatToken[string] {
	expandedTokenMap := map[string]FormatToken[string]{}
	for token, tokenDef := range formatter.tokenDef {
		expandedTokenMap[token] = tokenDef
		for _, alias := range tokenDef.aliases {
			expandedTokenMap[alias] = tokenDef
		}
	}
	return expandedTokenMap
}

func (formatter DateFormatterPrefix) TokenDesc() string {
	var output strings.Builder
	for _, tokenChar := range slices.Sorted(maps.Keys(formatter.tokenDef)) {
		tokenDef := formatter.tokenDef[tokenChar]
		output.WriteString(fmt.Sprintf("%c%s: %s\n", formatter.Prefix, tokenChar, tokenDef.Desc))
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

type DateFormatterString struct {
	escapeChars []rune
	tokenDef    TokenMap
	tokenGraph  *TokenGraphNode[FormatToken[string]]
}

func (formatter *DateFormatterString) TokenMap() TokenMap {
	return expandTokenMap(&formatter.tokenDef)
}

func (formatter *DateFormatterString) Format(dt time.Time, locale locales.Translator, str *string) string {
	var formattedDate strings.Builder
	var tokens strings.Builder

	tokenNode := formatter.tokenGraph

	escapeSupport := len(formatter.escapeChars) > 0
	escapeMode := false
	var escapeStartChar rune
	var escapeEndChar rune

	if escapeSupport {
		escapeStartChar = formatter.escapeChars[0]
		escapeEndChar = escapeStartChar
		if len(formatter.escapeChars) > 1 {
			escapeEndChar = formatter.escapeChars[1]
		}
	}

	for _, char := range *str {
		if escapeSupport && char == escapeEndChar && escapeMode {
			escapeMode = false
			continue
		} else if escapeSupport && char == escapeStartChar {
			escapeMode = true
			continue
		}

		if node, hasToken := tokenNode.children[char]; hasToken && !escapeMode {
			tokens.WriteRune(char)
			tokenNode = node
			continue
		}

		if formatFunc := tokenNode.value.expand; formatFunc != nil {
			formattedDate.WriteString(formatFunc(dt, locale))
		} else {
			formattedDate.WriteString(tokens.String())
		}
		tokens.Reset()

		tokenNode = formatter.tokenGraph
		if node, hasToken := formatter.tokenGraph.children[char]; hasToken && !escapeMode {
			tokenNode = node
			continue
		}

		formattedDate.WriteRune(char)
	}

	if formatFunc := tokenNode.value.expand; formatFunc != nil {
		formattedDate.WriteString(formatFunc(dt, locale))
	} else {
		formattedDate.WriteString(tokens.String())
	}

	return formattedDate.String()
}

func (formatter DateFormatterString) TokenDesc() string {
	var output strings.Builder
	for _, tokenStr := range slices.Sorted(maps.Keys(formatter.tokenDef)) {
		tokenDef := formatter.tokenDef[tokenStr]
		output.WriteString(fmt.Sprintf("%s: %s\n", tokenStr, tokenDef.Desc))
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
