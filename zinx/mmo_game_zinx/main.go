package main

import (
	"fmt"
	"zinx/mmo_game_zinx/apis"
	"zinx/mmo_game_zinx/core"
	"zinx/ziface"
	"zinx/znet"
)

func OnConnectionAdd(conn ziface.IConnection) {
	// 创建player
	player := core.NewPlayer(conn)
	// 给客户端发送MsgID:1的消息
	player.SyncPid()
	// 给客户端发送MsgID:200的消息
	player.BroadCastStartPosition()
	// 将玩家加入worldMgr
	core.WorldMgr.AddPlayer(player)
	// 将玩家pId绑定到连接属性
	conn.SetProperty("pid", player.Pid)
	// 同步周边玩家 通知当前玩家上线 广播当前玩家位置
	player.SyncSurrounding()

	fmt.Println("---> player PID:", player.Pid, " is arrived")
}

func OnConnectionLost(conn ziface.IConnection) {
	// 1.得到当前玩家
	pid, _ := conn.GetProperty("pid")
	player := core.WorldMgr.GetPlayerByPid(pid.(int32))
	// 2.玩家下线
	player.Offline()
}

func main() {
	// 创建server
	s := znet.NewServer()

	// 客户端钩子
	s.SetOnConnStart(OnConnectionAdd)
	s.SetOnConnStop(OnConnectionLost)

	// 注册路由
	s.AddRouter(2, &apis.WorldChatApi{})
	s.AddRouter(3, &apis.MoveApi{})

	// 启动
	s.Run()
}
