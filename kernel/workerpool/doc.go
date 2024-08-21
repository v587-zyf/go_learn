package workerpool

import (
	"kernel/iface"
	"kernel/session/tpc_session"
	"kernel/session/ws_session"
	"sync"
)

type ITask interface {
	Do()
}

var defaultWorkPoll *WorkerPool
var once sync.Once

func Init(cfg ...*Config) (err error) {
	once.Do(func() {
		defaultWorkPoll, err = New(cfg...)
		defaultWorkPoll.Start()
	})

	return
}

func Assign(task ITask) error {
	return defaultWorkPoll.Assign(task)
}

func AssignTcpTask(fn tpc_session.Recv, ss iface.ITcpSession, data any) error {
	return defaultWorkPoll.Assign(&NetTask{
		Func:    fn,
		Session: ss,
		Data:    data,
	})
}

func AssignWsTask(fn ws_session.Recv, ss iface.IWsSession, data any) error {
	return defaultWorkPoll.Assign(&WsTask{
		Func:    fn,
		Session: ss,
		Data:    data,
	})
}
