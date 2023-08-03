package initialize

import (
	"gin-example/global"
	"gin-example/model"
)

func Initialize(reload bool) {
	if !reload {
		// VIPER
		global.VP = Viper()
	}
	// 日志
	global.Log = Zap()
	// 初始化DB
	global.DB = GormMysql()
	// 迁移
	global.DB.AutoMigrate(model.User{})
	// 初始化Redis
	global.Redis = Redis()
}
