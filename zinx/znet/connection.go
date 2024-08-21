package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"zinx/utils"
	"zinx/ziface"
)

type Connection struct {
	// 当前conn属于哪个server
	TcpServer ziface.IServer
	// 当前连接socket
	Conn *net.TCPConn
	// 当前连接ID
	ConnID uint32
	// 当前连接状态
	isClosed bool
	// 退出信号
	ExitChan chan struct{}
	// 消息管理 msgID对应业务API
	MsgHandle ziface.IMsgHandle
	// 无缓冲 读和写通信
	msgChan chan []byte
	// 连接属性集合
	property map[string]any
	// 保护连接属性的锁
	propertyLock sync.RWMutex
}

func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandle ziface.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer: server,
		Conn:      conn,
		ConnID:    connID,
		isClosed:  false,
		MsgHandle: msgHandle,
		ExitChan:  make(chan struct{}, 1),
		msgChan:   make(chan []byte),
		property:  make(map[string]any),
	}
	c.TcpServer.GetConnMgr().Add(c)
	return c
}

func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine Is Running]")
	defer fmt.Println("[Reader Exit] connID:", c.ConnID, " addr:", c.GetRemoteAddr().String())
	defer c.Stop()

	for {
		//buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		//_, err := c.Conn.Read(buf)
		//if err != nil {
		//	fmt.Println("connID:", c.ConnID, " read err:", err.Error())
		//	continue
		//}

		dp := NewDataPack()
		// 1.读 msg  二进制流 8字节
		headData := make([]byte, dp.GetHeadLen())
		_, err := io.ReadFull(c.GetTCPConnection(), headData)
		if err != nil {
			fmt.Println("read head err:", err)
			break
		}
		// 2.读 head
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack err:", err)
			break
		}
		// 3.根据dataLen 读data 放msg.Data
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			_, err = io.ReadFull(c.GetTCPConnection(), data)
			if err != nil {
				fmt.Println("read data err:", err)
				break
			}
		}
		msg.SetMsgData(data)

		// 得到当前conn的request数据
		req := Request{
			conn: c,
			msg:  msg,
		}

		// 执行注册的路由方法
		//go func(request ziface.IRequest) {
		//	c.Router.PreHandle(request)
		//	c.Router.Handle(request)
		//	c.Router.PostHandle(request)
		//}(&req)

		if utils.GlobalObject.WorkerPoolSize > 0 {
			// 开启了工作池 消息发送工作池即可
			go c.MsgHandle.SendMsgToTaskQueue(&req)
		} else {
			// 普通消息发送 根据msgID找到对应API并执行
			go c.MsgHandle.DoMsgHandle(&req)
		}
	}
}

func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine Is Running]")
	defer fmt.Println("[Write Exit] connID:", c.ConnID, " addr:", c.GetRemoteAddr().String())

	for {
		select {
		case msg := <-c.msgChan:
			_, err := c.Conn.Write(msg)
			if err != nil {
				fmt.Println("Send data err:", err)
				return
			}
		case <-c.ExitChan:
			return
		}
	}
}

func (c *Connection) Start() {
	fmt.Println("[conn start] connID:", c.ConnID)
	// 读数据业务
	go c.StartReader()
	// 写数据业务
	go c.StartWriter()

	// 执行OnConnStart的Hook函数
	c.TcpServer.CallOnConnStart(c)
}

func (c *Connection) Stop() {
	fmt.Println("[conn stop] connID:", c.ConnID)

	if c.isClosed == true {
		return
	}
	c.isClosed = true

	// 调用OnConnStop的Hook函数
	c.TcpServer.CallOnConnStop(c)

	// close conn
	c.Conn.Close()
	// 告诉写入程序退出
	c.ExitChan <- struct{}{}
	// 将当当前连接从ConnMgr摘除
	c.TcpServer.GetConnMgr().Remove(c)

	// close chan
	close(c.ExitChan)
	close(c.msgChan)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) GetRemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) SendMsg(msgID uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("connection closed when send msg")
	}

	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMsgPackage(msgID, data))
	if err != nil {
		fmt.Println("pack msgID:", msgID, " err:", err)
		return errors.New("pack msg err")
	}

	//_, err = c.Conn.Write(binaryMsg)
	//if err != nil {
	//	fmt.Println("Write msgID:", msgID, " err:", err)
	//	return errors.New("send msg err")
	//}

	c.msgChan <- binaryMsg

	return nil
}

// 设置连接属性
func (c *Connection) SetProperty(key string, value any) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value
}

// 获取连接属性
func (c *Connection) GetProperty(key string) (any, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	value, ok := c.property[key]
	if ok {
		return value, nil
	} else {
		return nil, fmt.Errorf("property %s not found", key)
	}
}

// 移除连接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}
