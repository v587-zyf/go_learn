package handle

import (
	"comm/t_data/db"
	"comm/t_data/db/db_user"
	enum "comm/t_enum"
	errCode "comm/t_errcode"
	model "comm/t_model"
	pb "comm/t_proto/out/client"
	"errors"
	"github.com/v587-zyf/gc/errcode"
	"github.com/v587-zyf/gc/gcnet/http_server"
	"github.com/v587-zyf/gc/log"
	"github.com/v587-zyf/gc/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

func Register(c *http_server.Ctx) (any, error) {
	var (
		uid uint64
		err error

		registerReq = new(model.RegisterReq)

		validateSlice = []string{"account", "password"}
	)

	if err = c.BodyParser(registerReq); err != nil {
		return nil, errcode.ERR_PARAM
	}
	if res := utils.ValidateColumn(registerReq, validateSlice); res {
		return nil, errcode.ERR_PARAM
	}

	dbUser := db_user.GetUserModel()
	channelInfo := &db.AccountChannelInfo{
		Channel:     enum.AccountT[pb.LoginType_password],
		AccountInfo: &db.AccPass{Account: registerReq.Account, Password: registerReq.Password},
	}

	var userInfo *db_user.User
	userInfo, err = dbUser.GetUserByChannelInfo(channelInfo)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		log.Error("get account info err", zap.Error(err), zap.Reflect("channelInfo", channelInfo))
		return nil, errcode.ERR_MONGO_FIND
	}
	if userInfo.ID != 0 {
		return nil, errCode.ERR_ACCOUNT_ALREADY_REGISTER
	}

	uid, err = db.GenUserIdSeq()
	if err != nil {
		log.Error("gen uid err", zap.Error(err))
		return nil, err
	}

	u := &db_user.User{
		ID:       uid,
		Basic:    db.NewBasic(),
		Accounts: db.NewPassAccounts(registerReq.Account, registerReq.Password),
		Telegram: new(db.Telegram),
		Card:     db.NewCard(),
		Invite:   db.NewInvite(registerReq.Invite),
	}

	if err = dbUser.NewUserUnique(channelInfo, u); err != nil {
		log.Error("accountRegisterUnique err", zap.Error(err), zap.Reflect("channelInfo", channelInfo))
		return nil, err
	}

	return nil, nil
}
