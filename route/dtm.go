package router

import (
	"database/sql"
	"gin-example/api"
	"gin-example/utils"
	"github.com/dtm-labs/client/dtmcli"
	"github.com/dtm-labs/dtm/dtmutil"
	"github.com/gin-gonic/gin"
)

type DtmRouter struct{}

var DtmRouterApp = new(UserRouter)

func (s *UserRouter) InitDtmRouter(router *gin.RouterGroup) {
	r := router.Group("/dtm")
	{
		// SAGA
		sageRouter := r.Group("/saga")
		{
			sageRouter.GET("", api.DtmSage)
			sageRouter.POST("/out", dtmutil.WrapHandler(func(c *gin.Context) interface{} {
				barrier := utils.MustBarrierFromGin(c)
				return barrier.CallWithDB(utils.DbGet(), func(tx *sql.Tx) error {
					return api.SageOut(c, tx)
				})
			}))
			sageRouter.POST("/outCompensate", dtmutil.WrapHandler(func(c *gin.Context) interface{} {
				barrier := utils.MustBarrierFromGin(c)
				return barrier.CallWithDB(utils.DbGet(), func(tx *sql.Tx) error {
					return api.SageOutCompensate(c, tx)
				})
			}))
			sageRouter.POST("/in", dtmutil.WrapHandler(func(c *gin.Context) interface{} {
				barrier := utils.MustBarrierFromGin(c)
				return barrier.CallWithDB(utils.DbGet(), func(tx *sql.Tx) error {
					return api.SageIn(c, tx)
				})
			}))
			sageRouter.POST("/inCompensate", dtmutil.WrapHandler(func(c *gin.Context) interface{} {
				barrier := utils.MustBarrierFromGin(c)
				return barrier.CallWithDB(utils.DbGet(), func(tx *sql.Tx) error {
					return api.SageInCompensate(c, tx)
				})
			}))
		}

		// MSG
		msgRouter := r.Group("/msg")
		{
			msgRouter.GET("", api.DtmMsg)
			msgRouter.POST("/in", dtmutil.WrapHandler(func(c *gin.Context) interface{} {
				barrier := utils.MustBarrierFromGin(c)
				return barrier.CallWithDB(utils.DbGet(), func(tx *sql.Tx) error {
					return api.MsgIn(c, tx)
				})
			}))
			msgRouter.POST("/in2", dtmutil.WrapHandler(func(c *gin.Context) interface{} {
				barrier := utils.MustBarrierFromGin(c)
				return barrier.CallWithDB(utils.DbGet(), func(tx *sql.Tx) error {
					return api.MsgIn2(c, tx)
				})
			}))
			msgRouter.GET("/QueryPreparedB", dtmutil.WrapHandler(func(c *gin.Context) interface{} {
				return utils.MustBarrierFromGin(c).QueryPrepared(utils.DbGet())
			}))
		}

		// TCC
		tccRouter := r.Group("/tcc")
		{
			tccRouter.GET("", api.DtmTcc)
			tccRouter.POST("/outTry", dtmutil.WrapHandler(func(c *gin.Context) interface{} {
				barrier := utils.MustBarrierFromGin(c)
				return barrier.CallWithDB(utils.DbGet(), func(tx *sql.Tx) error {
					return api.TccOutTry(c, tx)
				})
			}))
			tccRouter.POST("/outConfirm", dtmutil.WrapHandler(func(c *gin.Context) interface{} {
				barrier := utils.MustBarrierFromGin(c)
				return barrier.CallWithDB(utils.DbGet(), func(tx *sql.Tx) error {
					return api.TccOutConfirm(c, tx)
				})
			}))
			tccRouter.POST("/outCancel", dtmutil.WrapHandler(func(c *gin.Context) interface{} {
				barrier := utils.MustBarrierFromGin(c)
				return barrier.CallWithDB(utils.DbGet(), func(tx *sql.Tx) error {
					return api.TccOutCancel(c, tx)
				})
			}))
			tccRouter.POST("/inTry", dtmutil.WrapHandler(func(c *gin.Context) interface{} {
				barrier := utils.MustBarrierFromGin(c)
				return barrier.CallWithDB(utils.DbGet(), func(tx *sql.Tx) error {
					return api.TccInTry(c, tx)
				})
			}))
			tccRouter.POST("/inConfirm", dtmutil.WrapHandler(func(c *gin.Context) interface{} {
				barrier := utils.MustBarrierFromGin(c)
				return barrier.CallWithDB(utils.DbGet(), func(tx *sql.Tx) error {
					return api.TccInConfirm(c, tx)
				})
			}))
			tccRouter.POST("/inCancel", dtmutil.WrapHandler(func(c *gin.Context) interface{} {
				barrier := utils.MustBarrierFromGin(c)
				return barrier.CallWithDB(utils.DbGet(), func(tx *sql.Tx) error {
					return api.TccInCancel(c, tx)
				})
			}))
		}

		// XA
		xaRouter := r.Group("/xa")
		{
			xaRouter.GET("", api.DtmXa)
			xaRouter.POST("/out", dtmutil.WrapHandler(func(c *gin.Context) interface{} {
				return dtmcli.XaLocalTransaction(c.Request.URL.Query(), utils.DtmConf, func(db *sql.DB, xa *dtmcli.Xa) error {
					return api.XaOut(c, db)
				})
			}))
			xaRouter.POST("/in", dtmutil.WrapHandler(func(c *gin.Context) interface{} {
				return dtmcli.XaLocalTransaction(c.Request.URL.Query(), utils.DtmConf, func(db *sql.DB, xa *dtmcli.Xa) error {
					return api.XaIn(c, db)
				})
			}))
		}
	}
}
