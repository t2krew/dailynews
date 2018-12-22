package mail

import (
	"bytes"
	"html/template"
)

func ParseTemplate(name string, data interface{}) (content string, err error) {
	tpl, err := template.ParseFiles(name)
	if err != nil {
		return
	}
	buffer := new(bytes.Buffer)
	if err = tpl.Execute(buffer, data); err != nil {
		return
	}
	return buffer.String(), nil
}
