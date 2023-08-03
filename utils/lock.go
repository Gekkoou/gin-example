package utils

import (
	"context"
	"errors"
	"gin-example/global"
	"math/rand"
	"strconv"
	"time"
)

type redisLock struct{}

var RedisLock = new(redisLock)

func (lock *redisLock) Lock(key string, outtime int) (string, error) {
	randNumber := lock.getRand()
	ok, err := global.Redis.SetNX(context.Background(), key, randNumber, time.Duration(outtime)*time.Second).Result()
	if err != nil {
		return "", err
	}
	if !ok {
		return "", errors.New("操作繁忙")
	}
	return strconv.Itoa(randNumber), err
}

func (lock *redisLock) getRand() int {
	return rand.Intn(99999)
}

func (lock *redisLock) UnLock(key string, value string) bool {
	if key == "" {
		return false
	}
	randNumber := global.Redis.Get(context.Background(), key).Val()
	if randNumber != value {
		return false
	}
	global.Redis.Del(context.Background(), key).Val()
	return true
}
