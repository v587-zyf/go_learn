package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

type HelloRouter struct {
	znet.BaseRouter
}

func (p *HelloRouter) Handle(request ziface.IRequest) {
	//fmt.Println("Call Hello Handle...")
	//_, err := request.GetConn().GetTCPConnection().Write([]byte("now...\n"))
	//if err != nil {
	//	fmt.Println("Handle write err:", err)
	//}
	fmt.Println("recv msgID:", request.GetMsgID(),
		" data:", string(request.GetMsgData()))

	err := request.GetConn().SendMsg(201, []byte("hello...hello...hello"))
	if err != nil {
		fmt.Println("send msg err:", err)
	}
}

type PingRouter struct {
	znet.BaseRouter
}

//func (p *PingRouter) PreHandle(request ziface.IRequest) {
//	fmt.Println("PreHandle...")
//	_, err := request.GetConn().GetTCPConnection().Write([]byte("before...\n"))
//	if err != nil {
//		fmt.Println("PreHandle write err:", err)
//	}
//}

func (p *PingRouter) Handle(request ziface.IRequest) {
	//fmt.Println("Call Ping Handle...")
	//_, err := request.GetConn().GetTCPConnection().Write([]byte("now...\n"))
	//if err != nil {
	//	fmt.Println("Handle write err:", err)
	//}
	fmt.Println("recv msgID:", request.GetMsgID(),
		" data:", string(request.GetMsgData()))

	err := request.GetConn().SendMsg(200, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println("send msg err:", err)
	}
}

//func (p *PingRouter) PostHandle(request ziface.IRequest) {
//	fmt.Println("PostHandle...")
//	_, err := request.GetConn().GetTCPConnection().Write([]byte("after...\n"))
//	if err != nil {
//		fmt.Println("PostHandle write err:", err)
//	}
//}

func OnConnStart(conn ziface.IConnection) {
	fmt.Println("-----> OnConnStart Is Called...")
	err := conn.SendMsg(202, []byte("OnConnStart"))
	if err != nil {
		fmt.Println(err)
	}

	conn.SetProperty("Name", "this is test value")
	conn.SetProperty("Age", 25)
	conn.SetProperty("Info", struct {
		Name string
		Age  int
	}{
		"a",
		25,
	})
}

func OnConnStop(conn ziface.IConnection) {
	fmt.Println("-----> OnConnStop Is Called...")
	fmt.Println("connID:", conn.GetConnID(), " Is Exit...")

	val, err := conn.GetProperty("Name")
	if err == nil {
		fmt.Println("Name:", val)
	}
	val, err = conn.GetProperty("Age")
	if err == nil {
		fmt.Println("Age:", val)
	}
	val, err = conn.GetProperty("Info")
	if err == nil {
		fmt.Println("Info:", val)
	}
}

func main() {
	// 1.创建实例
	s := znet.NewServer()

	// 2.hook钩子函数
	s.SetOnConnStart(OnConnStart)
	s.SetOnConnStop(OnConnStop)

	// 3.自定义router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})

	// 4.启动服务
	s.Run()
}
