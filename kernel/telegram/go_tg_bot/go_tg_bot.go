package go_tg_bot

import (
	"context"
	"core/log"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"go.uber.org/zap"
	"time"
)

type TgBot struct {
	options *TgBotOption

	ctx    context.Context
	cancel context.CancelFunc

	bot        *gotgbot.Bot
	updater    *ext.Updater
	dispatcher *ext.Dispatcher
}

func NewTgBot() *TgBot {
	t := &TgBot{
		options: NewOption(),
	}

	return t
}

func (t *TgBot) Init(ctx context.Context, option ...any) (err error) {
	t.ctx, t.cancel = context.WithCancel(ctx)

	for _, opt := range option {
		opt.(Option)(t.options)
	}

	t.bot, err = gotgbot.NewBot(t.options.token, nil)
	if err != nil {
		log.Error("Failed to connect to Telegram", zap.Error(err))
		return err
	}
	t.dispatcher = ext.NewDispatcher(&ext.DispatcherOpts{
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			log.Error("an error occurred while handling update", zap.Error(err))
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
	})
	t.updater = ext.NewUpdater(t.dispatcher, nil)

	return nil
}

func (t *TgBot) Get() *gotgbot.Bot {
	return t.bot
}

func (t *TgBot) GetCtx() context.Context {
	return t.ctx
}

func (t *TgBot) AddHandle(handler ext.Handler) {
	t.dispatcher.AddHandler(handler)
}

func (t *TgBot) Start() {
	err := t.updater.StartPolling(t.bot, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			Timeout: 9,
			RequestOpts: &gotgbot.RequestOpts{
				Timeout: time.Second * 10,
			},
		},
	})
	if err != nil {
		panic("failed to start polling: " + err.Error())
	}

	t.updater.Idle()
}
