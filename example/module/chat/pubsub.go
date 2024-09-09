package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type (
	subscriber chan interface{}         // 订阅者管道
	topicFunc  func(v interface{}) bool // 主题过滤器
)

// 发布者对象
type Publisher struct {
	m           sync.RWMutex             // 读写锁
	buffer      int                      // 订阅队列缓存大小
	timeout     time.Duration            // 发布超时时间
	subscribers map[subscriber]topicFunc // 订阅者信息
}

func testPubSub() {
	p := NewPublisher(100*time.Millisecond, 10)
	defer p.Close()

	all := p.Subscribe()
	golang := p.SubscribeTopic(func(v interface{}) bool {
		if s, ok := v.(string); ok {
			return strings.Contains(s, "golang")
		}
		return false
	})

	p.Publish("hello world")
	p.Publish("hello golang")

	go func() {
		for msg := range all {
			fmt.Println("all:", msg)
		}
	}()

	go func() {
		for msg := range golang {
			fmt.Println("golang:", msg)
		}
	}()

	time.Sleep(3 * time.Second)
}

// 构建发布者对象 设置发布超时时间和缓存队列长度
func NewPublisher(publishTimeout time.Duration, buffer int) *Publisher {
	return &Publisher{
		buffer:      buffer,
		timeout:     publishTimeout,
		subscribers: make(map[subscriber]topicFunc),
	}
}

// 添加订阅者 订阅全部主题
func (this *Publisher) Subscribe() chan interface{} {
	return this.SubscribeTopic(nil)
}

// 添加新订阅者 订阅过滤器筛选后的主题
func (this *Publisher) SubscribeTopic(topic topicFunc) chan interface{} {
	ch := make(chan interface{}, this.buffer)
	this.m.Lock()
	this.subscribers[ch] = topic
	this.m.Unlock()
	return ch
}

// 退出订阅
func (this *Publisher) Evict(sub chan interface{}) {
	this.m.Lock()
	defer this.m.Unlock()

	delete(this.subscribers, sub)
	close(sub)
}

// 发布一个主题
func (this *Publisher) Publish(v interface{}) {
	this.m.RLock()
	defer this.m.RUnlock()

	var wg sync.WaitGroup
	for sub, topic := range this.subscribers {
		wg.Add(1)
		go this.sendTopic(sub, topic, v, &wg)
	}
	wg.Wait()
}

// 发送主题（可以容忍一定超时）
func (this *Publisher) sendTopic(sub subscriber, topic topicFunc, v interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	if topic != nil && !topic(v) {
		return
	}

	select {
	case sub <- v:
	case <-time.After(this.timeout):
	}
}

// 关闭发布者对象 同时关闭所有订阅者管道
func (this *Publisher) Close() {
	this.m.Lock()
	defer this.m.Unlock()

	for sub := range this.subscribers {
		delete(this.subscribers, sub)
		close(sub)
	}
}
