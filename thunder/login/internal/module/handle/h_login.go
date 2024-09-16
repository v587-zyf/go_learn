package handle

import (
	"comm/t_data/db"
	"comm/t_data/db/db_user"
	"comm/t_data/redis"
	"comm/t_enum"
	"comm/t_model"
	pb "comm/t_proto/out/client"
	"encoding/json"
	"errors"
	"github.com/v587-zyf/gc/errcode"
	"github.com/v587-zyf/gc/gcnet/http_server"
	"github.com/v587-zyf/gc/log"
	"github.com/v587-zyf/gc/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

func Login(c *http_server.Ctx) (any, error) {
	var (
		err error

		loginReq = new(model.LoginReq)

		validateSlice []string
	)

	if err = c.BodyParser(loginReq); err != nil {
		return nil, errcode.ERR_JSON_UNMARSHAL_ERR
	}

	switch loginReq.Types {
	case pb.LoginType_password:
		validateSlice = []string{"account", "password"}
	case pb.LoginType_telegram:
		validateSlice = []string{"init_data"}
	default:
		return nil, errcode.ERR_PARAM
	}
	if res := utils.ValidateColumn(loginReq, validateSlice); res {
		return nil, errcode.ERR_PARAM
	}

	switch loginReq.Types {
	case pb.LoginType_password:
		return login_password(loginReq)
	case pb.LoginType_telegram:
		return login_telegram(loginReq)
	default:
		return nil, errcode.ERR_PARAM
	}
}

func login_password(loginReq *model.LoginReq) (any, error) {
	var (
		err error
	)

	dbUser := db_user.GetUserModel()
	channelInfo := &db.AccountChannelInfo{
		Channel:     enum.AccountT[loginReq.Types],
		AccountInfo: &db.AccPass{Account: loginReq.Account, Password: loginReq.Password},
	}

	var userInfo *db_user.User
	userInfo, err = dbUser.GetUserByChannelInfo(channelInfo)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		log.Error("get account info err", zap.Error(err), zap.Reflect("channelInfo", channelInfo))
		return nil, errcode.ERR_MONGO_FIND
	}
	if userInfo.ID == 0 {
		return nil, errcode.ERR_PARAM
	}
	if !userInfo.Accounts.IsPassTrue(loginReq.Password) {
		return nil, errcode.ERR_SIGN
	}

	token := utils.Token()
	loginKey := redis.FormatUserLogin(userInfo.ID)
	gateAddr, gateId, errNo := redis.GetGate()
	if !errors.Is(errNo, errcode.ERR_SUCCEED) {
		return nil, errNo
	}

	loginMap := make(map[string]any)
	loginMap[enum.Login_Token] = token
	loginMap[enum.Login_Gate] = gateId
	loginMap[enum.Login_UID] = userInfo.ID
	if err = redis.SetUserLoginInfo(loginKey, loginMap); err != nil {
		log.Error("redis set err", zap.Error(err), zap.String("loginKey", loginKey), zap.Any("loginMap", loginMap))
		return nil, errcode.ERR_REDIS_UPDATE_USER
	}

	ack := &model.LoginAck{
		UserId:   userInfo.ID,
		Token:    token,
		LinkAddr: gateAddr,
	}

	return ack, nil
}

func login_telegram(loginReq *model.LoginReq) (any, error) {
	var (
		uid uint64
		err error
	)

	// 解析initData
	tgData, res := utils.TgCheck(loginReq.InitData, GetHandleOps().Tg_Login_token)
	if !res {
		return nil, errcode.ERR_PARAM
	}

	tgUser := &model.TgUser{}
	if err = json.Unmarshal([]byte(tgData.Get("user")), &tgUser); err != nil {
		log.Error("tgUser info err", zap.Error(err), zap.String("tgUser", tgData.Get("user")))
		return nil, errcode.ERR_PARAM
	}

	dbUser := db_user.GetUserModel()
	channelInfo := &db.AccountChannelInfo{
		Channel:     enum.AccountT[loginReq.Types],
		AccountInfo: &db.AccTelegram{UserID: tgUser.UserID},
	}

	var userInfo *db_user.User
	userInfo, err = dbUser.GetUserByChannelInfo(channelInfo)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		log.Error("get account info err", zap.Error(err), zap.Reflect("channelInfo", channelInfo))
		return nil, errcode.ERR_MONGO_FIND
	}
	if userInfo.ID == 0 {
		uid, err = db.GenUserIdSeq()
		if err != nil {
			log.Error("gen uid err", zap.Error(err))
			return nil, err
		}

		telegram := db.NewTelegram(tgUser)

		//tgPhotos, err := go_tg_bot.Get().GetUserProfilePhotos(tgUser.UserID, &gotgbot.GetUserProfilePhotosOpts{Limit: 1})
		//if err != nil {
		//	log.Error("tg get head err", zap.Error(err))
		//	return nil, err
		//}
		//if len(tgPhotos.Photos) > 0 && len(tgPhotos.Photos[0]) > 0 {
		//	headFile, err := go_tg_bot.Get().GetFile(tgPhotos.Photos[0][len(tgPhotos.Photos[0])-1].FileId, nil)
		//	if err == nil {
		//		telegram.Head = fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", go_tg_bot.Get().Token, headFile.FilePath)
		//	}
		//}

		u := &db_user.User{
			ID:       uid,
			Basic:    db.NewBasic(),
			Accounts: db.NewTgAccounts(tgUser.UserID),
			Telegram: telegram,
			Card:     db.NewCard(),
			Invite:   db.NewInvite(loginReq.Invite),
		}

		if err = dbUser.NewUserUnique(channelInfo, u); err != nil {
			log.Error("accountRegisterUnique err", zap.Error(err), zap.Reflect("channelInfo", channelInfo))
			return nil, err
		}

		userInfo, err = dbUser.GetUserById(uid)
		if err != nil {
			log.Error("get user info err", zap.Error(err), zap.Uint64("userID", uid))
			return nil, err
		}
	}

	token := utils.Token()
	loginKey := redis.FormatUserLogin(userInfo.ID)
	gateAddr, gateId, errNo := redis.GetGate()
	if !errors.Is(errNo, errcode.ERR_SUCCEED) {
		return nil, errNo
	}

	loginMap := make(map[string]any)
	loginMap[enum.Login_Token] = token
	loginMap[enum.Login_Gate] = gateId
	loginMap[enum.Login_UID] = userInfo.ID
	if err = redis.SetUserLoginInfo(loginKey, loginMap); err != nil {
		log.Error("redis set err", zap.Error(err), zap.String("loginKey", loginKey), zap.Any("loginMap", loginMap))
		return nil, errcode.ERR_REDIS_UPDATE_USER
	}

	ack := &model.LoginAck{
		UserId:   userInfo.ID,
		Token:    token,
		LinkAddr: gateAddr,
	}

	return ack, nil
}
