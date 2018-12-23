package dingding

import (
	"encoding/json"
	"fmt"
	"github.com/t2krew/daily/output"
	"github.com/t2krew/daily/util"
	"log"
	"time"
)

type dingding struct {
	robot string
}

func New(robot string) *dingding {
	return &dingding{robot: robot}
}

func init() {
	httpClient = util.NewClient()
}

var httpClient *util.Client

func (d *dingding) Robot() string {
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

	ret, err := httpClient.Post(d.Robot(), b, 5*time.Second, map[string]string{
		"Content-Type": "application/json",
	})
	if err != nil {
		return
	}
	log.Printf("[task][钉钉] done, ret: %s\n", string(ret))
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
