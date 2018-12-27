package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/x/bsonx"
	"github.com/t2krew/dailynews/config"
	"github.com/t2krew/dailynews/mgo"
	"github.com/t2krew/dailynews/output"
	"github.com/t2krew/dailynews/output/dtalk"
	"github.com/t2krew/dailynews/output/mail"
	"github.com/t2krew/dailynews/spider"
	"github.com/t2krew/dailynews/util"
)

func main() {
	conf, err := config.Configer("app")
	if err != nil {
		panic(err)
	}

	subConf, err := config.Configer("subscribe")
	if err != nil {
		panic(err)
	}

	var (
		email    = conf.GetString("mail.email")
		password = conf.GetString("mail.password")
		host     = conf.GetString("mail.host")
		port     = conf.GetInt("mail.port")
		nickname = conf.GetString("mail.nickname")

		dailyCollection   = conf.GetString("mongo.daily_collection")
		articleCollection = conf.GetString("mongo.article_collection")

		interval = 1 * time.Minute // 1分钟 不推荐修改，会影响内存数据判断今日是否推送
	)

	dCol := mgo.Client.Collection(dailyCollection)
	_, err = dCol.Indexs([]mongo.IndexModel{
		{Keys: bsonx.Doc{{"md5", bsonx.Int32(-1)}}},
		{Keys: bsonx.Doc{{"date", bsonx.Int32(-1)}}},
		{Keys: bsonx.Doc{{"source", bsonx.Int32(-1)}}},
	})
	if err != nil {
		log.Println(err)
		return
	}

	aCol := mgo.Client.Collection(articleCollection)
	_, err = aCol.Indexs([]mongo.IndexModel{
		{Keys: bsonx.Doc{{"md5", bsonx.Int32(-1)}}},
		{Keys: bsonx.Doc{{"title", bsonx.Int32(-1)}}},
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

	log.Println("[START TO RUN]")

	for {

		now := util.Now()
		// 22:00 - 10:00 区间不推送
		if now.Hour() < 10 || now.Hour() > 21 {
			continue
		}

		var wg sync.WaitGroup
		for _, s := range spider.Spiders {

			if s.IsDone() {
				log.Printf("读取内存，[%s] [今日已推送]\n", s.Name())
				continue
			}

			tResult, err := dCol.FindOne(bson.M{"source": s.Name(), "date": now.Format("2006-01-02")})
			if err != nil {
				return
			}

			if len(tResult) > 0 {
				log.Printf("读取数据库，[%s] [今日已推送]\n", s.Name())
				s.SetDone() // 设置今日已推
				continue
			}

			wg.Add(1)
			go func(s spider.Spider) {
				defer wg.Done()

				var sub = subConf.GetStringMapStringSlice(s.Name())

				ret, err := s.Handler()
				if err != nil {
					log.Println(err)
					return
				}

				if len(ret.List) == 0 {
					return
				}

				hash := util.Md5(ret.Url)
				date := util.Today().Format("2006-01-02")

				result, err := dCol.FindOne(bson.M{"md5": hash})
				if err != nil {
					log.Println(err)
					return
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
							return
						}
					}

					s.SetDone() // 设置今日已推

					go func() {
						var insertdata = map[string]interface{}{
							"md5":     hash,
							"date":    ret.Date,
							"url":     ret.Url,
							"source":  s.Name(),
							"content": ret.List,
						}

						_, err := dCol.InsertOne(insertdata)
						if err != nil {
							log.Println(err)
							return
						}
						log.Println("daily insert done")
					}()

					go func() {
						for _, item := range ret.List {
							var insertdata = map[string]interface{}{
								"url":   item["link"],
								"title": item["title"],
								"md5":   util.Md5(item["link"]),
							}
							_, err := aCol.InsertOne(insertdata)
							if err != nil {
								log.Println(err)
								return
							}
						}
						log.Println("artiles insert done")
					}()

				} else {
					log.Printf("读取数据库，[%s] [最新数据已爬取]\n", s.Name())
				}
			}(s)
		}
		wg.Wait()

		time.Sleep(interval)
	}
}
