package passwd_crack

import (
	"context"
	"github.com/gin-gonic/gin"
	"ias_tool_v2/api"
	"ias_tool_v2/model"
	"net/http"
)

type PasswdCrackHandler struct {
	api.Handler
}

func GetPasswdCrackHandler() *PasswdCrackHandler {
	return &PasswdCrackHandler{}
}

func (p *PasswdCrackHandler) Start(ctx *gin.Context) {
	var (
		err    error
		task   *model.PasswdCrackTask
		ctxMe  context.Context
		cancel context.CancelFunc
	)

	params := &model.PasswdCrackParams{
		ServiceType: model.GetServiceType(ctx.Request.URL.String()),
	}

	if err = ctx.BindJSON(params); err != nil {
		goto ERR
	}
	if err = params.IsValid(); err != nil {
		goto ERR
	}

	ctxMe, cancel = context.WithCancel(context.Background())
	task = model.NewPasswdCrackTask(params)
	task.ChangeTaskStatus(model.StatusEnum.Received)

	model.InsertCtx(task.TaskId, task.ServiceType, ctxMe, cancel)

	go task.RecordProgress()

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
