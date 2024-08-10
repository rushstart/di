package di

import (
	"github.com/rushstart/tid"
	"reflect"
)

type scope int

const (
	Singleton scope = iota
	Request
)

type Scope interface {
	Supply(tid.ID, reflect.Value) error
	Resolve(tid.ID) (reflect.Value, error)
	ResolveDefinition(definition Definition) (reflect.Value, error)
	Invoke(function any, opts ...InvokeOption) ([]reflect.Value, error)
	sign() int64
}
