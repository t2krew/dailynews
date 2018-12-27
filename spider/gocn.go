package spider

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/t2krew/dailynews/util"
	"regexp"
)

const SpiderGocn = "gocn"

type Gocn struct {
	done bool
	name string
	root string
}

func init() {
	gocn := NewGocn("https://gocn.vip/explore/category-14")
	register("gocn", gocn)
}

func NewGocn(root string) *Gocn {
	return &Gocn{
		name: SpiderGocn,
		root: root,
	}
}

func (g *Gocn) Name() string {
	return g.name
}

func (g *Gocn) IsDone() bool {
	now := util.Now()
	if now.Hour() == 0 && now.Minute() == 0 {
		g.done = false
	}
	return g.done
}

func (g *Gocn) SetDone() {
	g.done = true
}

func (g *Gocn) Handler() (ret *Data, err error) {
	dateReg, err := regexp.Compile("\\((\\d{4}-\\d{1,2}-\\d{1,2})\\)$")
	if err != nil {
		return
	}

	var dailyList []DailyItem
	dailyCollector := colly.NewCollector()
	dailyCollector.OnHTML("div.aw-item", func(e *colly.HTMLElement) {
		e.ForEach("h4>a", func(i int, element *colly.HTMLElement) {
			match := dateReg.FindStringSubmatch(element.Text)
			if len(match) == 0 {
				return
			}
			date := match[1]
			item := DailyItem{date, element.Attr("href")}
			dailyList = append(dailyList, item)
		})
	})

	err = dailyCollector.Visit(g.root)
	if err != nil {
		return
	}

	var articleList []ArticleItem
	articleCollector := colly.NewCollector()
	articleCollector.OnHTML(".content", func(e *colly.HTMLElement) {
		e.ForEach("ol li", func(i int, element *colly.HTMLElement) {
			var item ArticleItem
			element.DOM.Contents().Each(func(i int, selection *goquery.Selection) {
				if i == 0 {
					item.Title = selection.Text()
				}
				if i == 1 {
					if href, ok := selection.Attr("href"); ok {
						item.Link = href
					}
				}
			})
			articleList = append(articleList, item)
		})
	})

	var target = dailyList[0]

	err = articleCollector.Visit(target.Link)
	if err != nil {
		return
	}

	b, err := json.Marshal(articleList)
	if err != nil {
		return
	}

	var m []map[string]string
	err = json.Unmarshal(b, &m)
	if err != nil {
		return
	}

	return &Data{
		List:   m,
		Spider: g.Name(),
		Title:  "Golang",
		Date:   target.Date,
		Url:    target.Link,
	}, nil
}
