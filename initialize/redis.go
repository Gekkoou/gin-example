package initialize

import (
	"gin-example/global"
	"github.com/redis/go-redis/v9"
)

func Redis() *redis.Client {
	config := global.Config.Redis
	return redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password, // 没有密码，默认值
		DB:       config.DB,       // 默认DB 0
	})
}
