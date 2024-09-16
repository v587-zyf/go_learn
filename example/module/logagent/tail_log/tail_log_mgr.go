package tail_log

import (
	"example/module/logagent/etcd"
	"fmt"
	"time"
)

var tskMgr *TailLogMgr

type TailLogMgr struct {
	LogEntryList []*etcd.LogEntry
	TskMap       map[string]*TailTask
	NewConfChan  chan []*etcd.LogEntry
}

func Init(logEntries []*etcd.LogEntry) {
	tskMgr = &TailLogMgr{
		LogEntryList: logEntries,
		TskMap:       make(map[string]*TailTask, 16),
		NewConfChan:  make(chan []*etcd.LogEntry),
	}

	for _, entry := range logEntries {
		tsk, err := NewTailTask(entry.Path, entry.Topic)
		if err != nil {
			fmt.Println("tail log err:", err)
			continue
		}
		key := fmt.Sprintf("%s_%s", entry.Path, entry.Topic)
		tskMgr.TskMap[key] = tsk
	}

	go tskMgr.Run()
}

func GetLogMgr() *TailLogMgr {
	return tskMgr
}

func GetNewConfChan() chan []*etcd.LogEntry {
	return tskMgr.NewConfChan
}

func (mgr *TailLogMgr) Run() {
	for {
		select {
		case newConf := <-mgr.NewConfChan:
			for _, entry := range newConf {
				key := fmt.Sprintf("%s_%s", entry.Path, entry.Topic)
				if _, ok := mgr.TskMap[key]; ok {
					continue
				}

				tsk, err := NewTailTask(entry.Path, entry.Topic)
				if err != nil {
					fmt.Println("tail log err:", err)
					continue
				}
				tskMgr.TskMap[key] = tsk
			}

			for _, tEntry := range mgr.LogEntryList {
				isDel := true
				for _, nEntry := range newConf {
					if tEntry.Path == nEntry.Path &&
						tEntry.Topic == nEntry.Topic {
						isDel = false
						continue
					}
				}
				if isDel {
					key := fmt.Sprintf("%s_%s", tEntry.Path, tEntry.Topic)
					tskMgr.TskMap[key].Close()
					delete(tskMgr.TskMap, key)
				}
			}

			//fmt.Println("---newConf:", newConf)
			//mgr.handleNewConf(newConf)
		default:
			time.Sleep(time.Second)
		}
	}
}

func (mgr *TailLogMgr) handleNewConf(newConf []*etcd.LogEntry) {
	for _, entry := range newConf {
		_, err := NewTailTask(entry.Path, entry.Topic)
		if err != nil {
			fmt.Println("tail log err:", err)
			continue
		}
	}
}
