package workerpool

import (
	"kernel/errcode"
	"kernel/log"
	"runtime"
	"sync"
	"time"

	"go.uber.org/zap"
)

type WorkerPool struct {
	maxWorkerCnt int
	minWorkerCnt int

	idleCleanTime time.Duration

	curWorkerCnt int

	pool sync.Pool

	mu sync.Mutex

	ready []*worker

	stopCh chan struct{}

	mustStop bool
}

func New(cfg ...*Config) (*WorkerPool, error) {
	p := &WorkerPool{
		maxWorkerCnt:  256 * 1024,
		minWorkerCnt:  10,
		idleCleanTime: 5 * 60 * time.Second,

		pool: sync.Pool{
			New: func() interface{} {
				return &worker{
					task: make(chan ITask, workerCap),
				}
			},
		},

		ready: make([]*worker, 0),
	}

	if len(cfg) > 0 {
		wpCfg := cfg[0]
		if wpCfg.MaxCount > 0 {
			p.maxWorkerCnt = int(wpCfg.MaxCount)
		}

	}

	return p, nil
}

var workerCap = func() int {
	// Use blocking workerChan if GOMAXPROCS=1.
	// This immediately switches Serve to WorkerFunc, which results
	// in higher performance (under go1.5 at least).
	if runtime.GOMAXPROCS(0) == 1 {
		return 0
	}

	// Use non-blocking workerChan if GOMAXPROCS>1,
	// since otherwise the Serve caller (Acceptor) may lag accepting
	// new connections if WorkerFunc is CPU-bound.
	return 1
}()

func (p *WorkerPool) Start() {
	if p.stopCh != nil {
		return
	}

	p.stopCh = make(chan struct{})
	stopCh := p.stopCh

	p.mustStop = false

	go func() {
		timer := time.NewTicker(10 * time.Second)
		var scratch []*worker
	LOOP:
		for {
			select {
			case <-stopCh:
				break LOOP
			case <-timer.C:
				p.clean(&scratch)
			}
		}

		timer.Stop()
		timer = nil
	}()
}

func (p *WorkerPool) Stop() {
	if p.stopCh == nil {
		return
	}
	close(p.stopCh)
	p.stopCh = nil

	p.mu.Lock()
	ready := p.ready
	for i := range ready {
		ready[i].task <- nil
		// close(ready[i].task)
		ready[i] = nil
	}
	p.ready = ready[:0]
	p.mustStop = true
	p.mu.Unlock()
}

func (p *WorkerPool) Assign(task ITask) error {
	w := p.getWorker()
	if w == nil {
		log.Error("WorkPool Assign Failed", zap.Int("maxWorkerCnt", p.maxWorkerCnt), zap.Any("task", task))
		return errcode.ERR_WP_TOO_MANY_WORKER
	}
	w.task <- task
	return nil
}

func (p *WorkerPool) clean(scratch *[]*worker) {
	criticalTime := time.Now().Add(-p.idleCleanTime)

	p.mu.Lock()
	ready := p.ready
	n := len(ready)
	if n <= p.minWorkerCnt {
		p.mu.Unlock()
		return
	}

	// Use binary-search algorithm to find out the index of the least recently worker which can be cleaned up.
	l, r, mid := 0, n-1, 0
	for l <= r {
		mid = (l + r) / 2
		if criticalTime.After(p.ready[mid].lastUseTime) {
			l = mid + 1
		} else {
			r = mid - 1
		}
	}
	i := r
	if i == -1 {
		p.mu.Unlock()
		return
	}

	*scratch = append((*scratch)[:0], ready[:i+1]...)
	m := copy(ready, ready[i+1:])
	for i = m; i < n; i++ {
		ready[i] = nil
	}
	p.ready = ready[:m]
	p.mu.Unlock()

	// Notify obsolete workers to stop.
	// This notification must be outside the wp.lock, since ch.ch
	// may be blocking and may consume a lot of time if many workers
	// are located on non-local CPUs.
	tmp := *scratch
	for i := range tmp {
		tmp[i].task <- nil
		tmp[i] = nil
	}
	log.Info("tmp clean", zap.Int("len", len(tmp)))
}

func (p *WorkerPool) getWorker() *worker {
	var w *worker
	canCreate := false

	p.mu.Lock()
	ready := p.ready
	n := len(ready) - 1
	if n < 0 {
		if p.curWorkerCnt < p.maxWorkerCnt {
			canCreate = true
			p.curWorkerCnt++
		}
	} else {
		w = ready[n]
		ready[n] = nil
		p.ready = ready[:n]
	}
	p.mu.Unlock()

	if w == nil {
		if !canCreate {
			return nil
		}
		v := p.pool.Get()
		w = v.(*worker)
		go func() {
			w.run(p)

			p.mu.Lock()
			p.curWorkerCnt--
			p.mu.Unlock()

			p.pool.Put(v)
		}()
	}
	return w
}

func (p *WorkerPool) release(w *worker) bool {
	w.lastUseTime = time.Now()
	p.mu.Lock()
	if p.mustStop {
		p.mu.Unlock()
		return false
	}
	p.ready = append(p.ready, w)
	p.mu.Unlock()
	return true
}
