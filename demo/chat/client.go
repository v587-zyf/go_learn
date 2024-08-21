package main

import (
	"flag"
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", client.ServerIp, client.ServerPort))
	if err != nil {
		fmt.Println("net.Dial err:", err)
		return nil
	}

	client.conn = conn

	return client
}

func (c *Client) Menu() bool {
	var flag int

	fmt.Println("0.exit")
	fmt.Println("1.public chat")
	fmt.Println("2.private chat")
	fmt.Println("3.rename")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		c.flag = flag
		return true
	} else {
		fmt.Println("input error")
		return false
	}
}

func (c *Client) UpdateName() bool {
	fmt.Println("input your name:")
	fmt.Scanln(&c.Name)

	sendMsg := "rename|" + c.Name
	if _, err := c.conn.Write([]byte(sendMsg)); err != nil {
		fmt.Println("c.conn.Write err:", err)
		return false
	}

	return true
}

func (c *Client) PublicChat() {
	var chatMsg string

	fmt.Println("input msg. input 'exit' to exit")
	fmt.Scanln(&chatMsg)

	for chatMsg != "exit" {
		// 输入不为空发送给服务器
		if len(chatMsg) != 0 {
			sendMsg := chatMsg
			if _, err := c.conn.Write([]byte(sendMsg)); err != nil {
				fmt.Println("c.conn.Write err:", err)
				break
			}
		}

		chatMsg = ""
		fmt.Println("input msg. input 'exit' to exit")
		fmt.Scanln(&chatMsg)
	}
}

func (c *Client) GetAllUser() {
	sendMsg := "who"
	if _, err := c.conn.Write([]byte(sendMsg)); err != nil {
		fmt.Println("c.conn.Write err:", err)
		return
	}
}

func (c *Client) PrivateChat() {
	var remoteName string

	c.GetAllUser()
	fmt.Println("input remote name. input 'exit' to exit")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		var chatMsg string
		fmt.Println("input msg. input 'exit' to exit")
		fmt.Scanln(&chatMsg)

		for {
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg
				if _, err := c.conn.Write([]byte(sendMsg)); err != nil {
					fmt.Println("c.conn.Write err:", err)
					break
				}
			}
			chatMsg = ""
			fmt.Println("input msg. input 'exit' to exit")
			fmt.Scanln(&chatMsg)
		}
	}
}

func (c *Client) Run() {
	for c.flag != 0 {
		for !c.Menu() {

		}
		switch c.flag {
		case 0:
			fmt.Println("exit")
		case 1:
			c.PublicChat()
			//fmt.Println("public chat")
		case 2:
			c.PrivateChat()
			//fmt.Println("private chat")
		case 3:
			c.UpdateName()
			//fmt.Println("rename")
		}
	}
}

func (c *Client) DealResponse() {
	// 下面一段可简化成这样
	//io.Copy(os.Stdout, c.conn)

	for {
		buf := make([]byte, 4096)
		n, err := c.conn.Read(buf)
		if err != nil {
			fmt.Println("c.conn.Read err:", err)
			continue
		}
		fmt.Println(string(buf[:n-1]))
	}
}

var (
	serverIp   string
	serverPort int
)

func init() {
	flag.StringVar(&serverIp, "ip", "localhost", "server ip addr")
	flag.IntVar(&serverPort, "port", 8888, "server port")
}

func main() {
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println("client link fail")
		return
	}

	fmt.Println("client link success")

	go client.DealResponse()

	client.Run()
}
