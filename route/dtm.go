package router

import (
	"database/sql"
	"gin-example/api"
	"gin-example/utils"
	"github.com/dtm-labs/dtm/dtmutil"
	"github.com/gin-gonic/gin"
)

type DtmRouter struct{}

var DtmRouterApp = new(UserRouter)

func (s *UserRouter) InitDtmRouter(router *gin.RouterGroup) {
	r := router.Group("/dtm")
	{
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
	}
}
