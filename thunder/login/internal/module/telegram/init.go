package telegram

import (
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
	"github.com/v587-zyf/gc/telegram/go_tg_bot"
)

func TelegramInit() {
	initCmd()

	go go_tg_bot.Start()
}

func initCmd() {
	go_tg_bot.AddHandle(handlers.NewCommand("start", start))
	go_tg_bot.AddHandle(handlers.NewCallback(callbackquery.Prefix("@"), guild))
}
