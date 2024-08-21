package ziface

type IMsgHandle interface {
	// 调度router
	DoMsgHandle(request IRequest)
	// 添加处理逻辑
	AddRouter(msgID uint32, router IRouter)
	// 启动Worker工作池
	StartWorkerPool()
	// 将消息交给TaskQueue 由router处理
	SendMsgToTaskQueue(request IRequest)
}
