# About

**A dependency injection toolkit based on Generics**

## Getting started

Compatible with Go 1.22 or more.

Import package:

```
go get -u github.com/rushstart/di
```

### Create a DI container

The simplest way to start is to use the default options:

```go
import "github.com/rushstart/di"

container, cleanup := di.New()
defer cleanup()
```

### Service registration and invocation

Engine:
```go
// Provider
func NewEngine() (*Engine, error) {
    return &Engine{
        Started: false,
    }, nil
}

type Engine struct {
    Started bool
}

// Shutdown called on cleanup
func (e *Engine) Shutdown() error {
    e.Started = false
    return nil
}
```
Car:
```go
// Provider
func NewCar(engine *Engine) (*Car, error) {
    return &Car{
        Engine: engine,
    }, nil
}

type Car struct {
    Engine *Engine
}

func (c *Car) Start() {
    c.Engine.Started = true
    println("vroooom")
}
```
### Register services using individual declaration

```go
func main() {
    // create DI container
    container, cleanup := di.New()
    defer cleanup()
	
    container.Bind(di.Define[*Engine](NewEngine))
    container.Bind(di.Define[*Car](NewCar))
    // Before invoking the services, you need to bootstrap the container
    container.Bootstrap()

    // invoking car will instantiate Car services and its Engine dependency
    car, err := di.Get[*Car](container)
    if err != nil {
        log.Fatal(err.Error())
    }

    car.Start()
}
```
### Register services using package declaration

```go
var definitions = []di.Definition{
    di.Define[*Engine](NewEngine),
    di.Define[*Car](NewCar),
}

func main() {
    // create DI container
    container, cleanup := di.New()
    defer cleanup()
    
    container.Bind(definitions...)
    // Before invoking the services, you need to bootstrap the container
    container.Bootstrap()
    
    // invoking car will instantiate Car services and its Engine dependency
    car, err := di.Get[*Car](container)
    if err != nil {
        log.Fatal(err.Error())
    }
    
    car.Start()
}
```
### Tagged definitions

```go
var definitions = []di.Definition{
    di.Define[*Engine](NewEngine, di.WithTag("v8")),
    di.Define[*Engine](NewEngine, di.WithTag("v12")),
	
    di.Define[*Car](func(
        s struct {
            engine *Engine `tag:"v8"`
        },
    ) (*Car, error) {
        return &Car{Engine: s.engine}, nil
    }, di.WithTag("audi")),
	
    di.Define[*Car](func(
        s struct {
            engine *Engine `tag:"v12"`
        },
    ) (*Car, error) {
        return &Car{Engine: s.engine}, nil
    }, di.WithTag("porsche")),
	
    di.DefineValue[int](4242, di.WithTag("port"))
}
```
## Service invocation

```go
myService, err := di.Get[*MyService](container)
myService, err := di.Get[*MyService](container, "some-tag")
myService := di.MustGet[*MyService](container)
myService := di.MustGet[*MyService](container, "some-tag")

myService, err := container.Call(func(int count, service *Service) *MyService {
    return &MyService{}
}, 55)
myService := container.MustCall(func(int count, service *Service) *MyService {
    return &MyService{}
}, 55)
```

## Context scope

```go
// create DI container
container, cleanup := di.New()
defer cleanup()

container.Bind(di.Define[*Engine](NewEngine))
container.Bind(di.Define[*Car](NewCar), di.WithScope(di.Context))
// Before invoking the services, you need to bootstrap the container
container.Bootstrap()
go func() {
    cc := container.AcquireContainer()
    defer container.ReleaseContainer(cc)
    car := di.MustGet[*Car](cc)
    car.Start()
}()
```

## Trigger health check

```go
// returns the status (error or nil) for each service
container.HealthCheck() map[string]error
container.HealthCheckWithContext(context.Context) map[string]error
// on a single service
di.HealthCheck[T any](container, tag ...string) error
di.HealthCheckWithContext[T any](context.Context, container, tag ...string) error
```

## Two-step function call

```go
// prepare
callable, err := di.MakeCallable(container, func(int count, service *Service) *MyService {
    return &MyService{}
})
// ...
go func() {
    cc := container.AcquireContainer()
    defer container.ReleaseContainer(cc)
	// call
    callable.Call(cc, 55)
}()

```