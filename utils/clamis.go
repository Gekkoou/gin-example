package utils

import (
	"gin-example/global"
	"gin-example/model/common/request"
	"github.com/gin-gonic/gin"
)

func GetClaims(c *gin.Context) (*request.CustomClaims, error) {
	token := c.Request.Header.Get(global.Config.JWT.TokenKey)
	j := NewJWT()
	claims, err := j.ParseToken(token)
	return claims, err
}

func GetUserID(c *gin.Context) uint64 {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return 0
		} else {
			return cl.BaseClaims.Id
		}
	} else {
		customClaims := claims.(*request.CustomClaims)
		return customClaims.BaseClaims.Id
	}
}

func GetUserName(c *gin.Context) string {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return ""
		} else {
			return cl.Name
		}
	} else {
		customClaims := claims.(*request.CustomClaims)
		return customClaims.BaseClaims.Name
	}
}
