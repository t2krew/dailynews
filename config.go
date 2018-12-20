package main

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"strings"
	"sync"
)

var lock sync.Mutex
var viperLists map[string]*viper.Viper

func init() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	viper.SetEnvPrefix("ENV")
	_ = viper.BindEnv("app")

	suffix := ""
	env := viper.GetString("app")
	env = strings.TrimSpace(env)
	if len(env) > 0 {
		suffix = "-" + env
	}

	viper.SetDefault("configPath", fmt.Sprintf("%s/config%s", wd, suffix))
	viper.SetDefault("version", "1.0.0")

	viperLists = make(map[string]*viper.Viper)
}

func Configer(name string) (v *viper.Viper, err error) {
	v, ok := viperLists[name]
	if !ok {
		lock.Lock()
		defer lock.Unlock()
		v, ok = viperLists[name]
		if !ok {
			v, err = newViper(name)
			if err != nil {
				return
			}
			viperLists[name] = v
		}
	}
	return
}

func newViper(name string) (v *viper.Viper, err error) {
	v = viper.New()
	v.SetConfigType("json")
	v.AddConfigPath(viper.GetString("configPath"))
	v.SetConfigName(name)

	err = v.ReadInConfig()
	return
}
