package utils

import (
	"bytes"
	"html/template"
)

func StrTemplate(text []byte, data interface{}) (string, error) {
	tmpl, _ := template.New("tmpl").Parse(string(text))

	return render(tmpl, data)
}

func render(tmpl *template.Template, data interface{}) (string, error) {
	var body bytes.Buffer

	_ = tmpl.Execute(&body, data)

	return body.String(), nil
}
