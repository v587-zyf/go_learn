package utils

import (
	"strconv"
	"time"
)

func GetNowUTC() time.Time {
	loc, _ := time.LoadLocation("UTC")
	return time.Now().In(loc)
}

// 20240625
func GetYearMonthDay(t time.Time) int {
	date, _ := strconv.Atoi(t.Format("20060102"))
	return date
}

// 202447
func GetYearWeek(t time.Time) int {
	year, week := t.ISOWeek()
	date := year*100 + week
	return date
}

// 用字符串格式化时间
func GetTimeByData(dateStr string) (time.Time, error) {
	loc, _ := time.LoadLocation("Local")
	return time.ParseInLocation("2006-01-02 15:04:05", dateStr, loc)
}

// 获取指定时间到现在多少天
func GetTheDays(startTime time.Time) int {
	serverOpenTime := startTime
	serverOpenZeroTime := GetZeroTime(serverOpenTime).Unix()
	nowTime := time.Now()
	nowZeroTime := GetZeroTime(nowTime).Unix()
	openDays := (nowZeroTime-serverOpenZeroTime)/(24*60*60) + 1
	return int(openDays)
}

// 获取0点时间
func GetZeroTime(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.Local)
}

// 获取传入时间的零点时刻的时间戳
func GetZeroTimeInt64(t time.Time) int64 {
	ts, _ := time.Parse("2006-01-02", t.Format("2006-01-02"))
	return ts.Unix()
}

// 时间转换为int 如(20230101)
func GetDateInt(t time.Time) int {
	y, m, d := t.Date()
	date := y*10000 + int(m)*100 + d
	return date
}

// 获取两个时间差多少小时
func getHourDiffer(startTime, endTime string) int64 {
	var hour int64
	t1, err := time.ParseInLocation("2006-01-02 15:04:05", startTime, time.Local)
	t2, err := time.ParseInLocation("2006-01-02 15:04:05", endTime, time.Local)
	if err == nil && t1.Before(t2) {
		diff := t2.Unix() - t1.Unix()
		hour = diff / 3600
		return hour
	} else {
		return hour
	}
}
