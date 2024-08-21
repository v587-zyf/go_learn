package tcp_server

import (
	"context"
	"go.uber.org/zap"
	"kernel/log"
	"kernel/session/tpc_session"
	"net"
	"sync"
)

type TcpServer struct {
	options *TcpOption

	ctx    context.Context
	cancel context.CancelFunc

	listener net.Listener

	wg sync.WaitGroup
}

func NewTcpServer() *TcpServer {
	s := &TcpServer{
		options: NewTcpOption(),
	}

	return s
}

func (s *TcpServer) Init(ctx context.Context, option ...any) (err error) {
	s.ctx, s.cancel = context.WithCancel(ctx)

	for _, opt := range option {
		opt.(Option)(s.options)
	}

	s.listener, err = net.Listen("tcp", s.options.listenAddr)
	if err != nil {
		log.Error("net listen err", zap.Error(err))
		return
	}

	return nil
}

func (s *TcpServer) Start() {
	s.wg.Add(1)
	go func(svr *TcpServer) {
		defer func() {
			svr.wg.Done()
		}()

	LOOP:
		for {
			c, err := s.listener.Accept()
			if err != nil {
				log.Error("tcp listen err", zap.Error(err))
				break LOOP
			}
			ss := tpc_session.NewSession(context.Background(), c)
			ss.Hooks().OnMethod(s.options.method)
			ss.Start()
		}

		log.Debug("server end", zap.String("addr", s.options.listenAddr))
	}(s)

	s.Wait()
}

func (s *TcpServer) Stop() {
	s.listener.Close()
}

func (s *TcpServer) Wait() {
	s.wg.Wait()
}
