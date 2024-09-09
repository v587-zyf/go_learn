package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

var cli *clientv3.Client

type LogEntry struct {
	Path  string `json:"path"`
	Topic string `json:"topic"`
}

func Init(address string, timeOut time.Duration) (err error) {
	cli, err = clientv3.New(clientv3.Config{
		// 连接节点
		Endpoints: []string{address},
		/// 超时时间
		DialTimeout: timeOut,
	})
	if err != nil {
		panic(err)
	}

	return nil
}

func GetConf(key string) (logEntries []*LogEntry, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp, err := cli.Get(ctx, key)
	if err != nil {
		fmt.Println("get from etcd failed, err:", err)
		return
	}
	for _, ev := range resp.Kvs {
		if err = json.Unmarshal(ev.Value, &logEntries); err != nil {
			fmt.Println("unmarshal failed, err:", err)
			return
		}
	}

	return
}

func WatchConf(key string, newConfCh chan<- []*LogEntry) {
	rch := cli.Watch(context.Background(), key)
	for wresp := range rch {
		for _, ev := range wresp.Events {
			fmt.Printf("Type:%s key:%v val:%v\n",
				ev.Type, string(ev.Kv.Key), string(ev.Kv.Value))
			var newConf []*LogEntry
			if ev.Type != clientv3.EventTypeDelete {
				if err := json.Unmarshal(ev.Kv.Value, &newConf); err != nil {
					fmt.Println("unmarshal failed, err:", err)
					continue
				}
			}
			newConfCh <- newConf
		}
	}
}
