//go:build ignore

package main

import (
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"
	"time"
)

type Locale struct {
	Name      string
	AliasList []string
}

func main() {
	pkgName := "localiser"
	err := os.MkdirAll(pkgName, 0750)
	if err != nil {
		panic(err)
	}

	file, err := os.Create(path.Join(pkgName, "parser.go"))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	packageTemplate.Execute(file, struct {
		Timestamp   time.Time
		PackageName string
		Locales     []Locale
	}{
		Timestamp:   time.Now(),
		PackageName: pkgName,
		Locales: []Locale{
			{
				Name:      "en_US",
				AliasList: []string{"en", "us"},
			},
			{
				Name: "fr",
			},
			{
				Name: "es",
			},
		},
	})
	fmt.Println("Done")
}

var funcMap = template.FuncMap{
	"lower": strings.ToLower,
}

var packageTemplate = template.Must(template.New("").Funcs(funcMap).Parse(`// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots at
// {{ .Timestamp }}
package {{ .PackageName }}

import (
	"fmt"
	"strings"

	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en_GB"
	{{- range .Locales }}
	"github.com/go-playground/locales/{{ printf "%s" .Name }}"
	{{- end }}
)

func Parse(localeStr string) (locales.Translator, error) {
	selectedLocale := en_GB.New()
	var err error

	switch strings.ToLower(localeStr) {
	case "en_gb", "gb", "uk":
		selectedLocale = en_GB.New()

	{{- range .Locales }}
	case "{{ printf "%s" .Name | lower }}" {{- range .AliasList }}{{ printf ", %q" . }}{{- end }}:
		selectedLocale = {{ printf "%s" .Name }}.New()
	{{- end }}
	default:
		err = fmt.Errorf("Unsupported locale: '%s'", localeStr)
	}

	return selectedLocale, err
}
`))
