package di

import "reflect"

type instance struct {
	value   reflect.Value
	cleanup func() error
	isset   bool
}

func instanceFromOut(outs ...reflect.Value) (instance, error) {
	return instance{}, nil
}
