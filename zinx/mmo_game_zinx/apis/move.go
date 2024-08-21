package apis

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"zinx/mmo_game_zinx/core"
	"zinx/mmo_game_zinx/pb/pb"
	"zinx/ziface"
	"zinx/znet"
)

type MoveApi struct {
	znet.BaseRouter
}

func (m *MoveApi) Handle(request ziface.IRequest) {
	// 1.解析客户端proto协议
	protoMsg := &pb.Position{}
	if err := proto.Unmarshal(request.GetMsgData(), protoMsg); err != nil {
		fmt.Println("proto unmarshal err:", err)
		return
	}
	// 2.得到发送位置的玩家
	pid, err := request.GetConn().GetProperty("pid")
	if err != nil {
		fmt.Println("GetProperty pid err:", err)
		return
	}
	//fmt.Printf("Player Pid:%d move(%f,%f,%f,%f)\n",
	//	pid, protoMsg.GetX(), protoMsg.GetY(), protoMsg.GetZ(), protoMsg.GetV())
	// 3.向其他玩家广播当前玩家位置
	player := core.WorldMgr.GetPlayerByPid(pid.(int32))
	// 3.1.广播并更新当前玩家坐标
	player.UpdatePos(protoMsg.GetX(), protoMsg.GetY(), protoMsg.GetZ(), protoMsg.GetV())
}
