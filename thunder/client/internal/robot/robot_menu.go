package robot

import (
	"fmt"
	"github.com/v587-zyf/gc/errcode"
	"github.com/v587-zyf/gc/gcnet/ws_session"
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
	"reflect"
	"sort"
)

type SendCell struct {
	Des string        // 指令
	Fn  reflect.Value // 方法
}

func (c *SendCell) Do() {

}

type RecvCell struct {
	Fn reflect.Value // 方法
}

func (r *Robot) RegisterSend(actionID int32, fn interface{}, des string) error {
	_, has := r.sendMap[actionID]
	if has {
		log.Error("msg already exist", zap.Int32("actionID", actionID))
		return errcode.ERR_STANDARD_ERR
	}

	fnVal := reflect.ValueOf(fn)
	if fnVal.Kind() != reflect.Func {
		log.Error("func is err")
		return errcode.ERR_STANDARD_ERR
	}

	cell := &SendCell{
		Des: des,
		Fn:  fnVal,
	}
	r.sendMap[actionID] = cell

	//log.Info("menu registerSend succ", zap.Int32("msgId", msgId), zap.String("des", des))

	return nil
}

func (r *Robot) RegisterRecv(msgId uint16, fn ws_session.Recv) error {
	_, has := r.recvMap[msgId]
	if has {
		log.Error("msg already exist", zap.Uint16("msgID", msgId))
		return errcode.ERR_STANDARD_ERR
	}

	fnVal := reflect.ValueOf(fn)
	if fnVal.Kind() != reflect.Func {
		log.Error("func is err")
		return errcode.ERR_STANDARD_ERR
	}

	cell := &RecvCell{
		Fn: fnVal,
	}
	r.recvMap[msgId] = cell

	r.RegisterHandler(msgId, fn)

	//log.Info("menu registerRecv succ", zap.Uint16("msgId", msgId))

	return nil
}

func (r *Robot) ShowDes() {
	if r.menuStr == "" {
		keys := []int(nil)
		for k := range r.sendMap {
			keys = append(keys, int(k))
		}
		sort.Ints(keys)

		outStr := ""
		for _, key := range keys {
			cell, has := r.sendMap[int32(key)]
			if !has {
				continue
			}
			outStr += fmt.Sprintf("%d -> %s\n", key, cell.Des)
		}
		//log.Debug("--", zap.String("outStr", outStr))
		outStr = outStr[:len(outStr)-1]
		r.menuStr = outStr
	}
	fmt.Println(r.menuStr)
}

func (r *Robot) DoSend(action int32) bool {
	sendCell, ok := r.sendMap[action]
	if !ok {
		log.Error("send not found", zap.Int32("action", action))
		return false
	}

	ret := sendCell.Fn.Call(nil)
	if len(ret) <= 0 {
		log.Error("send fn call err")
		return false
	}
	e := ret[len(ret)-1].Interface()
	if e != nil {
		err, ok := e.(error)
		if ok {
			log.Error("send err", zap.String("err", err.Error()), zap.Int32("action", action))
			return false
		}
	}

	return true
}
