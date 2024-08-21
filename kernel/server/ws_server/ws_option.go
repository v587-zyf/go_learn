package ws_server

import (
	"kernel/iface"
)

type WsOption struct {
	addr string
	pem  string
	key  string

	dev bool

	method iface.IWsSessionMethod
}

type Option func(opts *WsOption)

func NewWsOption() *WsOption {
	o := &WsOption{}

	return o
}

func WithAddr(addr string) Option {
	return func(opts *WsOption) {
		opts.addr = addr
	}
}

func WithPem(pem string) Option {
	return func(opts *WsOption) {
		opts.pem = pem
	}
}

func WithKey(key string) Option {
	return func(opts *WsOption) {
		opts.key = key
	}
}

func WithMethod(m iface.IWsSessionMethod) Option {
	return func(opts *WsOption) {
		opts.method = m
	}
}

func WithDev(dev bool) Option {
	return func(opts *WsOption) {
		opts.dev = dev
	}
}
