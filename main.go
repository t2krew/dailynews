package main

import (
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/x/bsonx"
	"github.com/t2krew/dailynews/mgo"
	"github.com/t2krew/dailynews/output"
	"github.com/t2krew/dailynews/output/dtalk"
	"github.com/t2krew/dailynews/output/mail"
	"github.com/t2krew/dailynews/spider"
	"github.com/t2krew/dailynews/util"
	"log"
	"time"
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
		ddrobot  = conf.GetStringSlice("dingding.robot")
		receiver = conf.GetStringSlice("mail.receiver")

		musername = conf.GetString("mongo.username")
		mpassword = conf.GetString("mongo.password")
		mhost     = conf.GetString("mongo.host")
		mport     = conf.GetInt("mongo.port")
		mdatabase = conf.GetString("mongo.database")

		interval = time.Duration(conf.GetInt("interval")) * time.Second
	)

	cli, err := mgo.New(mhost, mport, musername, mpassword, mdatabase)
	if err != nil {
		log.Println(err)
		return
	}

	col := cli.Collection("dailynews")
	_, err = col.Indexs([]mongo.IndexModel{
		{Keys: bsonx.Doc{{"md5", bsonx.Int32(-1)}}},
		{Keys: bsonx.Doc{{"date", bsonx.Int32(-1)}}},
	})
	if err != nil {
		log.Println(err)
		return
	}

	var (
		dd      = dtalk.New(ddrobot)                              // 钉钉
		mailbox = mail.New(email, password, nickname, host, port) // 邮件
	)

	output.AddAdapter(dd)      // 钉钉发送实现
	output.AddAdapter(mailbox) // 邮件发送实现

	for {
		var list []*spider.Data
		for _, s := range spider.Spiders {
			ret, err := s.Handler()
			if err != nil {
				log.Println(err)
				continue
			}
			list = append(list, ret)
		}

		date := util.Today().Format("2006-01-02")

		hash := util.Md5(list[0].Url)

		ret, err := col.FindOne(bson.M{"md5": hash})
		if err != nil {
			time.Sleep(interval)
			continue
		}

		if len(ret) == 0 {
			data := output.Content{
				Subject: fmt.Sprintf("今日推荐 (%s)", date),
				Data: &output.Data{
					Date: date,
					List: list[0].List,
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

			var insertdata = map[string]interface{}{
				"md5":     hash,
				"date":    list[0].Date,
				"url":     list[0].Url,
				"content": list[0].List,
			}

			ret, err := col.InsertOne(insertdata)
			if err != nil {
				log.Println(err)
			}

			log.Printf("insert result: %v\n", ret)
		} else {
			log.Println("最新数据已存在")
		}

		time.Sleep(interval)
	}
}
