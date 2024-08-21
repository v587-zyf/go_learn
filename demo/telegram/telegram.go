package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"kernel/log"
)

func TelegramDo() {
	bot, err := tgbotapi.NewBotAPI("7492635198:AAFTNaXiRCpERqAQBaXOC3q2hRHAoqrUvpQ")
	if err != nil {
		log.Error("Failed to connect to Telegram", zap.Error(err))
		return
	}

	bot.Debug = true

	log.Debug("Authorized on account", zap.String("name", bot.Self.UserName))

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Debug("", zap.String("name", update.Message.From.UserName),
				zap.String("text", update.Message.Text), zap.Reflect("update", update))

			if update.Message.NewChatMembers != nil {
				// 新成员加入事件
				for _, newMember := range update.Message.NewChatMembers {
					log.Info("New member joined", zap.String("name", newMember.UserName))
				}
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}
	}

	//chatCfg := tgbotapi.ChatInfoConfig{}
	//chat, err := bot.GetChat(chatCfg)
	//log.Debug("chat", zap.Reflect("chat", chat), zap.Error(err))

	//memCntCfg := tgbotapi.ChatMemberCountConfig{
	//	ChatConfig: tgbotapi.ChatConfig{
	//		ChatID: -4252592771,
	//	},
	//}
	//count, err := bot.GetChatMembersCount(memCntCfg)
	//if err != nil {
	//	log.Error("get chat members count", zap.Error(err))
	//	return
	//}
	//log.Debug("count", zap.Int("count", count))
	//
	//chatCfg := tgbotapi.ChatInfoConfig{
	//	ChatConfig: tgbotapi.ChatConfig{
	//		ChatID: -4252592771,
	//	},
	//}
	//chat, err := bot.GetChat(chatCfg)
	//if err != nil {
	//	log.Error("get chat", zap.Error(err))
	//	return
	//}
	//log.Debug("chat", zap.Reflect("chat", chat))
}
