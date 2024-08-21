package znet

import (
	"fmt"
	"sync"
	"zinx/ziface"
)

/*
*
连接管理模块
*/
type ConnManager struct {
	connections map[uint32]ziface.IConnection
	connLock    sync.RWMutex
}

func NewConnManager() *ConnManager {
	cm := &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
	return cm
}

func (connMgr *ConnManager) Add(conn ziface.IConnection) {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	// conn加入map
	connMgr.connections[conn.GetConnID()] = conn
	fmt.Println("[Add ConnID:", conn.GetConnID(), " Succ] All connLen:", connMgr.Len())
}

func (connMgr *ConnManager) Remove(conn ziface.IConnection) {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	// 删除连接
	delete(connMgr.connections, conn.GetConnID())
	fmt.Println("[Remove ConnID:", conn.GetConnID(), " Succ] All connLen:", connMgr.Len())
}

func (connMgr *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	conn, ok := connMgr.connections[connID]
	if ok {
		return conn, nil
	} else {
		return nil, fmt.Errorf("connID:%d not found", connID)
	}
}

func (connMgr *ConnManager) Len() int {
	return len(connMgr.connections)
}

func (connMgr *ConnManager) ClearConn() {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	// 删除conn并停止工作
	for connID, conn := range connMgr.connections {
		// stop
		conn.Stop()
		// remove
		delete(connMgr.connections, connID)
	}
	fmt.Println("[Clear Conn Succ] All connLen:", connMgr.Len())
}
