package router

import (
	"ias_tool_v2/api/port"
	"ias_tool_v2/middlewares"

	"github.com/gin-gonic/gin"
)

func InitPortRouter(Router *gin.RouterGroup) {
	p := Router.Group("port").Use(middlewares.CostTime()).Use(middlewares.BasicAuth())
	{
		p.POST("/start", port.Start)
		p.POST("/progress", port.Progress)
		p.POST("/result", port.Res)
		p.POST("/stop", port.Stop)
	}
}
