package handler

import (
	"encoding/binary"
	"fmt"
	"github.com/v587-zyf/gc/enums"
	"github.com/v587-zyf/gc/gcnet/ws_session"
	"github.com/v587-zyf/gc/iface"
)

type ClientWsHandlerUnit struct {
	msgID   uint16          // 协议id
	handler ws_session.Recv // 协议具体方法
}

type ClientWsHandler struct {
	handlers map[uint16]*ClientWsHandlerUnit
}

var cWHandler = &ClientWsHandler{handlers: make(map[uint16]*ClientWsHandlerUnit)}

func GetClientWsHandler() *ClientWsHandler {
	return cWHandler
}

func (h *ClientWsHandler) Register(msgID uint16, handler ws_session.Recv) {
	h.handlers[msgID] = &ClientWsHandlerUnit{
		msgID:   msgID,
		handler: handler,
	}
}

func (h *ClientWsHandler) GetHandler(msgID uint16) ws_session.Recv {
	if h, ok := h.handlers[msgID]; ok {
		return h.handler
	}
	return nil
}

func (h *ClientWsHandler) HasHandler(msgID uint16) bool {
	_, ok := h.handlers[msgID]
	return ok
}

//func (h *ClientWsHandler) UnmarshalClient(data []byte) (*iface.MessageFrame, error) {
//	if len(data) < enums.MSG_HEADER_SIZE {
//		return nil, errors.New("packet has a wrong header, data too long")
//	}
//
//	buffer := bytes.NewBuffer(data)
//	frame := new(iface.MessageFrame)
//
//	binary.Read(buffer, binary.BigEndian, &frame.Len)
//	if frame.Len > enums.MSG_MAX_PACKET_SIZE-enums.MSG_HEADER_SIZE {
//		log.Error("err msg len", zap.Uint32("bodyLen", frame.Len))
//		return nil, fmt.Errorf("msg len too long")
//	}
//	binary.Read(buffer, binary.BigEndian, &frame.MsgID)
//	binary.Read(buffer, binary.BigEndian, &frame.Tag)
//	binary.Read(buffer, binary.BigEndian, &frame.UserID)
//
//	if msgPrototype := pb.GetMsgProtoType(frame.MsgID); msgPrototype != nil {
//		body := reflect.New(msgPrototype).Interface().(iface.IProtoMessage)
//		if err := body.Unmarshal(data[enums.MSG_HEADER_SIZE:]); err != nil {
//			return nil, err
//		}
//		frame.Body = body
//		return frame, nil
//	}
//
//	return nil, fmt.Errorf("unmarshl error, cmdId: %d, dataLen: %d", frame.MsgID, len(data))
//}
//func (h *ClientWsHandler) UnmarshalServer(data []byte) (*iface.MessageFrame, error) {
//	if len(data) < enums.MSG_HEADER_SIZE {
//		return nil, errors.New("packet has a wrong header, data too long")
//	}
//
//	buffer := bytes.NewBuffer(data)
//	frame := new(iface.MessageFrame)
//
//	binary.Read(buffer, binary.BigEndian, &frame.Len)
//	if frame.Len > enums.MSG_MAX_PACKET_SIZE-enums.MSG_HEADER_SIZE {
//		log.Error("err msg len", zap.Uint32("bodyLen", frame.Len))
//		return nil, fmt.Errorf("msg len too long")
//	}
//	binary.Read(buffer, binary.BigEndian, &frame.MsgID)
//	binary.Read(buffer, binary.BigEndian, &frame.Tag)
//	binary.Read(buffer, binary.BigEndian, &frame.UserID)
//
//	if msgPrototype := server.GetMsgProtoType(frame.MsgID); msgPrototype != nil {
//		body := reflect.New(msgPrototype).Interface().(iface.IProtoMessage)
//		if err := body.Unmarshal(data[enums.MSG_HEADER_SIZE:]); err != nil {
//			return nil, err
//		}
//		frame.Body = body
//		return frame, nil
//	}
//
//	return nil, fmt.Errorf("unmarshl error, cmdId: %d, dataLen: %d", frame.MsgID, len(data))
//}

func (h *ClientWsHandler) Marshal(msgID uint16, Tag uint32, userID uint64, msg iface.IProtoMessage) ([]byte, error) {
	size := msg.Size()
	data := make([]byte, enums.MSG_HEADER_SIZE+size)
	n, err := msg.MarshalTo(data[enums.MSG_HEADER_SIZE:])
	if err != nil {
		return nil, err
	}

	binary.BigEndian.PutUint32(data[0:4], uint32(n))
	binary.BigEndian.PutUint16(data[4:6], uint16(msgID))
	binary.BigEndian.PutUint32(data[6:10], uint32(Tag))
	binary.BigEndian.PutUint64(data[10:18], uint64(userID))
	dataLen := enums.MSG_HEADER_SIZE + size
	if dataLen <= enums.MSG_MAX_PACKET_SIZE {
		return data[:dataLen], nil
	} else {
		return nil, fmt.Errorf("msg len too long")
	}
}
