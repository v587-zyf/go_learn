package pb
import (
	"reflect"
)

var msgProtoTypes = make(map[uint16]reflect.Type)
var msgNames = make(map[uint16]string)

func init() {
	msgProtoTypes[MsgID_HeartbeatId]=reflect.TypeOf((*Heartbeat)(nil)).Elem()
	msgProtoTypes[MsgID_Kick_NtfId]=reflect.TypeOf((*KickNtf)(nil)).Elem()
	msgProtoTypes[MsgID_Err_NtfId]=reflect.TypeOf((*ErrNtf)(nil)).Elem()
	msgProtoTypes[MsgID_Gm_ReqId]=reflect.TypeOf((*GmReq)(nil)).Elem()
	msgProtoTypes[MsgID_Gm_AckId]=reflect.TypeOf((*GmAck)(nil)).Elem()
	msgProtoTypes[MsgID_TestMsgId]=reflect.TypeOf((*TestMsg)(nil)).Elem()
	msgProtoTypes[MsgID_Login_ReqId]=reflect.TypeOf((*LoginReq)(nil)).Elem()
	msgProtoTypes[MsgID_Login_AckId]=reflect.TypeOf((*LoginAck)(nil)).Elem()
	msgProtoTypes[MsgID_Verify_ReqId]=reflect.TypeOf((*VerifyReq)(nil)).Elem()
	msgProtoTypes[MsgID_Verify_AckId]=reflect.TypeOf((*VerifyAck)(nil)).Elem()
	msgProtoTypes[MsgID_Reconnect_ReqId]=reflect.TypeOf((*ReconnectReq)(nil)).Elem()
	msgProtoTypes[MsgID_Reconnect_AckId]=reflect.TypeOf((*ReconnectAck)(nil)).Elem()
	msgProtoTypes[MsgID_Enter_ReqId]=reflect.TypeOf((*EnterReq)(nil)).Elem()
	msgProtoTypes[MsgID_Enter_NtfId]=reflect.TypeOf((*EnterNtf)(nil)).Elem()
	msgProtoTypes[MsgID_Move_ReqId]=reflect.TypeOf((*MoveReq)(nil)).Elem()
	msgProtoTypes[MsgID_Move_AckId]=reflect.TypeOf((*MoveAck)(nil)).Elem()
	msgProtoTypes[MsgID_OpenWall_ReqId]=reflect.TypeOf((*OpenWallReq)(nil)).Elem()
	msgProtoTypes[MsgID_OpenWall_AckId]=reflect.TypeOf((*OpenWallAck)(nil)).Elem()
	msgProtoTypes[MsgID_GetTreasure_ReqId]=reflect.TypeOf((*GetTreasureReq)(nil)).Elem()
	msgProtoTypes[MsgID_GetTreasure_AckId]=reflect.TypeOf((*GetTreasureAck)(nil)).Elem()
	msgProtoTypes[MsgID_OpenTreasure_ReqId]=reflect.TypeOf((*OpenTreasureReq)(nil)).Elem()
	msgProtoTypes[MsgID_OpenTreasure_AckId]=reflect.TypeOf((*OpenTreasureAck)(nil)).Elem()
	msgProtoTypes[MsgID_Revive_ReqId]=reflect.TypeOf((*ReviveReq)(nil)).Elem()
	msgProtoTypes[MsgID_Revive_AckId]=reflect.TypeOf((*ReviveAck)(nil)).Elem()
	msgProtoTypes[MsgID_ResetMap_ReqId]=reflect.TypeOf((*ResetMapReq)(nil)).Elem()
	msgProtoTypes[MsgID_ResetMap_AckId]=reflect.TypeOf((*ResetMapAck)(nil)).Elem()
	msgProtoTypes[MsgID_UpLv_ReqId]=reflect.TypeOf((*UpLvReq)(nil)).Elem()
	msgProtoTypes[MsgID_UpLv_AckId]=reflect.TypeOf((*UpLvAck)(nil)).Elem()
	msgProtoTypes[MsgID_Strength_NtfId]=reflect.TypeOf((*StrengthNtf)(nil)).Elem()
	msgProtoTypes[MsgID_Income_NtfId]=reflect.TypeOf((*IncomeNtf)(nil)).Elem()
	msgProtoTypes[MsgID_Diamond_NtfId]=reflect.TypeOf((*DiamondNtf)(nil)).Elem()
	msgProtoTypes[MsgID_Card_ReqId]=reflect.TypeOf((*CardReq)(nil)).Elem()
	msgProtoTypes[MsgID_Card_AckId]=reflect.TypeOf((*CardAck)(nil)).Elem()
	msgProtoTypes[MsgID_Card_NtfId]=reflect.TypeOf((*CardNtf)(nil)).Elem()
	msgProtoTypes[MsgID_CardUpLv_ReqId]=reflect.TypeOf((*CardUpLvReq)(nil)).Elem()
	msgProtoTypes[MsgID_CardUpLv_AckId]=reflect.TypeOf((*CardUpLvAck)(nil)).Elem()
	msgProtoTypes[MsgID_ShopBuy_ReqId]=reflect.TypeOf((*ShopBuyReq)(nil)).Elem()
	msgProtoTypes[MsgID_ShopBuy_AckId]=reflect.TypeOf((*ShopBuyAck)(nil)).Elem()
	msgProtoTypes[MsgID_Hasten_ReqId]=reflect.TypeOf((*HastenReq)(nil)).Elem()
	msgProtoTypes[MsgID_Hasten_AckId]=reflect.TypeOf((*HastenAck)(nil)).Elem()
	msgProtoTypes[MsgID_Invite_ReqId]=reflect.TypeOf((*InviteReq)(nil)).Elem()
	msgProtoTypes[MsgID_Invite_AckId]=reflect.TypeOf((*InviteAck)(nil)).Elem()
	msgProtoTypes[MsgID_InviteReward_ReqId]=reflect.TypeOf((*InviteRewardReq)(nil)).Elem()
	msgProtoTypes[MsgID_InviteReward_AckId]=reflect.TypeOf((*InviteRewardAck)(nil)).Elem()
	msgProtoTypes[MsgID_Invite_NtfId]=reflect.TypeOf((*InviteNtf)(nil)).Elem()
	msgProtoTypes[MsgID_RedPoint_ReqId]=reflect.TypeOf((*RedPointReq)(nil)).Elem()
	msgProtoTypes[MsgID_RedPoint_AckId]=reflect.TypeOf((*RedPointAck)(nil)).Elem()
	msgProtoTypes[MsgID_RedPoint_NtfId]=reflect.TypeOf((*RedPointNtf)(nil)).Elem()
	msgProtoTypes[MsgID_Rank_ReqId]=reflect.TypeOf((*RankReq)(nil)).Elem()
	msgProtoTypes[MsgID_Rank_AckId]=reflect.TypeOf((*RankAck)(nil)).Elem()
	msgProtoTypes[MsgID_GuildList_ReqId]=reflect.TypeOf((*GuildListReq)(nil)).Elem()
	msgProtoTypes[MsgID_GuildList_AckId]=reflect.TypeOf((*GuildListAck)(nil)).Elem()
	msgProtoTypes[MsgID_GuildRank_ReqId]=reflect.TypeOf((*GuildRankReq)(nil)).Elem()
	msgProtoTypes[MsgID_GuildRank_AckId]=reflect.TypeOf((*GuildRankAck)(nil)).Elem()
	msgProtoTypes[MsgID_GuildJoin_ReqId]=reflect.TypeOf((*GuildJoinReq)(nil)).Elem()
	msgProtoTypes[MsgID_GuildJoin_AckId]=reflect.TypeOf((*GuildJoinAck)(nil)).Elem()
	msgProtoTypes[MsgID_GuildLeave_ReqId]=reflect.TypeOf((*GuildLeaveReq)(nil)).Elem()
	msgProtoTypes[MsgID_GuildLeave_AckId]=reflect.TypeOf((*GuildLeaveAck)(nil)).Elem()
	msgNames[MsgID_HeartbeatId]="Heartbeat"
	msgNames[MsgID_Kick_NtfId]="KickNtf"
	msgNames[MsgID_Err_NtfId]="ErrNtf"
	msgNames[MsgID_Gm_ReqId]="GmReq"
	msgNames[MsgID_Gm_AckId]="GmAck"
	msgNames[MsgID_TestMsgId]="TestMsg"
	msgNames[MsgID_Login_ReqId]="LoginReq"
	msgNames[MsgID_Login_AckId]="LoginAck"
	msgNames[MsgID_Verify_ReqId]="VerifyReq"
	msgNames[MsgID_Verify_AckId]="VerifyAck"
	msgNames[MsgID_Reconnect_ReqId]="ReconnectReq"
	msgNames[MsgID_Reconnect_AckId]="ReconnectAck"
	msgNames[MsgID_Enter_ReqId]="EnterReq"
	msgNames[MsgID_Enter_NtfId]="EnterNtf"
	msgNames[MsgID_Move_ReqId]="MoveReq"
	msgNames[MsgID_Move_AckId]="MoveAck"
	msgNames[MsgID_OpenWall_ReqId]="OpenWallReq"
	msgNames[MsgID_OpenWall_AckId]="OpenWallAck"
	msgNames[MsgID_GetTreasure_ReqId]="GetTreasureReq"
	msgNames[MsgID_GetTreasure_AckId]="GetTreasureAck"
	msgNames[MsgID_OpenTreasure_ReqId]="OpenTreasureReq"
	msgNames[MsgID_OpenTreasure_AckId]="OpenTreasureAck"
	msgNames[MsgID_Revive_ReqId]="ReviveReq"
	msgNames[MsgID_Revive_AckId]="ReviveAck"
	msgNames[MsgID_ResetMap_ReqId]="ResetMapReq"
	msgNames[MsgID_ResetMap_AckId]="ResetMapAck"
	msgNames[MsgID_UpLv_ReqId]="UpLvReq"
	msgNames[MsgID_UpLv_AckId]="UpLvAck"
	msgNames[MsgID_Strength_NtfId]="StrengthNtf"
	msgNames[MsgID_Income_NtfId]="IncomeNtf"
	msgNames[MsgID_Diamond_NtfId]="DiamondNtf"
	msgNames[MsgID_Card_ReqId]="CardReq"
	msgNames[MsgID_Card_AckId]="CardAck"
	msgNames[MsgID_Card_NtfId]="CardNtf"
	msgNames[MsgID_CardUpLv_ReqId]="CardUpLvReq"
	msgNames[MsgID_CardUpLv_AckId]="CardUpLvAck"
	msgNames[MsgID_ShopBuy_ReqId]="ShopBuyReq"
	msgNames[MsgID_ShopBuy_AckId]="ShopBuyAck"
	msgNames[MsgID_Hasten_ReqId]="HastenReq"
	msgNames[MsgID_Hasten_AckId]="HastenAck"
	msgNames[MsgID_Invite_ReqId]="InviteReq"
	msgNames[MsgID_Invite_AckId]="InviteAck"
	msgNames[MsgID_InviteReward_ReqId]="InviteRewardReq"
	msgNames[MsgID_InviteReward_AckId]="InviteRewardAck"
	msgNames[MsgID_Invite_NtfId]="InviteNtf"
	msgNames[MsgID_RedPoint_ReqId]="RedPointReq"
	msgNames[MsgID_RedPoint_AckId]="RedPointAck"
	msgNames[MsgID_RedPoint_NtfId]="RedPointNtf"
	msgNames[MsgID_Rank_ReqId]="RankReq"
	msgNames[MsgID_Rank_AckId]="RankAck"
	msgNames[MsgID_GuildList_ReqId]="GuildListReq"
	msgNames[MsgID_GuildList_AckId]="GuildListAck"
	msgNames[MsgID_GuildRank_ReqId]="GuildRankReq"
	msgNames[MsgID_GuildRank_AckId]="GuildRankAck"
	msgNames[MsgID_GuildJoin_ReqId]="GuildJoinReq"
	msgNames[MsgID_GuildJoin_AckId]="GuildJoinAck"
	msgNames[MsgID_GuildLeave_ReqId]="GuildLeaveReq"
	msgNames[MsgID_GuildLeave_AckId]="GuildLeaveAck"
}

