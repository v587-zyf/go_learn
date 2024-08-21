package errcode

import "kernel/enums"

var (
	ERR_SUCCEED           = CreateErrCode(1000, "成功", enums.LANG_CN)
	ERR_STANDARD_ERR      = CreateErrCode(1001, "失败", enums.LANG_CN)
	ERR_SIGN              = CreateErrCode(1002, "验证未通过", enums.LANG_CN)
	ERR_PARAM             = CreateErrCode(1003, "参数错误", enums.LANG_CN)
	ERR_RECONNECT_MAX     = CreateErrCode(1004, "重连已达最大次数", enums.LANG_CN)
	ERR_KICK_OUT          = CreateErrCode(1005, "被踢出", enums.LANG_CN)
	ERR_SQUEEZE           = CreateErrCode(1006, "该账号在别处登录", enums.LANG_CN)
	ERR_RECONNECT_TIMEOUT = CreateErrCode(1007, "重连超时数", enums.LANG_CN)

	// net
	ERR_NET_BODY_LEN_INVALID = CreateErrCode(11, "数据格式错误", enums.LANG_CN)
	ERR_NET_TCP_CLOSED       = CreateErrCode(12, "网络已经关闭", enums.LANG_CN)
	ERR_NET_TIMEOUT          = CreateErrCode(13, "网络数据超时", enums.LANG_CN)
	ERR_NET_SEND_TIMEOUT     = CreateErrCode(14, "发送数据超时", enums.LANG_CN)
	ERR_NET_RECV_TIMEOUT     = CreateErrCode(15, "接收数据超时", enums.LANG_CN)
	ERR_NET_PKG_LEN_LIMIT    = CreateErrCode(16, "数据包长度限制", enums.LANG_CN)
	ERR_NET_SESSION_EMPTY    = CreateErrCode(17, "没有建立会话连接", enums.LANG_CN)
	ERR_NET_CONNECT_FAILED   = CreateErrCode(18, "网络连接失败", enums.LANG_CN)

	// server
	ERR_SERVER_INTERNAL          = CreateErrCode(21, "服务器内部错误", enums.LANG_CN)
	ERR_SERVER_NOT_FOUND         = CreateErrCode(22, "服务器未找到", enums.LANG_CN)
	ERR_SERVER_REG_TIMEOUT       = CreateErrCode(23, "服务器注册超时", enums.LANG_CN)
	ERR_SERVER_REG_FAILED        = CreateErrCode(24, "服务器注册失败", enums.LANG_CN)
	ERR_SERVER_REG_PARAM_INVALID = CreateErrCode(25, "服务注册参数错误", enums.LANG_CN)
	ERR_SERVER_REG_DUPLICATE     = CreateErrCode(26, "服务重复注册", enums.LANG_CN)
	ERR_WP_TOO_MANY_WORKER       = CreateErrCode(27, "工作池任务太多", enums.LANG_CN)

	// event
	ERR_EVENT_LISTENER_LIMIT    = CreateErrCode(31, "监听事件达到上限", enums.LANG_CN)
	ERR_EVENT_LISTENER_EMPTY    = CreateErrCode(32, "监听事件队列为空", enums.LANG_CN)
	ERR_EVENT_LISTENER_NOT_FIND = CreateErrCode(33, "监听事件未找到", enums.LANG_CN)
	ERR_EVENT_PARAM_INVALID     = CreateErrCode(34, "输入参数无效", enums.LANG_CN)

	ERR_JSON_MARSHAL_ERR   = CreateErrCode(101, "json打包错误", enums.LANG_CN)
	ERR_JSON_UNMARSHAL_ERR = CreateErrCode(102, "json解包错误", enums.LANG_CN)

	ERR_USER_DATA_NOT_FOUND = CreateErrCode(201, "用户信息未找到", enums.LANG_CN)
	ERR_USER_DATA_INVALID   = CreateErrCode(202, "用户信息错误", enums.LANG_CN)

	// http
	ERR_HTTP_METHOD      = CreateErrCode(301, "请求类型错误", enums.LANG_CN)
	ERR_HTTP_PARAM       = CreateErrCode(302, "参数错误", enums.LANG_CN)
	ERR_HTTP_REQUEST_ERR = CreateErrCode(303, "http请求错误", enums.LANG_CN)

	// redis
	ERR_REDIS_UPDATE_USER    = CreateErrCode(401, "redis更新玩家数据错误", enums.LANG_CN)
	ERR_REDIS_DATA_NIL       = CreateErrCode(402, "redis数据为空", enums.LANG_CN)
	ERR_REDIS_LOGIN_DATA_NIL = CreateErrCode(403, "redis登陆数据数据为空", enums.LANG_CN)

	// server
	ERR_SERVER_GATE_NIL = CreateErrCode(501, "Gate服务器为空", enums.LANG_CN)

	// mongodb
	ERR_MONGO_UPSERT     = CreateErrCode(701, "upsert错误", enums.LANG_CN)
	ERR_MONGO_FIND       = CreateErrCode(702, "未找到数据", enums.LANG_CN)
	ERR_MONGO_INSERT_ONE = CreateErrCode(703, "插入一条数据错误", enums.LANG_CN)
	ERR_MONGO_UPDATE_ONE = CreateErrCode(704, "更新一条数据错误", enums.LANG_CN)
	ERR_MONGO_QUERY      = CreateErrCode(705, "数据库请求错误", enums.LANG_CN)
	ERR_MONGO_KEY_DUP    = CreateErrCode(706, "数据库主键重复", enums.LANG_CN)
	ERR_MONGO_DEL        = CreateErrCode(707, "删除数据失败", enums.LANG_CN)

	ERR_MQ_REPLY_HEAD_LEN      = CreateErrCode(3001, "mq回复头长度错误", enums.LANG_CN)
	ERR_MQ_BUFF_WRITE          = CreateErrCode(3002, "mq写入缓冲区错误", enums.LANG_CN)
	ERR_MQ_REPLY_EMPTY         = CreateErrCode(3003, "mq回复数据为空", enums.LANG_CN)
	ERR_MQ_REPLY_PB            = CreateErrCode(3004, "mq回复数据解析错误", enums.LANG_CN)
	ERR_MQ_RECV_DATA_UNMARSHAL = CreateErrCode(3005, "mq接收数据解析错误", enums.LANG_CN)
	ERR_MQ_MSG_ID_NOT_REGISTER = CreateErrCode(3006, "mq消息未注册", enums.LANG_CN)
	ERR_MQ_CONNECT_FAIL        = CreateErrCode(3007, "mq连接失败", enums.LANG_CN)
	ERR_MQ_REQ_TIMEOUT         = CreateErrCode(3008, "mq请求超时", enums.LANG_CN)
	ERR_MQ_SERVER_NOT_FOUND    = CreateErrCode(3009, "mq服务器未找到", enums.LANG_CN)
)

//func CreateErrCode(code int32, desc string) ErrCode {
//	errCode := ErrCode(code)
//	if _, ok := defaultErrs[errCode]; ok {
//		msg := fmt.Sprintf("duplicate create err code, code:%d msg:%s", code, desc)
//		panic(msg)
//	}
//
//	defaultErrs[errCode] = desc
//	return errCode
//}
//
//func (code ErrCode) Error() string {
//	if v, ok := defaultErrs[code]; !ok {
//		return fmt.Sprintf("未知代号[%d]", code)
//	} else {
//		return v
//	}
//}
//
//func (code ErrCode) Int() int {
//	return int(code)
//}
//func (code ErrCode) Int32() int32 {
//	return int32(code)
//}
//
//var defaultErrs map[ErrCode]string = map[ErrCode]string{}
