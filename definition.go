package di

import (
	"github.com/rushstart/tid"
	"reflect"
)

type DefinitionOption func(definitionOptions) definitionOptions

type definitionOptions struct {
	scope   scope
	tag     string
	asValue bool
}

func WithScope(scope scope) DefinitionOption {
	return func(opts definitionOptions) definitionOptions {
		opts.scope = scope
		return opts
	}
}

func WithTag(t string) DefinitionOption {
	return func(opts definitionOptions) definitionOptions {
		opts.tag = t
		return opts
	}
}

type Definition struct {
	id       tid.ID
	provider reflect.Value
	asValue  bool
	ctor     InvokeInfo
	scope    scope
	order    int
	sign     int64
}

func (d Definition) ValidateBinding(s Scope) error {
	if d.sign != s.sign() {
		return ErrInvalidDefinitionBinding
	}
	return nil
}

func (d Definition) resolve(s Scope) (instance, error) {
	err := d.ValidateBinding(s)
	if err != nil {
		return instance{}, err
	}

	if d.asValue {
		return instance{
			value:   d.provider,
			cleanup: nil,
			isset:   true,
		}, nil
	}

	out, err := d.ctor.Call(s)
	if err != nil {
		return instance{}, err
	}
	return newInstance(out)
}

func Define[T any](constructor any, opts ...DefinitionOption) Definition {
	var options definitionOptions
	for _, opt := range opts {
		options = opt(options)
	}

	return define(tid.From[T](options.tag), reflect.ValueOf(constructor), options)
}

func DefineValue[T any](value T, opts ...DefinitionOption) Definition {
	var options definitionOptions
	for _, opt := range opts {
		options = opt(options)
	}
	options.asValue = true

	return define(tid.From[T](options.tag), reflect.ValueOf(value), options)
}

func define(id tid.ID, provider reflect.Value, opts definitionOptions) Definition {
	return Definition{
		id:       id,
		provider: provider,
		scope:    opts.scope,
		asValue:  opts.asValue,
	}
}
