package telegram

import (
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
	"github.com/v587-zyf/gc/telegram/go_tg_bot"
)

const (
	GameConst_ShareTreasureCmd = "/treasure#"
)

func TgInit() {
	initCmd()

	go go_tg_bot.Start()
}
func initCmd() {
	go_tg_bot.AddHandle(handlers.NewMessage(message.HasPrefix(GameConst_ShareTreasureCmd), TgCmd))
	go_tg_bot.AddHandle(handlers.NewCallback(callbackquery.Prefix(GameConst_ShareTreasureCmd), TgTreasureCallBack))
}
