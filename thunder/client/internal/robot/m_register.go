package robot

import (
	"encoding/json"
	"github.com/v587-zyf/gc/log"
	"github.com/v587-zyf/gc/utils"
	"go.uber.org/zap"
)

type RegisterResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func (r *Robot) Register() (err error) {
	registerReq := new(LoginReq)
	registerReq.Account = InputString("input account")
	registerReq.Password = InputString("input password")

	reqBytes, err := json.Marshal(registerReq)
	if err != nil {
		log.Error("marshal err", zap.Error(err))
		return
	}

	var resp []byte
	if r.cfg.Https {
		resp, err = utils.PostJsonDiyClient(r.cfg.HttpAddr+"/register", reqBytes, r.cfg.Pem)
	} else {
		resp, err = utils.PostJson(r.cfg.HttpAddr+"/register", reqBytes)
	}

	if err != nil {
		log.Error("http post err", zap.Error(err))
		return err
	}

	registerResp := new(RegisterResp)
	if err = json.Unmarshal(resp, &registerResp); err != nil {
		log.Error("loginResp unmarshal err", zap.Error(err))
		return
	}
	log.Debug("register resp", zap.Any("registerResp", registerResp))

	return
}
