package model

type TaskStatus struct {
	Received  int //接收
	Started   int //开始
	Finished  int //完成
	Errored   int //执行错误
	EarlyExit int //提前退出
}

var StatusEnum = TaskStatus{
	Received: 1,
	Started:  2,
	Finished: 3,
}

//GetStatusName 通过status获取状态中文名称
func GetStatusName(statusNum int) (statusName string) {
	if statusNum == StatusEnum.Received {
		return "Received"
	} else if statusNum == StatusEnum.Started {
		return "Started"
	} else if statusNum == StatusEnum.Finished {
		return "Finished"
	} else {
		return "EarlyExit"
	}
}
