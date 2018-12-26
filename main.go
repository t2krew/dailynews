package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/x/bsonx"
	"github.com/t2krew/dailynews/mgo"
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

	subConf, err := Configer("subscribe")
	if err != nil {
		panic(err)
	}

	var (
		email    = conf.GetString("mail.email")
		password = conf.GetString("mail.password")
		host     = conf.GetString("mail.host")
		port     = conf.GetInt("mail.port")
		nickname = conf.GetString("mail.nickname")

		musername   = conf.GetString("mongo.username")
		mpassword   = conf.GetString("mongo.password")
		mhost       = conf.GetString("mongo.host")
		mport       = conf.GetInt("mongo.port")
		mdatabase   = conf.GetString("mongo.database")
		mcollection = conf.GetString("mongo.collection")

		interval = time.Duration(conf.GetInt("interval")) * time.Second
	)

	cli, err := mgo.New(mhost, mport, musername, mpassword, mdatabase)
	if err != nil {
		log.Println(err)
		return
	}

	col := cli.Collection(mcollection)
	_, err = col.Indexs([]mongo.IndexModel{
		{Keys: bsonx.Doc{{"md5", bsonx.Int32(-1)}}},
		{Keys: bsonx.Doc{{"date", bsonx.Int32(-1)}}},
	})
	if err != nil {
		log.Println(err)
		return
	}

	var (
		dingding = dtalk.New()
		mailbox  = mail.New(email, password, nickname, host, port)
	)

	output.AddAdapter(mailbox)
	output.AddAdapter(dingding)

	for {
		var wg sync.WaitGroup
		for _, s := range spider.Spiders {
			wg.Add(1)
			go func(s spider.Spider) {
				defer wg.Done()

				var sub = subConf.GetStringMapStringSlice(s.Name())

				ret, err := s.Handler()
				if err != nil {
					log.Println(err)
				}

				hash := util.Md5(ret.Url)
				date := util.Today().Format("2006-01-02")

				result, err := col.FindOne(bson.M{"md5": hash})
				if err != nil {
					log.Println(err)
				}

				if len(result) == 0 {
					data := output.Content{
						Subject: fmt.Sprintf("%s 今日推荐 (%s)", ret.Title, date),
						Data: &output.Data{
							Date:  date,
							List:  ret.List,
							Title: ret.Title,
						},
					}

					tplName := "daily"
					for _, sender := range output.Outputers {
						receiver, ok := sub[sender.Name()]
						if !ok {
							continue
						}
						err = sender.Send(tplName, receiver, data)
						if err != nil {
							fmt.Println(err)
						}
					}

					var insertdata = map[string]interface{}{
						"md5":     hash,
						"date":    ret.Date,
						"url":     ret.Url,
						"source":  s.Name(),
						"content": ret.List,
					}

					ret, err := col.InsertOne(insertdata)
					if err != nil {
						log.Println(err)
					}

					log.Printf("insert result: %v\n", ret)
				} else {
					log.Println("最新数据已存在")
				}
			}(s)
		}
		wg.Wait()

		time.Sleep(interval)
	}
}
