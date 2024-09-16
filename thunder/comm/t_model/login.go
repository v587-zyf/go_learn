package model

import pb "comm/t_proto/out/client"

type RegisterReq struct {
	Account  string `json:"account" validate:"account"`         // 账号
	Password string `json:"password" validate:"password"`       // 密码
	Invite   uint64 `json:"invite,omitempty" validate:"invite"` // 邀请人id
}

type LoginReq struct {
	Types    pb.LoginType `json:"type" validate:"type"`                     // 登陆类型 枚举
	InitData string       `json:"init_data,omitempty" validate:"init_data"` // tg数据
	Invite   uint64       `json:"invite,omitempty" validate:"invite"`       // 邀请人id

	Account  string `json:"account,omitempty" validate:"account"`   // 账号
	Password string `json:"password,omitempty" validate:"password"` // 密码
}
type LoginAck struct {
	UserId   uint64 `json:"user_id" validate:"user_id"`     // 用户id
	Token    string `json:"token" validate:"token"`         // token 进入或重连游戏校验
	LinkAddr string `json:"link_addr" validate:"link_addr"` // 连接游戏地址
}

type TgUser struct {
	UserID          int64  `url:"id" json:"id"`
	FirstName       string `url:"first_name" json:"first_name"`
	LastName        string `url:"last_name" json:"last_name"`
	UserName        string `url:"user_name" json:"user_name"`
	LanguageCode    string `url:"language_code" json:"language_code"`
	AllowsWriteToPM bool   `url:"allows_write_to_pm" json:"allows_write_to_pm"`
}
type TgData struct {
	QueryID  string `url:"query_id" json:"query_id"`
	TgUser   TgUser `url:"user" json:"user"`
	AuthDate string `url:"auth_date" json:"auth_date"`
	Hash     string `url:"hash" json:"hash"`
}
