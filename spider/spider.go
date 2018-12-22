package spider

import "sync"

type Spider interface {
	Handler() ([]map[string]string, error)
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
