package mq

import (
	"errors"
	"fmt"
	"kernel/errcode"
	"kernel/log"
	"reflect"
	"runtime"

	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
)

func RegisterHandler[T1 any, T2 any](h *Handler, msgID int32, f func(*Session, *Msg, int32, *T1) (*T2, error)) error {
	pbType := reflect.TypeOf((*T1)(nil)).Elem()

	h.handlers[msgID] = func(c2 *Session, msg *Msg, msgID2 int32, data []byte) {

		pbVal := reflect.New(pbType)

		var ack *T2
		var err error
		defer func() {
			if r := recover(); r != nil {
				buf := make([]byte, 1<<10)
				runtime.Stack(buf, true)
				if err2, ok := r.(error); ok {
					log.Error("handler core dump", zap.String("subj", msg.Subject), zap.Int32("msgID", msgID2), zap.String("err", err2.Error()), zap.ByteString("core", buf))
					err = err2
				} else if err2, ok := r.(string); ok {
					log.Error("handler core dump", zap.String("subj", msg.Subject), zap.Int32("msgID", msgID2), zap.String("err", err2), zap.ByteString("core", buf))
					err = errors.New(err2)
				} else {
					log.Error("handler core dump", zap.String("subj", msg.Subject), zap.Int32("msgID", msgID2), zap.Reflect("err", err2), zap.ByteString("core", buf))
					err = fmt.Errorf("core dump, %+v", err2)
				}
			}

			if msg.Reply != "" {
				out, err := MarshalReply(err, ack)
				if err != nil {
					log.Error("marshal reply err", zap.String("subj", msg.Subject), zap.Int32("msgID", msgID), zap.Reflect("ack", ack), zap.String("err", err.Error()))
					return
				}
				err = c2.conn.Publish(msg.Reply, out)
				if err != nil {
					log.Error("mq reply, publish err", zap.String("subj", msg.Subject), zap.Int32("msgID", msgID), zap.String("err", err.Error()))
				}
			}

		}()

		err = proto.Unmarshal(data, pbVal.Interface().(proto.Message))
		if err != nil {
			log.Error("req pb unmarsal failed", zap.String("subj", msg.Subject), zap.Int32("msgID", msgID2), zap.String("err", err.Error()))
			err = errcode.ERR_MQ_RECV_DATA_UNMARSHAL
			return
		}
		ack, err = f(c2, msg, msgID2, pbVal.Interface().(*T1))

		// refV := reflect.ValueOf(ack)
		// out, err := proto.Marshal(refV.Interface().(proto.Message))

	}

	return nil
}

func RegisterHandlerEmptyReply[T1 any](h *Handler, msgID int32, f func(*Session, *Msg, int32, *T1) error) error {
	pbType := reflect.TypeOf((*T1)(nil)).Elem()

	h.handlers[msgID] = func(c2 *Session, msg *Msg, msgID2 int32, data []byte) {

		pbVal := reflect.New(pbType)

		var err error
		defer func() {
			if r := recover(); r != nil {
				buf := make([]byte, 1<<10)
				runtime.Stack(buf, true)
				if err2, ok := r.(error); ok {
					log.Error("handler core dump", zap.String("subj", msg.Subject), zap.Int32("msgID", msgID2), zap.String("err", err2.Error()), zap.ByteString("core", buf))
					err = err2
				} else if err2, ok := r.(string); ok {
					log.Error("handler core dump", zap.String("subj", msg.Subject), zap.Int32("msgID", msgID2), zap.String("err", err2), zap.ByteString("core", buf))
					err = errors.New(err2)
				} else {
					log.Error("handler core dump", zap.String("subj", msg.Subject), zap.Int32("msgID", msgID2), zap.Reflect("err", err2), zap.ByteString("core", buf))
					err = fmt.Errorf("core dump, %+v", err2)
				}
			}

			if msg.Reply != "" {
				out, err := MarshalReplyEmpty(err)
				if err != nil {
					log.Error("marshal reply err", zap.String("subj", msg.Subject), zap.Int32("msgID", msgID), zap.String("err", err.Error()))
					return
				}
				err = c2.conn.Publish(msg.Reply, out)
				if err != nil {
					log.Error("mq reply, publish err", zap.String("subj", msg.Subject), zap.Int32("msgID", msgID), zap.String("err", err.Error()))
				}
			}

		}()

		err = proto.Unmarshal(data, pbVal.Interface().(proto.Message))
		if err != nil {
			log.Error("req pb unmarsal failed", zap.String("subj", msg.Subject), zap.Int32("msgID", msgID2), zap.String("err", err.Error()))
			err = errcode.ERR_MQ_RECV_DATA_UNMARSHAL
			return
		}
		err = f(c2, msg, msgID2, pbVal.Interface().(*T1))

		// refV := reflect.ValueOf(ack)
		// out, err := proto.Marshal(refV.Interface().(proto.Message))

	}

	return nil
}

func NewHandler() *Handler {
	return &Handler{
		handlers: map[int32]HandlerFn{},
	}
}

type Handler struct {
	handlers map[int32]HandlerFn

	// onErrHandler ErrHandlerFn
}

func (h *Handler) HandleMsg(c *Session, msg *Msg, msgID int32, body []byte) {
	// 执行命令
	handler, ok := h.handlers[msgID]
	if !ok {
		log.Warn("recv msg no handler", zap.String("subj", msg.Subject), zap.Int32("msgID", msgID))

		if msg.Reply != "" {
			out, err := MarshalReplyEmpty(errcode.ERR_MQ_MSG_ID_NOT_REGISTER)
			if err != nil {
				log.Error("marshal reply err", zap.String("subj", msg.Subject), zap.Int32("msgID", msgID), zap.String("err", err.Error()))
				return
			}
			err = c.conn.Publish(msg.Reply, out)
			if err != nil {
				log.Error("mq reply, publish err", zap.String("subj", msg.Subject), zap.Int32("msgID", msgID), zap.String("err", err.Error()))
			}
		}

		return
	}

	handler(c, msg, msgID, body)

}
