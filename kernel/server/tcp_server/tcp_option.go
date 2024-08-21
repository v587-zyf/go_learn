package tcp_server

import (
	"kernel/iface"
)

type TcpOption struct {
	listenAddr string

	method iface.ITpcSessionMethod
}

type Option func(opts *TcpOption)

func NewTcpOption() *TcpOption {
	o := &TcpOption{}

	return o
}

func WithListenAddr(addr string) Option {
	return func(opts *TcpOption) {
		opts.listenAddr = addr
	}
}

func WithMethod(m iface.ITpcSessionMethod) Option {
	return func(opts *TcpOption) {
		opts.method = m
	}
}
