package parser

import (
	"fmt"
	"maps"
	"slices"
	"strings"
	"time"

	"github.com/go-playground/locales"
)

// Date handler for string based tokens of varying lengths with support for escaping
// characters
type DateHandlerString struct {
	escapeChars []rune
	tokenDef    TokenMap
	tokenGraph  *TokenGraphNode[FormatToken[string]]
}

func (formatter *DateHandlerString) TokenMap() TokenMap {
	return expandTokenMap(&formatter.tokenDef)
}

func (formatter *DateHandlerString) Format(dt time.Time, locale locales.Translator, str *string) string {
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

func (formatter *DateHandlerString) Parse(input, format string) (time.Time, error) {
	return time.Now(), nil
}

func (formatter *DateHandlerString) TokenDescTokenFormatter(tokenFmt func(format string, a ...any) string) string {
	var output strings.Builder
	for _, tokenStr := range slices.Sorted(maps.Keys(formatter.tokenDef)) {
		tokenDef := formatter.tokenDef[tokenStr]
		output.WriteString(tokenFmt("%s: ", tokenStr))
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

func (formatter *DateHandlerString) TokenDesc() string {
	return formatter.TokenDescTokenFormatter(fmt.Sprintf)
}
