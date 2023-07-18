package router

import (
	"ias_tool_v2/api/port"
	"ias_tool_v2/middlewares"

	"github.com/gin-gonic/gin"
)

func InitPortRouter(Router *gin.RouterGroup) {
	p := Router.Group("port").Use(middlewares.CostTime())
	{
		p.POST("/start", port.Start)
		p.GET("/progress", port.Progress)
		p.GET("/res", port.Res)
	}
}
