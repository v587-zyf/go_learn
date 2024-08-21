package znet

import (
	"fmt"
	"zinx/utils"
	"zinx/ziface"
)

type MsgHandle struct {
	Apis           map[uint32]ziface.IRouter // msgID:router
	TaskQueue      []chan ziface.IRequest    // worker工作池消息队列
	WorkerPoolSize uint32                    // worker工作池数量
}

func NewMsgHandle() *MsgHandle {
	mh := &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}

	return mh
}

// 调度router
func (mh *MsgHandle) DoMsgHandle(request ziface.IRequest) {
	// 1.取msgID
	handle, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Printf("msgID:%d not found\n", request.GetMsgID())
		return
	}
	// 2.根据msgID调度
	handle.PreHandle(request)
	handle.Handle(request)
	handle.PostHandle(request)
}

// 添加处理逻辑
func (mh *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	// 1.是否已存在
	_, ok := mh.Apis[msgID]
	if ok {
		panic(fmt.Sprintf("repeat add msgID:%d", msgID))
	}
	// 2.添加
	mh.Apis[msgID] = router

	fmt.Println("add api msgID:", msgID)
}

// 启动一个Worker工作池(只能发生一次)
func (mh *MsgHandle) StartWorkerPool() {
	var i uint32
	for ; i < mh.WorkerPoolSize; i++ {
		// 1.给每个消息队列分配空间
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		// 2.启动到当前worker 阻塞等channel传递消息
		go mh.startOneWorker(i, mh.TaskQueue[i])
	}
}

// 启动一个Worker工作流程
func (mh *MsgHandle) startOneWorker(workerID uint32, taskQueue chan ziface.IRequest) {
	fmt.Println("[Worker Start] workerID:", workerID)

	// 阻塞等待队列消息
	for {
		select {
		// 接到消息 执行request绑定业务
		case request := <-taskQueue:
			mh.DoMsgHandle(request)
		}
	}
}

// 将消息交给TaskQueue 由router处理
func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	// 1.消息平均分配给worker
	// 根据客户端connID分配
	// todo 这里只是简单取模 后续可完善
	workerID := request.GetConn().GetConnID() % mh.WorkerPoolSize

	// 2.消息发给worker的TaskQueue
	mh.TaskQueue[workerID] <- request
}
