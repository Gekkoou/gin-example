package utils

import (
	"database/sql"
	"fmt"
	"gin-example/global"
	"github.com/dtm-labs/client/dtmcli"
	"github.com/dtm-labs/dtm/dtmutil"
	"github.com/gin-gonic/gin"
)

var DtmConf = dtmcli.DBConf{
	Driver:   "mysql",
	Host:     "127.0.0.1",
	Port:     3306,
	User:     "root",
	Password: "root",
}

func MustBarrierFromGin(c *gin.Context) *dtmcli.BranchBarrier {
	ti, err := dtmcli.BarrierFromQuery(c.Request.URL.Query())
	fmt.Println(err)
	return ti
}

func DbGet() *sql.DB {
	db, _ := global.DB.DB()
	return db
}

func DtmBarrierHandler(fn func(c *gin.Context) interface{}) gin.HandlerFunc {
	return dtmutil.WrapHandler(fn)
}

// DtmBarrierBusiFunc type for busi func
type DtmBarrierBusiFunc func(c *gin.Context, db dtmcli.DB) error

func Barrier(c *gin.Context, fn DtmBarrierBusiFunc) error {
	barrier := MustBarrierFromGin(c)
	return barrier.CallWithDB(DbGet(), func(tx *sql.Tx) error {
		return fn(c, tx)
	})
}
