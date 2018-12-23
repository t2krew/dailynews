package spider

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"regexp"
)

type Gocn struct {
	root string
}

type DailyItem struct {
	Date string `json:"date"`
	Link string `json:"link"`
}

type ArticleItem struct {
	Link  string `json:"link"`
	Title string `json:"title"`
}

func init() {
	gocn := New("https://gocn.vip/explore/category-14")
	register("gocn", gocn)

}

func New(root string) *Gocn {
	return &Gocn{root: root}
}

func (g *Gocn) Handler() (ret []map[string]string, err error) {
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
			item := DailyItem{match[1], element.Attr("href")}
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

	err = articleCollector.Visit(dailyList[0].Link)
	if err != nil {
		return
	}

	b, err := json.Marshal(articleList)
	if err != nil {
		return
	}

	err = json.Unmarshal(b, &ret)
	return
}
