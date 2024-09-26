package main

import (
	"flag"
	"fmt"
	"genProto/genProtoMsg/conf"
	"genProto/genProtoMsg/gen"
	"os"
)

var (
	source   = flag.String("s", "", "pb存放目录")
	out      = flag.String("o", "", "生成目录")
	showHelp = flag.Bool("h", false, "show help")
)

func main() {
	flag.Parse()
	if *showHelp {
		flag.Usage()
		return
	}
	if *source == "" {
		fmt.Println("source is nil")
		return
	}
	if *out == "" {
		fmt.Println("out is nil")
		return
	}
	//fmt.Println("source dir:", *source)
	//fmt.Println("out dir:", *out)
	conf.Init(*source, *out)

	os.Mkdir(*out, 0755)

	genSlice := map[string]func(string){
		*out + "msgMap.go": gen.GenMsgMap,
	}
	for k, v := range genSlice {
		os.Create(k)
		v(k)
	}
}
