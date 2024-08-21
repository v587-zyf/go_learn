package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

func TestDataPack(test *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("listen err:", err)
		return
	}

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("accept err:", err)
				continue
			}

			go func(conn net.Conn) {
				dp := NewDataPack()
				for {
					// 1.读head
					headData := make([]byte, dp.GetHeadLen())
					_, err = io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head err:", err)
						break
					}
					// 2.读dataLen
					msgHead, err := dp.Unpack(headData)
					if err != nil {
						fmt.Println("unpack head err:", err)
						return
					}
					// 3.读data
					if msgHead.GetMsgLen() > 0 {
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen())

						_, err = io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("read data err:", err)
							return
						}
						fmt.Println("---> Recv MsgID:", msg.Id, " len:",
							msg.DataLen, " data:", string(msg.Data))
					}
				}
			}(conn)
		}
	}()

	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client dial err:", err)
		return
	}

	dp := NewDataPack()
	msg1 := &Message{
		Id:      1,
		DataLen: 4,
		Data:    []byte{'z', 'i', 'n', 'x'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("pack msg1 err:", err)
		return
	}

	msg2 := &Message{
		Id:      2,
		DataLen: 7,
		Data:    []byte{'n', 'i', 'h', 'a', 'o', '!', '!'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("pack msg2 err:", err)
		return
	}

	sendData1 = append(sendData1, sendData2...)

	conn.Write(sendData1)

	select {}
}
