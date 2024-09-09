package example

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
)

type HelloService struct{}

func (this *HelloService) Hello(request string, reply *string) error {
	*reply = "hello:" + request
	return nil
}

func testRpcService() {
	rpc.RegisterName("HelloService", new(HelloService))

	listener, err := net.Listen("tcp", ":999")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := listener.Accept()
	if err != nil {
		log.Fatal(err)
	}
	rpc.ServeConn(conn)

	fmt.Println("rpc service start")
}

func testRpcClient() {
	client, err := rpc.Dial("tcp", ":999")
	if err != nil {
		log.Fatal(err)
	}

	var reply string
	err = client.Call("HelloService.Hello", "hello", &reply)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(reply)

	fmt.Println("rpc client start")
}
