package api

import (
	"context"
	"errors"
	"fmt"
	"gin-example/global"
	"github.com/dtm-labs/client/dtmcli"
	"github.com/dtm-labs/dtm/dtmutil"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/lithammer/shortuuid"
	"strconv"
	"time"
)

type DtmReq struct {
	Amount int    `json:"amount"`
	User   string `json:"user"`
}

type DtmTccReq struct {
	Amount       int    `json:"amount"`
	FrozenAmount int    `json:"frozen_amount"`
	User         string `json:"user"`
}

var (
	dtmServer = "http://127.0.0.1:36789/api/dtmsvr"
	outServer = "http://127.0.0.1:8888/dtm"
	inServer  = "http://127.0.0.1:8888/dtm"

	reqOut = &DtmReq{Amount: 30, User: "a"}
	reqIn  = &DtmReq{Amount: 30, User: "b"}

	reqTccOut = &DtmTccReq{Amount: 30, User: "a"}
	reqTccIn  = &DtmTccReq{Amount: 30, User: "b"}
)

func DtmSage(c *gin.Context) {

	// 先安装运行dtm, 参考文档: https://dtm.pub/guide/install.html

	gid := shortuuid.New()

	global.Redis.Set(context.Background(), "dtm-saga-"+reqOut.User, 100, 5*time.Minute).Err()
	global.Redis.Set(context.Background(), "dtm-saga-"+reqIn.User, 100, 5*time.Minute).Err()

	saga := dtmcli.NewSaga(dtmServer, gid).
		Add(outServer+"/saga/out", outServer+"/saga/outCompensate", reqOut).
		Add(inServer+"/saga/in", inServer+"/saga/inCompensate", reqIn)

	err := saga.Submit()
	if err != nil {
		fmt.Println(err)
	}
	c.JSON(200, saga)
}

func DtmMsg(c *gin.Context) {
	gid := shortuuid.New()

	global.Redis.Set(context.Background(), "dtm-msg-"+reqOut.User, 100, 5*time.Minute).Err()
	global.Redis.Set(context.Background(), "dtm-msg-"+reqIn.User, 100, 5*time.Minute).Err()

	// 二阶段消息
	// msg := dtmcli.NewMsg(dtmServer, gid).
	// 	Add(inServer+"/msg/in", reqIn)
	// err := msg.DoAndSubmitDB(inServer+"/msg/QueryPreparedB", utils.DbGet(), func(tx *sql.Tx) error {
	// 	return MsgOut(tx, reqOut)
	// })

	// 普通消息
	msg := dtmcli.NewMsg(dtmServer, gid).
		Add(inServer+"/msg/in", reqIn).
		Add(inServer+"/msg/in2", reqIn)
	err := msg.Submit()

	if err != nil {
		fmt.Println(err)
	}
	c.JSON(200, msg)
}

func DtmTcc(c *gin.Context) {

	// Try 阶段：尝试执行，完成所有业务检查（一致性）, 预留必须业务资源（准隔离性）
	// Confirm 阶段：如果所有分支的Try都成功了，则走到Confirm阶段。Confirm真正执行业务，不作任何业务检查，只使用 Try 阶段预留的业务资源
	// Cancel 阶段：如果所有分支的Try有一个失败了，则走到Cancel阶段。Cancel释放 Try 阶段预留的业务资源。

	gid := shortuuid.New()

	global.Redis.HMSet(context.Background(), "dtm-tcc-"+reqTccOut.User, "amount", 30, "frozen_amount", 0)
	global.Redis.HMSet(context.Background(), "dtm-tcc-"+reqTccIn.User, "amount", 30, "frozen_amount", 0)
	global.Redis.Expire(context.Background(), "dtm-tcc-"+reqTccOut.User, 5*time.Minute)
	global.Redis.Expire(context.Background(), "dtm-tcc-"+reqTccIn.User, 5*time.Minute)

	err := dtmcli.TccGlobalTransaction(dtmServer, gid, func(tcc *dtmcli.Tcc) (*resty.Response, error) {
		resp, err := tcc.CallBranch(reqTccOut, outServer+"/tcc/outTry", outServer+"/tcc/outConfirm", outServer+"/tcc/outCancel")
		if err != nil {
			return resp, err
		}
		return tcc.CallBranch(reqTccIn, inServer+"/tcc/inTry", inServer+"/tcc/inConfirm", inServer+"/tcc/inCancel")
	})

	if err != nil {
		fmt.Println(err)
	}
	c.JSON(200, "dtm tcc ok")
}

