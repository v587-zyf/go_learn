package mq

import (
	"context"
	"encoding/binary"
	"errors"
	"kernel/errcode"
	"kernel/iface"
	"kernel/log"
	"net"
	"strings"
	"time"
	"unsafe"

	"github.com/golang/protobuf/proto"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

func NewSession(cfg *Config) (*Session, error) {
	conn, err := nats.Connect(cfg.Url,
		nats.RetryOnFailedConnect(true),
		nats.ReconnectWait(2*time.Second),
		nats.DisconnectErrHandler(func(c *nats.Conn, err error) {
			if err != nil {
				log.Debug("mq disconnected ", zap.String("SType", cfg.SType), zap.Int32("SId", cfg.SId), zap.String("err", err.Error()))
			} else {
				log.Debug("mq disconnected ", zap.String("SType", cfg.SType), zap.Int32("SId", cfg.SId))
			}

		}),
		nats.ReconnectHandler(func(c *nats.Conn) {
			log.Debug("mq reconnect to ", zap.String("SType", cfg.SType), zap.Int32("SId", cfg.SId), zap.String("url", c.ConnectedUrl()))
		}),
		nats.ClosedHandler(func(c *nats.Conn) {
			err := c.LastError()
			if err != nil {
				log.Debug("mq connection closed ", zap.String("SType", cfg.SType), zap.Int32("SId", cfg.SId), zap.String("err", c.LastError().Error()))
			} else {
				log.Debug("mq connection closed ", zap.String("SType", cfg.SType), zap.Int32("SId", cfg.SId))
			}
		}),
		nats.ConnectHandler(func(c *nats.Conn) {
			log.Debug("mq connect succ ", zap.String("SType", cfg.SType), zap.Int32("SId", cfg.SId), zap.String("url", c.ConnectedUrl()))
		}),
	)
	if err != nil {
		return nil, err
	}

	isConnected := conn.IsConnected()
	if !isConnected {
		return nil, errcode.ERR_MQ_CONNECT_FAIL
	}

	subj := FormatServerSubject(cfg.SType, cfg.SId)
	ch := make(chan *nats.Msg, 128)
	subSingle, err := conn.ChanSubscribe(subj, ch)
	if err != nil {
		return nil, err
	}

	subBroadcast, err := conn.ChanSubscribe(FormatBroadcastSubject(cfg.SType), ch)
	if err != nil {
		return nil, err
	}

	subRandom, err := conn.ChanQueueSubscribe(FormatRandomSubject(cfg.SType), FormatRandomGroup(cfg.SType), ch)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	client := &Session{
		conn: conn,

		id:     int32(cfg.SId),
		method: cfg.Method,

		cfg: cfg,

		ctx:    ctx,
		cancel: cancel,

		ch:           ch,
		subSingle:    subSingle,
		subBroadcast: subBroadcast,
		subRandom:    subRandom,
	}

	client.startPump()

	return client, nil
}

type Session struct {
	conn *nats.Conn
	id   int32

	cfg *Config

	ctx    context.Context
	cancel context.CancelFunc

	ch chan *nats.Msg

	subSingle    *nats.Subscription
	subBroadcast *nats.Subscription
	subRandom    *nats.Subscription

	// handler *Handler
	method IMethod
}

func (s *Session) GetID() int32 {
	return s.id
}
func (s *Session) SetID(id int32) {
	s.id = id
}

func (s *Session) GetLoginTime() int64 {
	return 0
}

func (s *Session) SetLoginTime(now int64) {

}

func (s *Session) GetConn() net.Conn {
	return nil
}

func (s *Session) GetCtx() context.Context {
	return s.ctx
}

func (s *Session) Serve(method iface.IWsSessionMethod) error {
	return nil
}

func (s *Session) Close() error {
	s.cancel()
	return nil
}

func (s *Session) Send(msgID, tag, userID int32, data []byte) error {
	return nil
}

func (s *Session) SendMsg(msgID, tag, userID int32, pb proto.Message) error {
	return nil
}

func (s *Session) GetCache(key string) (interface{}, bool) {
	return nil, false
}

func (s *Session) SetCache(key string, value interface{}) {

}

func (s *Session) GetConfig() *Config {
	return s.cfg
}

func (s *Session) RequestOne(SType string, SId int32, msgID int32, msg proto.Message, reply proto.Message) error {
	return s.Request(FormatServerSubject(SType, SId), msgID, msg, reply)
}

func (c *Session) RequestRand(SType string, msgID int32, msg proto.Message, reply proto.Message) error {
	return c.Request(FormatRandomSubject(SType), msgID, msg, reply)
}

func (s *Session) SendOne(SType string, SId int32, msgID int32, msg proto.Message) error {
	return s.SendBase(FormatServerSubject(SType, SId), msgID, msg)
}

func (s *Session) Broadcast(SType string, msgID int32, msg proto.Message) error {
	return s.SendBase(FormatBroadcastSubject(SType), msgID, msg)
}

func (s *Session) SendRand(SType string, msgID int32, msg proto.Message) error {
	return s.SendBase(FormatRandomSubject(SType), msgID, msg)
}

func (s *Session) Request(subj string, msgID int32, msg proto.Message, reply proto.Message) error {

	out, err := MarshalSend(msgID, msg)
	if err != nil {
		log.Error("mq request, marshal err", zap.String("subj", subj), zap.Int32("msgID", msgID), zap.String("err", err.Error()))
		return err
	}

	retMsg, err := s.conn.Request(subj, out, 5*time.Second)
	if err != nil {
		if strings.HasSuffix(err.Error(), "timeout") {
			return errcode.ERR_MQ_REQ_TIMEOUT
		} else if strings.HasSuffix(err.Error(), "no responders available for request") {
			return errcode.ERR_MQ_SERVER_NOT_FOUND
		}
		log.Error("mq request, conn err", zap.String("subj", subj), zap.Int32("msgID", msgID), zap.String("err", err.Error()))
		return err
	}

	head, err := UnmarshalHead(retMsg.Data)
	if err != nil {
		return err
	}
	if !errors.Is(head.Code, errcode.ERR_SUCCEED) {
		return head.Code
	}

	if len(head.Body) > 0 {
		err = proto.Unmarshal(head.Body, reply)
		if err != nil {
			log.Error("mq request, unmarshal reply err", zap.String("subj", subj), zap.Int32("msgID", msgID), zap.String("err", err.Error()))
			return err
		}
	}

	return nil
}

func (s *Session) SendBase(subj string, msgID int32, msg proto.Message) error {

	out, err := MarshalSend(msgID, msg)
	if err != nil {
		return err
	}

	err = s.conn.Publish(subj, out)
	return err
}

func (s *Session) startPump() {

	go func() {

		s.method.Start(s)
	LOOP:
		for {
			select {
			case msg := <-s.ch:
				s.onRecv(msg)
			case <-s.ctx.Done():
				break LOOP
			}
		}

		s.method.Stop(s)

		s.subSingle.Drain()
		s.subSingle = nil

		s.subBroadcast.Drain()
		s.subBroadcast = nil

		s.subRandom.Drain()
		s.subRandom = nil

		s.conn.Close()
		s.conn = nil

		close(s.ch)
		s.ch = nil

		log.Debug("mq client close", zap.String("SType", s.cfg.SType), zap.Int32("SId", s.cfg.SId))
	}()
}

func (s *Session) onRecv(msg *nats.Msg) {

	if len(msg.Data) < 6 {
		log.Error("invalid msg data", zap.String("subj", msg.Subject), zap.Int("dataLen", len(msg.Data)))
		return
	}

	m := (*Msg)(unsafe.Pointer(msg))

	// len := binary.LittleEndian.Uint16(msg.Data[:2])
	msgID := binary.LittleEndian.Uint32(msg.Data[2:6])
	s.method.Recv(s, m, int32(msgID), msg.Data[6:])
}
