package redis

import (
	errCode "comm/t_errcode"
	"github.com/redis/go-redis/v9"
	"github.com/v587-zyf/gc/enums"
	"github.com/v587-zyf/gc/errcode"
	"github.com/v587-zyf/gc/log"
	"github.com/v587-zyf/gc/rdb/rdb_cluster"
	"go.uber.org/zap"
)

func GetGate() (gateAddr string, gateID int32, errNo errcode.ErrCode) {
	rc := rdb_cluster.Get()
	rCtx := rdb_cluster.GetCtx()

	errNo = errcode.ERR_SUCCEED

	gateSlice := []string(nil)
	if err := rc.ZRange(rCtx, enums.RDB_KEY_SER_GATE, 0, 0).ScanSlice(&gateSlice); err != nil {
		errNo = errCode.ERR_SERVER_GATE_NIL
		log.Error("redis GetGate err", zap.Error(err))
		return
	}
	if len(gateSlice) == 0 {
		errNo = errCode.ERR_SERVER_GATE_NIL
		return
	}

	gateAddr, gateID = ParseGateData(gateSlice[0])
	if gateAddr == "" || gateID == 0 {
		errNo = errCode.ERR_SERVER_GATE_NIL
		return
	}

	return
}

func SetGate(score float64, member any) {
	rc := rdb_cluster.Get()
	rCtx := rdb_cluster.GetCtx()
	if err := rc.ZAdd(rCtx, enums.RDB_KEY_SER_GATE, redis.Z{Score: score, Member: member}).Err(); err != nil {
		log.Error("redis SetGate err", zap.Error(err))
	}
}

func DelGate(member any) {
	rc := rdb_cluster.Get()
	rCtx := rdb_cluster.GetCtx()

	if err := rc.ZRem(rCtx, enums.RDB_KEY_SER_GATE, member).Err(); err != nil {
		log.Error("redis DelGate err", zap.Error(err))
	}
}
