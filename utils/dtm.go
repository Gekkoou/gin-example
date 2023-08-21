package utils

import (
	"database/sql"
	"fmt"
	"gin-example/global"
	"github.com/dtm-labs/client/dtmcli"
	"github.com/gin-gonic/gin"
)

func MustBarrierFromGin(c *gin.Context) *dtmcli.BranchBarrier {
	ti, err := dtmcli.BarrierFromQuery(c.Request.URL.Query())
	fmt.Println(err)
	return ti
}

func DbGet() *sql.DB {
	db, _ := global.DB.DB()
	return db
}
