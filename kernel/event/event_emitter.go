package event

import (
	"kernel/errcode"
	"kernel/log"
	"reflect"
	"strings"
	"sync"

	"go.uber.org/zap"
)

type EventEmitter struct {
	events map[string][]ListenerFn

	mu sync.RWMutex

	max int
}

func (e *EventEmitter) On(eventName string, fns ...ListenerFn) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	v, ok := e.events[eventName]
	if !ok {
		v = make([]ListenerFn, 0, 10)
	}

	if len(v)+len(fns) > e.max {
		return errcode.ERR_EVENT_LISTENER_LIMIT
	}

	e.events[eventName] = append(v, fns...)

	return nil
}

func (e *EventEmitter) Off(eventName string, fn ListenerFn) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	fns, ok := e.events[eventName]
	if !ok {
		return errcode.ERR_EVENT_LISTENER_EMPTY
	}

	var idx int = -1
	fnp := reflect.ValueOf(fn).Pointer()
	for k, v := range fns {
		vp := reflect.ValueOf(v).Pointer()
		if fnp == vp {
			idx = k
			break
		}
	}
	if idx < 0 {
		return errcode.ERR_EVENT_LISTENER_NOT_FIND
	}

	newListerners := make([]ListenerFn, 0, len(fns)-1)

	newListerners = append(newListerners, fns[0:idx]...)
	newListerners = append(newListerners, fns[idx+1:]...)

	e.events[eventName] = newListerners
	return nil
}

func (e *EventEmitter) Once(eventName string, fns ...ListenerFn) error {
	v, ok := e.events[eventName]
	if !ok {
		v = make([]ListenerFn, 0, 10)
	}

	if len(v)+len(fns) > e.max {
		return errcode.ERR_EVENT_LISTENER_LIMIT
	}

	wrapFns := make([]ListenerFn, 0, len(fns))
	for _, fn := range fns {
		var wrapFn ListenerFn
		wrapFn = func(params ...interface{}) {
			fn(params...)
			e.Off(eventName, wrapFn)
		}
		wrapFns = append(wrapFns, wrapFn)
	}

	e.On(eventName, wrapFns...)

	return nil
}

func (e *EventEmitter) Emit(eventName string, params ...interface{}) error {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				log.Error("event emit err", zap.String("eventName", eventName), zap.String("err", err.Error()))
			} else if err, ok := r.(string); ok {
				if strings.HasPrefix(err, "reflect") {
					err = "Emit" + err[7:]
				}
				log.Error("event emit err", zap.String("eventName", eventName), zap.String("err", err))
			} else {
				log.Error("event emit err", zap.String("eventName", eventName), zap.Reflect("err", err))
			}
		}
	}()

	e.mu.RLock()
	v, ok := e.events[eventName]
	if !ok {
		return errcode.ERR_EVENT_LISTENER_EMPTY
	}
	listeners := make([]ListenerFn, len(v))
	copy(listeners, v)
	e.mu.RUnlock()

	for _, v := range listeners {
		v(params...)
	}

	return nil
}
