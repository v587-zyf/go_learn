package core

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"math/rand/v2"
	"sync"
	"zinx/mmo_game_zinx/pb/pb"
	"zinx/ziface"
)

var (
	PidGen int32 = 1
	IDlOCK sync.Mutex
)

type Player struct {
	Pid  int32              // 玩家ID
	Conn ziface.IConnection // 玩家连接
	X    float32            // 平面X坐标
	Y    float32            // 高度
	Z    float32            // 平面Y坐标(注意不是Y)
	V    float32            // 旋转的0-360角度
}

func NewPlayer(conn ziface.IConnection) *Player {
	IDlOCK.Lock()
	defer IDlOCK.Unlock()
	id := PidGen
	PidGen++

	p := &Player{
		Pid:  id,
		Conn: conn,
		X:    float32(160 + rand.IntN(10)), // 随机160坐标 基于X若干偏移
		Y:    0,
		Z:    float32(140 + rand.IntN(20)), // 随机140坐标 基于Y若干偏移
		V:    0,                            // 角度0
	}

	return p
}

// 把protobuf数据序列化 调用zinx的sendMsg
func (p *Player) SendMsg(msgID uint32, data proto.Message) {
	// 1.proto message序列化 转二进制
	msg, err := proto.Marshal(data)
	if err != nil {
		fmt.Println("marshal msg err:", err)
		return
	}
	// 2.将二进制 通过zinx的sendMsg发送客户端
	if p.Conn == nil {
		fmt.Println("connection in player is nil")
		return
	}

	if err = p.Conn.SendMsg(msgID, msg); err != nil {
		fmt.Println("player sendMsg err")
		return
	}
}

// 同步玩家上线
func (p *Player) SyncPid() {
	// 组建消息MsgID:0消息
	protoMsg := &pb.SyncPid{Pid: p.Pid}
	// 发送客户端
	p.SendMsg(1, protoMsg)
}

// 广播玩家出生地点
func (p *Player) BroadCastStartPosition() {
	// 组建消息MsgID:200消息
	protoMsg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			}},
	}
	// 发送客户端
	p.SendMsg(200, protoMsg)
}

// 玩家广播世界聊天
func (p *Player) Talk(content string) {
	// 组建msgID：200 proto数据
	protoMsg := &pb.BroadCast{
		Pid:  p.Pid,
		Tp:   1,
		Data: &pb.BroadCast_Content{Content: content},
	}
	// 得到当前所有在线玩家
	players := WorldMgr.GetAllPlayer()
	// 向所有玩家（包括自己）发送msgID：200消息
	for _, player := range players {
		player.SendMsg(200, protoMsg)
	}
}

// 同步玩家上线的位置信息
func (p *Player) SyncSurrounding() {
	// 1.获取当前玩家周围玩家信息（九宫格）
	pIds := WorldMgr.AOIManager.GetElemsSByPos(p.X, p.Z)
	players := make([]*Player, 0, len(pIds))
	for _, pId := range pIds {
		players = append(players, WorldMgr.GetPlayerByPid(int32(pId)))
	}
	// 2.将当前玩家位置通过MsgID：200消息发送给周围玩家
	// 2.1. 组建MsgID：200消息
	protoMsg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			}},
	}
	// 2.2.周围玩家各自发送200消息
	for _, player := range players {
		player.SendMsg(200, protoMsg)
	}
	// 3.将周围玩家位置通过MsgID：200消息发送给当前玩家客户端
	// 3.1.组建MsgID：200消息
	playersProtoMsg := make([]*pb.Player, 0, len(players))
	for _, player := range players {
		pbPlayer := &pb.Player{
			Pid: player.Pid,
			P: &pb.Position{
				X: player.X,
				Y: player.Y,
				Z: player.Z,
				V: player.V,
			},
		}
		playersProtoMsg = append(playersProtoMsg, pbPlayer)
	}
	// 3.2.将组建数据发给当前玩家
	syncPlayersProtoMsg := &pb.SyncPlayers{
		Players: playersProtoMsg[:],
	}
	p.SendMsg(202, syncPlayersProtoMsg)
}

func (p *Player) UpdatePos(x, y, z, v float32) {
	// 1.更新当前玩家坐标
	p.X = x
	p.Y = y
	p.Z = z
	p.V = v
	// 2.组件MsgID:200 广播协议
	protoMsg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  3,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			}},
	}
	// 3.获取当前玩家周边玩家（九宫格）
	players := p.GetSurroundingPlayers()
	// 4.依次给周边玩家发送坐标消息
	for _, player := range players {
		player.SendMsg(200, protoMsg)
	}
}

// 获取当前玩家周边玩家（九宫格）
func (p *Player) GetSurroundingPlayers() []*Player {
	// 得到当前AOI九宫格所有玩家PID
	elems := WorldMgr.AOIManager.GetElemsSByPos(p.X, p.Z)
	// 将所有pid对应player放切片中
	players := make([]*Player, 0, len(elems))
	for _, elem := range elems {
		players = append(players, WorldMgr.GetPlayerByPid(int32(elem)))
	}

	return players
}

// 玩家下线
func (p *Player) Offline() {
	// 1.得到当前玩家周边九宫格玩家
	players := p.GetSurroundingPlayers()
	// 2.给周围玩家广播MsgID:201消息
	protoMsg := &pb.SyncPid{Pid: p.Pid}
	for _, player := range players {
		player.SendMsg(201, protoMsg)
	}
	// 3.将当前玩家从世界管理器和AOI删除
	WorldMgr.AOIManager.RemoveFromPos(int(p.Pid), p.X, p.Z)
	WorldMgr.RemovePlayer(p.Pid)
}
