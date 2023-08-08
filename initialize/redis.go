package initialize

import (
	"context"
	"gin-example/global"
	"github.com/redis/go-redis/v9"
)

func Redis() (client *redis.Client, err error) {
	config := global.Config.Redis
	client = redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password, // 没有密码，默认值
		DB:       config.DB,       // 默认DB 0
	})
	_, err = client.Ping(context.Background()).Result()
	return
}
