package cron

import (
	"sync"
	"time"
)

const (
	WORK_TYPE_ONCE = iota // 执行一次
	WORK_TYPE_LOOP        // 轮询执行
)

type TimeCronUnit struct {
	Fn    func()        // 具体要做的方法
	Timer *time.Timer   // 定时器
	Dur   time.Duration // 定时时间
	Type  int           // 类型

	CloseChan chan struct{}
}

func (tc *TimeCronUnit) Start() {
	defer func() {
		if tc.Timer != nil {
			tc.Timer.Stop()
		}
	}()

	tc.Timer = time.NewTimer(tc.Dur)
	for {
		select {
		case <-tc.Timer.C:
			tc.Fn()
			if tc.Type == WORK_TYPE_LOOP {
				tc.Timer.Reset(tc.Dur)
			}
		case <-tc.CloseChan:
			return
		}
	}
}

func (tc *TimeCronUnit) Stop() {
	tc.CloseChan <- struct{}{}
}

func (tc *TimeCronUnit) Reset(d time.Duration) {
	if tc.Timer == nil {
		return
	}

	tc.Dur = d
	tc.Timer.Reset(d)
}

type TimeCron struct {
	CronList map[any]*TimeCronUnit
	Mu       sync.Mutex
}

func NewTimeCron() *TimeCron {
	tw := &TimeCron{
		CronList: make(map[any]*TimeCronUnit),
	}
	return tw
}

func (tc *TimeCron) Start(key any, tcUnit *TimeCronUnit) {
	tc.Mu.Lock()
	defer tc.Mu.Unlock()

	// delete old cron
	cron, has := tc.CronList[key]
	if has {
		cron.Stop()
		delete(tc.CronList, key)
	}

	// do new cron
	tc.CronList[key] = tcUnit

	go tcUnit.Fn()
}

func (tc *TimeCron) StopByKey(key any) {
	tc.Mu.Lock()
	defer tc.Mu.Unlock()

	cron, has := tc.CronList[key]
	if !has {
		return
	}

	cron.Stop()
	delete(tc.CronList, key)
}

func (tc *TimeCron) Close() {
	tc.Mu.Lock()
	defer tc.Mu.Unlock()

	for _, cron := range tc.CronList {
		cron.Stop()
	}
	tc.CronList = make(map[any]*TimeCronUnit)
}

func (tc *TimeCron) ResetByKey(key any, d time.Duration) {
	tc.Mu.Lock()
	defer tc.Mu.Unlock()

	cron, has := tc.CronList[key]
	if !has {
		return
	}
	cron.Reset(d)
}

func (tc *TimeCron) Once() {

}
