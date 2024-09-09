package ws_session

import (
	"bytes"
	"context"
	"encoding/binary"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"kernel/enums"
	"kernel/errcode"
	"kernel/iface"
	"kernel/log"
	"runtime"
	"sync"
	"time"
)

type Session struct {
	id   uint64
	conn *websocket.Conn

	ctx    context.Context
	cancel context.CancelFunc

	cache     map[string]any
	cacheLock sync.RWMutex

	inChan  chan []byte
	outChan chan []byte

	hooks *Hooks

	method iface.ITpcSessionMethod

	reconnectTimes int

	once sync.Once
}

func NewSession(ctx context.Context, conn *websocket.Conn) *Session {
	ctx, cancel := context.WithCancel(ctx)
	s := &Session{
		ctx:    ctx,
		cancel: cancel,

		inChan:  make(chan []byte, 1024),
		outChan: make(chan []byte, 1024),

		cache: make(map[string]any),

		hooks: NewHooks(),
	}
	s.conn = conn

	return s
}

func (s *Session) Start() {
	go s.readPump()
	go s.IOPump()
	//go func() {
	//	s.parsePump()
	//}()
	//
	//go func() {
	//	s.writePump()
	//}()
}

func (s *Session) Hooks() *Hooks {
	return s.hooks
}

func (s *Session) Set(key string, value any) {
	s.cacheLock.Lock()
	defer s.cacheLock.Unlock()

	s.cache[key] = value
}
func (s *Session) Get(key string) (any, bool) {
	s.cacheLock.RLock()
	defer s.cacheLock.RUnlock()

	v, ok := s.cache[key]
	return v, ok
}
func (s *Session) Remove(key string) {
	s.cacheLock.Lock()
	defer s.cacheLock.Unlock()

	delete(s.cache, key)
}

func (s *Session) GetID() uint64 {
	return s.id
}
func (s *Session) SetID(id uint64) {
	if id <= 0 {
		id = 0
	}
	s.id = id
}

func (s *Session) Close() error {
	// return s.Conn.Close()
	//log.Info("session close", zap.Int32("sessID", s.GetID()))
	s.once.Do(func() {
		s.cancel()
		s.conn.Close()
	})

	return nil
}

func (s *Session) GetConn() *websocket.Conn {
	return s.conn
}

func (s *Session) GetCtx() context.Context {
	return s.ctx
}

func (s *Session) Send(msgID uint16, tag uint32, userID uint64, msg iface.IProtoMessage) error {
	//log.Debug("2---------------------", zap.Any("msg", msg), zap.Uint16("msgID", msgID))
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.BigEndian, uint32(msg.Size()))
	binary.Write(buf, binary.BigEndian, msgID)
	binary.Write(buf, binary.BigEndian, tag)
	binary.Write(buf, binary.BigEndian, userID)

	if msg.Size()+enums.MSG_HEADER_SIZE > enums.MSG_MAX_PACKET_SIZE {
		return errcode.ERR_NET_PKG_LEN_LIMIT
	}

	data, err := msg.Marshal()
	if err != nil {
		log.Error("msg marshal err", zap.Uint64("userID", userID), zap.Any("msg", msg))
		return err
	}
	//log.Debug("3---------------------", zap.Any("msg", msg), zap.Uint16("msgID", msgID), zap.ByteString("data", data))

	buf.Write(data)

	select {
	case s.outChan <- buf.Bytes():
		return nil
	default:
		return errcode.ERR_NET_SEND_TIMEOUT
	}
}

func (s *Session) Send2User(msgID uint16, msg iface.IProtoMessage) error {
	return s.Send(msgID, 0, s.GetID(), msg)
}

func (s *Session) GetReconnectTimes() int {
	return s.reconnectTimes
}
func (s *Session) AddReconnectTimes() {
	s.reconnectTimes++
}

func (s *Session) readPump() {
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 1<<10)
			runtime.Stack(buf, true)
			if err, ok := r.(error); ok {
				log.Error("core dump", zap.Uint64("sessID", s.GetID()),
					zap.String("err", err.Error()), zap.ByteString("core", buf))
			} else if err, ok := r.(string); ok {
				log.Error("core dump", zap.Uint64("sessID", s.GetID()),
					zap.String("err", err), zap.ByteString("core", buf))
			} else {
				log.Error("core dump", zap.Uint64("sessID", s.GetID()),
					zap.Reflect("err", err), zap.ByteString("core", buf))
			}
		}
	}()

	s.hooks.ExecuteStart(s)

LOOP:
	for {
		_, message, err := s.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err,
				websocket.CloseGoingAway,
				websocket.CloseNoStatusReceived,
				websocket.CloseAbnormalClosure) {
				log.Info("ws_server read err", zap.Error(err))
			}
			break LOOP
		}
		if message != nil && len(message) > 0 {
			//log.Debug("1------------------recv msg")
			dataCopy := make([]byte, len(message))
			copy(dataCopy, message)
			s.inChan <- dataCopy
		}
	}

	s.cancel()
}

