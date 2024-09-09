package go_tg_bot

type TgBotOption struct {
	token string
}

type Option func(opts *TgBotOption)

func NewOption() *TgBotOption {
	o := &TgBotOption{}

	return o
}

func WithToken(token string) Option {
	return func(opts *TgBotOption) {
		opts.token = token
	}
}
