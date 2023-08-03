package initialize

import (
	"gin-example/api"
	router "gin-example/route"
	"github.com/gin-gonic/gin"
)

func Routes() *gin.Engine {
	r := gin.Default()
	r.GET("/", api.Index)
	routerGroup := r.Group("")
	{
		router.UserRouterApp.InitUserRouter(routerGroup)
	}
	return r
}
