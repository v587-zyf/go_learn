package enum

import "time"

const (
	DUMP_USER_EXPIRE_TIME = 60 * 60 * 12 * time.Second
)

const SER_NAME = "{Thunder}"

const (
	REDIS_KEY_USER        = SER_NAME + "_User:%d"
	REDIS_KEY_RANK_DAILY  = SER_NAME + "_RankDaily"
	REDIS_KEY_RANK_WEEKLY = SER_NAME + "_RankWeekly"

	REDIS_KEY_GUILD                = SER_NAME + "_Guild"
	REDIS_KEY_GUILD_RANK_DAILY     = SER_NAME + "_Guild_RankDaily"
	REDIS_KEY_GUILD_RANK_WEEKLY    = SER_NAME + "_Guild_RankWeekly"
	REDIS_KEY_GUILD_IN_RANK_DAILY  = SER_NAME + "_Guild_In_RankDaily"
	REDIS_KEY_GUILD_IN_RANK_WEEKLY = SER_NAME + "_Guild_In_RankWeekly"
)