func DtmXa(c *gin.Context) {
	gid := shortuuid.New()

	global.Redis.Set(context.Background(), "dtm-xa-"+reqOut.User, 100, 5*time.Minute).Err()
	global.Redis.Set(context.Background(), "dtm-xa-"+reqIn.User, 100, 5*time.Minute).Err()

	err := dtmcli.XaGlobalTransaction(dtmutil.DefaultHTTPServer, gid, func(xa *dtmcli.Xa) (*resty.Response, error) {
		resp, err := xa.CallBranch(reqOut, outServer+"/xa/out")
		if err != nil {
			return resp, err
		}
		return xa.CallBranch(reqIn, inServer+"/xa/in")
	})

	if err != nil {
		fmt.Println(err)
	}
	c.JSON(200, "dtm xa ok")
}

func SageOut(c *gin.Context, db dtmcli.DB) error {
	var req DtmReq
	if err := c.BindJSON(&req); err != nil {
		return dtmcli.ErrFailure
	}

	if err := global.Redis.DecrBy(context.Background(), "dtm-saga-"+req.User, int64(req.Amount)).Err(); err != nil {
		return dtmcli.ErrFailure
	}

	return nil
}

func SageOutCompensate(c *gin.Context, db dtmcli.DB) error {
	var req DtmReq
	if err := c.BindJSON(&req); err != nil {
		return dtmcli.ErrFailure
	}

	if err := global.Redis.IncrBy(context.Background(), "dtm-saga-"+req.User, int64(req.Amount)).Err(); err != nil {
		return dtmcli.ErrFailure
	}

	return nil
}

func SageIn(c *gin.Context, db dtmcli.DB) error {
	var req DtmReq
	if err := c.BindJSON(&req); err != nil {
		return dtmcli.ErrFailure
	}

	if err := global.Redis.IncrBy(context.Background(), "dtm-saga-"+req.User, int64(req.Amount)).Err(); err != nil {
		return dtmcli.ErrFailure
	}

	return nil
}

func SageInCompensate(c *gin.Context, db dtmcli.DB) error {
	var req DtmReq
	if err := c.BindJSON(&req); err != nil {
		return dtmcli.ErrFailure
	}

	if err := global.Redis.DecrBy(context.Background(), "dtm-saga-"+req.User, int64(req.Amount)).Err(); err != nil {
		return dtmcli.ErrFailure
	}

	return nil
}

func MsgOut(db dtmcli.DB, req *DtmReq) error {
	if err := global.Redis.DecrBy(context.Background(), "dtm-msg-"+req.User, int64(req.Amount)).Err(); err != nil {
		return dtmcli.ErrFailure
	}
	return nil
}

func MsgIn(c *gin.Context, db dtmcli.DB) error {
	var req DtmReq
	if err := c.BindJSON(&req); err != nil {
		return dtmcli.ErrFailure
	}

	if err := global.Redis.IncrBy(context.Background(), "dtm-msg-"+req.User, int64(req.Amount)).Err(); err != nil {
		return dtmcli.ErrFailure
	}

	return nil
}

func MsgIn2(c *gin.Context, db dtmcli.DB) error {
	var req DtmReq
	if err := c.BindJSON(&req); err != nil {
		return dtmcli.ErrFailure
	}

	if err := global.Redis.IncrBy(context.Background(), "dtm-msg-"+req.User, int64(req.Amount)).Err(); err != nil {
		return dtmcli.ErrFailure
	}

	return nil
}

func getTccUserAmount(req *DtmTccReq) (amount int, frozenAmount int, err error) {
	res, err := global.Redis.HMGet(context.Background(), "dtm-tcc-"+req.User, "amount", "frozen_amount").Result()
	if err != nil {
		return 0, 0, errors.New("get tcc user amount error")
	}
	amount, err = strconv.Atoi(res[0].(string))
	frozenAmount, err = strconv.Atoi(res[1].(string))
	if err != nil {
		return 0, 0, errors.New("get tcc user amount error")
	}
	return
}

