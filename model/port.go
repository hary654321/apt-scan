package model

import (
	"fmt"
	"ias_tool_v2/core/scanner"
	"time"

	"github.com/go-playground/validator/v10"
)

type PortTask struct {
	Task
	ChReq       chan ReqParams `json:"request_tcp_param"` //待执行tcp参数的channel
	Payloads    []Port         `json:"payload,omitempty"` //报文信息
	Timeout     int            `json:"timeout"`           //读超时配置
	ChMaxThread chan struct{}
}

type Port struct {
	PortName string  `json:"Port_name" validate:"required,required"`
	Payload  *string `json:"payload" validate:"required,required"`
}

type PortReqParam struct {
	TaskId    string   `json:"task_id" validate:"required"`                //task id
	ScanAddrs []string `json:"addrs" validate:"required,dive,required"`    //待扫描地址列表
	Timeout   int      `json:"timeout,omitempty" validate:"required"`      //发起请求读超时
	Threads   int      `json:"threads" validate:"min=1,max=8192,required"` //最大执行goroutine数量
}

// IsValid 校验参数
func (request *PortReqParam) IsValid() (err error) {

	validate := validator.New()
	if err = validate.Struct(request); err != nil {
		fmt.Println("字段校验出错", err.Error())
		return err
	}
	return nil
}

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

func NewPortTask(p *PortReqParam) *scanner.PortClient {
	PortConfig := scanner.DefaultConfig()
	PortConfig.Threads = p.Threads
	PortConfig.Timeout = time.Duration(p.Timeout) * time.Second // getTimeout(len(app.Setting.Port))
	EngineArr[p.TaskId] = scanner.NewPortScanner(PortConfig)
	return EngineArr[p.TaskId]
}
