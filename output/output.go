package output

type Outputer interface {
	Name() string
	Send(string, []string, Content) error
}

type Output struct {
	OutputName string
}

var Outputers []Outputer

func init() {
	Outputers = []Outputer{}
}

func AddAdapter(o Outputer) {
	Outputers = append(Outputers, o)
}

type Data struct {
	Title string              `json:"title"`
	Date  string              `json:"date"`
	List  []map[string]string `json:"list"`
}

type Content struct {
	Subject string `json:"subject"`
	Data    *Data  `json:"data"`
	Mime    string `json:"content_type"`
}
