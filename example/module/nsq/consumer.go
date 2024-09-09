package nsq

import (
	"fmt"
	"github.com/nsqio/go-nsq"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type MyHandler struct {
	Title string
}

func (m *MyHandler) HandleMessage(msg *nsq.Message) (err error) {
	fmt.Printf("%s recv from:%s msg:%v\n",
		m.Title, msg.NSQDAddress, msg.Body)

	return
}

func initConsumer(topic, channel, address string) (err error) {
	config := nsq.NewConfig()
	config.LookupdPollInterval = 15 * time.Second
	c, err := nsq.NewConsumer(topic, channel, config)
	if err != nil {
		return err
	}
	consumer := &MyHandler{
		Title: "test",
	}
	c.AddHandler(consumer)

	if err = c.ConnectToNSQLookupd(address); err != nil {
		return err
	}

	return nil
}

func ConsumerDo() {
	address := "127.0.0.1:4161"
	if err := initConsumer("topic_demo", "first", address); err != nil {
		fmt.Println("initConsumer err:", err)
		return
	}
	c := make(chan os.Signal)
	// 转发键盘信号到 c
	signal.Notify(c, syscall.SIGINT)
	<-c
}
