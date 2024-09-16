package robot

import (
	pb "comm/t_proto/out/client"
	"fmt"
	"github.com/v587-zyf/gc/iface"
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
)

func (r *Robot) Hasten() (err error) {
	tStr := "请输入加速类型\n"
	for k, v := range pb.HastenType_name {
		tStr += fmt.Sprintf("%d->%s\n", k, v)
	}
	t := InputInt32(tStr)

	req := &pb.HastenReq{
		HastenType: pb.HastenType(t),
	}
	if err = r.SendMsg(pb.MsgID_Hasten_ReqId, req); err != nil {
		log.Error("send hasten err", zap.String("err", err.Error()))
	}
	return nil
}
func (r *Robot) HastenAck(conn iface.IWsSession, data any) {
	msg := data.(*iface.MessageFrame).Body.(*pb.HastenAck)

	log.Debug("hastenAck", zap.Any("msg", msg))
}
