package cron

import (
	"center/internal/module"
	"comm/t_data/redis"
	"comm/t_enum"
)

func threeSecond() (err error) {
	if err = redis.DbDump(module.GetClientModuleMgrOptions().SID); err != nil {
		return err
	}

	return nil
}

func oneSecond() (err error) {
	if err = module.GetClientModuleMgr().GetModule(enum.C_M_GOLD).(*module.GoldMgr).Auto(); err != nil {
		return err
	}
	if err = module.GetClientModuleMgr().GetModule(enum.C_M_STRENGTH).(*module.StrengthMgr).Auto(); err != nil {
		return err
	}

	return nil
}
