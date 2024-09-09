package example

import "fmt"

type Queue struct {
	Array    []int
	Head     int
	Rear     int
	Capacity int
}

func NewQueue(capacity int) *Queue {
	return &Queue{
		Head:     -1,
		Rear:     -1,
		Capacity: capacity,
	}
}

// 判断队列是否为空
func (this *Queue) IsEmpty() bool {
	return this.Head == -1
}

// 判断队列是否已满
func (this *Queue) IsFull() bool {
	return (this.Rear+1)%this.Capacity == this.Head
}

// 获取队列长度
func (this *Queue) GetQueueSize() int {
	if this.Head == -1 {
		return 0
	}
	return (this.Rear + 1 - this.Head + this.Capacity) % this.Capacity
}

// 从尾部入队列
func (this *Queue) EnQueue(data int) {
	if this.IsFull() {
		fmt.Println("队列已满")
	} else {
		this.Rear = (this.Rear + 1) % this.Capacity
		this.Array[this.Rear] = data
		if this.Head == -1 {
			this.Head = this.Rear
		}
	}
}

// 从头部取数据
func (this *Queue) DeQueue() int {
	var data int
	if this.IsEmpty() {
		fmt.Println("队列为空")
		return -1
	} else {
		data = this.Array[this.Head]
		if this.Head == this.Rear {
			this.Head = -1
			this.Rear = -1
		} else {
			this.Head = (this.Head + 1) % this.Capacity
		}
		return data
	}
}
