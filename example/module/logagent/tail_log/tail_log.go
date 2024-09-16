package tail_log

import (
	"context"
	"example/module/logagent/kafka"
	"fmt"
	"github.com/hpcloud/tail"
)

/**
 * go install github.com/hpcloud/tail
 */

type TailTask struct {
	Path     string
	Topic    string
	Instance *tail.Tail

	ctx    context.Context
	cancel context.CancelFunc
}

func NewTailTask(path, topic string) (*TailTask, error) {
	ctx, cancel := context.WithCancel(context.Background())
	t := &TailTask{
		Path:  path,
		Topic: topic,

		ctx:    ctx,
		cancel: cancel,
	}
	if err := t.Init(); err != nil {
		fmt.Println("tail task init err:", err)
		return nil, err
	}

	return t, nil
}

func (t *TailTask) Init() (err error) {
	config := tail.Config{
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, // 从哪个地方开始读
		ReOpen:    true,                                 // 重新打开
		MustExist: false,                                // 文件不存在是否报错
		Follow:    true,                                 // 跟随文件
		Poll:      true,
	}
	t.Instance, err = tail.TailFile(t.Path, config)
	if err != nil {
		fmt.Println("tail file err:", err)
		return
	}
	go t.Run()

	return
}

func (t *TailTask) Run() {
	for {
		select {
		case line := <-t.Instance.Lines:
			//kafka.SendToKafka(t.Topic, line.Text)
			//fmt.Println(line.Text)
			kafka.SendToChan(t.Topic, line.Text)
		case <-t.ctx.Done():
			fmt.Printf("tail task %s_%s exit\n", t.Path, t.Topic)
			return
		}
	}
}

func (t *TailTask) Close() {
	t.cancel()
	t.Instance.Stop()
}

// example
//var (
//	tailObj *tail.Tail
//)
//func Init(fileName string) (err error) {
//	config := tail.Config{
//		Location: &tail.SeekInfo{
//			Offset: 0,
//			Whence: 2,
//		}, // 从哪个地方开始读
//		ReOpen:    true,  // 重新打开
//		MustExist: false, // 文件不存在是否报错
//		Follow:    true,  // 跟随文件
//		Poll:      false,
//	}
//	tailObj, err = tail.TailFile(fileName, config)
//	if err != nil {
//		fmt.Println("tail file err:", err)
//		return
//	}
//
//	return
//}
//func TailLogDo() {
//	fileName := "./log/my.log"
//	config := tail.Config{
//		Location: &tail.SeekInfo{
//			Offset: 0,
//			Whence: 2,
//		}, // 从哪个地方开始读
//		ReOpen:    true,  // 重新打开
//		MustExist: false, // 文件不存在是否报错
//		Follow:    true,  // 跟随文件
//		Poll:      false,
//	}
//	tails, err := tail.TailFile(fileName, config)
//	if err != nil {
//		fmt.Println("tail file err:", err)
//		return
//	}
//	var (
//		line *tail.Line
//		ok   bool
//	)
//	for {
//		line, ok = <-tails.Lines
//		if !ok {
//			fmt.Printf("tail file close reopen. fileName:%s\n",
//				tailObj.Filename)
//			time.Sleep(time.Second)
//			continue
//		}
//		fmt.Println(line.Text)
//	}
//}
