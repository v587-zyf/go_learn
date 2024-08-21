package tpc_session

import (
	"kernel/iface"
)

// RECV 接受函数
type Recv func(conn iface.ITcpSession, data any)

// CALL 开始结束回调
type Call func(ss iface.ITcpSession)
