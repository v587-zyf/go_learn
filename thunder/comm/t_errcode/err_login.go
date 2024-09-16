package errCode

import (
	"github.com/v587-zyf/gc/enums"
	"github.com/v587-zyf/gc/errcode"
)

var (
	ERR_ACCOUNT_ALREADY_REGISTER = errcode.CreateErrCode(1301, errcode.NewCodeLang("账号已注册", enums.LANG_CN), errcode.NewCodeLang("The account is registered", enums.LANG_EN))
)
