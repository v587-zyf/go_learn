package main

import (
	"bytes"
	"fmt"
	"log"
)

func (mgr *UploadMgr) CmdDo() {
	// 执行远程命令
	var b bytes.Buffer
	mgr.session.Stdout = &b
	if err := mgr.session.Run(mgr.options.cmdInfo); err != nil {
		log.Fatalf("unable to run command: %v", err)
	}

	// 打印命令输出
	fmt.Println(b.String())
}
