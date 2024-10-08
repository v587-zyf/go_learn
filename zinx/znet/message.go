package znet

type Message struct {
	Id      uint32 // 消息id
	DataLen uint32 // 消息长度
	Data    []byte // 消息内容
}

func NewMsgPackage(msgID uint32, data []byte) *Message {
	msg := &Message{
		Id:      msgID,
		DataLen: uint32(len(data)),
		Data:    data,
	}
	return msg
}

func (m *Message) GetMsgId() uint32 {
	return m.Id
}
func (m *Message) GetMsgLen() uint32 {
	return m.DataLen
}
func (m *Message) GetMsgData() []byte {
	return m.Data
}

func (m *Message) SetMsgId(id uint32) {
	m.Id = id
}
func (m *Message) SetMsgLen(len uint32) {
	m.DataLen = len
}
func (m *Message) SetMsgData(data []byte) {
	m.Data = data
}
