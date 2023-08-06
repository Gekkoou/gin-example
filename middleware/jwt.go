package middleware

import (
	"gin-example/global"
	"gin-example/model/common/response"
	"gin-example/service"
	"gin-example/utils"
	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get(global.Config.JWT.TokenKey)
		if token == "" {
			response.FailWithMessage("请登录", c)
			c.Abort()
			return
		}
		j := utils.NewJWT()
		claims, err := j.ParseToken(token)
		if err != nil {
			response.FailWithMessage("token无效, 请登录", c)
			c.Abort()
			return
		}

		jwtStr, err := service.JwtServiceApp.GetRedisJWT(int(claims.BaseClaims.Id))
		if err != nil || jwtStr != token {
			response.FailWithMessage("token无效, 请登录", c)
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}
