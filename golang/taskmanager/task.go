package taskmanager

import (
	"fmt"
	"reflect"
	"time"
)

type Task interface {
	Call() (out interface{}, err error)
	StartTime() time.Time
}

type task struct {
	Fn    interface{}
	Args  []interface{}
	Start time.Time
}

func NewTask(start time.Time, fn interface{}, args ...interface{}) Task {
	return &task{Fn: fn, Args: args, Start: start}
}

func (t *task) Call() (interface{}, error) {
	fn := reflect.ValueOf(t.Fn)
	if k := fn.Type().Kind(); k != reflect.Func {
		return nil,
			fmt.Errorf("%s: expecting type Func found %s",
				fn.Type().Name(), k.String())
	}
	if len(t.Args) != fn.Type().NumIn() {
		return nil,
			fmt.Errorf("%s: expecting %d arguments found %d",
				fn.Type().Name(), fn.Type().NumIn(), len(t.Args))
	}
	in := make([]reflect.Value, len(t.Args))
	for i, a := range t.Args {
		in[i] = reflect.ValueOf(a)
	}
	return fn.Call(in), nil
}

func (t *task) StartTime() time.Time {
	return t.Start
}
