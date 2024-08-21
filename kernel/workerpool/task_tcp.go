package workerpool

import (
	"kernel/iface"
	"kernel/session/tpc_session"
)

func (p *WorkerPool) AssignNetTask(fn tpc_session.Recv, ss iface.ITcpSession, data any) error {
	return Assign(&NetTask{
		Func:    fn,
		Session: ss,
		Data:    data,
	})
}

type NetTask struct {
	Func    tpc_session.Recv
	Session iface.ITcpSession
	Data    any
}

func (t *NetTask) Do() {
	t.Func(t.Session, t.Data)
}
