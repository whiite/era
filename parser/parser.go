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
	aliases []T
}

type DateFormatter interface {
	TokenDesc() string
	Parse(dt time.Time, locale locales.Translator, str *string) string
}

type DateFormatterWrapper struct {
	format   func(dt time.Time, formatStr string) string
	tokenMap map[string]FormatToken[string]
}

func (formatter DateFormatterWrapper) Parse(dt time.Time, locale locales.Translator, str *string) string {
	return formatter.format(dt, *str)
}

func (formatter DateFormatterWrapper) TokenDesc() string {
	var output strings.Builder
	for _, tokenChar := range slices.Sorted(maps.Keys(formatter.tokenMap)) {
		tokenDef := formatter.tokenMap[tokenChar]
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
	tokenMap   map[string]FormatToken[string]
	tokenGraph *TokenGraphNode[FormatToken[string]]
}

func (formatter DateFormatterPrefix) Parse(dt time.Time, locale locales.Translator, str *string) string {
	var formattedDate strings.Builder
	var tokens strings.Builder

	if formatter.tokenGraph == nil {
		tokenMap := formatter.TokenMapExpanded()
		formatter.tokenGraph = createTokenGraph(&tokenMap)
	}
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
	for token, tokenDef := range formatter.tokenMap {
		expandedTokenMap[token] = tokenDef
		for _, alias := range tokenDef.aliases {
			expandedTokenMap[alias] = tokenDef
		}
	}
	return expandedTokenMap
}

func (formatter DateFormatterPrefix) TokenDesc() string {
	var output strings.Builder
	for _, tokenChar := range slices.Sorted(maps.Keys(formatter.tokenMap)) {
		tokenDef := formatter.tokenMap[tokenChar]
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
	tokenMap    map[string]FormatToken[string]
	tokenGraph  *TokenGraphNode[FormatToken[string]]
}

func (formatter *DateFormatterString) Parse(dt time.Time, locale locales.Translator, str *string) string {
	var formattedDate strings.Builder
	var tokens strings.Builder

	if formatter.tokenGraph == nil {
		tokenMap := formatter.TokenMapExpanded()
		formatter.tokenGraph = createTokenGraph(&tokenMap)
	}
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

// Expands the inner token map to include aliases
func (formatter DateFormatterString) TokenMapExpanded() map[string]FormatToken[string] {
	expandedTokenMap := map[string]FormatToken[string]{}
	for token, tokenDef := range formatter.tokenMap {
		expandedTokenMap[token] = tokenDef
		for _, alias := range tokenDef.aliases {
			expandedTokenMap[alias] = tokenDef
		}
	}
	return expandedTokenMap
}

func (formatter DateFormatterString) TokenDesc() string {
	var output strings.Builder
	for _, tokenStr := range slices.Sorted(maps.Keys(formatter.tokenMap)) {
		tokenDef := formatter.tokenMap[tokenStr]
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
