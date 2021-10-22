package email

import (
	"bytes"
	"html/template"
)

func RenderString(text string, data interface{}) (string, error) {
	tmpl, _ := template.New("tmpl").Parse(text)

	return render(tmpl, data)
}

func RenderFile(filePath string, data interface{}) (string, error) {
	//var body bytes.Buffer

	t1 := template.Must(template.ParseFiles(filePath))
	return render(t1, data)

	//_ = t1.Execute(&body, data)
	//
	//return body.String(), nil

	//t, _ := template.New("tmpl").ParseFiles(filePath)
	//_ = t1.ExecuteTemplate(&body, "verify_code.tmpl", data)
	//
	//return body.String(), nil
}

func render(tmpl *template.Template, data interface{}) (string, error) {
	var body bytes.Buffer

	_ = tmpl.Execute(&body, data)

	return body.String(), nil
}
