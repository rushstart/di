package di

import (
	"errors"
	"fmt"
	"github.com/rushstart/tid"
	"reflect"
	"sync"
	"sync/atomic"
)

var _ Scope = (*Container)(nil)

// New constructs a Container.
func New() (*Container, func()) {
	c := &Container{
		definitions: make(map[tid.ID]Definition),
	}

	return c, func() { c.cleanup() }
}

type Container struct {
	id          int64
	loaded      atomic.Bool
	mu          sync.Mutex
	definitions map[tid.ID]Definition
	logger      Logger
	instances   []instance
	// order of definitions in each the scope
	order [2][]tid.ID
	// request scope pool
	pool sync.Pool

	invokeInfoCache syncMap[tid.ID, InvokeInfo]
}

func (c *Container) sign() int64 {
	return c.id
}

// Load is called only once
// after loading and before cleanup, the container becomes immutable
func (c *Container) Load() {
	if !c.loaded.Load() {
		c.mu.Lock()
		if !c.loaded.Load() {
			c.doLoad()
			c.loaded.Store(true)
		}
		c.mu.Unlock()
	}
}

func (c *Container) doLoad() {
	c.pool = sync.Pool{New: func() any {
		return &RequestScope{
			container: c,
			instances: make([]instance, len(c.order[Request])),
		}
	}}

	c.instances = make([]instance, len(c.order[Singleton]))
	for _, id := range c.order[Singleton] {
		_, err := c.Resolve(id)
		if err != nil {
			panic(err)
		}
	}
}

func (c *Container) Supply(id tid.ID, value reflect.Value) error {
	return errors.New("DI: not allowed set a value in a singleton scope")
}

func (c *Container) Resolve(id tid.ID) (reflect.Value, error) {
	d, err := c.Definition(id)
	if err != nil {
		return reflect.Value{}, err
	}

	return c.ResolveDefinition(d)
}

func (c *Container) ResolveDefinition(d Definition) (val reflect.Value, err error) {
	if c.instances == nil {
		return val, errors.New("DI: container must be loaded before use")
	}

	inst := c.instances[d.order]
	if inst.isset {
		return inst.value, nil
	}

	if c.loaded.Load() {
		return val, errors.New("DI: container already loaded")
	}

	inst, err = d.resolve(c)
	if err != nil {
		return val, err
	}

	c.instances[d.order] = inst
	return inst.value, nil
}

func (c *Container) Invoke(function any, opts ...InvokeOption) ([]reflect.Value, error) {
	invokeInfo, ok := function.(InvokeInfo)
	if !ok {
		var err error
		fnVal := reflect.ValueOf(function)
		invokeInfo, ok = c.invokeInfoCache.Load(tid.FromType(fnVal.Type()))
		if !ok {
			invokeInfo, err = NewInvokeInfo(c, fnVal)
			if err != nil {
				return nil, err
			}
			c.invokeInfoCache.Store(tid.FromType(fnVal.Type()), invokeInfo)
		}
	}
	return invokeInfo.Call(c)
}

func (c *Container) cleanup() {
	for i := len(c.instances) - 1; i >= 0; i-- {
		if inst := c.instances[i]; inst.cleanup != nil {
			if err := inst.cleanup(); err != nil {
				c.logger.Error(err)
			}
		}
	}
	clear(c.instances)
}

func (c *Container) AcquireRequestScope() *RequestScope {
	if !c.loaded.Load() {
		panic("DI: container must be loaded before use")
	}
	return c.pool.Get().(*RequestScope)
}

func (c *Container) ReleaseRequestScope(s *RequestScope) {
	s.release()
	c.pool.Put(s)
}

func (c *Container) Definition(id tid.ID) (Definition, error) {
	definition, ok := c.definitions[id]
	if !ok {
		return Definition{}, fmt.Errorf("%w (%s)", ErrDefinitionNotFound, id.String())
	}

	return definition, nil
}

func (c *Container) Bind(definitions ...Definition) {
	if c.loaded.Load() {
		panic("DI: binding to a loaded container is not allowed")
	}

	for _, definition := range definitions {
		if err := c.canBind(definition); err != nil {
			panic(err)
		}

		definition.order = len(c.order[definition.scope])
		if !definition.asValue {
			ctor, err := NewInvokeInfo(c, definition.provider)
			if err != nil {
				panic(err)
			}
			definition.ctor = ctor
		}
		c.order[definition.scope] = append(c.order[definition.scope], definition.id)
		definition.sign = c.sign()
		c.definitions[definition.id] = definition
	}
}

func (c *Container) canBind(definition Definition) error {
	return nil
}
