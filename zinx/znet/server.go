package znet

import (
	"fmt"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

type Server struct {
	// 服务器名称
	Name string
	// 服务器IP
	Ip string
	// 服务器端口
	Port int
	// 服务器IP版本
	IPVersion string
	// 当前Server消息管理模块 绑定msgID和业务API
	MsgHandle ziface.IMsgHandle
	// 连接管理器
	ConnMgr ziface.IConnManager
	// server创建连接后自动调用函数
	OnConnStart func(conn ziface.IConnection)
	// server删除连接后自动调用函数
	OnConnStop func(conn ziface.IConnection)
}

func NewServer() *Server {
	s := &Server{
		Name:      utils.GlobalObject.Name,
		Ip:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.TcpPort,
		IPVersion: "tcp4",
		MsgHandle: NewMsgHandle(),
		ConnMgr:   NewConnManager(),
	}
	return s
}

func (s *Server) Start() {
	fmt.Println("[Server Start] Serer ",
		s.Name, utils.GlobalObject.Version,
		" Listen on :", s.Ip, ":", s.Port)

	go func() {
		// 0.开启消息队列worker工作池
		s.MsgHandle.StartWorkerPool()

		// 1.获取tpc的addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.Ip, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr err:", err)
			return
		}
		// 2.监听
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen tcp ", s.Ip, ":", s.Port, " err:", err)
			return
		}

		var cid uint32
		cid = 0

		// 3.阻塞等待客户端连接，进行业务
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("accept tcp conn err:", err)
				continue
			}

			// 最大连接数量
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				// todo 给客户端响应超出最大连接的错误包
				fmt.Println("Too Many Connection MaxConn:", utils.GlobalObject.MaxConn)

				conn.Close()
				continue
			}

			dealConn := NewConnection(s, conn, cid, s.MsgHandle)
			go dealConn.Start()

			cid++
		}

	}()
}
func (s *Server) Stop() {
	// 停止|回收
	fmt.Println("[Server Stop] ", s.Name)

	s.ConnMgr.ClearConn()
}

func (s *Server) Run() {
	s.Start()

	// todo 额外业务

	// 阻塞
	select {}
}

func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandle.AddRouter(msgID, router)
	//fmt.Println("add router success!")
}

func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

// 注册OnConnStart 钩子函数的方法
func (s *Server) SetOnConnStart(hookFunc func(conn ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

// 注册OnConnStop 钩子函数的方法
func (s *Server) SetOnConnStop(hookFunc func(conn ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

// 调用OnConnStart 钩子函数的方法
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("---> Call OnConnStart() ...")
		s.OnConnStart(conn)
	}
}

// 调用OnConnStop 钩子函数的方法
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("---> Call OnConnStop() ...")
		s.OnConnStop(conn)
	}
}
