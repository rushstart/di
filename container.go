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
	c := &Container{}

	return c, func() { c.cleanup() }
}

type Container struct {
	id          int64
	loaded      atomic.Bool
	mu          sync.Mutex
	order       [2][]tid.ID
	definitions map[tid.ID]Definition
	pool        sync.Pool
	logger      Logger
	instances   []instance
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
	inst := c.instances[d.order]
	if inst.isset {
		return inst.value, nil
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
		invokeInfo, err = NewInvokeInfo(c, function)
		if err != nil {
			return nil, err
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
		panic("DI: cannot acquire a request scope in an unloaded container")
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
