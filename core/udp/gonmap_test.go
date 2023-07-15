package udp

import (
	"testing"
	"zrWorker/pkg/utils"
)

// 文件是   https://svn.nmap.org/nmap/nmap-service-probes

func TestScan(t *testing.T) {

	//slog.Println(slog.DEBUG, "inited")
	//nmapbanner, err := GoNmap.GoNmapScan("172.16.130.138", "9200", "tcp") //端口不存货  会阻塞 1分钟
	nmapbanner, _ := UdpInfo("115.45.105.240", 177) //支持udp

	utils.WriteJson("nmapbanner.json", nmapbanner)
}
