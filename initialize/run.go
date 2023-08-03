package initialize

import (
	"fmt"
	"gin-example/global"
	"github.com/gin-gonic/gin"
)

func Run(r *gin.Engine) error {
	config := global.Config.Serve
	return r.Run(fmt.Sprintf("%s:%s", config.Ip, config.Addr))
}
