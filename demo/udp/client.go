package udp

import (
	"fmt"
	"net"
)

func Client() {
	client, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 8100,
	})
	if err != nil {
		fmt.Println("dailUdp err:", err)
		return
	}
	defer client.Close()
	fmt.Println("udp dail success...")

	sendData := []byte("hello world")
	_, err = client.Write(sendData)
	if err != nil {
		fmt.Println("client write err:", err)
		return
	}

	data := make([]byte, 4096)
	n, addr, err := client.ReadFromUDP(data)
	if err != nil {
		fmt.Println("client read err:", err)
		return
	}
	fmt.Printf("recv:%v addr:%v count:%v\n",
		string(data[:n]), addr, n)
}
