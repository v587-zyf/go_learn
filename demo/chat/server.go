package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Addr string
	Port int

	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	Message chan string
}

func NewServer(addr string, port int) *Server {
	server := &Server{
		Addr:      addr,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}

	return server
}

func (s *Server) Broadcast(user *User, msg string) {
	sendMsg := fmt.Sprintf("[%s]%s:%s", user.Addr, user.Name, msg)

	s.Message <- sendMsg
}

func (s *Server) Handler(conn net.Conn) {
	user := NewUser(conn, s)

	user.Online()

	liveChan := make(chan struct{})
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Exit()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("read err:", err)
				return
			}
			msg := string(buf[:n])

			user.DoMessage(msg)

			liveChan <- struct{}{}
		}
	}()

	for {
		select {
		case <-liveChan:

		case <-time.After(time.Second * 300):
			user.SendMsg("you are kick out")

			user.Exit()
			return
		}
	}
}

func (s *Server) ListenMessage() {
	for {
		msg := <-s.Message

		s.mapLock.Lock()
		for _, user := range s.OnlineMap {
			user.C <- msg
		}
		s.mapLock.Unlock()
	}
}

func (s *Server) Start() {
	// 连接
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Addr, s.Port))
	if err != nil {
		fmt.Println("net listen err:", err)
		return
	}
	fmt.Printf("listen %s:%d...\n", s.Addr, s.Port)

	// 关闭
	defer listener.Close()

	// 启动监听message
	go s.ListenMessage()

	// 接收并操作
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}

		go s.Handler(conn)
	}

}
