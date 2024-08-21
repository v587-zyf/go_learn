package mq

type CommonMethod struct {
	handler *Handler
}

func (m *CommonMethod) Start(c *Session) {

}

func (m *CommonMethod) Recv(c *Session, msg *Msg, msgID int32, body []byte) {
	m.handler.HandleMsg(c, msg, msgID, body)
}

func (m *CommonMethod) Stop(c *Session) {

}
