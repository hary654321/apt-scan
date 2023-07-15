package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"ias_tool_v2/initialize"
	"ias_tool_v2/logger"
	"ias_tool_v2/middlewares"
)

func main() {
	var err error
	var Router *gin.Engine
	//加载配置
	if err := initialize.InitConfig(); err != nil {
		goto ERR
	}
	//加载路由
	Router = initialize.Routers()
	if initialize.GConfig.HttpsServer {
		Router.Use(middlewares.TlsHandler())
		if err = Router.RunTLS(fmt.Sprintf(":%d", initialize.GConfig.ApiPort), ".pem/.cert.pem", ".pem/.key.pem"); err != nil {
			goto ERR
		}
	} else {
		if err = Router.Run(fmt.Sprintf(":%d", initialize.GConfig.ApiPort)); err != nil {
			goto ERR
		}
	}
ERR:
	logger.Warningf(err.Error())
}
