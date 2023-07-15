package udp

import (
	"context"
	"embed"
	"fmt"
	"net"
	"time"
	"zrWorker/core/slog"
	"zrWorker/pkg/utils"
)

var UdpPort = []int{20, 21, 53, 67, 68, 69, 123, 137, 138, 139, 161, 162, 500, 514, 520, 631, 1194, 1701, 1718, 1719, 1720, 1723, 1812, 1813, 1863, 1900, 4500, 2049}

// go test -v -run TestUdpClient port_test.go
func UpdStatus(ip net.IP, port int) string {
	//address := "192.168.56.1:9090"
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	ch := make(chan string)
	go readRes(ctx, ip, port, ch)

	for {
		select {
		case v, ok := <-ch:
			slog.Println(slog.DEBUG, "have", v, ok)
			return v
		case <-ctx.Done():
			//slog.Println(slog.DEBUG, "not have")
			return ""
		}
	}
}

// go test -v -run TestUdpClient port_test.go
func readRes(ctx context.Context, ip net.IP, port int, ch chan string) {
	// 连接服务器
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   ip,
		Port: port,
	})

	if err != nil {
		//slog.Println(slog.DEBUG,"Connect to udp server failed,err:", err)
		return
	}

	// 发送数据
	_, err1 := conn.Write([]byte(fmt.Sprintf("udp testing:%v", "hi")))
	if err1 != nil {
		//slog.Println(slog.DEBUG,"Send data failed,err:", err)
		return
	}
	//接收数据
	result := make([]byte, 1024)

	//slog.Println(slog.DEBUG,"Listening...")
	//read from UDPConn here
	n, remoteAddr, err := conn.ReadFromUDP(result)
	if err != nil {
		//slog.Println(slog.DEBUG,"Read from udp server failed ,err:", err)
		return
	}
	slog.Println(slog.DEBUG, "Recived msg from  ", remoteAddr, string(result[:n]))

	ch <- string(result[:n])
	select {
	case <-ctx.Done():
		close(ch)
		slog.Println(slog.DEBUG, ctx.Err())
		return
	default:
		//slog.Println(slog.DEBUG,n, "default")
	}

}

var cli *VScan
var f embed.FS

func init() {
	cli = GoNmapInit()
}
func UdpInfo(ip string, port int) (*Result, error) {

	nmapbanner, err := cli.GoNmapScan(ip, utils.GetInterfaceToString(port), "udp") //支持udp

	if err != nil {
		//slog.Println(slog.DEBUG, "udpInfo", err)
	}

	return nmapbanner, err
}
