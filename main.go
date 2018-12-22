package main

import (
	"fmt"
	"log"

	"github.com/t2krew/daily/mail"
	"github.com/t2krew/daily/spider"
	"github.com/t2krew/daily/util"
)

func main() {
	mailConf, err := Configer("email")
	if err != nil {
		panic(err)
	}

	var list []map[string]string
	for _, s := range spider.Spiders {
		ret, err := s.Handler()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(ret)
		list = append(list, ret...)
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

	to := []string{"xxx@qq.com"}

	date := util.Today().Format("2006-01-02")

	data := mail.Content{
		Subject: fmt.Sprintf("Daily Articles (%s)", date),
		Data: &mail.Data{
			Date: date,
			List: list,
		},
		Mime: "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n",
	}

	err = mailbox.Send("template/daily.html", to, data)
	if err != nil {
		log.Panicln(err)
	}

	log.Println("success")
}
