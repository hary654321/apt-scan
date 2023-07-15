package hostinfo

import (
	"context"
	"ias_tool_v2/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Start(ctx *gin.Context) {
	var (
		err    error
		task   *model.ProbeTask
		ctxMe  context.Context
		cancel context.CancelFunc
	)

	params := &model.PortReqParam{}
	if err = ctx.BindJSON(params); err != nil {
		goto ERR
	}
	if err = params.IsValid(); err != nil {
		goto ERR
	}

	ctxMe, cancel = context.WithCancel(context.Background())
	task = model.NewPortTask(params)
	task.ChangeTaskStatus(model.StatusEnum.Received)

	model.InsertCtx(task.TaskId, task.ServiceType, ctxMe, cancel)

	go task.RecordProgress()

	go task.RecordResult()

	go task.Custom(ctxMe)

	go task.Product(ctxMe, params)

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
