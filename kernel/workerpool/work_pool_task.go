package workerpool

import (
	"kernel/iface"
	"time"
)

type DelaySendTask struct {
	Delay  time.Duration
	MsgID  uint16
	Tag    uint32
	UserID uint64
	Msg    iface.IProtoMessage
	Func   func(mid uint16, tag uint32, uid uint64, pb iface.IProtoMessage)
}

func (t *DelaySendTask) Do() {
	time.Sleep(t.Delay)
	t.Func(t.MsgID, t.Tag, t.UserID, t.Msg)
}

func AssignDelaySendTask(delay time.Duration,
	fn func(mid uint16, tag uint32, uid uint64, pb iface.IProtoMessage),
	msgID uint16, tag uint32, userID uint64, msg iface.IProtoMessage) error {
	return Assign(&DelaySendTask{
		Delay:  delay,
		MsgID:  msgID,
		Tag:    tag,
		UserID: userID,
		Msg:    msg,
		Func:   fn,
	})
}
