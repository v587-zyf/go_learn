package udp

import (
	"fmt"
	"net"
)

func Server() {
	udpAddr := &net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 8100,
	}
	listen, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("udp listen error:", err)
		return
	}
	defer listen.Close()
	fmt.Println("udp listen...")

	for {
		var data [1024]byte
		n, addr, err := listen.ReadFromUDP(data[:])
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Printf("data:%v addr:%v count%v\n",
			string(data[:n]), addr, n)
		_, err = listen.WriteToUDP(data[:n], addr)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
}
