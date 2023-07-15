package scanner

import (
	"embed"
	"fmt"
	"github.com/hary654321/gonmap"
	"log"
	"net"
	"net/http"
	"testing"
	"time"
	"zrWorker/core/slog"
	"zrWorker/lib/appfinger"
	"zrWorker/lib/ping"
	"zrWorker/lib/simplehttp"
	"zrWorker/lib/uri"
	"zrWorker/pkg/utils"
)

func TestNmap(t *testing.T) {
	nmap := gonmap.New()
	nmap.SetTimeout(time.Second * 10)

	//具体进行端口扫描
	status, response := nmap.ScanTimeout("127.0.0.1", 3306, time.Second*10)
	slog.Println(slog.DEBUG, "端口状态：", status, response)
}

func TestGetBannerWithURL(t *testing.T) {

	URL := uri.URLParse("http://www.zorelworld.com")
	req, _ := simplehttp.NewRequest(http.MethodGet, URL.String(), nil)

	cli := simplehttp.NewClient()
	banner, err := appfinger.GetBannerWithURL(URL, req, cli)
	slog.Println(slog.DEBUG, "banner", banner, err)
}

func TestUdpServer(t *testing.T) {
	// 监听UDP服务
	udpConn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 9090,
	})

	if err != nil {
		log.Fatal("Listen failed,", err)
		return
	}

	// 循环读取消息
	for {
		var data [1024]byte
		n, addr, err := udpConn.ReadFromUDP(data[:])
		if err != nil {
			log.Printf("Read from udp server:%s failed,err:%s", addr, err)
			break
		}
		go func() {
			// 返回数据
			fmt.Printf("Addr:%s,data:%v count:%d \n", addr, string(data[:n]), n)
			_, err := udpConn.WriteToUDP([]byte("msg recived."), addr)
			if err != nil {
				fmt.Println("write to udp server failed,err:", err)
			}
		}()
	}

}

// go test -v -run TestUdpClient port_test.go
func TestUdpClient(t *testing.T) {
	// 连接服务器
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 9090,
	})

	if err != nil {
		log.Println("Connect to udp server failed,err:", err)
		return
	}

	println(1)

	// 发送数据
	_, err1 := conn.Write([]byte(fmt.Sprintf("udp testing:%v", "hi")))
	if err1 != nil {
		log.Println("Send data failed,err:", err)
		return
	}

	println(2)
	//接收数据
	result := make([]byte, 1024)

	println(3)
	n, remoteAddr, err := conn.ReadFromUDP(result)
	println(4)
	if err != nil {
		log.Println("Read from udp server failed ,err:", err)
		return
	}
	println(5)
	fmt.Printf("Recived msg from %s, data:%s \n", remoteAddr, string(result[:n]))

}

//go:embed fingerprint.txt
var fingerprintEmbed embed.FS

func TestUrlFinger(t *testing.T) {
	fs, _ := fingerprintEmbed.Open("fingerprint.txt")
	if n, err := appfinger.InitDatabaseFS(fs); err != nil {
		slog.Println(slog.ERROR, "指纹库加载失败，请检查【fingerprint.txt】文件", err)
	} else {
		slog.Printf(slog.INFO, "成功加载HTTP指纹:[%d]条", n)
	}

	client := simplehttp.NewClient()

	value := foo2{uri.URLParse("https://stage.akagishizenen.jp"), nil, nil, client}

	URL := value.URL
	req := value.req
	cli := value.client
	if appfinger.SupportCheck(URL.Scheme) == false {

		return
	}
	var banner *appfinger.Banner
	//var finger *appfinger.FingerPrint
	var err error

	banner, err = appfinger.GetBannerWithURL(URL, req, cli)
	utils.WriteJson("FoundDomain.json", banner)
	if err != nil {

		return
	}

	//slog.Printf(slog.DEBUG, "response == nil")
	//finger = appfinger.Search(URL, banner)
	//
	//m := misc.ToMap(finger)
	//
	//utils.WriteJson("a", m)

}

func TestPortFinger(t *testing.T) {
	//fs, _ := fingerprintEmbed.Open("fingerprint.txt")
	//if n, err := appfinger.InitDatabaseFS(fs); err != nil {
	//	slog.Println(slog.ERROR, "指纹库加载失败，请检查【fingerprint.txt】文件", err)
	//} else {
	//	slog.Printf(slog.INFO, "成功加载HTTP指纹:[%d]条", n)
	//}
	//

	nmap := gonmap.New()
	nmap.SetTimeout(time.Second * 10)
	_, response := nmap.ScanTimeout("172.16.130.138", 9200, time.Second*10)

	utils.WriteJson("tcp", response.FingerPrint)

}

func TestPing(t *testing.T) {
	res := ping.Pinger("1.0.125.62", 1)
	slog.Println(slog.DEBUG, "res", res)
}
