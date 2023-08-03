package initialize

import (
	"gin-example/global"
	"gin-example/model"
)

// 数据库迁移
func RegisterTables() error {
	return global.DB.AutoMigrate(model.User{})
}
