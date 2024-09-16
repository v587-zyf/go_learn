package cron

import (
	"github.com/robfig/cron/v3"
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
	"runtime/debug"
)

type CronJob struct {
	spec   string
	worker func() error
	name   string
}

func NewCron() *cron.Cron {
	diyJobs := []CronJob{
		{"0 0 0 * * *", daily, "DailyZero"},           // 定时每天0点执行
		{"0 */1 * * * ?", oneMinute, "OneMinute"},     // 1m
		{"55 59 23 * * * ", dailyZero, "dailyZero"},   // 每日 23.59
		{"59 59 23 * * 0", sundayZero, "SundayZero"},  // 每周日 23.59
		{"0/3 * * * * ?", threeSecond, "ThreeSecond"}, // 3s
		{"0/1 * * * * ?", oneSecond, "OneSecond"},     // 1s

		//{"*/10 * * * * ?", fiveSecond, "fiveSecond"}, // 5s
		//{"0 0 5 * * ?", Hour5Reset, "Hour5Reset"},       // 定时每天5点重围
		//{"0 0 23 * * ?", Hour23Reset, "Hour23Reset"},    // 定时每天23点
		//{"0 0 */1 * * ?", OneHour, "OneHour"},           // 1h
	}

	// 设置时区为UTC+0，也就是零时区
	//loc, _ := time.LoadLocation("UTC")
	c := cron.New(cron.WithSeconds())
	//c := cron.New()
	for _, job := range diyJobs {
		_, err := c.AddFunc(job.spec, wrap(job.name, job.worker))
		if err != nil {
			log.Error("add cron job err", zap.String("name", job.name), zap.Error(err))
			continue
		}
	}

	return c
}

func wrap(name string, f func() error) func() {
	return func() {
		defer func() {
			if err := recover(); err != nil {
				log.Error("job panic", zap.String("name", name),
					zap.Any("err", err), zap.ByteString("stack", debug.Stack()))
			}
		}()
		if err := f(); err != nil {
			log.Error("job err", zap.String("name", name), zap.Error(err))
		}
		//log.Info("job do succ", zap.String("name", name))
	}
}
