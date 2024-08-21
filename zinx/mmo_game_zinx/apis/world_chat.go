package apis

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"zinx/mmo_game_zinx/core"
	"zinx/mmo_game_zinx/pb/pb"
	"zinx/ziface"
	"zinx/znet"
)

// 世界聊天
type WorldChatApi struct {
	znet.BaseRouter
}

func (wc *WorldChatApi) Handle(request ziface.IRequest) {
	// 解析proto协议
	protoMsg := &pb.Talk{}
	if err := proto.Unmarshal(request.GetMsgData(), protoMsg); err != nil {
		fmt.Println("proto unmarshal error:", err)
		return
	}
	// 获取当前玩家
	pid, err := request.GetConn().GetProperty("pid")
	if err != nil {
		fmt.Println()
		return
	}
	// 根据pId找对应player
	player := core.WorldMgr.GetPlayerByPid(pid.(int32))
	// 将消息广播给在线玩家
	player.Talk(protoMsg.Content)
}
