package model

import (
	"context"
)

type Ctx struct {
	CtxCtx    context.Context
	CtxCancel context.CancelFunc
}

var TaskCtxMgr = make(map[string]Ctx)

func genTaskCtxMgrKey(taskId, serviceType string) string {
	return "ctx-" + taskId + serviceType
}

//InsertCtx 将task的context存储到TaskCtxMgr
func InsertCtx(taskId, serviceType string, ctx context.Context, cancel context.CancelFunc) {
	ct := Ctx{CtxCtx: ctx, CtxCancel: cancel}
	TaskCtxMgr[genTaskCtxMgrKey(taskId, serviceType)] = ct
}

//GetCtx 取出TaskCtxMgr中的context
func GetCtx(taskId, serviceType string) (ctx context.Context, cancel context.CancelFunc, ok bool) {
	ctxMap, ok := TaskCtxMgr[genTaskCtxMgrKey(taskId, serviceType)]
	if ok {
		return ctxMap.CtxCtx, ctxMap.CtxCancel, true
	}
	return nil, nil, false
}