func TccOutTry(c *gin.Context, db dtmcli.DB) error {
	var req DtmTccReq
	if err := c.BindJSON(&req); err != nil {
		return dtmcli.ErrFailure
	}

	_, frozenAmount, err := getTccUserAmount(reqTccOut)
	if err != nil {
		return dtmcli.ErrFailure
	}

	if err = global.Redis.HMSet(context.Background(), "dtm-tcc-"+reqTccOut.User, "frozen_amount", frozenAmount-reqTccOut.Amount).Err(); err != nil {
		return dtmcli.ErrFailure
	}

	return nil
}

func TccOutConfirm(c *gin.Context, db dtmcli.DB) error {
	var req DtmTccReq
	if err := c.BindJSON(&req); err != nil {
		return dtmcli.ErrFailure
	}

	amount, frozenAmount, err := getTccUserAmount(reqTccOut)
	if err != nil {
		return dtmcli.ErrFailure
	}

	if err = global.Redis.HMSet(context.Background(), "dtm-tcc-"+reqTccOut.User, "amount", amount-reqTccOut.Amount, "frozen_amount", frozenAmount+reqTccOut.Amount).Err(); err != nil {
		return dtmcli.ErrFailure
	}

	return nil
}

func TccOutCancel(c *gin.Context, db dtmcli.DB) error {
	var req DtmTccReq
	if err := c.BindJSON(&req); err != nil {
		return dtmcli.ErrFailure
	}

	_, frozenAmount, err := getTccUserAmount(reqTccOut)
	if err != nil {
		return dtmcli.ErrFailure
	}

	if err = global.Redis.HMSet(context.Background(), "dtm-tcc-"+reqTccOut.User, "frozen_amount", frozenAmount+reqTccOut.Amount).Err(); err != nil {
		return dtmcli.ErrFailure
	}

	return nil
}

func TccInTry(c *gin.Context, db dtmcli.DB) error {
	var req DtmTccReq
	if err := c.BindJSON(&req); err != nil {
		return dtmcli.ErrFailure
	}

	_, frozenAmount, err := getTccUserAmount(reqTccIn)
	if err != nil {
		return dtmcli.ErrFailure
	}

	if err = global.Redis.HMSet(context.Background(), "dtm-tcc-"+reqTccIn.User, "frozen_amount", frozenAmount+reqTccIn.Amount).Err(); err != nil {
		return dtmcli.ErrFailure
	}

	return nil
}

func TccInConfirm(c *gin.Context, db dtmcli.DB) error {
	var req DtmTccReq
	if err := c.BindJSON(&req); err != nil {
		return dtmcli.ErrFailure
	}

	amount, frozenAmount, err := getTccUserAmount(reqTccIn)
	if err != nil {
		return dtmcli.ErrFailure
	}

	if err = global.Redis.HMSet(context.Background(), "dtm-tcc-"+reqTccIn.User, "amount", amount+reqTccIn.Amount, "frozen_amount", frozenAmount-reqTccIn.Amount).Err(); err != nil {
		return dtmcli.ErrFailure
	}

	return nil
}

func TccInCancel(c *gin.Context, db dtmcli.DB) error {
	var req DtmTccReq
	if err := c.BindJSON(&req); err != nil {
		return dtmcli.ErrFailure
	}

	_, frozenAmount, err := getTccUserAmount(reqTccIn)
	if err != nil {
		return dtmcli.ErrFailure
	}

	if err = global.Redis.HMSet(context.Background(), "dtm-tcc-"+reqTccIn.User, "frozen_amount", frozenAmount-reqTccIn.Amount).Err(); err != nil {
		return dtmcli.ErrFailure
	}

	return nil
}

func XaOut(c *gin.Context, db dtmcli.DB) error {
	var req DtmReq
	if err := c.BindJSON(&req); err != nil {
		return dtmcli.ErrFailure
	}

	if err := global.Redis.DecrBy(context.Background(), "dtm-xa-"+req.User, int64(req.Amount)).Err(); err != nil {
		return dtmcli.ErrFailure
	}

	return nil
}

func XaIn(c *gin.Context, db dtmcli.DB) error {
	var req DtmReq
	if err := c.BindJSON(&req); err != nil {
		return dtmcli.ErrFailure
	}

	if err := global.Redis.IncrBy(context.Background(), "dtm-xa-"+req.User, int64(req.Amount)).Err(); err != nil {
		return dtmcli.ErrFailure
	}

	return nil
}
