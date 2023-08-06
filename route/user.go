package router

import (
	"gin-example/api"
	"gin-example/middleware"
	"github.com/gin-gonic/gin"
)

type UserRouter struct{}

var UserRouterApp = new(UserRouter)

func (s *UserRouter) InitUserRouter(route *gin.RouterGroup) {
	r := route.Group("/user")
	{
		r.GET("/:id", api.GetUser)
		r.GET("/list", api.GetUserList)
	}
	middlewareRouter := route.Group("/user").Use(middleware.JWTAuth())
	{
		middlewareRouter.POST("/changePassword", api.ChangePassword)
		middlewareRouter.POST("/update", api.UpdateUser)
		middlewareRouter.POST("/delete", api.DeleteUser)
	}
	route.Use(middleware.JWTAuth()).GET("/logout", api.Logout)
}
