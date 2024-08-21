package main

import (
	"fmt"
	"net"
	"strings"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	connAddr := conn.RemoteAddr().String()

	user := &User{
		Name:   connAddr,
		Addr:   connAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}

	go user.ListenMessage()

	return user
}

func (u *User) ListenMessage() {
	for {
		msg := <-u.C

		u.conn.Write([]byte(msg + "\r\n"))
	}
}

func (u *User) Online() {
	u.server.mapLock.Lock()
	u.server.OnlineMap[u.Name] = u
	u.server.mapLock.Unlock()

	u.server.Broadcast(u, "online")
}

func (u *User) Exit() {
	u.server.mapLock.Lock()
	delete(u.server.OnlineMap, u.Name)
	u.server.mapLock.Unlock()

	u.server.Broadcast(u, "exit")

	if _, ok := <-u.C; !ok {
		close(u.C)
	}
	u.conn.Close()
}

func (u *User) SendMsg(msg string) {
	u.conn.Write([]byte(msg))
}

func (u *User) DoMessage(msg string) {
	if msg == "who" {
		u.server.mapLock.Lock()
		for _, user := range u.server.OnlineMap {
			sendMsg := fmt.Sprintf("[%s]%s online...\r\n", user.Addr, user.Name)
			u.SendMsg(sendMsg)
		}
		u.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		newName := strings.Split(msg, "|")[1]
		if _, ok := u.server.OnlineMap[newName]; !ok {
			u.server.mapLock.Lock()
			u.server.OnlineMap[newName] = u
			delete(u.server.OnlineMap, u.Name)
			u.server.mapLock.Unlock()

			u.Name = newName
			u.SendMsg(fmt.Sprintf("rename success, your name is %s\r\n", u.Name))
		} else {
			u.SendMsg("this name is used early, please try another name\r\n")
		}
	} else if len(msg) > 3 && msg[:3] == "to|" {
		// 1.判断格式
		msgArr := strings.Split(msg, "|")
		if len(msgArr) != 3 {
			u.SendMsg("format err.need write: to|username|content\r\n")
			return
		}
		// 2.判断用户名
		userName := msgArr[1]
		if userName == "" {
			u.SendMsg("username is empty\r\n")
			return
		}
		user, ok := u.server.OnlineMap[userName]
		if !ok {
			u.SendMsg(fmt.Sprintf("userName:%s is not found\r\n", userName))
			return
		}
		// 3.判断发送内容
		content := msgArr[2]
		if content == "" {
			u.SendMsg("content is empty\r\n")
			return
		}
		user.SendMsg(fmt.Sprintf("%s send to you:%s\r\n", u.Name, content))
	} else {
		u.server.Broadcast(u, msg)
	}
}
