package dtalk

import (
	"encoding/json"
	"fmt"
	"github.com/t2krew/dailynews/output"
	"github.com/t2krew/dailynews/util"
	"log"
	"time"
)

type dingding struct {
	robot []string
}

func New(robot []string) *dingding {
	return &dingding{robot: robot}
}

func init() {
	httpClient = util.NewClient()
}

var httpClient *util.Client

func (d *dingding) Robot() []string {
	return d.robot
}

func (d *dingding) Send(tplname string, receiver []string, content output.Content) (err error) {
	initdata := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"title": content.Subject,
			"text":  mdfill(content),
		},
	}

	b, err := json.Marshal(initdata)
	if err != nil {
		return
	}

	for idx, url := range d.Robot() {
		ret, err := httpClient.Post(url, b, 5*time.Second, map[string]string{
			"Content-Type": "application/json",
		})
		if err != nil {
			log.Println(err)
			continue
		}
		log.Printf("[task][钉钉][%d] done, ret: %s\n", idx+1, string(ret))
	}
	return

}

func mdfill(content output.Content) (markdown string) {
	sub := content.Subject
	list := content.Data.List
	markdown = fmt.Sprintf("## %s\r\n", sub)
	for idx, item := range list {
		markdown += fmt.Sprintf("#### %d) [%s](%s)\r\n", idx+1, item["title"], item["link"])
	}
	return
}
