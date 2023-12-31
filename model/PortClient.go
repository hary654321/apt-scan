package model

import (
	"encoding/json"
	"fmt"
	"ias_tool_v2/core/scanner"
	"ias_tool_v2/core/slog"
	"ias_tool_v2/core/udp"
	"ias_tool_v2/core/utils"
	"ias_tool_v2/define"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
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
var TakData map[string]*ProbeReqParam        //任务数据

func init() {
	if EngineArr == nil {
		EngineArr = make(map[string]*scanner.PortClient)
	}
	if TakData == nil {
		TakData = make(map[string]*ProbeReqParam)
	}
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
	TakData[p.TaskId] = p
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

	m = utils.Dealdata(m)

	utils.WriteJsonString(runTaskID+".json", m)
}

func WatchDog(p *scanner.PortClient) {
	time.Sleep(time.Second * 1)
	for {
		if p.Stoped == true {
			break
		}
		slog.Println(slog.WARN, p.TaskId, "total--", p.Total, "done ---", p.DoneCount())
		if p.Total-p.DoneCount() <= 1 {
			//进行探针扫描
			if len(TakData[p.TaskId].Payloads) != 0 {
				time.Sleep(time.Second * 3)
				ToProbeScan(TakData[p.TaskId])
			}
			delete(EngineArr, p.TaskId)
			delete(TakData, p.TaskId)
			break
		}
		time.Sleep(time.Second * 1)
	}
}

func StopEn(p *scanner.PortClient) {
	p.Stop()
	p.Stoped = true
	delete(EngineArr, p.TaskId)
	delete(TakData, p.TaskId)
	os.Remove(p.TaskId + ".json")
}

func GetPortClient(taskId string) *scanner.PortClient {
	return EngineArr[taskId]
}

func ToProbeScan(p *ProbeReqParam) {

	slog.Println(slog.DEBUG, "开始进行探针扫描")

	path := p.TaskId + ".json"

	res, _ := utils.ReadLineData(path)

	var addrArr []string
	for _, v := range res {
		var portInfo define.Portres
		if err := json.Unmarshal([]byte(v), &portInfo); err != nil {
			slog.Println(slog.DEBUG, "json读取失败==", err)
			return
		}
		if portInfo.Service == "http" || portInfo.Service == "https" || strings.Contains(portInfo.Service, "ssl") {
			slog.Println(slog.DEBUG, portInfo.IP+":"+portInfo.Port)
			addrArr = append(addrArr, portInfo.IP+":"+portInfo.Port)
		}

		slog.Println(slog.DEBUG, portInfo.IP+":"+portInfo.Port, "==", portInfo.Service)
		if portInfo.Service == "" {
			slog.Println(slog.DEBUG, portInfo.IP+":"+portInfo.Port)
			addrArr = append(addrArr, portInfo.IP+":"+portInfo.Port)
		}
	}

	//进行探针扫描
	p.ScanAddrs = addrArr
	p.ServiceType = "probe"
	ProbeScan(p)
}
