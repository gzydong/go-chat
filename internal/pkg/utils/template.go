package utils

import (
	"bytes"
	"html/template"
)

func RenderString(text string, data interface{}) (string, error) {
	tmpl, _ := template.New("tmpl").Parse(text)

	return render(tmpl, data)
}

func RenderFile(filePath string, data interface{}) (string, error) {

	t1 := template.Must(template.ParseFiles(filePath))
	return render(t1, data)
}

func render(tmpl *template.Template, data interface{}) (string, error) {
	var body bytes.Buffer

	_ = tmpl.Execute(&body, data)

	return body.String(), nil
}
