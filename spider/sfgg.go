package spider

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/t2krew/dailynews/config"
	"github.com/t2krew/dailynews/mgo"
	"github.com/t2krew/dailynews/util"
	"log"
)

const SpiderSfgg = "sfgg"

type Sfgg struct {
	done bool
	name string
	root string
}

const COUNT = 5

var articleCollection string

func init() {
	sfgg := NewSfgg("https://segmentfault.com/hottest")
	register("sfgg", sfgg)

	conf, err := config.Configer("app")
	if err != nil {
		panic(err)
	}

	articleCollection = conf.GetString("mongo.article_collection")

}

func NewSfgg(root string) *Sfgg {
	return &Sfgg{
		name: SpiderSfgg,
		root: root,
	}
}

func (sf *Sfgg) Name() string {
	return sf.name
}

func (sf *Sfgg) IsDone() bool {
	now := util.Now()
	if now.Hour() == 0 && now.Minute() == 0 {
		sf.done = false
	}
	return sf.done
}

func (sf *Sfgg) SetDone() {
	sf.done = true
}

func (sf *Sfgg) Handler() (ret *Data, err error) {
	var articleList []ArticleItem
	var date = util.Today().Format("2006-01-02")

	var cnt = 0
	var aCol = mgo.Client.Collection(articleCollection)

	dailyCollector := colly.NewCollector()
	dailyCollector.OnHTML("div.news-list", func(e *colly.HTMLElement) {
		e.ForEach(".news__item-info a+a", func(i int, element *colly.HTMLElement) {
			href := element.Attr("href")
			href = fmt.Sprintf("https://segmentfault.com%s", href)
			title := element.DOM.Find("h4").Eq(0).Text()

			hash := util.Md5(href)

			result, err := aCol.FindOne(bson.M{"md5": hash})
			if err != nil {
				log.Println(err)
				return
			}

			if len(result) > 0 {
				return
			}

			if cnt >= COUNT {
				return
			}
			cnt++

			item := ArticleItem{href, title}
			articleList = append(articleList, item)
		})
	})

	err = dailyCollector.Visit(sf.root)
	if err != nil {
		log.Println(err)
		return
	}

	b, err := json.Marshal(articleList)
	if err != nil {
		log.Println(err)
		return
	}

	var m []map[string]string
	err = json.Unmarshal(b, &m)
	if err != nil {
		log.Println(err)
		return
	}

	return &Data{
		List:   m,
		Spider: sf.Name(),
		Title:  "Segmentfault",
		Date:   date,
		Url:    fmt.Sprintf("%s?date=%s", sf.root, date),
	}, nil
}
