package errcode

import (
	"fmt"
	"kernel/enums"
)

var (
	language    = enums.DEFAULT_LANGUAGE
	defaultErrs = make(ErrGroup)
)

type ErrCode int32
type ErrGroup map[ErrCode]map[enums.LANGUAGE]string

func SetLanguage(lang enums.LANGUAGE) {
	language = lang
}

func CreateErrCode(code int32, desc string, lang enums.LANGUAGE) ErrCode {
	errCode := ErrCode(code)
	_, ok := defaultErrs[errCode]
	if !ok {
		defaultErrs[errCode] = make(map[enums.LANGUAGE]string)
	}

	if _, ok = defaultErrs[errCode][lang]; ok {
		msg := fmt.Sprintf("duplicate create err code, code:%d msg:%s lang:%s", code, desc, lang)
		panic(msg)
	}
	defaultErrs[errCode][lang] = desc

	return errCode
}

func (code ErrCode) Error() string {
	if v, ok := defaultErrs[code][language]; !ok {
		return fmt.Sprintf("UNKNOW_ERR_CODE[%d]", code)
	} else {
		return v
	}
}

func (code ErrCode) Int() int {
	return int(code)
}
func (code ErrCode) Int32() int32 {
	return int32(code)
}
