package drive

import (
	"context"
	"gin-example/config/config"
	"github.com/redis/go-redis/v9"
)

type _redis struct {
	Produce *redis.Client
	Comsume *redis.Client
	Name    string
	Prefix  string
	Conf    config.Redis
}

func (r *_redis) NewRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     r.Conf.Addr,
		Password: r.Conf.Password, // 没有密码，默认值
		DB:       r.Conf.DB,       // 默认DB 0
	})
	return client
}

func NewRedis(topic string, conf config.Redis, prefix string) Interface {
	r := &_redis{Conf: conf}
	r.Produce = r.NewRedis()
	r.Comsume = r.Produce
	r.Name = topic
	r.Prefix = prefix
	return r
}

func (r *_redis) Push(ctx context.Context, message string) (err error) {
	if err = r.Produce.LPush(ctx, r.Prefix+r.Name, message).Err(); err != nil {
		return
	}
	return nil
}
func (r *_redis) GetMessage(ctx context.Context) (string, error) {
	msg, err := r.Produce.BRPop(ctx, 0, r.Prefix+r.Name).Result()
	return msg[1], err
}

func (r *_redis) CommitMessage(context.Context) error {
	return nil
}
