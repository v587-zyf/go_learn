package ziface

type IRequest interface {
	// 获取当前连接
	GetConn() IConnection
	// 获取当前消息数据
	GetMsgData() []byte
	// 获取当前消息id
	GetMsgID() uint32
}
