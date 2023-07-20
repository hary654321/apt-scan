package port

import (
	"ias_tool_v2/config"
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

	portScanner.Total = len(params.ScanAddrs)
	for _, addr := range params.ScanAddrs {
		netloc, port := utils.SplitWithNetlocPort(addr)
		slog.Println(slog.DEBUG, "ip"+netloc, "port", port)
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

func Progress(ctx *gin.Context) {

	params := &model.GetTaskID{}

	if err := ctx.BindJSON(params); err != nil {
		slog.Println(slog.DEBUG, err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}
	taskid := params.TaskId
	portScanner := model.GetPortClient(taskid)

	total := 0
	if portScanner == nil && !utils.PathExists(config.CoreConf.ResPath+taskid+".json") {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "错误的任务ID",
		})
		return
	}

	var ok int
	if portScanner == nil {
		ok = portScanner.Total
	} else {
		ok = int(portScanner.Ok)
		total = portScanner.Total
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":     200,
		"msg":      "success",
		"all_addr": total,
		"ok_addr":  ok,
		"type":     "port",
	})
	return

}

func Res(ctx *gin.Context) {

	params := &model.GetTaskID{}

	if err := ctx.BindJSON(params); err != nil {
		slog.Println(slog.DEBUG, err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}
	taskid := params.TaskId

	path := config.CoreConf.ResPath + taskid + ".json"
	if !utils.PathExists(path) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "错误的任务ID",
		})
		return
	}

	data, _ := utils.ReadLineData(path)
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"res":  data,
	})
	return
}
