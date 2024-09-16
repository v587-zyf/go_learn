package module

import (
	"comm/comm"
	"comm/t_data/redis"
	pb "comm/t_proto/out/client"
	"comm/t_tdb"
	"github.com/v587-zyf/gc/log"
	"github.com/v587-zyf/gc/module"
	"github.com/v587-zyf/gc/rdb/rdb_cluster"
	"go.uber.org/zap"
	"sync"
	"time"
)

type StrengthMgr struct {
	module.DefModule

	userMap map[uint64]struct{}
	uLock   sync.RWMutex
}

func NewStrengthMgr() *StrengthMgr {
	return &StrengthMgr{
		userMap: make(map[uint64]struct{}),
	}
}

func (m *StrengthMgr) Add(userID uint64) {
	m.uLock.Lock()
	defer m.uLock.Unlock()

	m.userMap[userID] = struct{}{}
}
func (m *StrengthMgr) Remove(userID uint64) {
	m.uLock.Lock()
	defer m.uLock.Unlock()

	delete(m.userMap, userID)
}

func (m *StrengthMgr) Auto() error {
	m.uLock.RLock()
	defer m.uLock.RUnlock()

	if len(m.userMap) <= 0 {
		return nil
	}

	var (
		timeNow  = time.Now()
		addTimes = time.Duration(t_tdb.Conf().Add_strengthen_second)
		addNum   = t_tdb.Conf().Add_strengthen_num

		err         error
		ret         bool
		maxStrength int
		rdbUser     *redis.User
		locker      *rdb_cluster.Locker
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

		maxStrength = t_tdb.GetUserUserCfg(rdbUser.Basic.Lv).Strength_max
		if rdbUser.Basic.Strength >= maxStrength && rdbUser.Basic.ExtraStrength >= rdbUser.Shop.GetAllShopStrength() {
			unlock()
			continue
		}

		diff := timeNow.Sub(rdbUser.Basic.LastStrengthAt)
		if diff >= addTimes*time.Second {
			if rdbUser.Basic.Strength < maxStrength {
				rdbUser.Basic.Strength += addNum
				if rdbUser.Basic.Strength > maxStrength {
					rdbUser.Basic.Strength = maxStrength
				}
			} else if rdbUser.Basic.ExtraStrength <= rdbUser.Shop.GetAllShopStrength() {
				rdbUser.Basic.ExtraStrength += addNum
				if rdbUser.Basic.ExtraStrength > rdbUser.Shop.GetAllShopStrength() {
					rdbUser.Basic.ExtraStrength = rdbUser.Shop.GetAllShopStrength()
				}
			}
			rdbUser.Basic.LastStrengthAt = timeNow

			if err = redis.SetUser(rdbUser, locker); err != nil {
				unlock()
				log.Error("add strength err", zap.Uint64("userID", userID))
				continue
			}

			comm.Send2User(userID, pb.MsgID_Strength_NtfId, &pb.StrengthNtf{Strength: int32(rdbUser.Basic.Strength), ExtraStrength: int32(rdbUser.Basic.ExtraStrength)})
		}

		unlock()
	}

	return nil
}
