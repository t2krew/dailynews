package main

import (
	"fmt"
	"log"

	"github.com/t2krew/dailynews/output"
	"github.com/t2krew/dailynews/output/dtalk"
	"github.com/t2krew/dailynews/output/mail"
	"github.com/t2krew/dailynews/spider"
	"github.com/t2krew/dailynews/util"
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

		dd      = dtalk.New(ddrobot)                              // 钉钉
		mailbox = mail.New(email, password, nickname, host, port) // 邮件
	)

	fmt.Println(email, password, nickname, host, port)

	output.AddAdapter(dd)      // 钉钉发送实现
	output.AddAdapter(mailbox) // 邮件发送实现

	var list []map[string]string
	for _, s := range spider.Spiders {
		ret, err := s.Handler()
		if err != nil {
			log.Println(err)
			continue
		}
		list = append(list, ret...)
	}

	date := util.Today().Format("2006-01-02")
	data := output.Content{
		Subject: fmt.Sprintf("每日推荐 (%s)", date),
		Data: &output.Data{
			Date: date,
			List: list,
		},
		Mime: "text/html",
	}

	tplName := "daily" // 模板文件名
	for _, sender := range output.Outputers {
		err = sender.Send(tplName, receiver, data)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
}
