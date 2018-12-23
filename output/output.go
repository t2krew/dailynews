package output

type Outputer interface {
	Send(string, []string, Content) error
}

var Outputers []Outputer

func init() {
	Outputers = []Outputer{}
}

func AddAdapter(o Outputer) {
	Outputers = append(Outputers, o)
}

type Data struct {
	Date string              `json:"date"`
	List []map[string]string `json:"list"`
}

type Content struct {
	Subject string `json:"subject"`
	Data    *Data  `json:"data"`
	Mime    string `json:"content_type"`
}
