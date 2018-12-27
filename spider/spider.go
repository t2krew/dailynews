package spider

import "sync"

type Spider interface {
	Name() string
	IsDone() bool
	SetDone()
	Handler() (*Data, error)
}

var mu sync.Mutex
var Spiders = make(map[string]Spider)

func register(name string, spider Spider) {
	_, ok := Spiders[name]
	if !ok {
		mu.Lock()
		defer mu.Unlock()
		_, ok = Spiders[name]
		if !ok {
			Spiders[name] = spider
		}
	}
}

type Data struct {
	Url    string              `json:"url"`
	Date   string              `json:"date"`
	Title  string              `json:"title"`
	Spider string              `json:"spider"`
	List   []map[string]string `json:"list"`
}

type DailyItem struct {
	Date string `json:"date"`
	Link string `json:"link"`
}

type ArticleItem struct {
	Link  string `json:"link"`
	Title string `json:"title"`
}
