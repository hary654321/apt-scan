package model

import (
	"fmt"
	"ias_tool_v2/core/scanner"
	"ias_tool_v2/core/slog"
	"ias_tool_v2/core/udp"
	"ias_tool_v2/core/utils"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/lcvvvv/appfinger"
	"github.com/lcvvvv/gonmap"
)

type Config struct {
	DeepInspection bool
	Timeout        time.Duration
	Threads        int
	Interval       time.Duration
}

func DefaultConfig() *Config {
	return &Config{
		DeepInspection: false,
		Timeout:        time.Second * 3,
		Threads:        800,
		Interval:       time.Millisecond * 300,
	}
}

// NewPortTask 根据用户请求的参数生成Port任务
// tlsport 如果 addrs 元素port部分存在在tlsport下，则执行https逻辑

var EngineArr map[string]*scanner.PortClient //引擎数组

func init() {
	EngineArr = make(map[string]*scanner.PortClient)

}

func GetRunTasks() string {
	ids := ""
	for k, _ := range EngineArr {
		ids += k
	}

	return ids
}

func NewPortTask(p *ProbeReqParam) *scanner.PortClient {
	PortConfig := scanner.DefaultConfig()
	PortConfig.Threads = p.Threads
	PortConfig.Timeout = time.Duration(p.Timeout) * time.Second // getTimeout(len(app.Setting.Port))
	EngineArr[p.TaskId] = scanner.NewPortScanner(PortConfig, p.TaskId)
	client := EngineArr[p.TaskId]
	runTaskID := p.TaskId
	client.HandlerClosed = func(addr net.IP, port int) {

	}
	client.HandlerOpen = func(addr net.IP, port int) {
		outputOpenResponse(runTaskID, addr, port)
	}
	client.HandlerNotMatched = func(addr net.IP, port int, response string) {
		outputUnknownResponse(runTaskID, addr, port, response)
	}
	client.HandlerMatched = func(addr net.IP, port int, response *gonmap.Response) {
		//slog.Println(slog.DEBUG, "HandlerMatched：", response.FingerPrint.Service, addr.String(), port)
		URLRaw := fmt.Sprintf("%s://%s:%d", response.FingerPrint.Service, addr.String(), port)
		URL, _ := url.Parse(URLRaw)
		if appfinger.SupportCheck(URL.Scheme) == true {
			//在这里处理http请求
			//return
		}

		//继续探针扫描
		outputNmapFinger(runTaskID, URL, response)

	}

	client.HandlerError = func(addr net.IP, port int, err error) {
		slog.Println(slog.DEBUG, "PortScanner Error: ", fmt.Sprintf("%s:%d", addr.String(), port), err)
	}

	client.OutputUdpResponse = func(addr net.IP, port int, res *udp.Result) {
		//输出结果
		protocol := gonmap.GuessProtocol(port)
		target := fmt.Sprintf("%s://%s:%d", protocol, addr.String(), port)
		URL, _ := url.Parse(target)

		m := map[string]interface{}{
			"IP":        URL.Hostname(),
			"Port":      strconv.Itoa(port),
			"Keyword":   res.Service.Name,
			"ProbeName": "UDP",
			"UdpInfo":   res,
		}

		m["runTaskID"] = runTaskID
		utils.WriteJsonAny(runTaskID+".json", m)
	}

	return EngineArr[p.TaskId]
}

func outputOpenResponse(runTaskID string, addr net.IP, port int) {
	//输出结果
	protocol := gonmap.GuessProtocol(port)
	target := fmt.Sprintf("%s://%s:%d", protocol, addr.String(), port)
	URL, _ := url.Parse(target)

	m := map[string]string{
		"IP":      URL.Hostname(),
		"Port":    strconv.Itoa(port),
		"Keyword": "response is empty",
	}
	m["runTaskID"] = runTaskID
	utils.WriteJsonString(runTaskID+".json", m)
}

func outputUnknownResponse(runTaskID string, addr net.IP, port int, response string) {

	//输出结果
	target := fmt.Sprintf("unknown://%s:%d", addr.String(), port)
	URL, _ := url.Parse(target)

	m := map[string]string{
		"Response": response,
		"IP":       URL.Hostname(),
		"Port":     strconv.Itoa(port),
		"Keyword":  "无法识别该协议",
	}
	m["runTaskID"] = runTaskID
	utils.WriteJsonString(runTaskID+".json", m)
}

func outputNmapFinger(runTaskID string, URL *url.URL, resp *gonmap.Response) {

	finger := resp.FingerPrint
	m := utils.ToMap(finger)
	m["Response"] = resp.Raw
	m["IP"] = URL.Hostname()
	m["Port"] = URL.Port()
	m["runTaskID"] = runTaskID

	utils.WriteJsonString(runTaskID+".json", m)
}

func WatchDog(p *scanner.PortClient) {
	time.Sleep(time.Second * 1)
	for {
		slog.Println(slog.WARN, "run --", p.RunningThreads())
		if p.RunningThreads() == 0 {
			delete(EngineArr, p.TaskId)
			break
		}
		time.Sleep(time.Second * 1)
	}
}

func GetPortClient(taskId string) *scanner.PortClient {
	return EngineArr[taskId]
}
