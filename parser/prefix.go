package parser

import (
	"fmt"
	"maps"
	"slices"
	"strings"
	"time"

	"github.com/go-playground/locales"
)

// Date handler for tokens that always start with a predefined prefix rune
//
// e.g. C based `strptime` function tokens: '%c', '%Oy' etc.
type DateHandlerPrefix struct {
	Prefix     rune
	tokenDef   TokenMap
	tokenGraph *TokenGraphNode[FormatToken[string]]
}

func (formatter *DateHandlerPrefix) Parse(input, format string) (time.Time, error) {
	dt := time.Unix(0, 0).UTC()

	tokenNode := formatter.tokenGraph
	interpretMode := false
	offsetIdx := 0
	var startIdx int

	for idx, char := range format {
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
			slice := input[startIdx:]
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
		slice := input[startIdx:]
		_, err := parseFunc(&dt, &slice)
		if err != nil {
			return dt, err
		}
	}

	return dt, nil
}

func (formatter *DateHandlerPrefix) TokenMap() TokenMap {
	return expandTokenMap(&formatter.tokenDef)
}

func (formatter DateHandlerPrefix) Format(dt time.Time, locale locales.Translator, str *string) string {
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
func (formatter DateHandlerPrefix) TokenMapExpanded() map[string]FormatToken[string] {
	expandedTokenMap := map[string]FormatToken[string]{}
	for token, tokenDef := range formatter.tokenDef {
		expandedTokenMap[token] = tokenDef
		for _, alias := range tokenDef.aliases {
			expandedTokenMap[alias] = tokenDef
		}
	}
	return expandedTokenMap
}

func (formatter DateHandlerPrefix) TokenDescTokenFormatter(tokenFmt func(format string, a ...any) string) string {
	var output strings.Builder
	for _, tokenChar := range slices.Sorted(maps.Keys(formatter.tokenDef)) {
		tokenDef := formatter.tokenDef[tokenChar]
		output.WriteString(tokenFmt("%c%s: ", formatter.Prefix, tokenChar))
		output.WriteString(fmt.Sprintf("%s\n", tokenDef.Desc))
		if len(tokenDef.aliases) > 0 {
			output.WriteString("  aliases:")
			for _, alias := range tokenDef.aliases {
				output.WriteString(" ")
				output.WriteString(tokenFmt("%s", alias))
			}
			output.WriteString("\n")
		}
	}
	return output.String()
}

func (formatter DateHandlerPrefix) TokenDesc() string {
	return formatter.TokenDescTokenFormatter(fmt.Sprintf)
}
