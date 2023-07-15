package hostinfo

import (
	"ias_tool_v2/core/hostinfo"
	"ias_tool_v2/logger"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/net"
)

var (
	nowtime    string
	hostip     string
	hostinfos  *host.InfoStat
	parts      hostinfo.Parts
	cpuinfos   hostinfo.CpuInfo
	mempercent float64
	meminfos   hostinfo.MemInfo
	netinfos   []net.IOCountersStat
	netspeed   []hostinfo.SpeedInfo

	//nodesstate  []*db.NodeState
	//tasksstate  []db.TaskState

)

// InfoGet 本机信息
func InfoGet(c *gin.Context) {
	var err error
	nowtime = time.Now().Format("2006-01-02 15:04:05")

	cpuinfos, err = hostinfo.GetCpuPercent()
	if err != nil {
		logger.Warningf(err.Error())
	}
	meminfos = hostinfo.GetMemInfo()
	netinfos, err = hostinfo.GetNetInfo()
	if err != nil {
		logger.Warningf(err.Error())
	}
	netspeed = hostinfo.GetNetSpeed()

	hostip = hostinfo.GetLocalIP()
	hostinfos, err = hostinfo.GetHostInfo()
	if err != nil {
		logger.Warningf(err.Error())
	}

	parts, err = hostinfo.GetDiskInfo()
	if err != nil {
		logger.Warningf(err.Error())
	}
	c.JSON(http.StatusOK,
		gin.H{"hostip": hostip, "hostinfos": hostinfos, "parts": parts,
			"cpuinfos": cpuinfos, "mempercent": mempercent, "meminfos": meminfos,
			"netinfos": netinfos, "netspeed": netspeed, "nowtime": nowtime,
		})
}
