package service

import (
	"context"
	"gin-example/global"
	"gin-example/utils"
	"strconv"
)

type JwtService struct{}

var JwtServiceApp = new(JwtService)

func (jwtService *JwtService) GetRedisJWT(uid int) (string, error) {
	jwtStr, err := global.Redis.Get(context.Background(), global.AccessTokenPrefixKey+strconv.Itoa(uid)).Result()
	return jwtStr, err
}

func (jwtService *JwtService) SetRedisJWT(uid int, jwt string) error {
	expiresTime, err := utils.ParseDuration(global.Config.JWT.ExpiresTime)
	if err != nil {
		return err
	}
	err = global.Redis.Set(context.Background(), global.AccessTokenPrefixKey+strconv.Itoa(uid), jwt, expiresTime).Err()
	return err
}

func (jwtService *JwtService) DelRedisJWT(uid int) {
	global.Redis.Del(context.Background(), global.AccessTokenPrefixKey+strconv.Itoa(uid)).Result()
}
