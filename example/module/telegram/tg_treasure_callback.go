package telegram

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
)

func TgTreasureCallBack(bot *gotgbot.Bot, ctx *ext.Context) (err error) {
	//log.Debug("1-----", zap.Any("msg", ctx.EffectiveMessage),
	//	zap.Any("callBackQuery", ctx.CallbackQuery), zap.Any("data", ctx.CallbackQuery.Data))
	log.Debug("1-----", zap.Any("chatId", ctx.EffectiveChat.Id), zap.Any("msgId", ctx.EffectiveMessage.MessageId),
		zap.Any("callBackData", ctx.CallbackQuery))
	//log.Debug("1-----", zap.Any("chatId", ctx.EffectiveChat.Id), zap.Any("msgId", ctx.EffectiveMessage.MessageId))

	//txt := "alert"
	//bot.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: txt, ShowAlert: true})
	//bot.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: txt, ShowAlert: false})
	msg := ctx.CallbackQuery.Message.(gotgbot.Message)

	chatId := ctx.CallbackQuery.Message.GetChat().Id
	//msgId := ctx.CallbackQuery.Message.GetMessageId() - 1
	photo := "https://candydream.mokoko.vip/source/3d19ea50aa4af0328cfee5c3d7b7e40.jpg"
	//msgIdStr := strings.Split(ctx.CallbackQuery.Data, "#")[2]
	//msgId, err := strconv.ParseInt(msgIdStr, 10, 64)
	//if err != nil {
	//	log.Error("msgId turn err", zap.Error(err))
	//	return
	//}
	caption, keyboard := makeTreasureMsg(ctx.CallbackQuery.Data)

	if _, err = bot.DeleteMessage(chatId, int64(msg.MessageId), nil); err != nil {
		log.Error("Failed to delete message", zap.Error(err))
		return
	}
	_, err = bot.SendPhoto(chatId, photo, &gotgbot.SendPhotoOpts{
		Caption:     caption + "this is new msg",
		ParseMode:   "HTML",
		ReplyMarkup: keyboard,
		//ReplyParameters: &gotgbot.ReplyParameters{
		//	MessageId: ctx.Message.MessageId,
		//},
	})
	//log.Debug("---", zap.Error(err))
	//_, _, err = bot.EditMessageCaption(&gotgbot.EditMessageCaptionOpts{
	//	ChatId:      chatId,
	//	MessageId:   int64(msg.MessageId),
	//	Caption:     caption,
	//	ParseMode:   "HTML",
	//	ReplyMarkup: keyboard,
	//})
	if err != nil {
		log.Debug("---", zap.Error(err))
	}

	return
}
