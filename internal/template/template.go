package template

import (
	"bytes"
	"html/template"
	"path"
	"strings"
)

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
