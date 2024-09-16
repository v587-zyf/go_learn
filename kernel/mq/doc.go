package mq

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/v587-zyf/gc/errcode"
	"reflect"

	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
)

type IMethod interface {
	Start(c *Session)
	Recv(c *Session, msg *Msg, msgID int32, body []byte)
	Stop(c *Session)
}

type HandlerFn func(c *Session, msg *Msg, msgID int32, body []byte)

type Msg nats.Msg

type Head struct {
	Len   uint16
	MsgID int32
	Body  []byte
}

type ReplyHead struct {
	Code errcode.ErrCode
	Body []byte
}

func FormatServerSubject(serverType string, serverID int32) string {
	return fmt.Sprintf("ser.%s.%d", serverType, serverID)
}

func FormatBroadcastSubject(serverType string) string {
	return fmt.Sprintf("serbc.%s", serverType)
}

func FormatRandomSubject(serverType string) string {
	return fmt.Sprintf("serran.%s", serverType)
}

func FormatRandomGroup(serverType string) string {
	return fmt.Sprintf("grp.%s", serverType)
}

func UnmarshalHead(data []byte) (*ReplyHead, error) {
	if len(data) < 4 {
		return nil, errcode.ERR_MQ_REPLY_HEAD_LEN
	}

	head := &ReplyHead{}

	errNo := binary.LittleEndian.Uint32(data[0:4])
	head.Code = errcode.ErrCode(errNo)

	head.Body = data[4:]

	return head, nil
}

func MarshalSend(msgID int32, msg proto.Message) ([]byte, error) {
	body, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	buff := new(bytes.Buffer)

	err = binary.Write(buff, binary.LittleEndian, uint16(0))
	if err != nil {
		return nil, err
	}
	err = binary.Write(buff, binary.LittleEndian, uint32(msgID))
	if err != nil {
		return nil, err
	}

	err = binary.Write(buff, binary.LittleEndian, body)
	if err != nil {
		return nil, err
	}

	out := buff.Bytes()
	binary.LittleEndian.PutUint16(out[0:2], uint16(buff.Len()))

	return out, nil
}

func MarshalReplyEmpty(err error) ([]byte, error) {
	buff := new(bytes.Buffer)

	var retCode int32
	if errCode, ok := err.(errcode.ErrCode); ok {
		retCode = int32(errCode)
	} else {
		retCode = int32(errcode.ERR_STANDARD_ERR)
	}

	err = binary.Write(buff, binary.LittleEndian, retCode)
	if err != nil {
		return nil, errcode.ERR_MQ_BUFF_WRITE
	}

	out := buff.Bytes()

	return out, nil

}

func MarshalReply[T any](recvErr error, msg *T) (out []byte, err error) {

	buff := new(bytes.Buffer)

	defer func() {
		if err == nil {
			return
		}
		out, err = MarshalReplyEmpty(err)
	}()

	if recvErr != nil {
		errCode, ok := recvErr.(errcode.ErrCode)
		if !ok || (ok && !errors.Is(errCode, errcode.ERR_SUCCEED)) {
			return nil, recvErr
		}
	}

	// if reflect.ValueOf(msg).IsNil() {
	// 	return nil, errcode.ERR_MQ_REPLY_EMPTY
	// }
	if msg == nil {
		return nil, errcode.ERR_MQ_REPLY_EMPTY
	}

	body, err := proto.Marshal(reflect.ValueOf(msg).Interface().(proto.Message))
	if err != nil {
		return nil, errcode.ERR_MQ_REPLY_PB
	}

	err = binary.Write(buff, binary.LittleEndian, int32(0))
	if err != nil {
		return nil, errcode.ERR_MQ_BUFF_WRITE
	}

	err = binary.Write(buff, binary.LittleEndian, body)
	if err != nil {
		return nil, errcode.ERR_MQ_BUFF_WRITE
	}

	return buff.Bytes(), nil
}
