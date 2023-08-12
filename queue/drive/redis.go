package drive

import (
	"context"
	"gin-example/config/config"
	"github.com/redis/go-redis/v9"
)

type _redis struct {
	Drive
	Produce *redis.Client
	Comsume *redis.Client
	Cfg     config.Redis
}

func (r *_redis) NewRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     r.Cfg.Addr,
		Password: r.Cfg.Password, // 没有密码，默认值
		DB:       r.Cfg.DB,       // 默认DB 0
	})
	return client
}

func NewRedis(topic string, cfg config.Redis, prefix, failureSuffix string) (Interface, error) {
	r := &_redis{Cfg: cfg}
	r.Name, r.FailureName = getQueueNames(topic, prefix, failureSuffix)
	r.Produce = r.NewRedis()
	r.Comsume = r.Produce
	return r, nil
}

func (r *_redis) Push(ctx context.Context, message string) (err error) {
	err = r.Produce.LPush(ctx, r.Name, message).Err()
	return
}

func (r *_redis) PushFailure(ctx context.Context, message string) (err error) {
	err = r.Produce.LPush(ctx, r.FailureName, message).Err()
	return
}

func (r *_redis) GetMessage(ctx context.Context) (string, error) {
	msg, err := r.Produce.BRPop(ctx, 0, r.Name).Result()
	if err != nil {
		return msg[1], err
	}
	return "", err
}

func (r *_redis) CommitMessage(context.Context) error {
	return nil
}
