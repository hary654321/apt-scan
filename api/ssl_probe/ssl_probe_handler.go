package ssl_probe

import (
	"ias_tool_v2/api"
	"ias_tool_v2/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProbeHandler struct {
	api.Handler
}

func GetPasswdCrackHandler() *ProbeHandler {
	return &ProbeHandler{}
}

func (p *ProbeHandler) Start(ctx *gin.Context) {
	var (
		err error
	)

	params := &model.ProbeReqParam{
		ServiceType: model.GetServiceType(ctx.Request.URL.String()),
	}

	if err = ctx.BindJSON(params); err != nil {
		goto ERR
	}
	// slog.Println(slog.DEBUG, params)
	if err = params.IsValid(); err != nil {
		goto ERR
	}

	model.ProbeScan(params)

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
