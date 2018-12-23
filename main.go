package main

import (
	"fmt"
	"github.com/t2krew/daily/output"
	"github.com/t2krew/daily/output/dingding"
	"github.com/t2krew/daily/output/mail"
	"github.com/t2krew/daily/spider"
	"github.com/t2krew/daily/util"
)

func main() {
	conf, err := Configer("app")
	if err != nil {
		panic(err)
	}

	var (
		email    = conf.GetString("mail.email")
		password = conf.GetString("mail.password")
		host     = conf.GetString("mail.host")
		port     = conf.GetInt("mail.port")
		nickname = conf.GetString("mail.nickname")
		ddrobot  = conf.GetString("dingding.robot")
		receiver = conf.GetStringSlice("mail.receiver")

		dd      = dingding.New(ddrobot)                               // 钉钉
		mailbox = mail.NewMail(email, password, nickname, host, port) // 邮件
	)

	output.Add(dd)
	output.Add(mailbox)

	var list []map[string]string
	for _, s := range spider.Spiders {
		ret, err := s.Handler()
		if err != nil {
			fmt.Println(err)
		}
		list = append(list, ret...)
	}

	var (
		date = util.Today().Format("2006-01-02")
		data = output.Content{
			Subject: fmt.Sprintf("Daily Articles (%s)", date),
			Data: &output.Data{
				Date: date,
				List: list,
			},
			Mime: "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n",
		}
	)

	for _, sender := range output.Outputers {
		err = sender.Send("template/daily.html", receiver, data)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
}
