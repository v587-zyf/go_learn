package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/rpc"
	"sync"
	"time"
)

// 1.实现3节点选举
// 2.改造代码成分布式选举代码 加rpc
// 3.演示完整代码 自主选举 日志复制

const (
	raftCount = 3
)

type Leader struct {
	Term     int // 任期
	LeaderId int // 编号
}

type Raft struct {
	mu              sync.Mutex // 锁
	me              int        // 节点编号
	currentTerm     int        // 当前任期
	votedFor        int        // 为那个节点投票
	state           int        // 状态 0 follower 1 candidate 2 leader
	lastMessageTime int64      // 最后发送消息时间
	currentLeader   int        // 当前leader
	message         chan bool  // 消息通道
	electCh         chan bool  // 选举通道
	heartBeat       chan bool  // 心跳通道
	heartbeatRe     chan bool  // 返回心跳通道
	timeout         int        // 超时时间
}

// 0 还没上任 -1 没有编号
var leader = Leader{0, -1}

func main() {
	// 1.3个节点 最初都是follower
	// 2.若有candidate 进行拉票
	// 3.产生leader

	for i := 0; i < raftCount; i++ {
		// 创建节点
		MakeNode(i)
	}

	// 加入服务端
	rpc.Register(new(Raft))
	rpc.HandleHTTP()
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}

	for {
	}
}

func MakeNode(me int) *Raft {
	r := &Raft{
		me:              me,
		votedFor:        -1, // -1谁都不投
		state:           0,
		lastMessageTime: 0,
		currentLeader:   -1,
		message:         make(chan bool),
		electCh:         make(chan bool),
		heartBeat:       make(chan bool),
		heartbeatRe:     make(chan bool),
		timeout:         0,
	}
	r.setTerm(0)

	// 随机种子
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// 选举协程
	go r.election()
	// 心跳检查协程
	go r.sendLeaderHeartbeat()

	return r
}

func (r *Raft) setTerm(term int) {
	r.currentTerm = term
}

func (r *Raft) election() {
	// 标记 是否产生leader
	var result bool
	for {
		// 设置超时
		timeout := r.randRange(150, 300)
		r.lastMessageTime = r.millisecond()

		select {
		// 延迟等待
		case <-time.After(time.Duration(timeout) * time.Millisecond):
			fmt.Println("当前节点状态为:", r.state)

		}

		result = false
		for !result {
			// 选举
			result = r.electionOneRound(&leader)
		}
	}
}

func (r *Raft) randRange(min, max int64) int64 {
	return rand.Int63n(max-min) + min
}

func (r *Raft) millisecond() int64 {
	return time.Now().Unix() / int64(time.Millisecond)
}

func (r *Raft) electionOneRound(leader *Leader) bool {
	var timeout int64
	timeout = 100

	// 票数
	var vote int
	// 是否开始心跳信号检查
	var triggerHeartbeat bool
	// 时间
	last := r.millisecond()
	// 返回值
	success := false

	// 当前节点变candidate
	r.mu.Lock()
	r.becomeCandidate()
	r.mu.Unlock()
	fmt.Println("start election leader")

	for {
		// 遍历所有节点拉选票
		for i := 0; i < raftCount; i++ {
			if i != r.me {
				go func() {
					if leader.LeaderId < 0 {
						r.electCh <- true
					}
				}()
			}
		}
		// 设置投票数量
		vote = 1
		// 遍历节点
		for i := 0; i < raftCount; i++ {
			// 计算投票数量
			select {
			case ok := <-r.electCh:
				if ok {
					vote++
					// 大于节点数 / 2 则成为leader
					success = vote > raftCount/2
					if success && !triggerHeartbeat {
						// 开始心跳检测
						triggerHeartbeat = true
						// 变成leader
						r.mu.Lock()
						r.becomeLeader()
						r.mu.Unlock()
						// 向其他节点发心跳
						r.heartBeat <- true
						fmt.Println(r.me, " 号节点为leader")
						fmt.Println("leader开始发送信号")
					}
				}
			}
		}
		// 校验 不超时且票数大于一半 选举成功 break
		if timeout+last < r.millisecond() ||
			vote > raftCount/2 ||
			r.currentLeader > -1 {
			break
		} else {
			// 等待操作
			select {
			case <-time.After(time.Duration(10) * time.Millisecond):
			}
		}
	}

	return success
}

func (r *Raft) becomeCandidate() {
	r.state = 1
	r.setTerm(r.currentTerm + 1)
	r.votedFor = r.me
	r.currentLeader = -1
}

func (r *Raft) becomeLeader() {
	r.state = 2
	r.currentLeader = r.me
}

// leader 发送心跳 数据同步 看follower是否有回应
func (r *Raft) sendLeaderHeartbeat() {
	for {
		select {
		case <-r.heartBeat:
			r.sendAppendEntriesImpl()
		}
	}
}

// 返回leader的确认信号
func (r *Raft) sendAppendEntriesImpl() {
	// leader不用管
	if r.currentLeader == r.me {
		// leader
		// 节点个数
		var successCount = 0
		for i := 0; i < raftCount; i++ {
			if i != r.me {
				go func() {
					//r.heartbeatRe <- true
					rpc, err := rpc.DialHTTP("tcp", "127.0.0.1:8080")
					if err != nil {
						log.Fatal(err)
					}
					// 接收服务器返回信息
					var ok = false
					err = rpc.Call("Raft.Communication",
						Param{Msg: "hello"}, &ok)
					if err != nil {
						log.Fatal(err)
					}
					if ok {
						r.heartbeatRe <- true
					}
				}()
			}
		}
		// 确认信号个数
		for i := 0; i < raftCount; i++ {
			select {
			case ok := <-r.heartbeatRe:
				if ok {
					successCount++
					if successCount > raftCount/2 {
						fmt.Println("选举成功 信号ok")
						log.Fatal("程序结束")
					}
				}
			}
		}
	} else {
		// follower
	}
}

// 首字母大写 RPC规范 用于分布式通信
type Param struct {
	Msg string
}

// 通信方法
func (r *Raft) Communication(p Param, a *bool) error {
	fmt.Println(p.Msg)
	*a = true

	return nil
}
