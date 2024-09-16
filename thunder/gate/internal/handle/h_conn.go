package handle

import (
	"comm/t_data/redis"
	errCode "comm/t_errcode"
	pb "comm/t_proto/out/client"
	"comm/t_proto/out/server"
	"github.com/v587-zyf/gc/enums"
	"github.com/v587-zyf/gc/errcode"
	"github.com/v587-zyf/gc/gcnet/grpc_msg"
	"github.com/v587-zyf/gc/gcnet/ws_session"
	"github.com/v587-zyf/gc/iface"
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
	"kernel/handler"
)

func Verify(ss iface.IWsSession, msg *iface.MessageFrame) {
	ack := &pb.VerifyAck{
		ErrNo: errcode.ERR_SUCCEED.Int32(),
	}
	defer ss.Send2User(pb.MsgID_Verify_AckId, ack)

	if msg.UserID == 0 {
		ack.ErrNo = errcode.ERR_SIGN.Int32()
		ack.ErrMsg = errcode.ERR_SIGN.Error()
		log.Info("invalid userID", zap.Uint64("userID", msg.UserID))
		return
	}

	req := msg.Body.(*pb.VerifyReq)

	loginKey := redis.FormatUserLogin(msg.UserID)
	loginInfo, err := redis.GetUserLoginInfo(loginKey)
	if err != nil {
		ack.ErrNo = errcode.ERR_REDIS_LOGIN_DATA_NIL.Int32()
		ack.ErrMsg = errcode.ERR_REDIS_LOGIN_DATA_NIL.Error()
		log.Error("GetUserLoginInfo err", zap.Error(err))
		return
	}
	redisData := redis.LoginInfo(loginInfo)

	if (req.Token != redisData.String(0)) || (msg.UserID != redisData.Uin64(2)) {
		ack.ErrNo = errcode.ERR_SIGN.Int32()
		ack.ErrMsg = errcode.ERR_SIGN.Error()
		log.Info("invalid token", zap.Uint64("reqUID", msg.UserID), zap.Uint64("rdbUID", redisData.Uin64(2)),
			zap.String("reqToken", req.Token), zap.String("redisToken", redisData.String(0)))
		return
	}

	if oldSS := ws_session.GetSessionMgr().GetOne(msg.UserID); oldSS != nil {
		oldSS.Send2User(pb.MsgID_Kick_NtfId, &pb.KickNtf{
			ErrNo:  errCode.ERR_KICK_OUT.Int32(),
			ErrMsg: errCode.ERR_KICK_OUT.Error(),
		})
		oldSS.Close()
		ws_session.GetSessionMgr().Disconnect(msg.UserID)
	}

	ss.SetID(msg.UserID)
	ws_session.GetSessionMgr().Add(ss)

	msgID, onlineMsg := makeUserOnlineMsg()
	go func() {
		SendToGame(msg.UserID, msgID, onlineMsg)
		SendToCenter(msg.UserID, msgID, onlineMsg)
	}()
}

func Reconnect(ss iface.IWsSession, msg *iface.MessageFrame) {
	ack := &pb.ReconnectAck{
		ErrNo: errcode.ERR_SUCCEED.Int32(),
	}
	defer ss.Send2User(pb.MsgID_Reconnect_AckId, ack)

	if msg.UserID == 0 {
		ack.ErrNo = errcode.ERR_SIGN.Int32()
		ack.ErrMsg = errcode.ERR_SIGN.Error()
		log.Info("invalid userID", zap.Uint64("userID", msg.UserID))
		return
	}

	req := msg.Body.(*pb.ReconnectReq)

	loginKey := redis.FormatUserLogin(msg.UserID)
	loginInfo, err := redis.GetUserLoginInfo(loginKey)
	if err != nil {
		ack.ErrNo = errcode.ERR_REDIS_LOGIN_DATA_NIL.Int32()
		ack.ErrMsg = errcode.ERR_REDIS_LOGIN_DATA_NIL.Error()
		log.Error("GetUserLoginInfo err", zap.Error(err))
		return
	}
	redisData := redis.LoginInfo(loginInfo)

	if (req.Token != redisData.String(0)) || (msg.UserID != redisData.Uin64(2)) {
		ack.ErrNo = errcode.ERR_SIGN.Int32()
		ack.ErrMsg = errcode.ERR_SIGN.Error()
		log.Info("invalid token", zap.Uint64("reqUID", msg.UserID), zap.Uint64("rdbUID", redisData.Uin64(1)),
			zap.String("reqToken", req.Token), zap.String("redisToken", redisData.String(0)))
		return
	}

	if oldSS := ws_session.GetSessionMgr().GetOne(msg.UserID); oldSS != nil {
		oldSS.Send2User(pb.MsgID_Kick_NtfId, &pb.KickNtf{
			ErrNo:  errCode.ERR_KICK_OUT.Int32(),
			ErrMsg: errCode.ERR_KICK_OUT.Error(),
		})
		oldSS.Close()
		ws_session.GetSessionMgr().Disconnect(msg.UserID)
	}

	if res := redis.CanReconnect(msg.UserID); !res {
		ack.ErrNo = errCode.ERR_RECONNECT_TIMEOUT.Int32()
		ack.ErrMsg = errCode.ERR_RECONNECT_TIMEOUT.Error()
		return
	}

	//ss.AddReconnectTimes()
	redis.DelReconnect(msg.UserID)
	ss.SetID(msg.UserID)
	ws_session.GetSessionMgr().Add(ss)

	msgID, onlineMsg := makeUserOnlineMsg()
	go func() {
		SendToGame(msg.UserID, msgID, onlineMsg)
		SendToCenter(msg.UserID, msgID, onlineMsg)
	}()
}

func SendToGame(userID uint64, msgID int, msg iface.IProtoMessage) {
	reqBytes, err := handler.GetClientWsHandler().Marshal(uint16(msgID), 0, userID, msg)
	if err != nil {
		panic(err)
		return
	}

	msgData := &server.MessageData{Sender: enums.SERVER_GATE, Receiver: enums.SERVER_GAME, Content: reqBytes, MsgType: server.MsgType_Server}
	grpc_msg.SendToMsg(msgData)
}
func SendToCenter(userID uint64, msgID int, msg iface.IProtoMessage) {
	reqBytes, err := handler.GetClientWsHandler().Marshal(uint16(msgID), 0, userID, msg)
	if err != nil {
		panic(err)
		return
	}

	msgData := &server.MessageData{Sender: enums.SERVER_GATE, Receiver: enums.SERVER_CENTER, Content: reqBytes, MsgType: server.MsgType_Server}
	grpc_msg.SendToMsg(msgData)
}

func makeUserOnlineMsg() (int, iface.IProtoMessage) {
	return server.MsgID_UserOnline_NtfId, &server.UserOnlineNtf{}
}
func makeUserOffMsg() (int, iface.IProtoMessage) {
	return server.MsgID_UserOff_NtfId, &server.UserOffNtf{}
}
