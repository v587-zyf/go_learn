package event

import (
	"kernel/errcode"
	"reflect"
)

type ListenerFn func(...interface{})

const (
	MAX_LISTENER_CNT = 100
)

func NewEventEmitter() *EventEmitter {
	return &EventEmitter{
		events: make(map[string][]ListenerFn),
		max:    1000,
	}
}

func NewPool() *Pool {
	return &Pool{}
}

func GenListener(fn interface{}) (ListenerFn, error) {
	refValue := reflect.ValueOf(fn)

	if refValue.Kind() != reflect.Func {
		return nil, errcode.ERR_EVENT_PARAM_INVALID
	}

	lfn := func(params ...interface{}) {
		paramRefs := make([]reflect.Value, len(params))

		for k, v := range params {
			paramRefs[k] = reflect.ValueOf(v)
		}

		refValue.Call(paramRefs)
	}

	return lfn, nil
}
