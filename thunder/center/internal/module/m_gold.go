package module

import (
	"comm/comm"
	"comm/t_data/redis"
	pb "comm/t_proto/out/client"
	"comm/t_proto/out/server"
	"comm/t_tdb"
	"github.com/v587-zyf/gc/enums"
	"github.com/v587-zyf/gc/gcnet/grpc_msg"
	"github.com/v587-zyf/gc/log"
	"github.com/v587-zyf/gc/module"
	"github.com/v587-zyf/gc/rdb/rdb_cluster"
	"github.com/v587-zyf/gc/utils"
	"go.uber.org/zap"
	"kernel/handler"
	"sync"
)

type GoldMgr struct {
	module.DefModule

	userMap map[uint64]struct{}
	uLock   sync.RWMutex
}

func NewGoldMgr() *GoldMgr {
	return &GoldMgr{
		userMap: make(map[uint64]struct{}),
	}
}

func (m *GoldMgr) Add(userID uint64) {
	m.uLock.Lock()
	defer m.uLock.Unlock()

	m.userMap[userID] = struct{}{}
}
func (m *GoldMgr) Remove(userID uint64) {
	m.uLock.Lock()
	defer m.uLock.Unlock()

	delete(m.userMap, userID)
}

func (m *GoldMgr) Auto() error {
	m.uLock.RLock()
	defer m.uLock.RUnlock()

	if len(m.userMap) <= 0 {
		return nil
	}

	var (
		err     error
		ret     bool
		rdbUser *redis.User
		locker  *rdb_cluster.Locker
	)
	for userID := range m.userMap {
		locker = redis.LockUser(userID, GetClientModuleMgrOptions().SID)
		if locker == nil {
			log.Error("get locker err")
			continue
		}
		var unlock = func() { locker.Unlock() }

		rdbUser, ret = GetRdbUser(userID, locker)
		if !ret {
			unlock()
			continue
		}

		// 每秒获得数量=总每小时挂机收益值/3600，向上保留5位小数
		allIncome := 0
		for id, lv := range rdbUser.Card.Data {
			if cardLvCfg := t_tdb.GetCardLv(id, lv); cardLvCfg != nil {
				allIncome += t_tdb.GetCardLv(id, lv).Income
			} else {
				log.Error("get card lv err", zap.Uint64("userID", userID), zap.Int("cardID", id), zap.Int("lv", lv))
			}
		}
		multiply := rdbUser.Hasten.GetMultiply()
		addIncome := utils.RoundFloat(float64(allIncome)/3600*multiply, 5)
		rdbUser.AddGold(addIncome)
		//log.Debug("---", zap.Float64("addIncome", addIncome), zap.Float64("multiply", multiply))
		if err = redis.SetUser(rdbUser, locker); err != nil {
			unlock()
			log.Error("add gold err", zap.Uint64("userID", userID))
			continue
		}

		comm.Send2User(userID, pb.MsgID_Income_NtfId, &pb.IncomeNtf{Gold: utils.RoundFloat(rdbUser.Basic.Gold, 5), IncomeType: pb.IncomeType_OnHook})

		{
			reqBytes, err := handler.GetClientWsHandler().Marshal(uint16(server.MsgID_UserIncome_NtfId), 0, userID, &server.UserIncomeNtf{})
			if err != nil {
				unlock()
				log.Error("marshal msg gold to game err", zap.Uint64("userID", userID), zap.Error(err))
				continue
			}
			msgData := &server.MessageData{Sender: enums.SERVER_CENTER, Receiver: enums.SERVER_GAME, Content: reqBytes, MsgType: server.MsgType_Server}
			grpc_msg.SendToMsg(msgData)
		}

		unlock()
	}

	return nil
}
