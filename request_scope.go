package di

import (
	"github.com/rushstart/tid"
	"reflect"
)

var _ Scope = (*RequestScope)(nil)

type RequestScope struct {
	container *Container
	instances []instance
}

func (s *RequestScope) sign() int64 {
	return s.container.sign()
}

func (s *RequestScope) Supply(id tid.ID, value reflect.Value) error {
	definition, err := s.container.Definition(id)
	if err != nil {
		return err
	}

	if definition.scope == Singleton {
		return s.container.Supply(id, value)
	}

	s.instances[definition.order] = instance{
		value:   value,
		cleanup: nil,
		isset:   true,
	}
	return nil
}

func (s *RequestScope) Resolve(id tid.ID) (reflect.Value, error) {
	d, err := s.container.Definition(id)
	if err != nil {
		return reflect.Value{}, err
	}

	return s.ResolveDefinition(d)
}

func (s *RequestScope) ResolveDefinition(d Definition) (val reflect.Value, err error) {
	if d.scope == Singleton {
		return s.container.ResolveDefinition(d)
	}

	inst := s.instances[d.order]
	if inst.isset {
		return inst.value, nil
	}

	inst, err = d.resolve(s)
	if err != nil {
		return val, err
	}

	s.instances[d.order] = inst
	return inst.value, nil
}

func (s *RequestScope) Invoke(function any, opts ...InvokeOption) ([]reflect.Value, error) {
	invokeInfo, ok := function.(InvokeInfo)
	if !ok {
		var err error
		fnVal := reflect.ValueOf(function)
		invokeInfo, ok = s.container.invokeInfoCache.Load(tid.FromType(fnVal.Type()))
		if !ok {
			invokeInfo, err = NewInvokeInfo(s.container, fnVal)
			if err != nil {
				return nil, err
			}
			s.container.invokeInfoCache.Store(tid.FromType(fnVal.Type()), invokeInfo)
		}
	}
	return invokeInfo.Call(s)
}

func (s *RequestScope) release() {
	for i := len(s.instances) - 1; i >= 0; i-- {
		if inst := s.instances[i]; inst.cleanup != nil {
			if err := inst.cleanup(); err != nil {
				s.container.logger.Error(err)
			}
		}
	}
	clear(s.instances)
}
