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
			sageRouter.POST("/out", utils.DtmBarrierHandler(func(c *gin.Context) interface{} {
				return utils.Barrier(c, api.SageOut)
			}))
			sageRouter.POST("/outCompensate", utils.DtmBarrierHandler(func(c *gin.Context) interface{} {
				return utils.Barrier(c, api.SageOutCompensate)
			}))
			sageRouter.POST("/in", utils.DtmBarrierHandler(func(c *gin.Context) interface{} {
				return utils.Barrier(c, api.SageIn)
			}))
			sageRouter.POST("/inCompensate", utils.DtmBarrierHandler(func(c *gin.Context) interface{} {
				return utils.Barrier(c, api.SageInCompensate)
			}))
		}

		// MSG
		msgRouter := r.Group("/msg")
		{
			msgRouter.GET("", api.DtmMsg)
			msgRouter.POST("/in", utils.DtmBarrierHandler(func(c *gin.Context) interface{} {
				return utils.Barrier(c, api.MsgIn)
			}))
			msgRouter.POST("/in2", utils.DtmBarrierHandler(func(c *gin.Context) interface{} {
				return utils.Barrier(c, api.MsgIn2)
			}))
			msgRouter.GET("/QueryPreparedB", dtmutil.WrapHandler(func(c *gin.Context) interface{} {
				return utils.MustBarrierFromGin(c).QueryPrepared(utils.DbGet())
			}))
		}

		// TCC
		tccRouter := r.Group("/tcc")
		{
			tccRouter.GET("", api.DtmTcc)
			tccRouter.POST("/outTry", utils.DtmBarrierHandler(func(c *gin.Context) interface{} {
				return utils.Barrier(c, api.TccOutTry)
			}))
			tccRouter.POST("/outConfirm", utils.DtmBarrierHandler(func(c *gin.Context) interface{} {
				return utils.Barrier(c, api.TccOutConfirm)
			}))
			tccRouter.POST("/outCancel", utils.DtmBarrierHandler(func(c *gin.Context) interface{} {
				return utils.Barrier(c, api.TccOutCancel)
			}))
			tccRouter.POST("/inTry", utils.DtmBarrierHandler(func(c *gin.Context) interface{} {
				return utils.Barrier(c, api.TccInTry)
			}))
			tccRouter.POST("/inConfirm", utils.DtmBarrierHandler(func(c *gin.Context) interface{} {
				return utils.Barrier(c, api.TccInConfirm)
			}))
			tccRouter.POST("/inCancel", utils.DtmBarrierHandler(func(c *gin.Context) interface{} {
				return utils.Barrier(c, api.TccInCancel)
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
