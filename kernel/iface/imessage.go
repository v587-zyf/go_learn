package iface

import (
	"github.com/golang/protobuf/proto"
)

type IProtoMessage interface {
	proto.Message
	Marshal() ([]byte, error)
	MarshalTo([]byte) (int, error)
	Unmarshal([]byte) error
	Size() int
}

type MessageFrame struct {
	Len    uint32
	MsgID  uint16
	Tag    uint32
	UserID uint64
	Body   IProtoMessage
}
