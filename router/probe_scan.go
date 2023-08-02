package router

import (
	"ias_tool_v2/api"
	"ias_tool_v2/api/ssl_probe"
	"ias_tool_v2/middlewares"

	"github.com/gin-gonic/gin"
)

func InitProbeScanRouter(Router *gin.RouterGroup) {
	PasswdCrackRouter := Router.Group("probe").Use(middlewares.CostTime()).Use(middlewares.BasicAuth())
	{
		PasswdCrackRouter.POST("/stop", middlewares.BaseParamsCheck(), api.InitHandler().Stop)
		PasswdCrackRouter.POST("/progress", middlewares.BaseParamsCheck(), ssl_probe.GetPasswdCrackHandler().Progress)
		PasswdCrackRouter.POST("/start", middlewares.BaseParamsCheck(), ssl_probe.GetPasswdCrackHandler().Start)
		PasswdCrackRouter.POST("/result", middlewares.BaseParamsCheck(), api.InitHandler().Result)
	}
}
