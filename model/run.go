package model

import (
	"context"
	"ias_tool_v2/core/slog"
)

func ProbeScan(params *ProbeReqParam) {
	var (
		task   *ProbeTask
		ctxMe  context.Context
		cancel context.CancelFunc
	)

	slog.Println(slog.DEBUG, params)
	ctxMe, cancel = context.WithCancel(context.Background())
	task = NewProbeTask(params)
	task.ChangeTaskStatus(StatusEnum.Received)

	InsertCtx(task.TaskId, task.ServiceType, ctxMe, cancel)

	go task.RecordProgress()

	go task.RecordResult()

	go task.Custom(ctxMe)

	go task.Product(ctxMe, params)

}
