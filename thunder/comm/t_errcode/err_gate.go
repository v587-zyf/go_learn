package errCode

import (
	"github.com/v587-zyf/gc/enums"
	"github.com/v587-zyf/gc/errcode"
)

var (
	ERR_KICK_OUT          = errcode.CreateErrCode(1204, errcode.NewCodeLang("账号在别处登录", enums.LANG_CN), errcode.NewCodeLang("The account is logged in elsewhere", enums.LANG_EN))
	ERR_RECONNECT_TIMEOUT = errcode.CreateErrCode(1205, errcode.NewCodeLang("重连超时数", enums.LANG_CN), errcode.NewCodeLang("The number of times out of reconnection", enums.LANG_EN))
	ERR_SERVER_GATE_NIL   = errcode.CreateErrCode(1206, errcode.NewCodeLang("gate server is nil", enums.LANG_CN), errcode.NewCodeLang("gate server is nil", enums.LANG_EN))
)
