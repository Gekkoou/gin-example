package router

import (
	"gin-example/api"
	"github.com/gin-gonic/gin"
)

type UserRouter struct{}

var UserRouterApp = new(UserRouter)

func (s *UserRouter) InitUserRouter(route *gin.RouterGroup) {
	r := route.Group("/user")
	{
		r.POST("/login", api.Login)
		r.POST("/register", api.Register)
		r.GET("/:id", api.GetUser)
		r.GET("/list", api.GetUserList)

		// TODO: 需要JWT中间件
		r.POST("/changePassword", api.ChangePassword)
		r.POST("/update", api.UpdateUser)
		r.POST("/delete", api.DeleteUser)
	}
}
