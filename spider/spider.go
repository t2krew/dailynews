package spider

import "sync"

type Spider interface {
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
	Date string `json:"date"`
	Url  string `json:"url"`
	List []map[string]string
}
