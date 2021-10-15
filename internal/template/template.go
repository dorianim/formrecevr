package template

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"
)

// CreateDefaultTemplate creates the default template
func CreateDefaultTemplate(path string) error {
	defaultTemplate := `New form submited using Formrecevr:
{{ range $key, $val := . }}{{ $key }}: {{ join $val "" }}{{ end }}`
	defaultTemplatePath := fmt.Sprintf("%s/default.html", path)

	_, err := os.Stat(defaultTemplatePath)
	if err == nil || !os.IsNotExist(err) {
		return err
	}
	os.MkdirAll(path, os.ModePerm)

	f, _ := os.OpenFile(defaultTemplatePath, os.O_CREATE|os.O_WRONLY, 0644)

	f.Truncate(0)
	f.Write([]byte(defaultTemplate))

	f.Close()
	return nil
}

// ExecuteTemplate parses a template, fills it and returns the result
func ExecuteTemplateFromFile(templateFilePath string, data interface{}) (string, error) {
	parsedTemplate, err := template.New(path.Base(templateFilePath)).Funcs(getFuncMap()).ParseFiles(templateFilePath)
	if err != nil {
		return "", err
	}

	return executeTemplate(parsedTemplate, data)
}

func ExecuteTemplateFromString(templateString string, data interface{}) (string, error) {
	parsedTemplate, err := template.New("t1").Funcs(getFuncMap()).Parse(templateString)
	if err != nil {
		return "", err
	}

	return executeTemplate(parsedTemplate, data)
}

func executeTemplate(template *template.Template, data interface{}) (string, error) {
	buf := new(bytes.Buffer)
	err := template.Execute(buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func getFuncMap() template.FuncMap {
	return template.FuncMap{
		"join": strings.Join,
		"print": func(stringList []string) string {
			return strings.Join(stringList, "")
		},
	}
}
