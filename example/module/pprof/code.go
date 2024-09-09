package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/pprof"
	"time"
)

/**
 * go tool pprof cpu.pprof
 * top 5 （查看占用最多的前5个函数）
 * quit 退出

 * 图形工具: graphviz https://graphviz.org/download/
 * go-torch: go get -v github.com/uber/go-torch
 * FlameGraph
 * 下载perl: https://www.perl.org/get.html
 * 下载FlameGraph: git clone https://github.com/brendangregg/FlameGraph.git
 * FlameGraph加入环境变量
 */

// 该方法有问题 c为通道 但未分配 未分配通道不管读还是写会阻塞
func logicCode() {
	var c chan int
	for {
		select {
		// 通道未分配 会一直阻塞
		case v := <-c:
			fmt.Println(v)
		// 没有写跳出循环或者方法 所以什么也不做 导致本循环一直阻塞运行
		default:
			//time.Sleep(500 * time.Millisecond)
		}
	}
}

func main() {
	var isCPUProf bool
	var isMemProf bool

	flag.BoolVar(&isCPUProf, "cpu", false, "return cpu pprof on")
	flag.BoolVar(&isMemProf, "mem", false, "return mem pprof on")
	flag.Parse()

	if isCPUProf {
		f, err := os.Create("./cpu.pprof")
		if err != nil {
			fmt.Println(err)
			return
		}
		pprof.StartCPUProfile(f)
		defer func() {
			pprof.StopCPUProfile()
			f.Close()
		}()
	}

	for i := 0; i < 8; i++ {
		go logicCode()
	}
	time.Sleep(20 * time.Second)

	if isMemProf {
		f, err := os.Create("./mem.pprof")
		if err != nil {
			fmt.Println(err)
			return
		}
		pprof.WriteHeapProfile(f)
		f.Close()
	}
}
