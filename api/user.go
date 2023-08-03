package api

import (
	"gin-example/model/common/request"
	"gin-example/model/common/response"
	userReq "gin-example/model/user/request"
	userRes "gin-example/model/user/response"
	"gin-example/service"
	"gin-example/utils"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var req userReq.Login
	if err := c.ShouldBind(&req); err != nil {
		response.FailWithMessage(utils.Translate(err), c)
		return
	}
	err := service.UserServiceApp.Login(req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("登录成功", c)
}

func Register(c *gin.Context) {
	var req userReq.Register
	if err := c.ShouldBind(&req); err != nil {
		response.FailWithMessage(utils.Translate(err), c)
		return
	}
	err := service.UserServiceApp.Register(req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("注册成功", c)
}

func ChangePassword(c *gin.Context) {
	var req userReq.ChangePassword
	if err := c.ShouldBind(&req); err != nil {
		response.FailWithMessage(utils.Translate(err), c)
		return
	}
	// TODO: 从 JWT 中提取 user id赋值到req.id
	req.Id = 0
	err := service.UserServiceApp.ChangePassword(req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("修改成功", c)
}

func GetUser(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindUri(&req); err != nil {
		response.FailWithMessage(utils.Translate(err), c)
		return
	}
	user, err := service.UserServiceApp.GetUser(req)
	if err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithData(gin.H{"userInfo": userRes.GetUser{
		ID:    user.Id,
		Name:  user.Name,
		Phone: user.Phone,
	}}, c)
}

func GetUserList(c *gin.Context) {
	var pageInfo request.PageInfo
	if err := c.ShouldBind(&pageInfo); err != nil {
		response.FailWithMessage(utils.Translate(err), c)
		return
	}
	list, total, err := service.UserServiceApp.GetUserList(pageInfo)
	if err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithData(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, c)
}

func UpdateUser(c *gin.Context) {
	var req userReq.UpdateUser
	if err := c.ShouldBind(&req); err != nil {
		response.FailWithMessage(utils.Translate(err), c)
		return
	}
	// TODO: 从 JWT 中提取 user id赋值到req.id
	req.Id = 0
	err := service.UserServiceApp.UpdateUser(req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

func DeleteUser(c *gin.Context) {
	var req userReq.DeleteUser
	if err := c.ShouldBind(&req); err != nil {
		response.FailWithMessage(utils.Translate(err), c)
		return
	}
	// TODO: 从 JWT 中提取 user id赋值到req.id
	req.Id = 0
	// TODO: 判断req.id 与 JWT 中提取 user id是否相等, 不相等则报错
	err := service.UserServiceApp.DeleteUser(req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("注销成功", c)
}
