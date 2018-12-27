package mgo

import (
	"github.com/t2krew/dailynews/config"
)

var Client *Mongo

func init() {
	conf, err := config.Configer("app")
	if err != nil {
		panic(err)
	}

	var (
		musername = conf.GetString("mongo.username")
		mpassword = conf.GetString("mongo.password")
		mhost     = conf.GetString("mongo.host")
		mport     = conf.GetInt("mongo.port")
		mdatabase = conf.GetString("mongo.database")
	)

	Client, err = New(mhost, mport, musername, mpassword, mdatabase)
	if err != nil {
		panic(err)
	}
}
