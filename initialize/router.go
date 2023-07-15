package initialize

import (
	"ias_tool_v2/router"

	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	Router := gin.Default()

	ApiGroup := Router.Group("/v1")
	router.InitHealthRouter(ApiGroup)
	router.InitProbeScanRouter(ApiGroup)
	router.InitPortRouter(ApiGroup)
	return Router
}
