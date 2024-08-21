package main

import (
	"flag"
	"fmt"
	"time"
)

var (
	up         = flag.Bool("up", true, "use upload file")
	uploadType = flag.String("t", "c", "c = client;s = server")
	connType   = flag.String("ct", CONN_TYPE_PASSWORD, "connection pw(password) or pk(public_key)")
	ip         = flag.String("ip", "", "remote ip")
	port       = flag.Int("port", 0, "remote port")
	user       = flag.String("user", "", "remote user")
	pass       = flag.String("pass", "", "pass")
	lPath      = flag.String("lp", "", "local path")
	rPath      = flag.String("rp", "", "remote path")
	cmd        = flag.Bool("c", false, "use cmd")
	cmdInfo    = flag.String("ci", "", "cmd info")
	showHelp   = flag.Bool("help", false, "show help")
	filter     = flag.String("f", "", "filter file or dir")
)

func main() {
	flag.Parse()
	if *up {
		if *lPath == "" {
			panic("must input lPath. can input help show list")
		}
		if *rPath == "" {
			panic("must input rPath. can input help show list")
		}
	}
	if *uploadType == "" {
		panic("must input uploadType. can input help show list")
	}
	if *ip == "" {
		panic("must input ip. can input help show list")
	}
	if *port == 0 {
		panic("must input port. can input help show list")
	}
	if *user == "" {
		panic("must input user. can input help show list")
	}
	if *pass == "" {
		panic("must input pass. can input help show list")
	}
	if *cmd && *cmdInfo == "" {
		panic("must input cmdInfo. can input help show list")
	}

	if *showHelp {
		flag.Usage()
		return
	}

	if err := Init(WithUp(*up), WithUpType(*uploadType), WithSshAddress(fmt.Sprintf("%s:%d", *ip, *port)),
		WithSshUser(*user), WithSshPass(*pass), WithLocalPath(*lPath), WithRemotePath(*rPath), WithCmd(*cmd),
		WithCmdInfo(*cmdInfo), WithFilter(*filter), WithConnType(*connType)); err != nil {
		panic(err)
	}

	mgr.run()

	fmt.Println("Wait 3 Seconds Auto Close!")

	time.Sleep(3 * time.Second)
}
