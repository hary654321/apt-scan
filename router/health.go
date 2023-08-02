package router

import (
	"ias_tool_v2/api/hostinfo"
	"ias_tool_v2/config"
	"ias_tool_v2/core/utils"
	"ias_tool_v2/middlewares"
	"ias_tool_v2/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

func GetMemPercent() float64 {
	memInfo, _ := mem.VirtualMemory()
	return memInfo.UsedPercent
}

func GetCpuPercent() float64 {
	percent, _ := cpu.Percent(time.Second, false)
	return percent[0]
}

func GetDiskPercent() float64 {
	parts, _ := disk.Partitions(true)
	diskInfo, _ := disk.Usage(parts[0].Mountpoint)
	return diskInfo.UsedPercent
}

func InitHealthRouter(Router *gin.RouterGroup) {
	c := Router.Group("check").Use(middlewares.CostTime()).Use(middlewares.BasicAuth())
	{
		c.GET("/health", func(context *gin.Context) {
			memPercent := int(GetMemPercent())
			if memPercent >= 80 {
				context.JSON(http.StatusBadRequest, gin.H{
					"code": 400,
					"msg":  "cur node mem not enough",
				})
				return
			}
			context.JSON(http.StatusOK, gin.H{
				"code": http.StatusOK,
				"msg":  "health",
			})
		})

		c.GET("/node_checker", func(context *gin.Context) {
			memPercent := int(GetMemPercent())
			cpuPercent := int(GetCpuPercent())
			diskPercent := int(GetDiskPercent())
			context.JSON(http.StatusOK, gin.H{
				"code":         http.StatusOK,
				"msg":          "health",
				"cpu_percent":  cpuPercent,
				"mem_percent":  memPercent,
				"dick_percent": diskPercent,
			})
		})

		c.GET("/hostinfo", hostinfo.InfoGet)

		c.GET("/heartbeat", func(c *gin.Context) {
			data := make(map[string]interface{})
			data["time"] = utils.GetTime()
			data["version"] = config.CoreConf.Version
			data["runningTasks"] = model.GetRunTasks()
			c.JSON(http.StatusOK, gin.H{
				"code": 200,
				"msg":  "",
				"data": data,
			})
		})

	}
}
