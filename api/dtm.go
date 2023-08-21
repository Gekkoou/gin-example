package api

import (
	"context"
	"fmt"
	"gin-example/global"
	"github.com/dtm-labs/client/dtmcli"
	"github.com/gin-gonic/gin"
	"github.com/lithammer/shortuuid"
	"time"
)

func SageOut(c *gin.Context, db dtmcli.DB) error {
	var req DtmReq
	if err := c.BindJSON(&req); err != nil {
		return dtmcli.ErrFailure
	}

	if err := global.Redis.DecrBy(context.Background(), "dtm-saga", int64(req.Amount)).Err(); err != nil {
		return dtmcli.ErrFailure
	}

	return nil
}

func SageOutCompensate(c *gin.Context, db dtmcli.DB) error {
	var req DtmReq
	if err := c.BindJSON(&req); err != nil {
		return dtmcli.ErrFailure
	}

	if err := global.Redis.IncrBy(context.Background(), "dtm-saga", int64(req.Amount)).Err(); err != nil {
		return dtmcli.ErrFailure
	}

	return nil
}

func SageIn(c *gin.Context, db dtmcli.DB) error {
	var req DtmReq
	if err := c.BindJSON(&req); err != nil {
		return dtmcli.ErrFailure
	}

	if err := global.Redis.IncrBy(context.Background(), "dtm-saga", int64(req.Amount)).Err(); err != nil {
		return dtmcli.ErrFailure
	}

	return nil
}

func SageInCompensate(c *gin.Context, db dtmcli.DB) error {
	var req DtmReq
	if err := c.BindJSON(&req); err != nil {
		return dtmcli.ErrFailure
	}

	if err := global.Redis.DecrBy(context.Background(), "dtm-saga", int64(req.Amount)).Err(); err != nil {
		return dtmcli.ErrFailure
	}

	return nil
}

func MsgOut(db dtmcli.DB, req DtmReq) error {
	if err := global.Redis.DecrBy(context.Background(), "dtm-msg", int64(req.Amount)).Err(); err != nil {
		return dtmcli.ErrFailure
	}
	return nil
}

func MsgIn(c *gin.Context, db dtmcli.DB) error {
	var req DtmReq
	if err := c.BindJSON(&req); err != nil {
		return dtmcli.ErrFailure
	}

	if err := global.Redis.IncrBy(context.Background(), "dtm-msg", int64(req.Amount)).Err(); err != nil {
		return dtmcli.ErrFailure
	}

	return nil
}

func MsgIn2(c *gin.Context, db dtmcli.DB) error {
	var req DtmReq
	if err := c.BindJSON(&req); err != nil {
		return dtmcli.ErrFailure
	}

	if err := global.Redis.IncrBy(context.Background(), "dtm-msg", int64(req.Amount)).Err(); err != nil {
		return dtmcli.ErrFailure
	}

	return nil
}

type DtmReq struct {
	Amount int `json:"amount"`
}

var (
	dtmServer = "http://127.0.0.1:36789/api/dtmsvr"
	outServer = "http://127.0.0.1:8888/dtm"
	inServer  = "http://127.0.0.1:8888/dtm"
)

func DtmSage(c *gin.Context) {

	// 先安装运行dtm, 参考文档: https://dtm.pub/guide/install.html

	gid := shortuuid.New()
	req := &DtmReq{Amount: 30}

	global.Redis.Set(context.Background(), "dtm-saga", 100, 5*time.Minute).Err()

	saga := dtmcli.NewSaga(dtmServer, gid).
		Add(outServer+"/saga/out", outServer+"/saga/outCompensate", req).
		Add(inServer+"/saga/in", inServer+"/saga/inCompensate", req)

	err := saga.Submit()
	if err != nil {
		fmt.Println(err)
	}
	c.JSON(200, saga)
}

func DtmMsg(c *gin.Context) {
	gid := shortuuid.New()
	req := &DtmReq{Amount: 30}

	global.Redis.Set(context.Background(), "dtm-msg", 100, 5*time.Minute).Err()

	// 二阶段消息
	// msg := dtmcli.NewMsg(dtmServer, gid).
	// 	Add(inServer+"/msg/in", req)
	// err := msg.DoAndSubmitDB(inServer+"/msg/QueryPreparedB", utils.DbGet(), func(tx *sql.Tx) error {
	// 	return MsgOut(tx, req.Amount)
	// })

	// 普通消息
	msg := dtmcli.NewMsg(dtmServer, gid).
		Add(inServer+"/msg/in", req).
		Add(inServer+"/msg/in2", req)
	err := msg.Submit()

	if err != nil {
		fmt.Println(err)
	}
	c.JSON(200, msg)
}
