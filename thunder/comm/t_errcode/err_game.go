package errCode

import (
	"github.com/v587-zyf/gc/enums"
	"github.com/v587-zyf/gc/errcode"
)

var (
	ERR_MOVE                  = errcode.CreateErrCode(1101, errcode.NewCodeLang("移动错误", enums.LANG_CN), errcode.NewCodeLang("Move error", enums.LANG_EN))
	ERR_OPEN_WALL             = errcode.CreateErrCode(1102, errcode.NewCodeLang("砸墙错误", enums.LANG_CN), errcode.NewCodeLang("Wall smashing error", enums.LANG_EN))
	ERR_STRENGTH_NOT_ENOUGH   = errcode.CreateErrCode(1103, errcode.NewCodeLang("体力不足", enums.LANG_CN), errcode.NewCodeLang("Lack of physical strength", enums.LANG_EN))
	ERR_POINT                 = errcode.CreateErrCode(1104, errcode.NewCodeLang("目标错误", enums.LANG_CN), errcode.NewCodeLang("Wrong target", enums.LANG_EN))
	ERR_NO_DEAD               = errcode.CreateErrCode(1105, errcode.NewCodeLang("玩家未死亡", enums.LANG_CN), errcode.NewCodeLang("The player is not dead", enums.LANG_EN))
	ERR_GM_PARAM              = errcode.CreateErrCode(1106, errcode.NewCodeLang("GM参数错误", enums.LANG_CN), errcode.NewCodeLang("The GM parameter is incorrect", enums.LANG_EN))
	ERR_DEAD                  = errcode.CreateErrCode(1107, errcode.NewCodeLang("玩家已死亡", enums.LANG_CN), errcode.NewCodeLang("The player is dead", enums.LANG_EN))
	ERR_GM_CLOSE              = errcode.CreateErrCode(1108, errcode.NewCodeLang("GM未开启", enums.LANG_CN), errcode.NewCodeLang("GM is not turned on", enums.LANG_EN))
	ERR_GOLD_NOT_ENOUGH       = errcode.CreateErrCode(1109, errcode.NewCodeLang("金币不足", enums.LANG_CN), errcode.NewCodeLang("Insufficient gold", enums.LANG_EN))
	ERR_LV_MAX                = errcode.CreateErrCode(1110, errcode.NewCodeLang("已达最大等级", enums.LANG_CN), errcode.NewCodeLang("The maximum level has been reached", enums.LANG_EN))
	ERR_CONF_NIL              = errcode.CreateErrCode(1111, errcode.NewCodeLang("配置未找到", enums.LANG_CN), errcode.NewCodeLang("Configuration not found", enums.LANG_EN))
	ERR_DIAMOND_NOT_ENOUGH    = errcode.CreateErrCode(1112, errcode.NewCodeLang("钻石不足", enums.LANG_CN), errcode.NewCodeLang("Diamonds are insufficient", enums.LANG_EN))
	ERR_BUY_NUM_MAX           = errcode.CreateErrCode(1113, errcode.NewCodeLang("已达最大购买次数", enums.LANG_CN), errcode.NewCodeLang("The maximum number of purchases has been reached", enums.LANG_EN))
	ERR_HASTEN_CD             = errcode.CreateErrCode(1114, errcode.NewCodeLang("该加速类型CD中", enums.LANG_CN), errcode.NewCodeLang("The acceleration type CD", enums.LANG_EN))
	ERR_HASTEN_MAX            = errcode.CreateErrCode(1115, errcode.NewCodeLang("加速已达上限", enums.LANG_CN), errcode.NewCodeLang("The acceleration has reached its upper limit", enums.LANG_EN))
	ERR_REWARD_ALREADY        = errcode.CreateErrCode(1116, errcode.NewCodeLang("奖励已领取", enums.LANG_CN), errcode.NewCodeLang("The reward has been claimed", enums.LANG_EN))
	ERR_INVITE_NUM_NOT_ENOUGH = errcode.CreateErrCode(1117, errcode.NewCodeLang("邀请人数不足", enums.LANG_CN), errcode.NewCodeLang("There are not enough inviters", enums.LANG_EN))
	ERR_GUILD_DID_NOT_JOIN    = errcode.CreateErrCode(1118, errcode.NewCodeLang("未加入公会", enums.LANG_CN), errcode.NewCodeLang("Not a member of a guild", enums.LANG_EN))
	ERR_GUILD_DATA_ERR        = errcode.CreateErrCode(1119, errcode.NewCodeLang("公会数据错误", enums.LANG_CN), errcode.NewCodeLang("Guild data error", enums.LANG_EN))
	ERR_GUILD_UPDATE_ERR      = errcode.CreateErrCode(1120, errcode.NewCodeLang("公会数据更新错误", enums.LANG_CN), errcode.NewCodeLang("Guild data update error", enums.LANG_EN))
	ERR_GUILD_ALREADY_JOIN    = errcode.CreateErrCode(1121, errcode.NewCodeLang("已加入公会", enums.LANG_CN), errcode.NewCodeLang("Joined a guild", enums.LANG_EN))
)
