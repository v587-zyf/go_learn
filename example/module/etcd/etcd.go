package etcd

import (
	"context"
	"fmt"
	"github.com/v587-zyf/gc/utils"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

func EtcdDo() {
	cli, err := clientv3.New(clientv3.Config{
		// 连接节点
		Endpoints: []string{"localhost:2379"},
		/// 超时时间
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("connect etcd success")
	defer cli.Close()

	ctxToDo := context.TODO()
	// 设置续期
	rsp, err := cli.Grant(ctxToDo, 5)
	if err != nil {
		panic(err)
	}

	// 存值
	key := "test"
	key = "/logagent/%s/colllect_config"
	ip, err := utils.GetLocalIp()
	if err != nil {
		fmt.Println("get local ip err:", err)
		return
	}
	key = fmt.Sprintf(key, ip)
	val := "123"
	val = `
[
   {
      "path" : "D:/code/golearn/demo/log/1.log",
      "topic" : "1_log"
   },
   {
      "path" : "D:/code/golearn/demo/log/2.log",
      "topic" : "2_log"
   },
   {
      "path" : "D:/code/golearn/demo/log/3.log",
      "topic" : "3_log"
   }
]
`

	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//defer cancel()
	// clientv3.WithLease 这个参数加上 设置的数据会有过期时间
	_, err = cli.Put(ctxToDo, key, val, clientv3.WithLease(rsp.ID))
	if err != nil {
		panic(err)
	}

	//ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	//defer cancel()
	resp, err := cli.Get(ctxToDo, key, clientv3.WithLease(rsp.ID))
	if err != nil {
		panic(err)
	}

	for _, ev := range resp.Kvs {
		fmt.Printf("%s: %s\n", ev.Key, ev.Value)
	}

	// 自动续期
	ch, err := cli.KeepAlive(ctxToDo, rsp.ID)
	if err != nil {
		panic(err)
	}

	for {
		c := <-ch
		fmt.Println(c)
	}

	//ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	//defer cancel()
	//_, err = cli.Delete(ctx, key)
	//if err != nil {
	//	panic(err)
	//}
	//
	//watchKey := "qwe"
	//rch := cli.Watch(context.Background(), watchKey)
	//for wresp := range rch {
	//	for _, ev := range wresp.Events {
	//		fmt.Printf("Type:%s key:%v val:%v\n",
	//			ev.Type, string(ev.Kv.Key), string(ev.Kv.Value))
	//	}
	//}

}
