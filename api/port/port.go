package port

import (
	"ias_tool_v2/core/scanner"
	"ias_tool_v2/core/slog"
	"ias_tool_v2/core/utils"
	"ias_tool_v2/model"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Start(ctx *gin.Context) {
	var (
		err         error
		portScanner *scanner.PortClient
	)

	params := &model.PortReqParam{}

	if err = ctx.BindJSON(params); err != nil {
		goto ERR
	}

	slog.Println(slog.DEBUG, params.ScanAddrs)

	if err = params.IsValid(); err != nil {
		goto ERR
	}

	portScanner = model.NewPortTask(params)

	go portScanner.Start()
	go model.WatchDog(portScanner)

	for _, addr := range params.ScanAddrs {
		netloc, port := utils.SplitWithNetlocPort(addr)
		go portScanner.Push(net.ParseIP(netloc), port)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
	})
	return
ERR:
	ctx.JSON(http.StatusBadRequest, gin.H{
		"code": 400,
		"msg":  err.Error(),
	})
}
