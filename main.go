package main

import (
	"fmt"
	"ias_tool_v2/config"
	"ias_tool_v2/initialize"
	"ias_tool_v2/logger"
	"ias_tool_v2/middlewares"

	"github.com/gin-gonic/gin"
)

func main() {
	var err error
	var Router *gin.Engine

	//加载配置
	config.Init("conf.toml")

	//加载路由
	Router = initialize.Routers()
	if config.CoreConf.HttpsServer {
		Router.Use(middlewares.TlsHandler())
		if err = Router.RunTLS(fmt.Sprintf(":%d", config.CoreConf.ApiPort), "pem/.cert.pem", "pem/.key.pem"); err != nil {
			goto ERR
		}
	} else {
		if err = Router.Run(fmt.Sprintf(":%d", config.CoreConf.ApiPort)); err != nil {
			goto ERR
		}
	}
ERR:
	logger.Warningf(err.Error())
}
