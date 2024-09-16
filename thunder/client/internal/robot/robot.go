package robot

import (
	"client/internal/enums"
	"fmt"
	"github.com/v587-zyf/gc/iface"
	"github.com/v587-zyf/gc/log"
	"github.com/v587-zyf/gc/workerpool"
	"go.uber.org/zap"
	"sync"
)

type RobotConf struct {
	HttpAddr string
	GateAddr string

	Pem string
	Key string

	Https bool
}

type Robot struct {
	Id     uint64
	userID uint64
	token  string

	cfg *RobotConf

	ss iface.IWsSession

	status int32

	once sync.Once
	done chan struct{}

	handlers map[uint16]*ClientWsHandlerUnit

	sendMap map[int32]*SendCell
	recvMap map[uint16]*RecvCell
	menuStr string
	ipStr   string
	gmStr   string
}

var ID uint64 = 1

func NewRobot(cfg *RobotConf) *Robot {
	r := &Robot{
		Id:     ID,
		userID: ID,

		done: make(chan struct{}),
		cfg:  cfg,

		handlers: make(map[uint16]*ClientWsHandlerUnit),
		sendMap:  make(map[int32]*SendCell),
		recvMap:  make(map[uint16]*RecvCell),
	}
	ID++
	return r
}

func (r *Robot) SetStatus(status int32) {
	r.status = status
}
func (r *Robot) SetGateAddr(gateAddr string) {
	r.cfg.GateAddr = gateAddr
}
func (r *Robot) SetUserID(userID uint64) {
	r.userID = userID
}
func (r *Robot) SetToken(token string) {
	r.token = token
}

func (r *Robot) GetSession() iface.IWsSession {
	return r.ss
}

func (r *Robot) SendMsg(msgID uint16, msg iface.IProtoMessage) error {
	if r.ss == nil {
		return nil
	}

	//log.Debug("send----", zap.Any("msg", msg), zap.Uint16("msgID", msgID))
	err := r.ss.Send(msgID, 0, r.userID, msg)
	if err != nil {
		return err
	}

	return err
}

func (r *Robot) SetHttps() {
	isDev := InputInt32("是否Https？1->是 0->否")
	r.cfg.Https = false
	if isDev == 1 {
		r.cfg.Https = true
	}
}
func (r *Robot) SelectIp() {
	ipStrMap := map[int32]string{
		0: "127.0.0.1:8101",
		1: "thunder.majyo.vip:8101",
	}
	if r.ipStr == "" {
		r.ipStr = "输入要选择的ip:\n"
		for n, s := range ipStrMap {
			r.ipStr += fmt.Sprintf("%d->%s \n", n, s)
		}
	}

	ipN := InputInt32(r.ipStr)
	r.cfg.HttpAddr = ipStrMap[ipN]
}

func (r *Robot) Run() {
	go func() {
	LOOP:
		for {
			select {
			case <-r.done:
				log.Info("robot done", zap.Uint64("userID", r.userID))
				break LOOP
			}
		}
		r.Stop(r.ss)
	}()

START:
	r.SelectIp()
	r.SetHttps()

	fmt.Println("--------------menu---------------")
	for {
		r.ShowDes()
		action := InputInt32("请输入行动id：")

		if action == 999 {
			goto START
		}

		r.DoSend(action)
		//time.Sleep(time.Millisecond * 200)
	}
}

func (r *Robot) Start(ss iface.IWsSession) {
}

func (r *Robot) Stop(ss iface.IWsSession) {
	log.Info("robot close", zap.Uint64("userID", r.userID))
	//r.once.Do(func() {
	//close(r.done)

	r.ss.Close()

	r.SetStatus(enums.STATUS_DISCONNECT)
	//})
}

func (r *Robot) Recv(ss iface.IWsSession, data any) {
	//log.Debug("2------------------msg recv")

	msg, err := r.UnmarshalClient(data.([]byte))
	if err != nil {
		log.Error("msg UnmarshalClient", zap.Error(err))
		return
	}

	workerpool.AssignWsTask(r.GetHandler(msg.MsgID), r.ss, msg)
	//log.Debug("2.1------------------msg recv", zap.Any("msg", msg))
}
