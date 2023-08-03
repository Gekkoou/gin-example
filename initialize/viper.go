package initialize

import (
	"fmt"
	"gin-example/global"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
)

func Viper(path ...string) *viper.Viper {
	env := os.Getenv("ENV")
	var config string
	switch env {
	case "debug":
		config = "./config/config.debug.yaml"
	case "":
		config = "./config/config.yaml"
	}
	v := viper.New()
	v.SetConfigFile(config)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err = v.Unmarshal(&global.Config); err != nil {
			fmt.Println(err)
		}
		Initialize(true)
	})
	if err = v.Unmarshal(&global.Config); err != nil {
		fmt.Println(err)
	}
	return v
}
