package server
import (
	"reflect"
)

var msgProtoTypes = make(map[uint16]reflect.Type)
var msgNames = make(map[uint16]string)

func init() {
	msgProtoTypes[MsgID_ZeroMsgId]=reflect.TypeOf((*ZeroMsg)(nil)).Elem()
	msgProtoTypes[MsgID_MessageDataId]=reflect.TypeOf((*MessageData)(nil)).Elem()
	msgProtoTypes[MsgID_Register_ReqId]=reflect.TypeOf((*RegisterReq)(nil)).Elem()
	msgProtoTypes[MsgID_Register_AckId]=reflect.TypeOf((*RegisterAck)(nil)).Elem()
	msgProtoTypes[MsgID_Send2UserId]=reflect.TypeOf((*Send2User)(nil)).Elem()
	msgProtoTypes[MsgID_BroadcastId]=reflect.TypeOf((*Broadcast)(nil)).Elem()
	msgProtoTypes[MsgID_UserOnline_NtfId]=reflect.TypeOf((*UserOnlineNtf)(nil)).Elem()
	msgProtoTypes[MsgID_UserOff_NtfId]=reflect.TypeOf((*UserOffNtf)(nil)).Elem()
	msgProtoTypes[MsgID_UserIncome_NtfId]=reflect.TypeOf((*UserIncomeNtf)(nil)).Elem()
	msgNames[MsgID_ZeroMsgId]="ZeroMsg"
	msgNames[MsgID_MessageDataId]="MessageData"
	msgNames[MsgID_Register_ReqId]="RegisterReq"
	msgNames[MsgID_Register_AckId]="RegisterAck"
	msgNames[MsgID_Send2UserId]="Send2User"
	msgNames[MsgID_BroadcastId]="Broadcast"
	msgNames[MsgID_UserOnline_NtfId]="UserOnlineNtf"
	msgNames[MsgID_UserOff_NtfId]="UserOffNtf"
	msgNames[MsgID_UserIncome_NtfId]="UserIncomeNtf"
}

func GetMsgProtoType(key uint16) reflect.Type {
	return msgProtoTypes[key]
}

func GetMsgName(key uint16) string {
	return msgNames[key]
}

const (
	MsgID_ZeroMsgId=0
	MsgID_MessageDataId=5001
	MsgID_Register_ReqId=5003
	MsgID_Register_AckId=5004
	MsgID_Send2UserId=5005
	MsgID_BroadcastId=5006
	MsgID_UserOnline_NtfId=5007
	MsgID_UserOff_NtfId=5008
	MsgID_UserIncome_NtfId=5009
)

func GetMsgIdFromType(i interface{}) uint16 {
	switch i.(type) {
	case *ZeroMsg:
		return 0
	case *MessageData:
		return 5001
	case *RegisterReq:
		return 5003
	case *RegisterAck:
		return 5004
	case *Send2User:
		return 5005
	case *Broadcast:
		return 5006
	case *UserOnlineNtf:
		return 5007
	case *UserOffNtf:
		return 5008
	case *UserIncomeNtf:
		return 5009
	default:
		return 0
	}
}
