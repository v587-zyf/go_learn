package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"zinx/znet"
)

func main() {
	fmt.Println("client start...")

	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client dial err:", err.Error())
		return
	}

	for {
		//if _, err = conn.Write([]byte("hello v0.1")); err != nil {
		//	fmt.Println("write err:", err.Error())
		//	return
		//}
		//
		//buf := make([]byte, 512)
		//n, err := conn.Read(buf)
		//if err != nil {
		//	fmt.Println("read err:", err.Error())
		//	return
		//}
		//
		//fmt.Println("server call back: ", string(buf[:n]))

		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMsgPackage(0, []byte("v0.5 test msg")))
		if err != nil {
			fmt.Println("pack err:", err)
			return
		}

		_, err = conn.Write(binaryMsg)
		if err != nil {
			fmt.Println("write err:", err)
			return
		}

		// 1.读msg
		binaryHead := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, binaryHead)
		if err != nil {
			fmt.Println("read head err:", err)
			return
		}
		// 2.读msgLen
		msgHead, err := dp.Unpack(binaryHead)
		if err != nil {
			fmt.Println("unpack msgHead err:", err)
			return
		}
		// 3.读msg
		if msgHead.GetMsgLen() > 0 {
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msgHead.GetMsgLen())

			_, err = io.ReadFull(conn, msg.Data)
			if err != nil {
				fmt.Println("read msg data err:", err)
				return
			}

			fmt.Println("---> Recv msgID:", msg.Id, " len:",
				msg.DataLen, " data:", string(msg.Data))
		}

		time.Sleep(time.Second)
	}

}
