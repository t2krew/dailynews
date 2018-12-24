package mail

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"html/template"
)

func ParseTemplate(name string, data interface{}) (content string, err error) {
	tplpath := fmt.Sprintf("%s/%s", viper.Get("apppath"), name)
	fmt.Println(tplpath)
	tpl, err := template.ParseFiles(tplpath)
	if err != nil {
		return
	}
	buffer := new(bytes.Buffer)
	if err = tpl.Execute(buffer, data); err != nil {
		return
	}
	return buffer.String(), nil
}
