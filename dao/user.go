package dao

import (
	"context"
	"errors"
	"gin-example/global"
	"gin-example/model"
	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type UserDao struct {
	group singleflight.Group
}

var UserDaoApp = new(UserDao)

func (userDao *UserDao) SetInfo(user model.User) {
	json, err := sonic.MarshalString(user)
	if err == nil {
		global.Redis.Set(context.Background(), userDao.GetKey(int(user.Id)), json, time.Duration(global.UserInfoDaoTtl)*time.Second).Result()
	}
}

func (userDao *UserDao) SetNA(uid int) {
	global.Redis.Set(context.Background(), userDao.GetKey(uid), "NA", time.Duration(global.UserInfoDaoTtl)*time.Second).Result()
}

func (userDao *UserDao) GetInfo(uid int) (model.User, error) {
	key := strconv.Itoa(uid)
	// 防缓存击穿
	r, err, _ := userDao.group.Do(key, func() (interface{}, error) {
		return userDao.DoGetInfo(uid)
	})
	return r.(model.User), err
}

func (userDao *UserDao) DoGetInfo(uid int) (user model.User, err error) {
	json, err := global.Redis.Get(context.Background(), userDao.GetKey(uid)).Result()
	if err == redis.Nil {
		// 查询mysql
		if err = global.DB.Model(&user).Where("id = ?", uid).First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 防缓存穿透
				userDao.SetNA(uid)
				return user, nil
			}
			return
		}
		userDao.SetInfo(user)
	} else if err == nil {
		if json == "NA" {
			return
		}
		err = sonic.UnmarshalString(json, &user)
	}
	return
}

func (userDao *UserDao) DelInfo(uid int) {
	global.Redis.Del(context.Background(), userDao.GetKey(uid)).Result()
}

func (userDao *UserDao) GetKey(uid int) string {
	return global.UserInfoDaoPrefixKey + strconv.Itoa(uid)
}
