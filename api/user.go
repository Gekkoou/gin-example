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

// @Tags 公共
// @Summary 用户登录
// @Description 用户登录
// @Accept mpfd
// @Produce json
// @Param userReq.Login formData userReq.Login true "用户登录"
// @Success 200 {string} string {}
// @Router /login [post]
func Login(c *gin.Context) {
	var req userReq.Login
	if err := c.ShouldBind(&req); err != nil {
		response.FailWithMessage(utils.Translate(err), c)
		return
	}
	user, err := service.UserServiceApp.Login(req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(user, "登录成功", c)
}

// @Tags 公共
// @Summary 用户注册
// @Description 用户注册
// @Accept mpfd
// @Produce json
// @Param userReq.Register formData userReq.Register true "用户注册"
// @Success 200 {string} string {}
// @Router /register [post]
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

// @Tags 用户
// @Summary 修改密码
// @Description 修改密码
// @Accept mpfd
// @Produce json
// @Security ApiKeyAuth
// @Param userReq.ChangePassword formData userReq.ChangePassword true "修改密码"
// @Success 200 {string} string {}
// @Router /user/changePassword [post]
func ChangePassword(c *gin.Context) {
	var req userReq.ChangePassword
	if err := c.ShouldBind(&req); err != nil {
		response.FailWithMessage(utils.Translate(err), c)
		return
	}
	req.Id = utils.GetUserID(c)
	err := service.UserServiceApp.ChangePassword(req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("修改成功", c)
}

// @Tags 用户
// @Summary 获取用户
// @Description 获取用户
// @Accept mpfd
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "id"
// @Success 200 {string} string {}
// @Router /user/{id} [get]
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

// @Tags 用户
// @Summary 获取用户列表
// @Description 获取用户列表
// @Accept mpfd
// @Produce json
// @Security ApiKeyAuth
// @Param request.PageInfo query request.PageInfo true "获取用户列表"
// @Success 200 {string} string {}
// @Router /user/GetUserList [get]
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

// @Tags 用户
// @Summary 更新用户
// @Description 更新用户
// @Accept mpfd
// @Produce json
// @Security ApiKeyAuth
// @Param  userReq.UpdateUser formData userReq.UpdateUser true "更新用户"
// @Success 200 {string} string {}
// @Router /user/UpdateUser [post]
func UpdateUser(c *gin.Context) {
	var req userReq.UpdateUser
	if err := c.ShouldBind(&req); err != nil {
		response.FailWithMessage(utils.Translate(err), c)
		return
	}
	req.Id = utils.GetUserID(c)
	err := service.UserServiceApp.UpdateUser(req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// @Tags 用户
// @Summary 删除用户
// @Description 删除用户
// @Accept mpfd
// @Produce json
// @Security ApiKeyAuth
// @Param userReq.DeleteUser formData userReq.DeleteUser true "删除用户"
// @Success 200 {string} string {}
// @Router /user/deleteUser [post]
func DeleteUser(c *gin.Context) {
	var req userReq.DeleteUser
	if err := c.ShouldBind(&req); err != nil {
		response.FailWithMessage(utils.Translate(err), c)
		return
	}
	if utils.GetUserID(c) != req.Id {
		response.Fail(c)
	}
	err := service.UserServiceApp.DeleteUser(req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("注销成功", c)
}

// @Tags 用户
// @Summary 登出
// @Description 登出
// @Accept mpfd
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {string} string {}
// @Router /logout [get]
func Logout(c *gin.Context) {
	uid := utils.GetUserID(c)
	service.UserServiceApp.Logout(int(uid))
	response.OkWithMessage("退出成功", c)
}
