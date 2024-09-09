package main

import (
	"fmt"
	"os"
	"runtime/trace"
	"time"
)

func main() {
	for i := 0; i < 5; i++ {
		time.Sleep(time.Second)
		fmt.Println("hello world!")
	}
}

// trace可视化试调gpm
// 1.创建文件
// 2.启动
// 3.停止
func main1() {
	f, err := os.Create("./trace.out")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = trace.Start(f)
	if err != nil {
		panic(err)
	}

	fmt.Println("hello world!")

	trace.Stop()
}
