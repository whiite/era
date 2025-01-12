package localiser

import (
	"fmt"
	"strings"

	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en_GB"
	"github.com/go-playground/locales/en_US"
	"github.com/go-playground/locales/es"
	"github.com/go-playground/locales/es_US"
	"github.com/go-playground/locales/fr"
)

// TODO: codegen this
func Parse(localeStr string) (locales.Translator, error) {
	selectedLocale := en_GB.New()
	var err error

	switch strings.ToLower(localeStr) {
	case "en_us", "en", "us":
		selectedLocale = en_US.New()
	case "es_us":
		selectedLocale = es_US.New()
	case "es":
		selectedLocale = es.New()
	case "fr":
		selectedLocale = fr.New()
	case "en_gb", "gb", "uk":
		selectedLocale = en_GB.New()
	default:
		err = fmt.Errorf("Unsupported locale: '%s'", localeStr)
	}
	return selectedLocale, err
}
