package cron

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/v587-zyf/gc/module"
	"math/rand"
	"runtime/debug"
)

type CronJob struct {
	spec   string       // 时间格式
	worker func() error // 要操作的方法
	name   string       // 任务名称
}

type Cron struct {
	module.DefModule
}

func NewCron() *Cron {
	return &Cron{}
}

func (c *Cron) Init() error {
	//randSecond := rand.Intn(90) % 60 // 随机秒
	randMinute := rand.Intn(5) + 15 // 随机分钟

	jobs := []CronJob{
		{"0 0 0 * * ?", c.DailyMidnight, "DailyMidnight"},                     // 每日0点
		{"0 0 5 * * ?", c.DailyFive, "DailyFive"},                             // 每日5点
		{"0 0 */1 * * ?", c.EveryOneHour, "EveryOneHour"},                     // 每1小时
		{"0 */1 0 * * ?", c.EveryOneMinute, "EveryOneMinute"},                 // 每1分钟
		{"0 0 1 1 * ?", c.EveryMonthFirst, "EveryMonthFirst"},                 // 每月1号1点
		{"0 26,29,33 * * * ?", c.ManyMinute, "ManyMinute"},                    // 多个分钟定时
		{"0 0 0,13,18,21 * * ?", c.ManyOClick, "ManyOClick"},                  // 多个定点
		{fmt.Sprintf("0 %d 0 * * ?", randMinute), c.DiyMinutes, "DiyMinutes"}, // 15-20分钟
	}

	secondParser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)
	cronTab := cron.New(cron.WithParser(secondParser), cron.WithChain())
	for _, job := range jobs {
		if _, err := cronTab.AddFunc(job.spec, c.wrap(job.name, job.worker)); err != nil {
			return err
		}
	}

	cronTab.Start()
	return nil
}

func (c *Cron) wrap(name string, f func() error) func() {
	return func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("cronjob %s panic: %v, %s \n", name, err, debug.Stack())
			}
		}()
		if err := f(); err != nil {
			fmt.Printf("cronjob %s error: %s \n", name, err.Error())
		}
	}
}

func (c *Cron) DailyMidnight() error {

	return nil
}

func (c *Cron) DailyFive() error {

	return nil
}

func (c *Cron) EveryOneHour() error {

	return nil
}

func (c *Cron) EveryOneMinute() error {

	return nil
}

func (c *Cron) EveryMonthFirst() error {

	return nil
}

func (c *Cron) ManyMinute() error {

	return nil
}

func (c *Cron) ManyOClick() error {

	return nil
}

func (c *Cron) DiyMinutes() error {

	return nil
}
