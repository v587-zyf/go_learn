package cron

import (
	"comm/t_data/redis"
	pb "comm/t_proto/out/client"
)

func sundayZero() error {
	if err := redis.UpdateRankByRankType(pb.RankType_RankWeekly); err != nil {
		return err
	}
	if err := redis.UpdateGuildRankByRankType(pb.RankType_RankWeekly); err != nil {
		return err
	}
	if err := redis.UpdateGuildInRankByRankType(pb.RankType_RankWeekly); err != nil {
		return err
	}

	return nil
}
