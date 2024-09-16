package redis

import (
	"comm/t_enum"
	"github.com/v587-zyf/gc/rdb/rdb_cluster"
	"github.com/v587-zyf/gc/utils"
)

type LoginInfo []any

// String 数据转string
func (t *LoginInfo) String(index int32) (out string) {
	if index < 0 || index >= int32(len(*t)) {
		return ""
	}

	data, ok := (*t)[index].(string)
	if !ok {
		return ""
	}

	return data
}

// Int32 数据转int32
func (t *LoginInfo) Uin64(index int32) uint64 {
	if index < 0 || index >= int32(len(*t)) {
		return 0
	}
	data, ok := (*t)[index].(string)
	if !ok {
		return 0
	}
	return utils.StrToUInt64(data)
}

// Int32 数据转int32
func (t *LoginInfo) Int64(index int32) int64 {
	if index < 0 || index >= int32(len(*t)) {
		return 0
	}
	data, ok := (*t)[index].(string)
	if !ok {
		return 0
	}

	return utils.StrToInt64(data)
}

// 保留时间比较久的数据
func SetUserLoginInfo(key string, value map[string]any) error {
	rc := rdb_cluster.Get()
	rCtx := rdb_cluster.GetCtx()

	if err := rc.HMSet(rCtx, key, value).Err(); err != nil {
		return err
	}

	return rc.Expire(rCtx, key, Second(TD_OneWeekSecond)).Err()
}

func GetUserLoginInfo(loginKey string) ([]any, error) {
	rc := rdb_cluster.Get()
	rCtx := rdb_cluster.GetCtx()

	return rc.HMGet(rCtx, loginKey, enum.Login_Token, enum.Login_Gate, enum.Login_UID).Result()
}
