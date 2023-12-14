package model

import (
	"context"
	"encoding/base64"
	"encoding/gob"
	"errors"
	"fmt"
	"ias_tool_v2/core/slog"
	"ias_tool_v2/core/utils"
	"ias_tool_v2/logger"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ProbeTask struct {
	Task
	ChReq       chan ReqParams `json:"request_tcp_param"` //待执行tcp参数的channel
	Payloads    []Probe        `json:"payload,omitempty"` //报文信息
	Timeout     int            `json:"timeout"`           //读超时配置
	ChMaxThread chan struct{}
}

type Probe struct {
	ProbeName     string  `json:"probe_name" validate:"required"`
	ProbeProtocol string  `json:"probe_protocol" validate:"required"`
	Payload       *string `json:"payload" validate:"required"`
	MT            string  `json:"probe_match_type"`
	Port          string  ` json:"port" `
}

// GetID get task id in post
type GetTaskID struct {
	TaskId string `json:"task_id" validate:"required"`
}

type ProbeReqParam struct {
	GetTaskID
	ScanAddrs   []string `json:"addrs" validate:"required,dive,required"`             //待扫描地址列表
	Payloads    []Probe  `json:"payload,omitempty" validate:"required,dive,required"` //探针详情,http报文或者tcp报文
	ServiceType string   `json:"service_type,required" validate:"required"`           //服务类型
	Timeout     int      `json:"timeout,omitempty" validate:"required"`               //发起请求读超时
	Threads     int      `json:"threads" validate:"min=1,max=50000,required"`         //最大执行goroutine数量
}

// ReqTcp tcp请求对象封装
type ReqParams struct {
	Addr          string `json:"addr"` //ip:port
	Timeout       int    `json:"-"`    //请求超时时间
	Payload       string `json:"-"`    //原请求的payload信息
	ProbeName     string `json:"probe_name"`
	MT            string `json:"probe_match_type"`
	ProbeProtocol string `json:"probe_protocol"`
}

// IsValid 校验参数
func (request *ProbeReqParam) IsValid() (err error) {

	validate := validator.New()
	if err = validate.Struct(request); err != nil {
		slog.Println(slog.DEBUG, "字段校验出错", err.Error())
		return err
	}
	return nil
}

// NewProbeTask 根据用户请求的参数生成probe任务
// tlsport 如果 addrs 元素port部分存在在tlsport下，则执行https逻辑
func NewProbeTask(p *ProbeReqParam) *ProbeTask {
	task := NewTask(p.TaskId, p.ServiceType, p.Timeout, p.Threads)
	probeTask := &ProbeTask{
		Task:        *task,
		Payloads:    p.Payloads,
		Timeout:     p.Timeout,
		ChReq:       make(chan ReqParams, p.Threads),
		ChMaxThread: make(chan struct{}, p.Threads),
	}
	probeTask.AllSampleNum = getTotal(p.Payloads) * len(p.ScanAddrs)
	return probeTask
}

func getTotal(Payloads []Probe) (total int) {

	for _, v := range Payloads {
		total += len(utils.GetPortArr(v.Port))
	}

	return
}

//func tlsRetryTcp(params ReqParams) (*PeerProbeResult, error) {
//	isTls := IsNotTls
//	return Scan(params, isTls)
//}

//ScanSchedule 探测调度
//func (task *ProbeTask) ScanSchedule(params ReqParams) {
//	var res = &SSLProbeResult{}
//
//	isTls := CheckIsTls(params.Addr)
//
//	resFlag, certResult := CheckIsTlsAndParseCert(params.Addr)
//
//	probeResult, err := Scan(params, isTls)
//	if isTls == IsTLS || resFlag == IsTls {
//		res.SslResult = &TlsResult{Cert: certResult}
//		if err == nil {
//			res.ProbeResult = probeResult
//		} else {
//			tlsRes := PeerProbeResult{
//				ReqInfo:  params,
//				ResPlain: "",
//				ResHex:   "",
//			}
//			res.ProbeResult = &tlsRes
//		}
//		task.chResult <- res
//	}
//	if isTls == IsTLS {
//		probeResult, err := tlsRetryTcp(params)
//		res.SslResult = &TlsResult{Cert: certResult}
//		if err == nil {
//			res.ProbeResult = probeResult
//		} else {
//			tlsRes := PeerProbeResult{
//				ReqInfo:  params,
//				ResPlain: "",
//				ResHex:   "",
//			}
//			res.ProbeResult = &tlsRes
//		}
//		task.chResult <- res
//	}
//}

func (task *ProbeTask) ScanSchedule(params ReqParams) {

	slog.Println(slog.DEBUG, "MT=====", params.MT)
	var res = &SSLProbeResult{}
	var isTls int
	var certData CertData

	if params.ProbeProtocol == "HTTP" {
		//TLS判断
		tlsFlagHalf, certDataHalf := CheckIsTlsAndParseCert(params.Addr)
		if tlsFlagHalf == IsTLS {
			isTls = IsTLS
			if len(certDataHalf.Thumbprint) > 0 {
				certData = certDataHalf
			}
		}
		if len(certDataHalf.Thumbprint) == 0 {
			tlsFlagFull, certDataFull := CheckIsTlsFullAndParseCert(params.Addr)
			if tlsFlagFull == IsTLS {
				isTls = IsTLS
				if len(certDataFull.Thumbprint) > 0 {
					certData = certDataFull
				}
			}
		}
	}
	//反正不管怎么样都要执行TCP探测
	tcpResult, err := Scan(params, isTls)
	// slog.Println(slog.DEBUG, "tcpResult:", tcpResult)
	// slog.Println(slog.DEBUG, "certData:", certData)
	res.SslResult = &TlsResult{Cert: certData}
	if err == nil {
		res.ProbeResult = tcpResult
	} else {
		res.ProbeResult = &PeerProbeResult{
			ReqInfo:  params,
			ResPlain: "",
			ResHex:   "",
		}
	}
	if res.ProbeResult.ResPlain != "" {
		task.chResult <- res
	}

	return
}

func checkIsTttp(payload string) bool {
	HttpMethods := []string{"GET ", "POST ", "PUT ", "DELETE ", "OPTIONS ", "TRACE ", "HEAD ", "CONNECT ", "PATCH "}
	payloadList := strings.Split(payload, "\r\n")
	for _, method := range HttpMethods {
		if strings.HasPrefix(payloadList[0], method) {
			return true
		}
	}
	return false
}

func escapeCharDel(decodedPayload string) string {
	decodedPayload = strings.Replace(decodedPayload, "\\r\\n", "\r\n", -1)
	decodedPayload = strings.Replace(decodedPayload, "\\\"", "\"", -1)
	decodedPayload = strings.Replace(decodedPayload, "\\/", "/", -1)
	return decodedPayload
}

func PayloadPreHandle(decodedPayload, midHost string) string {
	var build strings.Builder
	host := "Host: " + midHost + "\r\n"
	decodedPayload = escapeCharDel(decodedPayload)

	if !checkIsTttp(decodedPayload) {
		return decodedPayload
	}

	if strings.Contains(decodedPayload, "\r\n") {
		var midPayload strings.Builder
		var header string
		var body string
		var contentLength int
		payloadSplit := strings.Split(decodedPayload, "\r\n\r\n")
		if len(payloadSplit) > 1 {
			header = payloadSplit[0]
			body = strings.Join(payloadSplit[1:], "\r\n\r\n")
			contentLength = len(body)
		} else {
			header = decodedPayload
			contentLength = 0
		}

		var headerBuilder strings.Builder
		isHasHost := false
		var first string
		for idx, val := range strings.Split(header, "\r\n") {
			if idx == 0 {
				first = val + "\r\n"
				continue
			}
			if strings.HasPrefix(strings.ToLower(val), "content-length") {
				headerBuilder.WriteString(fmt.Sprintf("Content-Length: %d\r\n", contentLength))
				continue
			}
			if strings.HasPrefix(strings.ToLower(val), "host:") {
				isHasHost = true
			}
			headerBuilder.WriteString(val + "\r\n")
		}
		headerBuilder.WriteString("\r\n")
		if !isHasHost {
			header = fmt.Sprintf("%s%s%s", first, host, headerBuilder.String())
		} else {
			header = fmt.Sprintf("%s%s", first, headerBuilder.String())
		}
		midPayload.WriteString(header)
		midPayload.WriteString(body)
		return midPayload.String()
	}

	return build.String()
}

// Product 生产参数
func (task *ProbeTask) Product(ctx context.Context, p *ProbeReqParam) {

	ctxSub, cancel := context.WithCancel(ctx)
	defer cancel()
	task.ChangeTaskStatus(StatusEnum.Started)
	for _, addr := range p.ScanAddrs {
		select {
		case <-ctxSub.Done():
			cancel()
			logger.Infof("product: recv stop signal")
			goto EXIT
		default:
			//生产解析生成的结构体对象
			for _, payload := range task.Payloads {

				portStr := payload.Port
				portArr := utils.GetPortArr(portStr)
				for _, port := range portArr {
					decodePayload, err := base64.StdEncoding.DecodeString(*payload.Payload)
					strPayload := string(decodePayload)

					// slog.Println(slog.DEBUG, "payload:", strPayload)

					// slog.Println(slog.DEBUG, "PayloadPreHandle:", PayloadPreHandle(strPayload, addr))
					if err != nil {
						slog.Println(slog.DEBUG, err.Error())
						continue
					}
					// slog.Println(slog.DEBUG, "payload:", addr+":"+port)
					ReqHttpParam := ReqParams{
						Addr:          addr + ":" + port,
						Timeout:       p.Timeout,
						Payload:       PayloadPreHandle(strPayload, addr),
						ProbeName:     payload.ProbeName,
						MT:            payload.MT,
						ProbeProtocol: payload.ProbeProtocol,
					}
					task.ChReq <- ReqHttpParam
					task.ChMaxThread <- struct{}{}
				}

			}
		}
	}
	return
EXIT:
	logger.Infof("product: stop ok")
}

// Custom 消费者负责执行
func (task *ProbeTask) Custom(ctx context.Context) {
	ctxSub, cancel := context.WithCancel(ctx)
	defer cancel()
	for {
		select {
		case <-ctxSub.Done():
			cancel()
			logger.Infof("custom: recv stop signal")
			goto EXIT
		case reqParams, ok := <-task.ChReq:
			if !ok {
				task.ChangeTaskStatus(StatusEnum.Finished)
				goto FINISH
			}
			go func() {
				task.ScanSchedule(reqParams)
				<-task.ChMaxThread
				task.AddProgress()
			}()
		}
	}
FINISH:
	logger.Infof("custom: execute ok")

EXIT:
	task.ChangeTaskStatus(StatusEnum.EarlyExit)
	logger.Infof("custom: stop ok")
}

// TlsRes tls相关信息
type TlsRes struct {
	Subject string
	Before  string
	After   string
}

type PeerProbeResult struct {
	ReqInfo  ReqParams `json:"req_info"`
	ResPlain string    `json:"res_plain"`
	ResHex   string    `json:"res_hex"`
}

// SSLProbeResult 结果封装
type SSLProbeResult struct {
	ProbeResult *PeerProbeResult `json:"probe_result"`
	SslResult   *TlsResult       `json:"ssl_result,omitempty"`
}

// Scan 探测实现
func Scan(req ReqParams, isTls int) (res *PeerProbeResult, err error) {
	var resp string
	res = &PeerProbeResult{}

	// slog.Println(slog.DEBUG, "req.Payload", req.ProbeProtocol)

	if req.ProbeProtocol == "TCP" {
		resp, err = TcpSend("tcp", req.Addr, req.Payload, req.Timeout)

		if err != nil {
			slog.Println(slog.DEBUG, "tcp:", resp, err)
			return res, err
		}
		if len(resp) == 0 {
			err = errors.New("no result")
			//return res, err
		}
		res.ResPlain = resp

		// dump, _ := hex.DecodeString(resp)

		res.ResHex = "" // hex.Dump(dump)
	} else {
		if isTls == IsTLS {
			resp, err = HttpSend("tls", req.Addr, req.Payload, req.Timeout)

			if err != nil {
				slog.Println(slog.DEBUG, "resp:", resp, "err:", err)
				return res, err
			}
			if len(resp) == 0 {
				err = errors.New("no result")
				//return res, err
			}
			res.ResPlain = resp

			res.ResHex = "" // hex.Dump([]byte(resp))

		} else {
			resp, err = HttpSend("tcp", req.Addr, req.Payload, req.Timeout)
			if err != nil {
				// slog.Println(slog.DEBUG, "http", req.Addr, "====", req.Payload, err)
				return res, err
			}
			if len(resp) == 0 {
				err = errors.New("no result")
				//return res, err
			}
			res.ResPlain = resp

			res.ResHex = "" // hex.Dump([]byte(resp))
		}
	}

	res.ReqInfo = req

	return
}

// Encode 序列化任务结构体到文件
func (task *ProbeTask) Encode() {
	pickleFd, _ := os.OpenFile(task.PickleFilePath, os.O_WRONLY, 0666)
	defer func(pickleFd *os.File) {
		err := pickleFd.Close()
		if err != nil {
			logger.Errorf(err.Error())
		}
	}(pickleFd)
	enc := gob.NewEncoder(pickleFd)
	err := enc.Encode(&task)
	if err != nil {
		logger.Errorf(err.Error())
	}
}

// ChangeTaskStatus 修改任务状态
func (task *ProbeTask) ChangeTaskStatus(status int) {
	task.Status = status
	task.Encode()
}

// ProbeTaskDecode 反序列化任务结构体到文件
func ProbeTaskDecode(taskId, serviceType string) (*ProbeTask, error) {
	task := &ProbeTask{}
	picklePath := filepath.Join(PicklePathFolder(serviceType), taskId+".job")

	file, err := os.Open(picklePath)
	if err != nil {
		return nil, err
	}
	dec := gob.NewDecoder(file)
	err = dec.Decode(task)
	err = file.Close()
	if err != nil {
		return nil, err
	}
	return task, err
}
