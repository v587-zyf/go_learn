package db

import (
	enum "comm/t_enum"
	"go.mongodb.org/mongo-driver/bson"
)

type AccPass struct {
	Account  string `bson:"account"`
	Password string `bson:"password"`
}
type AccTelegram struct {
	UserID int64 `bson:"user_id"`
}
type AccountChannelInfo struct {
	Channel     string
	AccountInfo any
}

type Accounts struct {
	Password *AccPass     `bson:"password"`
	Telegram *AccTelegram `bson:"telegram"`
}

func NewTgAccounts(tgUID int64) *Accounts {
	return &Accounts{
		Telegram: &AccTelegram{
			UserID: tgUID,
		},
		Password: &AccPass{},
	}
}

func NewPassAccounts(account, password string) *Accounts {
	return &Accounts{
		Password: &AccPass{
			Account:  account,
			Password: password,
		},
		Telegram: &AccTelegram{},
	}
}

func (a *Accounts) IsPassTrue(password string) bool {
	return a.Password.Password == password
}

func MakeAccountFilter(accountInfo *AccountChannelInfo) bson.M {
	filter := bson.M{}
	switch accountInfo.Channel {
	case enum.ACCOUNT_TYPE_PASSWORD:
		accInfo := accountInfo.AccountInfo.(*AccPass)
		filter = bson.M{"accounts.password.account": accInfo.Account}
	case enum.ACCOUNT_TYPE_TELEGRAM:
		accInfo := accountInfo.AccountInfo.(*AccTelegram)
		filter = bson.M{"accounts.telegram.user_id": accInfo.UserID}
	}
	return filter
}
