package handle

import "github.com/v587-zyf/gc/gcnet/http_server"

type HandleOption struct {
	SID int64

	Tg_Login_token   string
	Tg_client_url    string
	Tg_start_photo   string
	Tg_start_caption string

	HttpServer *http_server.HttpServer
}

type Option func(opts *HandleOption)

func NewHandleOption() *HandleOption {
	o := &HandleOption{}

	return o
}
func WithSID(SID int64) Option {
	return func(opts *HandleOption) {
		opts.SID = SID
	}
}

func WithTgLoginToken(Tg_Login_token string) Option {
	return func(opts *HandleOption) {
		opts.Tg_Login_token = Tg_Login_token
	}
}

func WithTgClientUrl(Tg_client_url string) Option {
	return func(opts *HandleOption) {
		opts.Tg_client_url = Tg_client_url
	}
}
func WithTgStartPhoto(Tg_start_photo string) Option {
	return func(opts *HandleOption) {
		opts.Tg_start_photo = Tg_start_photo
	}
}
func WithTgStartCaption(Tg_start_caption string) Option {
	return func(opts *HandleOption) {
		opts.Tg_start_caption = Tg_start_caption
	}
}

func WithTgHttpServer(HttpServer *http_server.HttpServer) Option {
	return func(opts *HandleOption) {
		opts.HttpServer = HttpServer
	}
}
