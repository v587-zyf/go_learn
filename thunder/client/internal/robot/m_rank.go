package robot

import (
	pb "comm/t_proto/out/client"
	"fmt"
	"github.com/v587-zyf/gc/iface"
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
)

func (r *Robot) Rank() (err error) {
	rankTypeStr := "请输入消费类型\n"
	for k, v := range pb.RankType_name {
		if k == 0 {
			continue
		}
		rankTypeStr += fmt.Sprintf("%d->%s\n", k, v)
	}
	rankType := InputInt32(rankTypeStr)
	lv := InputInt32("请输入排行榜等级")
	req := &pb.RankReq{
		RankType: pb.RankType(rankType),
		Lv:       lv,
	}
	if err = r.SendMsg(pb.MsgID_Rank_ReqId, req); err != nil {
		log.Error("send rank err", zap.String("err", err.Error()))
	}
	return nil
}
func (r *Robot) RankAck(conn iface.IWsSession, data any) {
	msg := data.(*iface.MessageFrame).Body.(*pb.RankAck)

	log.Debug("RankAck", zap.Any("msg", msg))
}
