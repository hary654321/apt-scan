package initialize

import (
	"github.com/gin-gonic/gin"
	"ias_tool_v2/router"
)

func Routers() *gin.Engine {
	Router := gin.Default()

	ApiGroup := Router.Group("/v1")
	router.InitHealthRouter(ApiGroup)
	router.InitProbeScanRouter(ApiGroup)
	return Router
}
