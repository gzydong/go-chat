package utils

import (
	"bytes"
	"html/template"
)

func RenderTemplate(text []byte, data any) (string, error) {
	tmpl, _ := template.New("tmpl").Parse(string(text))

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return "", err
	}

	return body.String(), nil
}
