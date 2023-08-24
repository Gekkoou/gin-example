package utils

import (
	"database/sql"
	"fmt"
	"gin-example/global"
	"github.com/dtm-labs/client/dtmcli"
	"github.com/dtm-labs/dtm/dtmutil"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strconv"
)

func DbConf() dtmcli.DBConf {
	conf := global.Config.Mysql
	port, _ := strconv.ParseInt(conf.Port, 10, 64)
	return dtmcli.DBConf{
		Driver:   "mysql",
		Host:     conf.Path,
		Port:     port,
		User:     conf.Username,
		Password: conf.Password,
		Db:       conf.Dbname,
	}
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
type DtmBarrierBusiFunc func(c *gin.Context, db *gorm.DB) error

func Barrier(c *gin.Context, fn DtmBarrierBusiFunc) error {
	barrier := MustBarrierFromGin(c)
	return barrier.CallWithDB(DbGet(), func(tx *sql.Tx) error {
		gormDB, err := gorm.Open(mysql.New(mysql.Config{
			Conn: tx,
		}), &gorm.Config{})
		if err != nil {
			return dtmcli.ErrFailure
		}
		return fn(c, gormDB)
	})
}

func XaLocalTransaction(c *gin.Context, fn DtmBarrierBusiFunc) error {
	return dtmcli.XaLocalTransaction(c.Request.URL.Query(), DbConf(), func(db *sql.DB, xa *dtmcli.Xa) error {
		gormDB, err := gorm.Open(mysql.New(mysql.Config{
			Conn: db,
		}), &gorm.Config{})
		if err != nil {
			return dtmcli.ErrFailure
		}
		return fn(c, gormDB)
	})
}
