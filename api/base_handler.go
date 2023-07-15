package api

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"ias_tool_v2/model"
	"net/http"
	"strings"
)

type handler interface {
	Start(ctx *gin.Context)
	Stop(ctx *gin.Context)
	Progress(ctx *gin.Context)
	Result(ctx *gin.Context)
}

type Handler struct {
}

type Params struct {
	TaskId string `json:"task_id"`
}

func (h *Handler) Start(ctx *gin.Context) {
}

func GetServiceType(url string) string {
	urls := strings.Split(url, "/")
	return urls[2]
}

func (h *Handler) Result(ctx *gin.Context) {
	var (
		err       error
		res       []interface{}
		probeTask *model.ProbeTask
	)
	params := Params{}

	if err = ctx.BindJSON(&params); err != nil {
		goto ERR
	}

	switch GetServiceType(ctx.Request.URL.String()) {
	case "probe":
		probeTask, err = model.ProbeTaskDecode(params.TaskId, GetServiceType(ctx.Request.URL.String()))
		if err != nil {
			goto ERR
		}
		if res, err = probeTask.ReadResultFile(model.ProbeTask{}); err != nil {
			goto ERR
		}
		probeTask.TruncateFile()
	default:
		goto ERR
	}

	//通过url来获取serviceType
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": res,
	})
	return
ERR:
	ctx.JSON(http.StatusBadRequest, gin.H{
		"code": 400,
		"msg":  err.Error(),
	})
}

func (h *Handler) Progress(ctx *gin.Context) {
	var err error

	params := Params{}
	if err := ctx.BindJSON(&params); err != nil {
		goto ERR
	}
	switch GetServiceType(ctx.Request.URL.String()) {
	case "probe":
		if probeTask, err := model.ProbeTaskDecode(params.TaskId, GetServiceType(ctx.Request.URL.String())); err == nil {
			ok, all := probeTask.ReadProgressFile()
			ctx.JSON(http.StatusOK, gin.H{
				"code":     200,
				"msg":      "success",
				"ok_addr":  ok,
				"all_addr": all,
			})
			return
		} else {
			goto ERR
		}
	default:
		err = errors.New("StatusBadRequest")
		goto ERR
	}
ERR:
	err = errors.New("StatusBadRequest")
	ctx.JSON(http.StatusBadRequest, gin.H{
		"code": 400,
		"msg":  err.Error(),
	})
}

func (h *Handler) Stop(ctx *gin.Context) {
	var (
		err    error
		cancel context.CancelFunc
		ok     bool
	)
	params := Params{}
	if err = ctx.BindJSON(&params); err != nil {
		goto ERR
	}

	if _, cancel, ok = model.GetCtx(params.TaskId, GetServiceType(ctx.Request.URL.String())); ok {
		cancel()
		ctx.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "success",
		})
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "task has complete or not exists",
		})
	}
	return
ERR:
	ctx.JSON(http.StatusBadRequest, gin.H{
		"code": 400,
		"msg":  err.Error(),
	})
}

func InitHandler() *Handler {
	return &Handler{}
}
