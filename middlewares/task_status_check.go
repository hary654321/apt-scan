package middlewares

import (
	"errors"
	"ias_tool_v2/api"
	"ias_tool_v2/logger"
	"ias_tool_v2/model"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

// BaseParamsCheck 必要参数检查
func BaseParamsCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if _, ok := ctx.Get("task_id"); !ok {
			return
		}
		if _, ok := ctx.Get("service_type"); !ok {
			return
		}
		ctx.Next()
	}
}

// MustNotExecuted 必须没有执行过
func MustNotExecuted() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var err error
		var picklePath string
		serviceType := model.GetServiceType(ctx.Request.URL.String())
		task := &api.Params{}
		if err = ctx.BindJSON(task); err != nil {
			goto ERR
		}
		picklePath = filepath.Join(model.PicklePathFolder(serviceType), task.TaskId+".job")
		//通过taskId和serviceType获取.job文件是否存在
		if _, err = os.OpenFile(picklePath, os.O_WRONLY, 0666); err == nil {
			err = errors.New("task has executed")
			goto ERR
		}
		return
	ERR:
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err.Error(),
		})
	}
}

// MustExecuted 必须执行过
func MustExecuted() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		taskId, _ := ctx.Get("task_id")
		serviceType, _ := ctx.Get("service_type")
		//通过taskId和serviceType获取.job文件是否存在
		picklePath := model.PicklePathFolder(serviceType.(string)) + taskId.(string) + ".job"
		if _, err := os.OpenFile(picklePath, os.O_WRONLY, 0666); err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"code": http.StatusBadRequest,
				"msg":  "task has executed",
			})
		}

	}
}

func CostTime() gin.HandlerFunc {
	return func(c *gin.Context) {
		//请求前获取当前时间
		nowTime := time.Now()

		//请求处理
		c.Next()

		//处理后获取消耗时间
		costTime := time.Since(nowTime)
		//todo 当前url获取方法可能有用
		url := c.Request.URL.String()
		logger.Infof("the request URL %s cost %v\n", url, costTime)
	}
}
