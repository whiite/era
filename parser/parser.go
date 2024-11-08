package parser

import (
	"fmt"
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
}

type DateFormatterPrefix struct {
	Prefix   rune
	TokenMap map[rune]FormatToken[rune]
}

func (formatter DateFormatterPrefix) Parse(dt time.Time, locale locales.Translator, str *string) string {
	var formattedDate strings.Builder
	interpretToken := false

	for _, char := range *str {
		if char == formatter.Prefix {
			if interpretToken {
				formattedDate.WriteRune(formatter.Prefix)
			}
			interpretToken = true
			continue
		}

		if token, hasToken := formatter.TokenMap[char]; interpretToken && hasToken {
			formattedDate.WriteString(token.expand(dt, locale))
			interpretToken = false
			continue
		}

		formattedDate.WriteRune(char)
		interpretToken = false
	}

	return formattedDate.String()
}

func (formatter DateFormatterPrefix) TokenDesc() string {
	var output strings.Builder
	for tokenChar, tokenDef := range formatter.TokenMap {
		output.WriteString(fmt.Sprintf("%c%c: %s\n", formatter.Prefix, tokenChar, tokenDef.Desc))
	}
	return output.String()
}

type DateFormatterNoPrefix struct {
	escapeChars []rune
	TokenMap    map[string]FormatToken[string]
}

func (formatter DateFormatterNoPrefix) Parse(dt time.Time, locale locales.Translator, str *string) string {
	var formattedDate strings.Builder
	var tokens strings.Builder

	tokenNodeRoot := createTokenGraph(&formatter.TokenMap)
	tokenNode := tokenNodeRoot

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
		} else if escapeSupport && char == escapeStartChar {
			escapeMode = true
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
		formattedDate.WriteRune(char)

		tokenNode = tokenNodeRoot
		if node, hasToken := tokenNodeRoot.children[char]; hasToken && !escapeMode {
			tokenNode = node
		}
	}

	if formatFunc := tokenNode.value.expand; formatFunc != nil {
		formattedDate.WriteString(formatFunc(dt, locale))
	} else {
		formattedDate.WriteString(tokens.String())
	}

	return formattedDate.String()
}

func (formatter DateFormatterNoPrefix) TokenDesc() string {
	var output strings.Builder
	for tokenStr, tokenDef := range formatter.TokenMap {
		output.WriteString(fmt.Sprintf("%s: %s\n", tokenStr, tokenDef.Desc))
	}
	return output.String()
}

func numberSuffixed(num int) string {
	numstr := strconv.Itoa(num)
	switch numstr {
	case "1", "21", "31":
		return numstr + "st"
	case "2", "22":
		return numstr + "nd"
	case "3", "23":
		return numstr + "rd"
	default:
		return numstr + "th"
	}
}