func GetMsgProtoType(key uint16) reflect.Type {
	return msgProtoTypes[key]
}

func GetMsgName(key uint16) string {
	return msgNames[key]
}

const (
	MsgID_HeartbeatId=0
	MsgID_Kick_NtfId=1
	MsgID_Err_NtfId=2
	MsgID_Gm_ReqId=3
	MsgID_Gm_AckId=4
	MsgID_TestMsgId=999
	MsgID_Login_ReqId=1001
	MsgID_Login_AckId=1002
	MsgID_Verify_ReqId=1003
	MsgID_Verify_AckId=1004
	MsgID_Reconnect_ReqId=1005
	MsgID_Reconnect_AckId=1006
	MsgID_Enter_ReqId=1007
	MsgID_Enter_NtfId=1008
	MsgID_Move_ReqId=1009
	MsgID_Move_AckId=1010
	MsgID_OpenWall_ReqId=1011
	MsgID_OpenWall_AckId=1012
	MsgID_GetTreasure_ReqId=1013
	MsgID_GetTreasure_AckId=1014
	MsgID_OpenTreasure_ReqId=1015
	MsgID_OpenTreasure_AckId=1016
	MsgID_Revive_ReqId=1017
	MsgID_Revive_AckId=1018
	MsgID_ResetMap_ReqId=1019
	MsgID_ResetMap_AckId=1020
	MsgID_UpLv_ReqId=1021
	MsgID_UpLv_AckId=1022
	MsgID_Strength_NtfId=1023
	MsgID_Income_NtfId=1024
	MsgID_Diamond_NtfId=1025
	MsgID_Card_ReqId=1030
	MsgID_Card_AckId=1031
	MsgID_Card_NtfId=1032
	MsgID_CardUpLv_ReqId=1033
	MsgID_CardUpLv_AckId=1034
	MsgID_ShopBuy_ReqId=1040
	MsgID_ShopBuy_AckId=1041
	MsgID_Hasten_ReqId=1051
	MsgID_Hasten_AckId=1052
	MsgID_Invite_ReqId=1061
	MsgID_Invite_AckId=1062
	MsgID_InviteReward_ReqId=1063
	MsgID_InviteReward_AckId=1064
	MsgID_Invite_NtfId=1065
	MsgID_RedPoint_ReqId=1071
	MsgID_RedPoint_AckId=1072
	MsgID_RedPoint_NtfId=1073
	MsgID_Rank_ReqId=1081
	MsgID_Rank_AckId=1082
	MsgID_GuildList_ReqId=1101
	MsgID_GuildList_AckId=1102
	MsgID_GuildRank_ReqId=1103
	MsgID_GuildRank_AckId=1104
	MsgID_GuildJoin_ReqId=1105
	MsgID_GuildJoin_AckId=1106
	MsgID_GuildLeave_ReqId=1107
	MsgID_GuildLeave_AckId=1108
)

