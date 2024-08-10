package di_test

//var definitions = []di.Definition{
//	di.Define[*Engine](NewEngine, di.WithTag("v8")),
//	di.Define[*Engine](NewEngine, di.WithTag("v12")),
//
//	di.Define[*Car](func(
//		s struct {
//			engine *Engine `tag:"v8"`
//		},
//	) (*Car, error) {
//		return &Car{Engine: s.engine}, nil
//	}, di.WithTag("audi")),
//
//	di.Define[*Car](func(
//		s struct {
//			engine *Engine `tag:"v12"`
//		},
//	) (*Car, error) {
//		return &Car{Engine: s.engine}, nil
//	}, di.WithTag("porsche")),
//}
//
//func NewEngine() (*Engine, error) {
//	return &Engine{
//		Started: false,
//	}, nil
//}
//
//type Engine struct {
//	Started bool
//}
//
//func (e *Engine) Shutdown() error {
//	// called on injector shutdown
//	e.Started = false
//	return nil
//}
//
//func ExampleBind() {
//	container, cleanup := di.New(di.Config{logger: new(log.Logger)})
//	defer cleanup()
//
//	type test struct {
//		foobar string
//	}
//
//	type test2 struct {
//		*test
//	}
//	// without dependencies
//	container.Bind(
//		di.Define[*test](func() (*test, error) { return &test{foobar: "foobar"}, nil }),
//	)
//	// tagged without dependencies
//	container.Bind(
//		di.Define[*test](func() (*test, error) { return &test{foobar: "foobar"}, nil }, di.WithTag("some-tag")),
//	)
//
//	// with dependencies
//	container.Bind(
//		di.Define[*test2](func(test *test) (*test2, error) { return &test2{test}, nil }),
//	)
//
//	// tagged with tagged dependencies
//	container.Bind(
//		di.Define[*test2](
//			func(
//				s struct {
//					test *test `tag:"some-tag"`
//				},
//			) (*test2, error) {
//				return &test2{s.test}, nil
//			},
//			di.WithTag("some-tag"),
//		),
//	)
//
//	container.Bootstrap()
//
//	value, err := di.Get[*test2](container, "some-tag")
//
//	fmt.Println(value)
//	fmt.Println(err)
//}