func (s *Session) IOPump() {
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 1<<10)
			runtime.Stack(buf, true)
			if err, ok := r.(error); ok {
				log.Error("core dump", zap.Uint64("sessID", s.GetID()),
					zap.String("err", err.Error()), zap.ByteString("core", buf))
			} else if err, ok := r.(string); ok {
				log.Error("core dump", zap.Uint64("ssID", s.GetID()),
					zap.String("err", err), zap.ByteString("core", buf))
			} else {
				log.Error("core dump", zap.Uint64("ssID", s.GetID()),
					zap.Reflect("err", err), zap.ByteString("core", buf))
			}
		}
	}()

	heartbeatTicker := time.NewTicker(enums.HEARTBEAT_INTERVAL)
	defer heartbeatTicker.Stop()

LOOP:
	for {
		select {
		case data := <-s.inChan:
			s.hooks.ExecuteRecv(s, data)
		case data := <-s.outChan:
			s.conn.SetWriteDeadline(time.Now().Add(enums.CONN_WRITE_WAIT_TIME))

			err := s.conn.WriteMessage(websocket.BinaryMessage, data)
			if err != nil {
				msgID := binary.BigEndian.Uint16(data[0:2])
				log.Warn("conn write err", zap.Uint64("userID", s.id),
					zap.Uint16("msgID", msgID), zap.Int("len", len(data)), zap.Error(err))
				break LOOP
			}
		case <-heartbeatTicker.C:
			//msg := new(pb.Heartbeat)
			//s.Send2User(pb.MsgID_HeartbeatId, msg)
		case <-s.ctx.Done():
			break LOOP
		}
	}

	s.hooks.ExecuteStop(s)

	s.Close()
}

//func (s *Session) parsePump() {
//	defer func() {
//		if r := recover(); r != nil {
//			buf := make([]byte, 1<<10)
//			runtime.Stack(buf, true)
//			if err, ok := r.(error); ok {
//				log.Error("core dump", zap.Uint64("sessID", s.GetID()),
//					zap.String("err", err.Error()), zap.ByteString("core", buf))
//			} else if err, ok := r.(string); ok {
//				log.Error("core dump", zap.Uint64("sessID", s.GetID()),
//					zap.String("err", err), zap.ByteString("core", buf))
//			} else {
//				log.Error("core dump", zap.Uint64("sessID", s.GetID()),
//					zap.Reflect("err", err), zap.ByteString("core", buf))
//			}
//		}
//	}()
//
//LOOP:
//	for {
//		select {
//		case data := <-s.inChan:
//			s.hooks.ExecuteRecv(s, data)
//		case <-s.ctx.Done():
//			break LOOP
//		}
//	}
//}
//
//func (s *Session) writePump() {
//	defer func() {
//		if r := recover(); r != nil {
//			buf := make([]byte, 1<<10)
//			runtime.Stack(buf, true)
//			if err, ok := r.(error); ok {
//				log.Error("core dump", zap.Uint64("sessID", s.GetID()),
//					zap.String("err", err.Error()), zap.ByteString("core", buf))
//			} else if err, ok := r.(string); ok {
//				log.Error("core dump", zap.Uint64("sessID", s.GetID()),
//					zap.String("err", err), zap.ByteString("core", buf))
//			} else {
//				log.Error("core dump", zap.Uint64("sessID", s.GetID()),
//					zap.Reflect("err", err), zap.ByteString("core", buf))
//			}
//		}
//	}()
//
//LOOP:
//	for {
//		select {
//		case data := <-s.outChan:
//			s.conn.SetWriteDeadline(time.Now().Add(enums.CONN_WRITE_WAIT_TIME))
//
//			err := s.conn.WriteMessage(websocket.BinaryMessage, data)
//			if err != nil {
//				msgID := binary.BigEndian.Uint16(data[0:2])
//				log.Warn("conn write err", zap.Uint64("userID", s.id),
//					zap.Uint16("msgID", msgID), zap.Int("len", len(data)), zap.Error(err))
//				break LOOP
//			}
//		case <-s.ctx.Done():
//			break LOOP
//		}
//	}
//
//	s.conn.Close()
//	s.cancel()
//
//	s.hooks.ExecuteStop(s)
//}
//
//func (s *Session) split(data []byte, atEOF bool) (advance int, token []byte, err error) {
//	if atEOF && len(data) == 0 {
//		return 0, nil, nil
//	}
//	dataLen := len(data)
//	if dataLen < enums.MSG_HEADER_SIZE {
//		return 0, nil, nil
//	}
//
//	// body len
//	n := int(binary.BigEndian.Uint32(data[0:4]))
//	if n > enums.MSG_MAX_PACKET_SIZE-enums.MSG_HEADER_SIZE || n < 0 {
//		log.Error("body len invalid", zap.Uint64("sessID", s.id),
//			zap.Int("n", n), zap.String("addr", s.GetConn().RemoteAddr().String()))
//		return 0, nil, errcode.ERR_NET_BODY_LEN_INVALID
//	}
//	if dataLen < n+enums.MSG_HEADER_SIZE {
//		return 0, nil, nil
//	}
//	return n + enums.MSG_HEADER_SIZE, data[0 : n+enums.MSG_HEADER_SIZE], nil
//}
