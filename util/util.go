package util

import "time"

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
