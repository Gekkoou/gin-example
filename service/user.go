package service

import (
	"errors"
	"gin-example/dao"
	"gin-example/global"
	"gin-example/model"
	"gin-example/model/common/request"
	userReq "gin-example/model/user/request"
	userRes "gin-example/model/user/response"
	"gin-example/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
)

type UserService struct{}

var UserServiceApp = new(UserService)

// Login 用户登录
func (userService *UserService) Login(req userReq.Login) error {
	return nil
}

// Register 用户注册
func (userService *UserService) Register(req userReq.Register) (err error) {
	var u model.User
	if !errors.Is(global.DB.Where("name = ?", req.Name).First(&u).Error, gorm.ErrRecordNotFound) {
		return errors.New("用户名已存在")
	}
	var user = model.User{
		Name:     req.Name,
		Password: utils.BcryptHash(req.Password),
		Phone:    req.Phone,
	}
	if err = global.DB.Select("Name", "Password", "Phone").Create(&user).Error; err != nil {
		return err
	}
	return
}

// ChangePassword 修改密码
func (userService *UserService) ChangePassword(req userReq.ChangePassword) (err error) {
	var u model.User
	if errors.Is(global.DB.Where("id = ?", req.Id).First(&u).Error, gorm.ErrRecordNotFound) {
		return errors.New("用户不存在")
	}
	if ok := utils.BcryptCheck(u.Password, req.Password); !ok {
		return errors.New("原密码错误")
	}
	u.Password = utils.BcryptHash(req.NewPassword)
	if err = global.DB.Save(&u).Error; err != nil {
		return err
	}
	return
}

// GetUser 查询用户
func (userService *UserService) GetUser(req request.GetById) (model.User, error) {
	u, err := dao.UserDaoApp.GetInfo(int(req.Id))
	return u, err
}

// GetUserList 获取用户列表
func (userService *UserService) GetUserList(req request.PageInfo) (list interface{}, total int64, err error) {
	var userList []userRes.GetUserList
	limit := req.PageSize
	offset := req.PageSize * (req.Page - 1)
	db := global.DB.Model(model.User{})
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Find(&userList).Error
	return userList, total, err
}

// UpdateUser 更新用户
func (userService *UserService) UpdateUser(req userReq.UpdateUser) error {

	// 加锁
	key := "UserServiceUpdateUser:" + strconv.Itoa(int(req.Id))
	lockRand, err := utils.RedisLock.Lock(key, 30)
	if err != nil {
		global.Log.Error(err.Error())
		return err
	}
	defer utils.RedisLock.UnLock(key, lockRand)

	return global.DB.Transaction(func(tx *gorm.DB) error {
		var u model.User
		// for update 加锁
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Model(&u).Where("id = ?", req.Id).First(&u).Error; err != nil {
			return err
		}

		// 查询用户名是否已存在
		var u1 model.User
		if !errors.Is(tx.Where("name = ? AND Id <> ?", req.Name, u.Id).First(&u1).Error, gorm.ErrRecordNotFound) {
			return errors.New("用户名已存在")
		}

		// 修改用户信息
		upUserMap := map[string]any{
			"Name":  req.Name,
			"Phone": req.Phone,
		}
		if err := tx.Model(&u).Where("id = ?", req.Id).Updates(upUserMap).Error; err != nil {
			return err
		}
		return nil
	})
}

// DeleteUser 删除用户
func (userService *UserService) DeleteUser(req userReq.DeleteUser) (err error) {
	var u model.User
	if err = global.DB.Where("id = ?", req.Id).Delete(&u).Error; err != nil {
		return err
	}
	return
}
