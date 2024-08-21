package http_server

import (
	"errors"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"kernel/errcode"
	"kernel/log"
	"runtime"

	"go.uber.org/zap"
)

type Ctx struct {
	*fiber.Ctx
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

type ResponseHandlerFn func(*Ctx) (any, error)

type OriginHandlerFn func(*Ctx) error

func SendErrCode(c *fiber.Ctx, errCode errcode.ErrCode) error {
	resp := Response{
		Code: errCode.Int(),
		Msg:  errCode.Error(),
		Data: nil,
	}

	out, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	// return c.SendString(string(out))

	c.Response().SetStatusCode(200)
	c.Response().SetBodyRaw(out)
	return nil
}

func SendError(c *fiber.Ctx, err error) error {
	resp := Response{
		Code: errcode.ERR_STANDARD_ERR.Int(),
		Msg:  err.Error(),
		Data: nil,
	}

	out, err := json.Marshal(resp)
	if err != nil {
		return SendErrCode(c, errcode.ERR_JSON_MARSHAL_ERR)
	}

	c.Response().SetStatusCode(200)
	c.Response().SetBodyRaw(out)
	return nil
}

func SendResponse(c *fiber.Ctx, data any) error {
	resp := Response{
		Code: errcode.ERR_SUCCEED.Int(),
		Msg:  "ok",
		Data: data,
	}

	out, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	c.Response().SetStatusCode(200)
	c.Response().SetBodyRaw(out)

	return nil
}

func NewResponseHandlerFn(fn ResponseHandlerFn) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) (retErr error) {
		defer func() {
			if r := recover(); r != nil {
				buf := make([]byte, 1<<10)
				runtime.Stack(buf, true)
				if err, ok := r.(error); ok {
					log.Error("handler core dump", zap.String("clientIP", c.IP()),
						zap.String("err", err.Error()), zap.ByteString("core", buf))
				} else if err, ok := r.(string); ok {
					log.Error("handler core dump", zap.String("clientIP", c.IP()),
						zap.String("err", err), zap.ByteString("core", buf))
				} else {
					log.Error("handler core dump", zap.String("clientIP", c.IP()),
						zap.Reflect("err", err), zap.ByteString("core", buf))
				}
				retErr = SendErrCode(c, errcode.ERR_SERVER_INTERNAL)

				return
			}
		}()

		ctx := &Ctx{Ctx: c}
		resp, err := fn(ctx)
		if err != nil {
			var errCode errcode.ErrCode
			if errors.As(err, &errCode) && !errors.Is(errCode, errcode.ERR_SUCCEED) {
				return err
			}
		}
		return SendResponse(c, resp)
	}
}

func NewOriginHandlerFn(fn func(*Ctx) error) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) (retErr error) {
		defer func() {
			if r := recover(); r != nil {
				buf := make([]byte, 1<<10)
				runtime.Stack(buf, true)
				if err, ok := r.(error); ok {
					log.Error("handler core dump", zap.String("clientIP", c.IP()),
						zap.String("err", err.Error()), zap.ByteString("core", buf))
				} else if err, ok := r.(string); ok {
					log.Error("handler core dump", zap.String("clientIP", c.IP()),
						zap.String("err", err), zap.ByteString("core", buf))
				} else {
					log.Error("handler core dump", zap.String("clientIP", c.IP()),
						zap.Reflect("err", err), zap.ByteString("core", buf))
				}

				retErr = errcode.ERR_SERVER_INTERNAL
				return
			}
		}()

		ctx := &Ctx{Ctx: c}
		retErr = fn(ctx)

		return retErr
	}
}
