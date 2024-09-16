package handler

import (
	"encoding/binary"
	"fmt"
	"github.com/v587-zyf/gc/enums"
	"github.com/v587-zyf/gc/gcnet/tpc_session"
	"github.com/v587-zyf/gc/iface"
)

type ClientTcpHandlerUnit struct {
	msgID   uint16           // 协议id
	handler tpc_session.Recv // 协议具体方法
}

type ClientTcpHandler struct {
	handlers map[uint16]*ClientTcpHandlerUnit
}

var cTHandler = &ClientTcpHandler{handlers: make(map[uint16]*ClientTcpHandlerUnit)}

func GetClientTcpHandler() *ClientTcpHandler {
	return cTHandler
}

func (h *ClientTcpHandler) Register(msgID uint16, handler tpc_session.Recv) {
	h.handlers[msgID] = &ClientTcpHandlerUnit{
		msgID:   msgID,
		handler: handler,
	}
}

func (h *ClientTcpHandler) GetHandler(msgID uint16) tpc_session.Recv {
	if h, ok := h.handlers[msgID]; ok {
		return h.handler
	}
	return nil
}

func (h *ClientTcpHandler) HasHandler(msgID uint16) bool {
	_, ok := h.handlers[msgID]
	return ok
}

//func (h *ClientTcpHandler) Unmarshal(data []byte) (any, error) {
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

func (h *ClientTcpHandler) GetMsgID(data []byte) uint16 {
	if len(data) < enums.MSG_HEADER_SIZE {
		return 0
	}
	return binary.BigEndian.Uint16(data[0:2])
}

func (h *ClientTcpHandler) Marshal(msgID uint16, Tag uint32, userID uint64, msg iface.IProtoMessage) ([]byte, error) {
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
