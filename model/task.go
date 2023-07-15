package model

import (
	"bufio"
	"encoding/json"
	"github.com/wxnacy/wgo/file"
	"ias_tool_v2/config"
	"ias_tool_v2/logger"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type Tasker interface {
	WriteResult(content interface{}) (err error)   //修改任务结果
	WriteProgress(content interface{}) (err error) //修改任务结果
	AddProgress()                                  //修改任务进度
}

type Task struct {
	TaskId           string `json:"task_id" validate:"required"`               //task id
	Status           int    `json:"status" validate:"required"`                //任务状态
	Timeout          int    `json:"timeout" validate:"min=1,max=30"`           //每个请求最大超时时间
	Threads          int    `json:"threads,required" validate:"min=1,max=100"` //任务最大并行数量
	ServiceType      string `json:"service_type" validate:"required"`          //
	AllSampleNum     int    `json:"all_sample_num"`
	chResult         chan interface{}
	chProgress       chan struct{}
	ResultFilePath   string
	ProgressFilePath string
	PickleFilePath   string
}

type Progress struct {
	OkAddr  int `json:"ok_addr"`
	AllAddr int `json:"all_addr"`
}

func JsonProgress(ok, all int) []byte {
	p := &Progress{
		OkAddr:  ok,
		AllAddr: all,
	}
	jsonData, _ := json.Marshal(p)
	return jsonData
}

func NewTask(taskid, serviceType string, timeout, threads int) *Task {
	res, progress, pickle := GetCurTaskFile(serviceType, taskid)
	//创建进度和结果文件
	_, _ = os.Create(res)
	_, _ = os.Create(progress)
	_, _ = os.Create(pickle)

	return &Task{
		TaskId:      taskid,
		Status:      StatusEnum.Received,
		Timeout:     timeout,
		Threads:     threads,
		ServiceType: serviceType,

		ResultFilePath:   res,
		ProgressFilePath: progress,
		PickleFilePath:   pickle,

		chResult:   make(chan interface{}, 0),
		chProgress: make(chan struct{}, 0),
	}
}

// PicklePathFolder 根据传入的task_id和service_type 返回对应的.job文件夹
func PicklePathFolder(serviceType string) string {
	if runtime.GOOS == "windows" {
		return filepath.Join(config.GlobalPath, config.PicklePath, serviceType)
	} else {
		return filepath.Join(config.PicklePath, serviceType)
	}
}

//GetCurTaskFile 通过提交的服务类型和taskid返回对应的文件路径,!!!重要函数
func GetCurTaskFile(serviceType string, taskid string) (string, string, string) {
	//获取存放结果的总目录
	sysName := runtime.GOOS
	ServiceResMap := LoadServiceResMap()
	//获取当前任务的结果文件路径
	resPath := filepath.Join(ServiceResMap[serviceType][sysName], taskid)

	if !file.IsDir(resPath) {
		_ = os.MkdirAll(resPath, os.ModePerm)
	}

	picklePath := PicklePathFolder(serviceType)

	if !file.IsDir(picklePath) {
		_ = os.MkdirAll(picklePath, os.ModePerm)
	}

	res := filepath.Join(resPath, "result.json")
	progress := filepath.Join(resPath, "progress.dat")
	pickle := filepath.Join(picklePath, taskid+".job")
	return res, progress, pickle
}

// AddProgress 修改任务进度，给任务增加进度
func (t *Task) AddProgress() {
	t.chProgress <- struct{}{}
}

//RecordProgress 记录进度信息
func (t *Task) RecordProgress() {
	index := 0

	_ = t.WriteProgress(JsonProgress(0, t.AllSampleNum))

	for range t.chProgress {
		index += 1
		_ = t.WriteProgress(JsonProgress(index, t.AllSampleNum))
	}
}

// WriteProgress 封装写进度
func (t *Task) WriteProgress(content []byte) (err error) {
	var (
		fd *os.File
	)

	fd, err = os.OpenFile(t.ProgressFilePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		goto ERR
	}

	_, err = fd.Write(content)
	if err != nil {
		goto ERR
	}
	_, err = fd.Write([]byte("\n"))
	if err != nil {
		goto ERR
	}
	err = fd.Close()
	if err != nil {
		return err
	}

	return nil
ERR:
	logger.Errorf(err.Error())
	return err
}

func GetServiceType(url string) string {
	return strings.Split(url, "/")[2]
}

//RecordResult 记录进度信息
func (t *Task) RecordResult() {
	for result := range t.chResult {
		_ = t.WriteResult(result)
	}
}

//WriteResult 封装写结果文件
func (t *Task) WriteResult(content interface{}) (err error) {
	var (
		fd      *os.File
		jsonRes []byte
	)

	if fd, err = os.OpenFile(t.ResultFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644); err != nil {
		goto ERR
	}
	jsonRes, err = json.Marshal(content)
	if err != nil {
		goto ERR
	}
	_, err = fd.Write(jsonRes)
	if err != nil {
		goto ERR
	}
	_, err = fd.Write([]byte("\n"))
	if err != nil {
		goto ERR
	}
	err = fd.Close()
	if err != nil {
		return err
	}
	return nil
ERR:
	logger.Errorf(err.Error())
	return err
}

//truncateFile 清空结果文件
func (t *Task) TruncateFile() {
	_ = os.Truncate(t.ResultFilePath, 0)
}

//ReadResultFile 读取结果文件
func (t *Task) ReadResultFile(result interface{}) (res []interface{}, err error) {
	fi, _ := os.Open(t.ResultFilePath)

	defer fi.Close()
	br := bufio.NewReader(fi)

	for {
		byteStr, err := br.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		resLine := result

		if err = json.Unmarshal(byteStr, &resLine); err != nil {
			log.Println("ERROR", "parse failed json", err.Error())
			return nil, err
		} else {
			res = append(res, resLine)
		}
	}

	return res, nil
}

//ReadProgressFile 读取进度
func (t *Task) ReadProgressFile() (ok, all int) {
	var (
		buf []byte
		err error
	)
	progress := &Progress{}
	if buf, err = os.ReadFile(t.ProgressFilePath); err != nil {
		goto ERR
	}
	if err = json.Unmarshal(buf, progress); err != nil {
		goto ERR
	}
	return progress.OkAddr, progress.AllAddr
ERR:
	logger.Errorf("Error querying progress information: %v", err.Error())
	return
}
