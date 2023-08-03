package model

import (
	"fmt"
	"gin-example/global"
	"gin-example/utils"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

type User struct {
	Id        uint64                `gorm:"column:id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT" json:"id"`
	Name      string                `gorm:"column:name;type:varchar(50);comment:用户名;NOT NULL" json:"name"`
	Password  string                `gorm:"column:password;type:varchar(100);comment:密码;NOT NULL" json:"password"`
	Phone     string                `gorm:"column:phone;type:varchar(50);comment:手机号;NOT NULL" json:"phone"`
	CreatedAt int                   `gorm:"column:created_at;type:int(11) unsigned;comment:创建时间;NOT NULL;autoCreateTime" json:"created_at"`
	UpdatedAt int                   `gorm:"column:updated_at;type:int(11) unsigned;comment:更新时间;NOT NULL;autoUpdateTime" json:"updated_at"`
	DeletedAt soft_delete.DeletedAt `gorm:"column:deleted_at;type:int(11) unsigned;comment:删除时间;NOT NULL" json:"-"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	fmt.Println("BeforeSave User", u)
	return
}

func (u *User) AfterSave(tx *gorm.DB) (err error) {
	fmt.Println("AfterSave User", u)
	utils.DelCache(utils.GetCacheKeyById(global.UserInfoDaoPrefixKey, int(u.Id)))
	return
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	fmt.Println("BeforeUpdate User", u)
	return
}

func (u *User) AfterUpdate(tx *gorm.DB) (err error) {
	fmt.Println("AfterUpdate User", u)
	utils.DelCache(utils.GetCacheKeyById(global.UserInfoDaoPrefixKey, int(u.Id)))
	return
}

func (u *User) BeforeDelete(tx *gorm.DB) (err error) {
	fmt.Println("BeforeDelete User", u)
	return
}
func (u *User) AfterDelete(tx *gorm.DB) (err error) {
	fmt.Println("AfterDelete User", u)
	utils.DelCache(utils.GetCacheKeyById(global.UserInfoDaoPrefixKey, int(u.Id)))
	return
}
