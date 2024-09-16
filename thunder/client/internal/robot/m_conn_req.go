package robot

import (
	pb "comm/t_proto/out/client"
	"fmt"
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
)

func (r *Robot) Verify() (err error) {

	req := &pb.VerifyReq{
		Token: r.token,
	}
	if err = r.SendMsg(pb.MsgID_Verify_ReqId, req); err != nil {
		log.Error("send verify err", zap.String("err", err.Error()))
	}

	//r.SetStatus(enum.STATUS_IN_GAME)
	return nil
}

func (r *Robot) Reconnect() (err error) {
	req := &pb.ReconnectReq{
		Token: r.token,
	}
	if err = r.SendMsg(pb.MsgID_Reconnect_ReqId, req); err != nil {
		log.Error("send verify err", zap.Error(err))
		return
	}

	//r.SetStatus(enum.STATUS_IN_GAME)
	return nil
}

func (r *Robot) GM() (err error) {
	if r.gmStr == "" {
		r.gmStr = "请输入GM类型\n"
		for k, v := range pb.GmType_name {
			r.gmStr += fmt.Sprintf("%d->%s\n", k, v)
		}
	}
	gmType := InputInt32(r.gmStr)
	data := InputString("请输入值")

	req := &pb.GmReq{
		GmType: pb.GmType(gmType),
		Data:   data,
	}
	if err = r.SendMsg(pb.MsgID_Gm_ReqId, req); err != nil {
		log.Error("send open wall err", zap.String("err", err.Error()))
	}
	return nil
}
