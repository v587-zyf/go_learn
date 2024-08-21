package ziface

type IServer interface {
	// 启动
	Start()
	// 停止
	Stop()
	// 运行
	Run()
	// 路由
	AddRouter(msgID uint32, router IRouter)
	// 获取当前server连接管理器
	GetConnMgr() IConnManager
	// 注册OnConnStart 钩子函数的方法
	SetOnConnStart(func(conn IConnection))
	// 注册OnConnStop 钩子函数的方法
	SetOnConnStop(func(conn IConnection))
	// 调用OnConnStart 钩子函数的方法
	CallOnConnStart(conn IConnection)
	// 调用OnConnStop 钩子函数的方法
	CallOnConnStop(conn IConnection)
}
