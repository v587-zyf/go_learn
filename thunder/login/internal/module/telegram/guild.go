package telegram

import (
	"comm/t_data/db"
	"comm/t_data/db/db_guild"
	"comm/t_data/db/db_user"
	"comm/t_data/redis"
	enum "comm/t_enum"
	pb "comm/t_proto/out/client"
	"errors"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/v587-zyf/gc/errcode"
	"github.com/v587-zyf/gc/log"
	"github.com/v587-zyf/gc/rdb/rdb_cluster"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"login/internal/module/handle"
	"strings"
	"time"
)

func guild(b *gotgbot.Bot, ctx *ext.Context) (err error) {

	var (
		replyTxt   string
		guildName  string
		newGuildId uint64

		channelInfo *db.AccountChannelInfo
		dbUser      *db_user.User
		dbGuild     *db_guild.Guild
		rdbUser     *redis.User
		rdbGuild    *redis.Guild

		gLock *rdb_cluster.Locker
		uLock *rdb_cluster.Locker

		chatId = ctx.EffectiveChat.Id
		tgUId  = ctx.EffectiveUser.Id
		txt    = ctx.EffectiveMessage.Text
	)

	switch ctx.EffectiveChat.Type {
	case "private":
		// mongo user
		channelInfo = &db.AccountChannelInfo{
			Channel:     enum.AccountT[pb.LoginType_telegram],
			AccountInfo: &db.AccTelegram{UserID: tgUId},
		}
		dbUser, err = db_user.GetUserModel().GetUserByChannelInfo(channelInfo)
		if err != nil && !errors.Is(err, mongo.ErrNoDocuments) || dbUser.ID == 0 {
			replyTxt = fmt.Sprintf("Please Login First!")
			break
		}

		// search guild && check user joined guild
		guildName = strings.Replace(txt, "@", "", 1)
		dbGuild, err = db_guild.GetGuildModel().GetGuildByChatId(chatId)
		if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
			replyTxt = fmt.Sprintf("Guild 【%s】 Not Found!", guildName)
			break
		}
		if dbGuild.Members == nil {
			dbGuild.Members = db.NewGuildMembers()
		}

		if dbGuild.Members.Has(dbUser.ID) || dbUser.Basic.GuildID != 0 {
			replyTxt = fmt.Sprintf("Join Guild Already!")
			break
		}

		// add user to mongo guild
		dbGuild.Members.Add(dbUser.ID)
		if _, err = db_guild.GetGuildModel().Upsert(dbGuild); err != nil {
			log.Error("add guild err", zap.Error(err), zap.Uint64("userID", dbUser.ID), zap.String("guildName", guildName))
			replyTxt = fmt.Sprintf("Add Guild 【%s】 Faild! Try Again...", guildName)
			break
		}

		// add user to redis guild
		{
			gLock = redis.LockGuild(dbGuild.ID, handle.GetHandleOps().SID)
			if gLock == nil {
				log.Error("lock guild err")
				replyTxt = fmt.Sprintf("Add Guild 【%s】 Faild! Try Again...", guildName)
				break
			}
			var unLockFn = func() { gLock.Unlock() }

			rdbGuild, err = redis.GetGuild(dbGuild.ID, gLock)
			if err != nil {
				unLockFn()
				log.Error("get redis guild err", zap.Error(err))
				replyTxt = fmt.Sprintf("Add Guild 【%s】 Faild! Try Again...", guildName)
				break
			}
			rdbGuild.AddMember(dbUser.ID)
			if err = redis.SetGuild(rdbGuild, gLock); err != nil {
				unLockFn()
				log.Error("update redis guild err", zap.Error(err))
				replyTxt = fmt.Sprintf("Add Guild 【%s】 Faild! Try Again...", guildName)
				break
			}

			unLockFn()
		}

		// update redis user guild info
		{
			uLock = redis.LockUser(dbUser.ID, handle.GetHandleOps().SID)
			if uLock == nil {
				err = errcode.ERR_USER_DATA_INVALID
				return
			}
			uLock.Unlock()

			rdbUser, err = redis.GetUser(dbUser.ID, uLock)
			if err != nil {
				uLock.Unlock()
				replyTxt = fmt.Sprintf("Not Fount Your Game Data! Try Again...")
				break
			}
			rdbUser.Basic.GuildID = dbGuild.ID
			if err = redis.SetUser(rdbUser, uLock); err != nil {
				log.Error("save user fail", zap.Uint64("userID", dbUser.ID), zap.Error(err))
				replyTxt = fmt.Sprintf("Add Guild 【%s】 Faild! Try Again...", guildName)
				break
			}
			uLock.Unlock()
		}

		// send msg back
		replyTxt = fmt.Sprintf("Add Guild 【%s】 Succeed!", guildName)
	default:
		guildName = strings.Replace(txt, "@", "", 1)
		if guildName == ctx.EffectiveChat.Title {
			// mongo user
			channelInfo = &db.AccountChannelInfo{
				Channel:     enum.AccountT[pb.LoginType_telegram],
				AccountInfo: &db.AccTelegram{UserID: tgUId},
			}
			dbUser, err = db_user.GetUserModel().GetUserByChannelInfo(channelInfo)
			if err != nil && !errors.Is(err, mongo.ErrNoDocuments) || dbUser.ID == 0 {
				replyTxt = fmt.Sprintf("Please Login First!")
				break
			}
			if dbUser.Basic.GuildID != 0 {
				log.Error("user already add guild", zap.Uint64("userID", dbUser.Basic.GuildID))
				replyTxt = fmt.Sprintf("Create Guild 【%s】 Faild! Already Join Guild. PLease Change Account...", guildName)
				break
			}

			// check guild is existed
			dbGuild, err = db_guild.GetGuildModel().GetGuildByChatId(chatId)
			if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
				log.Error("get guild id err", zap.Error(err), zap.Int64("chatID", chatId))
				replyTxt = fmt.Sprintf("Create Guild 【%s】 Faild! Please Try Again...", guildName)
				break
			}
			if dbGuild.ID != 0 {
				replyTxt = fmt.Sprintf("Create Guild 【%s】 Faild! Already Create...", guildName)
				break
			}

			// create mongo guild
			newGuildId, err = db.GenGuildIdSeq()
			if err != nil {
				log.Error("create guild id err", zap.Error(err))
				replyTxt = fmt.Sprintf("Create Guild 【%s】 Faild! Please Try Again...", guildName)
				break
			}
			guildMembers := db.NewGuildMembers()
			guildMembers.Add(dbUser.ID)
			newDbGuild := &db_guild.Guild{
				ID:        newGuildId,
				CreateUID: dbUser.ID,
				ChatID:    chatId,
				Title:     guildName,
				Members:   guildMembers,
				CreateAt:  time.Now(),
			}
			if _, err = db_guild.GetGuildModel().Upsert(newDbGuild); err != nil {
				log.Error("add guild err", zap.Error(err), zap.Uint64("userID", newDbGuild.ID), zap.String("guildName", guildName))
				replyTxt = fmt.Sprintf("Create Guild 【%s】 Faild! Try Again...", guildName)
				break
			}

			// create redis guild
			{
				gLock = redis.LockGuild(newGuildId, handle.GetHandleOps().SID)
				if gLock == nil {
					log.Error("lock guild err")
					replyTxt = fmt.Sprintf("Create Guild 【%s】 Faild! Try Again...", guildName)
					break
				}
				newRdbGuild := &redis.Guild{
					ID:        newGuildId,
					CreateUID: dbUser.ID,
					ChatID:    chatId,
					Title:     guildName,
					Members:   guildMembers,
				}
				if err = redis.SetGuild(newRdbGuild, gLock); err != nil {
					gLock.Unlock()
					log.Error("redis create err", zap.Error(err))
					replyTxt = fmt.Sprintf("Create Guild 【%s】 Faild! Try Again...", guildName)
					break
				}
				gLock.Unlock()
			}

			// update redis user guild info
			{
				uLock = redis.LockUser(dbUser.ID, handle.GetHandleOps().SID)
				if uLock == nil {
					err = errcode.ERR_USER_DATA_INVALID
					return
				}
				var unLockFn = func() { uLock.Unlock() }

				rdbUser, err = redis.GetUser(dbUser.ID, uLock)
				if err != nil {
					unLockFn()
					replyTxt = fmt.Sprintf("Not Fount Your Game Data! Try Again...")
					break
				}
				rdbUser.Basic.GuildID = newGuildId
				if err = redis.SetUser(rdbUser, uLock); err != nil {
					unLockFn()
					log.Error("save user fail", zap.Uint64("userID", dbUser.ID), zap.Error(err))
					replyTxt = fmt.Sprintf("Create Guild 【%s】 Faild! Try Again...", guildName)
					break
				}
				unLockFn()
			}

			// send msg back
			replyTxt = fmt.Sprintf("Create Guild 【%s】 Succeed!", guildName)
		}
	}

	if replyTxt != "" {
		if _, err = b.SendMessage(chatId, replyTxt, &gotgbot.SendMessageOpts{ReplyParameters: &gotgbot.ReplyParameters{MessageId: ctx.EffectiveMessage.MessageId}}); err != nil {
			log.Error("send msg err", zap.Error(err))
		}
	}

	return nil
}
