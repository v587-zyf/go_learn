package nsq

import (
	"bufio"
	"fmt"
	"github.com/nsqio/go-nsq"
	"os"
	"strings"
)

/**
 * go install github.com/nsqio/go-nsq
 */

var producer *nsq.Producer

func initProducer(str string) (err error) {
	config := nsq.NewConfig()
	producer, err = nsq.NewProducer(str, config)
	if err != nil {
		return err
	}

	return nil
}

func ProducerDo() {
	nsqAddress := "127.0.0.1:4150"
	if err := initProducer(nsqAddress); err != nil {
		fmt.Println("init producer error:", err)
		return
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		data, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("read error:", err)
			continue
		}

		data = strings.TrimSpace(data)
		if strings.ToUpper(data) == "Q" {
			break
		}

		if err = producer.Publish("topic_demo", []byte(data)); err != nil {
			fmt.Println("publish error:", err)
			continue
		}

		fmt.Println("publish success")
	}
}
