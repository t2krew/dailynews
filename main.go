package main

import (
	"fmt"
	"github.com/t2krew/daily/mail"
	"log"
)

func main() {
	mailConf, err := Configer("email")
	if err != nil {
		panic(err)
	}

	var (
		email    = mailConf.GetString("email")
		password = mailConf.GetString("password")
		host     = mailConf.GetString("host")
		port     = mailConf.GetInt("port")
		nickname = mailConf.GetString("nickname")
	)

	fmt.Println(email, password, nickname, host, port)

	mailbox := mail.NewMail(email, password, nickname, host, port)

	to := []string{"648367227@qq.com"}

	data := mail.Content{
		Subject: "测试邮件",
		Message: "这是一封测试邮件，请查收后删除",
	}

	err = mailbox.Send(to, data)
	if err != nil {
		log.Panicln(err)
	}

	log.Println("success")
}
