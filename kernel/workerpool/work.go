package workerpool

import (
	"kernel/log"
	"runtime"
	"time"

	"go.uber.org/zap"
)

type worker struct {
	task        chan ITask
	lastUseTime time.Time
}

func (w *worker) run(p *WorkerPool) {
	defer func() {
		if r := recover(); r != nil {
			var errField zap.Field
			if err, ok := r.(error); ok {
				errField = zap.String("err", err.Error())
			} else if err, ok := r.(string); ok {
				errField = zap.String("err", err)
			} else {
				errField = zap.Any("err", r)
			}

			stackBuf := make([]byte, 512)
			stackLen := runtime.Stack(stackBuf, true)

			log.Error("worker panic", errField, zap.ByteString("stack", stackBuf[:stackLen]))
			return
		}
	}()

LOOP:
	for {
		select {
		case data := <-w.task:
			if data == nil {
				break LOOP
			}
			data.Do()
		}

		if !p.release(w) {
			break
		}
	}
}
