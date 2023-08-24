package model

import (
	"fmt"
	"gin-example/global"
	"gin-example/model/queue"
	"gin-example/utils"
	"gorm.io/gorm"
)

type User struct {
	global.MODEL
	Id       uint64 `gorm:"column:id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT" json:"id"`
	Name     string `gorm:"column:name;type:varchar(50);comment:用户名;NOT NULL" json:"name"`
	Password string `gorm:"column:password;type:varchar(100);comment:密码;NOT NULL" json:"password"`
	Phone    string `gorm:"column:phone;type:varchar(50);comment:手机号;NOT NULL" json:"phone"`
	Gold     int    `gorm:"column:gold;type:int(11);comment:金币;NOT NULL;default:0" json:"gold"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	fmt.Println("BeforeUpdate User", u)
	return
}

func (u *User) AfterUpdate(tx *gorm.DB) (err error) {
	fmt.Println("AfterUpdate User", u)
	u.DeleteCache(int(u.Id))
	return
}

func (u *User) BeforeDelete(tx *gorm.DB) (err error) {
	fmt.Println("BeforeDelete User", u)
	return
}
func (u *User) AfterDelete(tx *gorm.DB) (err error) {
	fmt.Println("AfterDelete User", u)
	u.DeleteCache(int(u.Id))
	return
}

func (u *User) DeleteCache(uid int) {
	key := utils.GetCacheKeyById(global.UserInfoDaoPrefixKey, uid)
	utils.DelCache(key)
	// 双删
	utils.AsynQueue(global.QueueDelCache, queue.DelCachePayload{Key: key})
}
