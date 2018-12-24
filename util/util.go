package util

import (
	"crypto/md5"
	"fmt"
	"net/url"
	"time"
)

var Loc *time.Location

func init() {
	var err error
	Loc, err = time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}
}

func Now() time.Time {
	return time.Now().In(Loc)
}

func Today() time.Time {
	n := Now()
	return time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, Loc)
}

func QueryStringToMap(query string) (ret map[string]string, err error) {

	ret = make(map[string]string)
	m, err := url.ParseQuery(query)

	if err != nil {
		return ret, err
	}

	for k, v := range m {
		if len(v) > 0 {
			ret[k] = v[0]
		}
	}

	return ret, nil
}

func Md5(str string) string {
	strByte := []byte(str)
	hash := fmt.Sprintf("%x", md5.Sum(strByte))
	return hash
}
