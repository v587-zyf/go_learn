package robot

import (
	"client/internal/enums"
	"comm/t_model"
	pb "comm/t_proto/out/client"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/v587-zyf/gc/gcnet/ws_session"
	"github.com/v587-zyf/gc/log"
	"github.com/v587-zyf/gc/utils"
	"go.uber.org/zap"
)

type LoginResp struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data *model.LoginAck `json:"data"`
}

func (r *Robot) Login() (err error) {
	tStr := "请输入登录类型\n"
	for k, v := range pb.LoginType_name {
		tStr += fmt.Sprintf("%d->%s\n", k, v)
	}

	t := InputInt32(tStr)

	loginReq := new(model.LoginReq)
	if t == int32(pb.LoginType_telegram) {
		loginReq.Types = pb.LoginType_telegram

		// 刘俊
		//loginReq.InitData = "query_id=AAEh9XR6AgAAACH1dHoWywt_&user=%7B%22id%22%3A6349452577%2C%22first_name%22%3A%22%E5%A4%A9%E5%A4%A9%22%2C%22last_name%22%3A%22%E5%90%91%E4%B8%8A%22%2C%22username%22%3A%22maowang0105%22%2C%22language_code%22%3A%22zh-hans%22%2C%22allows_write_to_pm%22%3Atrue%7D&auth_date=1721011277&hash=dded060bfa886a8ccdd4e2792f7364575208f26c93bf346f01deba2e53f2216c"

		// huihui
		//loginReq.InitData = "query_id=AAEhcA0aAwAAACFwDRpYRGDl&user=%7B%22id%22%3A6879539233%2C%22first_name%22%3A%22bi%22%2C%22last_name%22%3A%22shenghui%22%2C%22username%22%3A%22bishenghui66%22%2C%22language_code%22%3A%22zh-hans%22%2C%22allows_write_to_pm%22%3Atrue%7D&auth_date=1723518276&hash=98d275bdc739f347bb0aeb1d76fe0b9e62f826e94934841e9cee8209b85781a4"

		loginReq.InitData = "query_id=AAGUjvEAAwAAAJSO8QBDhoiZ&user=%7B%22id%22%3A6458281620%2C%22first_name%22%3A%22Redface%22%2C%22last_name%22%3A%22%22%2C%22username%22%3A%22redfacenine%22%2C%22language_code%22%3A%22zh-hans%22%2C%22allows_write_to_pm%22%3Atrue%7D&auth_date=1722246028&hash=6a15f4d507718d447dd41c20ee94b5d80535936bb2e8b5c09ede8fecc5af6d81"
		loginReq.Invite = 10000001
	} else {
		loginReq.Types = pb.LoginType_password

		loginReq.Account = InputString("input account")
		loginReq.Password = InputString("input password")
	}

	reqBytes, err := json.Marshal(loginReq)
	if err != nil {
		log.Error("marshal err", zap.Error(err))
		return
	}

	var resp []byte
	if r.cfg.Https {
		resp, err = utils.PostJsonDiyClient(r.cfg.HttpAddr+"/login", reqBytes, r.cfg.Pem)
	} else {
		resp, err = utils.PostJson(r.cfg.HttpAddr+"/login", reqBytes)
	}

	if err != nil {
		log.Error("http post err", zap.Error(err))
		return err
	}

	loginResp := new(LoginResp)
	if err = json.Unmarshal(resp, &loginResp); err != nil {
		log.Error("loginResp unmarshal err", zap.Error(err))
		return
	}
	log.Debug("login resp", zap.Any("loginResp", loginResp))

	if loginResp.Data != nil {
		r.SetToken(loginResp.Data.Token)
		r.SetUserID(loginResp.Data.UserId)
		r.SetGateAddr(loginResp.Data.LinkAddr)

		r.SetStatus(enums.STATUS_LOGIN)
	}

	return
}

func (r *Robot) LoginPass() (err error) {
	acc := InputString("input account")
	pass := InputString("input password")
	loginReq := &model.LoginReq{
		Types: pb.LoginType_password,

		Account:  acc,
		Password: pass,
	}

	reqBytes, err := json.Marshal(loginReq)
	if err != nil {
		log.Error("marshal err", zap.Error(err))
		return
	}

	var resp []byte
	if r.cfg.Https {
		resp, err = utils.PostJsonDiyClient(r.cfg.HttpAddr+"/login", reqBytes, r.cfg.Pem)
	} else {
		resp, err = utils.PostJson(r.cfg.HttpAddr+"/login", reqBytes)
	}

	if err != nil {
		log.Error("http post err", zap.Error(err))
		return err
	}

	loginResp := new(LoginResp)
	if err = json.Unmarshal(resp, &loginResp); err != nil {
		log.Error("loginResp unmarshal err", zap.Error(err))
		return
	}
	log.Debug("login resp", zap.Any("loginResp", loginResp))

	if loginResp.Data != nil {
		r.SetToken(loginResp.Data.Token)
		r.SetUserID(loginResp.Data.UserId)
		r.SetGateAddr(loginResp.Data.LinkAddr)

		r.SetStatus(enums.STATUS_LOGIN)
	}

	return
}

func (r *Robot) ConnectLs() (err error) {
	var urls string
	if r.cfg.Https {
		urls = fmt.Sprintf("ws://%s/ws", r.cfg.GateAddr)
	} else {
		urls = fmt.Sprintf("wss://%s/ws", r.cfg.GateAddr)
	}
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false, // 跳过证书验证
	}
	dialer := &websocket.Dialer{
		TLSClientConfig: tlsConfig,
	}

	conn, _, err := dialer.Dial(urls, nil)
	if err != nil {
		panic(err)
		return err
	}

	if r.ss != nil {
		r.ss.GetConn().Close()
	}

	ss := ws_session.NewSession(context.Background(), conn)
	if err != nil {
		log.Error("new session err", zap.Error(err))
		return
	}

	r.ss = ss
	r.ss.Start()
	r.ss.(*ws_session.Session).Hooks().OnMethod(r)

	return
}

func (r *Robot) Disconnect() (err error) {
	if r.ss.GetConn() != nil {
		r.ss.GetConn().Close()
	}

	return
}

func (r *Robot) StopMenu() (err error) {

	return
}
