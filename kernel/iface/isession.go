package iface

import (
	"context"
	"github.com/gorilla/websocket"
	"net"
)

type ITcpSession interface {
	Set(key string, value any)
	Get(key string) (any, bool)
	Remove(key string)

	GetID() uint64
	SetID(id uint64)

	Start()
	Close() error

	GetConn() net.Conn
	GetCtx() context.Context

	Send(msgID uint16, tag uint32, userID uint64, msg IProtoMessage) error
}

type IWsSession interface {
	Set(key string, value any)
	Get(key string) (any, bool)
	Remove(key string)

	GetID() uint64
	SetID(id uint64)

	Start()
	Close() error

	GetConn() *websocket.Conn
	GetCtx() context.Context

	Send(msgID uint16, tag uint32, userID uint64, msg IProtoMessage) error
	Send2User(msgID uint16, msg IProtoMessage) error

	GetReconnectTimes() int
	AddReconnectTimes()
}

type ITcpSessionMgr interface {
	Length() int
	GetOne(UID uint64) ITcpSession
	IsOnline(UID uint64) bool

	Add(ss ITcpSession)
	Disconnect(SID uint64)

	Once(UID uint64, fn func(mgr ITcpSession))
	Range(fn func(uint64, ITcpSession))
}

type ITpcSessionMethod interface {
	Start(ss ITcpSession)
	Recv(conn ITcpSession, data any)
	Stop(ss ITcpSession)
}
type IWsSessionMethod interface {
	Start(ss IWsSession)
	Recv(conn IWsSession, data any)
	Stop(ss IWsSession)
}
