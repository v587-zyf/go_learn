package redis

import (
	"comm/t_enum"
	"fmt"
	"github.com/v587-zyf/gc/utils"
)

// {UserID:1234567890}
func FormatKeyUserID(userID any) string {
	return fmt.Sprintf(enum.REDIS_KEY_USER, userID)
}
func FormatUserID(userID any) string {
	return fmt.Sprint("{UserID:", userID, "}")
}

// {UserID:1234567890}Reconnect
func FormatUserReconnect(userID any) string {
	return fmt.Sprint(enum.SER_NAME+FormatUserID(userID), "Reconnect")
}

// {GateAddr:127.0.0.1:8080}_{GateID:1234567890}
func FormatGateData(addr string, id int32) string {
	return fmt.Sprint(addr, "_", id)
}

// {UserID:1234567890}Login
func FormatUserLogin(userID any) string {
	return fmt.Sprint(enum.SER_NAME+FormatUserID(userID), "Login")
}

// {UserDump}
func FormatUserDumpKey() string {
	return fmt.Sprintf(enum.SER_NAME + "{UserDump}")
}

// RankDaily:1:20240625
func FormatRankDaily(lv int, d ...int) string {
	date := utils.GetYearMonthDay(utils.GetNowUTC())
	if len(d) > 0 {
		date = d[0]
	}
	return fmt.Sprintf("%s:%d:%d", enum.REDIS_KEY_RANK_DAILY, lv, date)
}

// RankWeekly:1:202447
func FormatRankWeekly(lv int, d ...int) string {
	date := utils.GetYearWeek(utils.GetNowUTC())
	if len(d) > 0 {
		date = d[0]
	}
	return fmt.Sprintf("%s:%d:%d", enum.REDIS_KEY_RANK_WEEKLY, lv, date)
}

// {UserID:1234567890}
func FormatGuildID(guildID any) string {
	return fmt.Sprintf("{GuildID:%d}", guildID)
}
func FormatGuild() string {
	return fmt.Sprintf("%s", enum.REDIS_KEY_GUILD)
}

// Guild_RankDaily:20240625
func FormatGuildRankDaily(d ...int) string {
	date := utils.GetYearMonthDay(utils.GetNowUTC())
	if len(d) > 0 {
		date = d[0]
	}
	return fmt.Sprintf("%s:%d", enum.REDIS_KEY_GUILD_RANK_DAILY, date)
}

// Guild_RankWeekly:202447
func FormatGuildRankWeekly(d ...int) string {
	date := utils.GetYearWeek(utils.GetNowUTC())
	if len(d) > 0 {
		date = d[0]
	}
	return fmt.Sprintf("%s:%d", enum.REDIS_KEY_GUILD_RANK_WEEKLY, date)
}

// Guild_In_RankDaily:20240625
func FormatGuildInRankDaily(guildID uint64, d ...int) string {
	date := utils.GetYearMonthDay(utils.GetNowUTC())
	if len(d) > 0 {
		date = d[0]
	}
	return fmt.Sprintf("%s:%d:%d", enum.REDIS_KEY_GUILD_IN_RANK_DAILY, guildID, date)
}

// Guild_In_RankWeekly:202447
func FormatGuildInRankWeekly(guildID uint64, d ...int) string {
	date := utils.GetYearWeek(utils.GetNowUTC())
	if len(d) > 0 {
		date = d[0]
	}
	return fmt.Sprintf("%s:%d:%d", enum.REDIS_KEY_GUILD_IN_RANK_WEEKLY, guildID, date)
}
