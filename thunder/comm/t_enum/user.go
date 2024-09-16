package enum

import pb "comm/t_proto/out/client"

var (
	AccountT = map[pb.LoginType]string{
		pb.LoginType_password: ACCOUNT_TYPE_PASSWORD,
		pb.LoginType_telegram: ACCOUNT_TYPE_TELEGRAM,
	}
)

const (
	ACCOUNT_TYPE_PASSWORD = "password"
	ACCOUNT_TYPE_TELEGRAM = "telegram"
)

var USER_FIELD_NO = map[int]string{
	0: USER_FIELD_ID,
	1: USER_FIELD_BASIC,
	2: USER_FIELD_TELEGRAM,
}

const (
	USER_FIELD_ID       = "ID"
	USER_FIELD_BASIC    = "Basic"
	USER_FIELD_TELEGRAM = "Telegram"
)
