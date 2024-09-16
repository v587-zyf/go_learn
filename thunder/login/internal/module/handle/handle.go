package handle

import (
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/v587-zyf/gc/gcnet/http_server/middleware"
	"time"
)

type Handle struct {
	options *HandleOption
}

var h *Handle

func init() {
	h = &Handle{
		options: NewHandleOption(),
	}
}

func GetHandle() *Handle {
	return h
}

func GetHandleOps() *HandleOption {
	return h.options
}

func HandleInit(opts ...Option) error {
	for _, opt := range opts {
		opt(h.options)
	}

	h.options.HttpServer.Use(middleware.NewErrHandler())
	h.options.HttpServer.UseOrigin(limiter.New(limiter.Config{
		Expiration: 10 * time.Second,
		Max:        20,
	}))

	h.options.HttpServer.Post("/login", Login)
	h.options.HttpServer.Post("/register", Register)
	return nil
}