func GetMsgIdFromType(i interface{}) uint16 {
	switch i.(type) {
	case *Heartbeat:
		return 0
	case *KickNtf:
		return 1
	case *ErrNtf:
		return 2
	case *GmReq:
		return 3
	case *GmAck:
		return 4
	case *TestMsg:
		return 999
	case *LoginReq:
		return 1001
	case *LoginAck:
		return 1002
	case *VerifyReq:
		return 1003
	case *VerifyAck:
		return 1004
	case *ReconnectReq:
		return 1005
	case *ReconnectAck:
		return 1006
	case *EnterReq:
		return 1007
	case *EnterNtf:
		return 1008
	case *MoveReq:
		return 1009
	case *MoveAck:
		return 1010
	case *OpenWallReq:
		return 1011
	case *OpenWallAck:
		return 1012
	case *GetTreasureReq:
		return 1013
	case *GetTreasureAck:
		return 1014
	case *OpenTreasureReq:
		return 1015
	case *OpenTreasureAck:
		return 1016
	case *ReviveReq:
		return 1017
	case *ReviveAck:
		return 1018
	case *ResetMapReq:
		return 1019
	case *ResetMapAck:
		return 1020
	case *UpLvReq:
		return 1021
	case *UpLvAck:
		return 1022
	case *StrengthNtf:
		return 1023
	case *IncomeNtf:
		return 1024
	case *DiamondNtf:
		return 1025
	case *CardReq:
		return 1030
	case *CardAck:
		return 1031
	case *CardNtf:
		return 1032
	case *CardUpLvReq:
		return 1033
	case *CardUpLvAck:
		return 1034
	case *ShopBuyReq:
		return 1040
	case *ShopBuyAck:
		return 1041
	case *HastenReq:
		return 1051
	case *HastenAck:
		return 1052
	case *InviteReq:
		return 1061
	case *InviteAck:
		return 1062
	case *InviteRewardReq:
		return 1063
	case *InviteRewardAck:
		return 1064
	case *InviteNtf:
		return 1065
	case *RedPointReq:
		return 1071
	case *RedPointAck:
		return 1072
	case *RedPointNtf:
		return 1073
	case *RankReq:
		return 1081
	case *RankAck:
		return 1082
	case *GuildListReq:
		return 1101
	case *GuildListAck:
		return 1102
	case *GuildRankReq:
		return 1103
	case *GuildRankAck:
		return 1104
	case *GuildJoinReq:
		return 1105
	case *GuildJoinAck:
		return 1106
	case *GuildLeaveReq:
		return 1107
	case *GuildLeaveAck:
		return 1108
	default:
		return 0
	}
}
