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
	fnVal, ok := function.(reflect.Value)
	if !ok {
		fnVal = reflect.ValueOf(function)
	}
	fnType := fnVal.Type()
	numIn := fnType.NumIn()
	var inputs []Input
	if numIn > 0 {
		inputs = make([]Input, 0, numIn)
		for i := range numIn {
			input, err := NewInput(c, tid.FromType(fnType.In(i)))
			if err != nil {
				return InvokeInfo{}, err
			}
			inputs = append(inputs, input)
		}
	}

	return InvokeInfo{fnVal: fnVal, inputs: inputs}, nil
}

type InvokeInfo struct {
	fnVal  reflect.Value
	inputs []Input
}

func (i InvokeInfo) Call(s Scope) ([]reflect.Value, error) {
	var in []reflect.Value
	l := len(i.inputs)
	if l > 0 {
		in = make([]reflect.Value, 0, l)
		for _, input := range i.inputs {
			inVal, err := input.Resolve(s)
			if err != nil {
				return nil, err
			}
			in = append(in, inVal)
		}
	}

	return i.fnVal.Call(in), nil
}
