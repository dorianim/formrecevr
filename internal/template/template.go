package template

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"path"
	"strings"
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
	log.Println(err)

	os.MkdirAll(path, os.ModePerm)

	f, err := os.OpenFile(defaultTemplatePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	f.Truncate(0)
	f.Write([]byte(defaultTemplate))

	f.Close()
	return nil
}

// ExecuteTemplate parses a template, fills it and returns the result
func ExecuteTemplate(templateFilePath string, data interface{}) (string, error) {
	parsedTemplate, err := template.New(path.Base(templateFilePath)).Funcs(template.FuncMap{
		"join": strings.Join,
	}).ParseFiles(templateFilePath)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	err = parsedTemplate.Execute(buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
