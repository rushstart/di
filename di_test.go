package di_test

import (
	"github.com/rushstart/di"
	"github.com/rushstart/di/tests/fixtures"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEndToEnd(t *testing.T) {
	container, cleanup := di.New()
	defer cleanup()

	t.Run("bind", func(t *testing.T) {
		container.Bind(
			di.Define[*fixtures.Wheel](func() *fixtures.Wheel { return fixtures.NewWheel(fixtures.WheelFrontLeft) }, di.WithTag("front-left")),
			di.Define[*fixtures.Wheel](func() *fixtures.Wheel { return fixtures.NewWheel(fixtures.WheelFrontRight) }, di.WithTag("front-right")),
			di.Define[*fixtures.Wheel](func() *fixtures.Wheel { return fixtures.NewWheel(fixtures.WheelBackLeft) }, di.WithTag("back-left")),
			di.Define[*fixtures.Wheel](func() *fixtures.Wheel { return fixtures.NewWheel(fixtures.WheelBackRight) }, di.WithTag("back-right")),
		)
		container.Bind(di.Define[*fixtures.Engine](func() *fixtures.Engine { return &fixtures.Engine{} }))
		container.Bind(di.Define[*fixtures.Car](func(s struct {
			engine *fixtures.Engine
			w1     *fixtures.Wheel `tag:"front-left"`
			w2     *fixtures.Wheel `tag:"front-right"`
			w3     *fixtures.Wheel `tag:"back-left"`
			w4     *fixtures.Wheel `tag:"back-right"`
		}) *fixtures.Car {
			return &fixtures.Car{
				Engine: s.engine,
				Wheels: [4]*fixtures.Wheel{s.w1, s.w2, s.w3, s.w4},
			}
		}))
		container.Load()
		car := di.MustGet[*fixtures.Car](container)
		car.Engine.Start()
		car.Engine.Stop()
	})

	t.Run("bind value", func(t *testing.T) {
		container.Bind(di.DefineValue[int](1))
		container.Bind(di.DefineValue[int](2, di.WithTag("two")))
		container.Load()
		assert.Equal(t, 1, di.MustGet[int](container))
		assert.Equal(t, 2, di.MustGet[int](container, "two"))
	})
}

func BenchmarkResolveValue(b *testing.B) {
	container, cleanup := di.New()
	defer cleanup()
	container.Bind(di.DefineValue[int](1))
	container.Bind(di.DefineValue[int](2, di.WithTag("two")))
	container.Load()
	ii, _ := di.NewInvokeInfo(container, func(s struct {
		i int `tag:"two"`
	}) int {
		return s.i
	})

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		container.Invoke(ii)
	}
}
