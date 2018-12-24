package main

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
	"sync"
)

var lock sync.Mutex
var viperLists map[string]*viper.Viper

func init() {
	viper.SetEnvPrefix("ENV")
	_ = viper.BindEnv("app")
	_ = viper.BindEnv("apppath")

	suffix := ""
	env := viper.GetString("app")
	env = strings.TrimSpace(env)
	if len(env) > 0 {
		suffix = "-" + env
	}

	configPath := fmt.Sprintf("config%s", suffix)

	apppath := viper.GetString("apppath")
	apppath = strings.TrimSpace(apppath)
	if len(apppath) > 0 {
		if string(apppath[len(apppath)-1:]) == "/" {
			apppath = apppath[:len(apppath)-1]
			viper.Set("apppath", apppath)
		}
		configPath = fmt.Sprintf("%s/%s", apppath, configPath)
	} else {
		configPath = fmt.Sprintf("./%s", configPath)
	}

	viper.SetDefault("configPath", configPath)
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
