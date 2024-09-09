package telegram

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"go.uber.org/zap"
	"kernel/log"
	"strings"
)

func TgCmd(bot *gotgbot.Bot, ctx *ext.Context) (err error) {
	//log.Debug("1-----", zap.Any("msg", ctx.EffectiveMessage),
	//	zap.Any("sender", ctx.EffectiveSender), zap.Any("chat", ctx.EffectiveChat),
	//	zap.Any("user", ctx.EffectiveUser))
	log.Debug("1-----", zap.Any("chatId", ctx.EffectiveChat.Id), zap.Any("msgId", ctx.EffectiveMessage.MessageId))
	txtSlice := strings.Split(ctx.EffectiveMessage.Text, "#")
	if len(txtSlice) != 2 || len(txtSlice[1]) < 11 {
		return nil
	}
	//
	//// 拆出playerId
	//playerIdStr := txtSlice[1][8 : len(txtSlice[1])-2]
	//log.Debug("2----", zap.Any("txtSlice", txtSlice), zap.String("playerIdStr", playerIdStr))
	//
	//playerIdUint, err := strconv.ParseUint(playerIdStr, 10, 64)
	//if err != nil {
	//	log.Error("handleTreasure", zap.Error(err))
	//	return
	//}
	//log.Debug("3----", zap.Uint64("playerIdUint", playerIdUint))

	//answerCallBackQuery := &gotgbot.AnswerCallbackQueryOpts{
	//	Text:      "alert",
	//	ShowAlert: true,
	//}
	//ret, err := bot.AnswerCallbackQuery(ctx.CallbackQuery.Id, answerCallBackQuery)
	//log.Debug("4----", zap.Any("ret", ret), zap.Error(err))

	photo := "https://candydream.mokoko.vip/source/3d19ea50aa4af0328cfee5c3d7b7e40.jpg"
	caption, keyboard := makeTreasureMsg(ctx.EffectiveMessage.Text)
	bot.SendPhoto(ctx.EffectiveChat.Id, photo, &gotgbot.SendPhotoOpts{
		Caption:     caption,
		ParseMode:   "HTML",
		ReplyMarkup: keyboard,
		//ReplyParameters: &gotgbot.ReplyParameters{
		//	MessageId: ctx.Message.MessageId,
		//},
	})
	log.Debug("---", zap.Error(err))
	//sendMsgOpts := &gotgbot.SendMessageOpts{
	//	ParseMode: "HTML",
	//}
	//ret, err := bot.SendMessage(ctx.EffectiveChat.Id, caption, sendMsgOpts)
	//log.Debug("5----", zap.Any("ret", ret), zap.Error(err))

	return nil
}

func makeTreasureMsg(callBackData string) (caption string, keyboard gotgbot.InlineKeyboardMarkup) {
	assistBtn := gotgbot.InlineKeyboardButton{
		Text:         "Assist",
		CallbackData: callBackData,
	}
	keyboard = gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
			{assistBtn},
		},
	}

	caption = "<b>Discover hidden treasure chests！！！</b>\n" +
		"help me open the treasure chest\n" +
		"🎁 Lv.%s\n" +
		"🔑 Turn on progress：%s%\n" +
		"⏰ The chest will open on %s\n" +
		"\n" +
		"<b>Rewards</b>\n" +
		"🌟 Gold：99999\n" +
		"💎 Diamond：5\n"

	return
}
