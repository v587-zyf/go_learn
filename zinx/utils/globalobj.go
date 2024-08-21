package utils

import (
	"encoding/json"
	"os"
	"zinx/ziface"
)

/**
 * 全局参数
 */
type GlobalObj struct {
	TcpServer ziface.IServer // 当前Server对象
	Host      string         // 当前Server监听IP
	TcpPort   int            // 当前Server监听端口
	Name      string         // 当前Server名称

	Version        string // 版本号
	MaxConn        int    // 当前允许最大连接数
	MaxPackageSize uint32 // 当前数据包最大值

	WorkerPoolSize   uint32 // 当前worker的大小
	MaxWorkerTaskLen uint32 // 每个worker消息队列数量最大值
}

var GlobalObject *GlobalObj

func (g *GlobalObj) Reload() {
	data, err := os.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

func init() {
	GlobalObject = &GlobalObj{
		Name:             "ZinxServerApp",
		Version:          "v0.1",
		TcpPort:          8999,
		Host:             "0.0.0.0",
		MaxConn:          1000,
		MaxPackageSize:   4096,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
	}

	GlobalObject.Reload()
}
