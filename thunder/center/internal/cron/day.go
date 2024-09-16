package cron

import (
	"comm/t_data/redis"
	pb "comm/t_proto/out/client"
)

func daily() error {

	return nil
}

func dailyZero() error {
	if err := redis.UpdateRankByRankType(pb.RankType_RankDaily); err != nil {
		return err
	}
	if err := redis.UpdateGuildRankByRankType(pb.RankType_RankDaily); err != nil {
		return err
	}
	if err := redis.UpdateGuildInRankByRankType(pb.RankType_RankDaily); err != nil {
		return err
	}

	return nil
}
