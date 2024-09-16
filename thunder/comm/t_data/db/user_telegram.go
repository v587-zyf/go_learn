package db

import model "comm/t_model"

type Telegram struct {
	TgID            int64  `bson:"tg_id" json:"tg_id,omitempty"`
	FirstName       string `bson:"first_name" json:"first_name,omitempty"`
	LastName        string `bson:"last_name" json:"last_name,omitempty"`
	UserName        string `bson:"user_name" json:"user_name,omitempty"`
	LanguageCode    string `bson:"language_code" json:"language_code,omitempty"`           // 语言
	AllowsWriteToPM bool   `bson:"allows_write_to_pm" json:"allows_write_to_pm,omitempty"` // 是否订阅机器人

	Head string `bson:"head" json:"head,omitempty"` // 头像
}

func NewTelegram(tgUser *model.TgUser) *Telegram {
	return &Telegram{
		TgID:            tgUser.UserID,
		FirstName:       tgUser.FirstName,
		LastName:        tgUser.LastName,
		UserName:        tgUser.UserName,
		LanguageCode:    tgUser.LanguageCode,
		AllowsWriteToPM: tgUser.AllowsWriteToPM,
	}
}
