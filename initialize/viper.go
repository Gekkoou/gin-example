package initialize

import (
	"gin-example/global"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
)

func Viper() (v *viper.Viper, err error) {
	env := os.Getenv("ENV")
	var config string
	switch env {
	case "debug":
		config = "./config/config.debug.yaml"
	case "":
		config = "./config/config.yaml"
	}
	v = viper.New()
	v.SetConfigFile(config)
	v.SetConfigType("yaml")
	err = v.ReadInConfig()
	if err != nil {
		return
	}
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		global.Log.Info("config file changed:" + e.Name)
		if err = v.Unmarshal(&global.Config); err != nil {
			global.Log.Error("config file changed:" + e.Name)
		}
		Initialize(true)
	})
	err = v.Unmarshal(&global.Config)
	return
}
