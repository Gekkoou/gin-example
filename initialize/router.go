package initialize

import (
	"gin-example/api"
	"gin-example/docs"
	router "gin-example/route"
	"github.com/gin-gonic/gin"

	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Routes() *gin.Engine {
	r := gin.Default()
	r.GET("/", api.Index)
	r.POST("/login", api.Login)
	r.POST("/register", api.Register)
	routerGroup := r.Group("")
	{
		router.UserRouterApp.InitUserRouter(routerGroup)
	}
	docs.SwaggerInfo.BasePath = "/"
	// 文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}
