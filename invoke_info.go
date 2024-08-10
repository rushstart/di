package di

import (
	"github.com/rushstart/tid"
	"reflect"
)

type InvokeOption func(invokeOptions) invokeOptions

type invokeOptions struct {
	scope   scope
	tag     string
	asValue bool
}

func NewInvokeInfo(c *Container, function any, opts ...InvokeOption) (InvokeInfo, error) {
	fnVal := reflect.ValueOf(function)
	fnType := fnVal.Type()
	inputs := make([]Input, 0, fnType.NumIn())
	for i := range fnType.NumIn() {
		input, err := NewInput(c, tid.FromType(fnType.In(i)))
		if err != nil {
			return InvokeInfo{}, err
		}
		inputs = append(inputs, input)
	}
	return InvokeInfo{fnVal: fnVal, inputs: inputs}, nil
}

type InvokeInfo struct {
	fnVal  reflect.Value
	inputs []Input
}

func (i InvokeInfo) Call(s Scope) ([]reflect.Value, error) {
	return []reflect.Value{}, nil
}
