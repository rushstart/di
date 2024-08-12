package di

import (
	"fmt"
	"reflect"
)

type instance struct {
	value   reflect.Value
	cleanup func() error
	isset   bool
}

func newInstance(outs []reflect.Value) (instance, error) {
	l := len(outs)
	if l < 1 || l > 3 {
		return instance{}, fmt.Errorf("DI: invalid number of out: got %d, want 1, 2 or 3", l)
	}
	inst := instance{value: outs[0]}
	if l >= 2 {
		err, ok := outs[l-1].Interface().(error)
		if ok {
			return instance{}, err
		}
		if !outs[1].IsNil() {
			cleanup, ok := outs[1].Interface().(func() error)
			if !ok {
				return instance{}, fmt.Errorf("DI: invalid type of constructor: got %T, want func() error", outs[1].Interface())
			}
			inst.cleanup = cleanup
		}
	}
	inst.isset = true
	return inst, nil
}
