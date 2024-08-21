package znet

import "zinx/ziface"

type Request struct {
	// 已与客户端建立的连接
	conn ziface.IConnection
	// 客户端请求的数据
	msg ziface.IMessage
}

func (r *Request) GetConn() ziface.IConnection {
	return r.conn
}

func (r *Request) GetMsgData() []byte {
	return r.msg.GetMsgData()
}

func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgId()
}
