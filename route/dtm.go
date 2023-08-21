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

func (s *UserRouter) InitDtmRouter(route *gin.RouterGroup) {
	r := route.Group("/dtm/saga")
	{
		r.GET("", api.DtmSage)
		r.POST("/out", dtmutil.WrapHandler(func(c *gin.Context) interface{} {
			barrier := utils.MustBarrierFromGin(c)
			return barrier.CallWithDB(utils.DbGet(), func(tx *sql.Tx) error {
				return api.SageOut(c, tx)
			})
		}))
		r.POST("/outCompensate", dtmutil.WrapHandler(func(c *gin.Context) interface{} {
			barrier := utils.MustBarrierFromGin(c)
			return barrier.CallWithDB(utils.DbGet(), func(tx *sql.Tx) error {
				return api.SageOutCompensate(c, tx)
			})
		}))
		r.POST("/in", dtmutil.WrapHandler(func(c *gin.Context) interface{} {
			barrier := utils.MustBarrierFromGin(c)
			return barrier.CallWithDB(utils.DbGet(), func(tx *sql.Tx) error {
				return api.SageIn(c, tx)
			})
		}))
		r.POST("/inCompensate", dtmutil.WrapHandler(func(c *gin.Context) interface{} {
			barrier := utils.MustBarrierFromGin(c)
			return barrier.CallWithDB(utils.DbGet(), func(tx *sql.Tx) error {
				return api.SageInCompensate(c, tx)
			})
		}))
	}
}
