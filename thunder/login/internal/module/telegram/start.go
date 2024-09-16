package telegram

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
	"login/internal/module/handle"
)

func start(b *gotgbot.Bot, ctx *ext.Context) (err error) {
	photo := handle.GetHandleOps().Tg_start_photo
	caption := handle.GetHandleOps().Tg_start_caption
	clientUrl := handle.GetHandleOps().Tg_client_url

	playButton := gotgbot.InlineKeyboardButton{
		Text: "🎮 Play Game",
		WebApp: &gotgbot.WebAppInfo{
			Url: clientUrl,
		},
	}

	sendPhotoOpts := &gotgbot.SendPhotoOpts{
		Caption: caption,
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{playButton},
			},
		},
	}
	if _, err = b.SendPhoto(ctx.EffectiveChat.Id, photo, sendPhotoOpts); err != nil {
		log.Error("tg bot send start err", zap.Error(err))
		return
	}

	return
}
