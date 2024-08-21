package ziface

import "net"

type IConnection interface {
	// 启动
	Start()
	// 停止
	Stop()
	// 获取连接socket
	GetTCPConnection() *net.TCPConn
	// 获取连接ID
	GetConnID() uint32
	// 获取连接地址
	GetRemoteAddr() net.Addr
	// 发送数据
	SendMsg(msgID uint32, data []byte) error
	// 设置连接属性
	SetProperty(key string, value any)
	// 获取连接属性
	GetProperty(key string) (any, error)
	// 移除连接属性
	RemoveProperty(key string)
}

type HandleFunc func(*net.TCPConn, []byte, int) error
