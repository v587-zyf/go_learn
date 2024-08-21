package tg_bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"kernel/log"
)

type TgBot struct {
	options *TgBotOption

	ctx    context.Context
	cancel context.CancelFunc

	Bot *tgbotapi.BotAPI
}

func NewTgBot() *TgBot {
	t := &TgBot{
		options: NewGrpcOption(),
	}

	return t
}

func (t *TgBot) Init(ctx context.Context, option ...any) (err error) {
	t.ctx, t.cancel = context.WithCancel(ctx)

	for _, opt := range option {
		opt.(Option)(t.options)
	}

	t.Bot, err = tgbotapi.NewBotAPI(t.options.token)
	if err != nil {
		log.Error("Failed to connect to Telegram", zap.Error(err))
		return err
	}

	return nil
}

func (t *TgBot) Get() *tgbotapi.BotAPI {
	return t.Bot
}

func (t *TgBot) GetCtx() context.Context {
	return t.ctx
}

//func (t *Telegram) RunListen(tgBot *TgBot) {
//	u := tgbotapi.NewUpdate(0)
//	u.AllowedUpdates = []string{"new_chat_members"}
//	//tgBot.Bot.Debug = true
//
//	// soSoValue
//	chatMember := tgbotapi.GetChatMemberConfig{
//		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
//			ChatID: -1001765198458,
//			UserID: 6458281620,
//		},
//	}
//	member, err := tgBot.Bot.GetChatMember(chatMember)
//	log.Debug("soSoValue -------------", zap.Reflect("member", member), zap.Error(err))
//	// safePal
//	chatMember = tgbotapi.GetChatMemberConfig{
//		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
//			ChatID: -1001937794848,
//			UserID: 6458281620,
//		},
//	}
//	member, err = tgBot.Bot.GetChatMember(chatMember)
//	log.Debug("safePal -------------", zap.Reflect("member", member), zap.Error(err))
//
//	tgBot.Bot.GetUserProfilePhotos()
//}
