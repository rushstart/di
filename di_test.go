package di_test

import (
	"github.com/rushstart/di"
	"github.com/rushstart/tid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEndToEnd(t *testing.T) {
	container, cleanup := di.New()
	defer cleanup()

	type Foo struct {
		Name string
	}
	type Bar struct {
		Name string
	}
	type FooBar struct {
		Foo
		*Bar
	}

	t.Run("bind", func(t *testing.T) {
		container.Bind(di.Define[Foo](func() Foo { return Foo{Name: "foo"} }))
		container.Bind(di.Define[Foo](func() Foo { return Foo{Name: "foo-tagged"} }, di.WithTag("tagged")))
		container.Load()
		assert.Equal(t, "foo", di.MustGet[Foo](container).Name)
		assert.Equal(t, "foo-tagged", di.MustGet[Foo](container, "tagged").Name)
	})

	t.Run("bind value", func(t *testing.T) {
		container.Bind(di.DefineValue(1))
		container.Bind(di.DefineValue(2, di.WithTag("two")))
		container.Load()
		assert.Equal(t, 1, di.MustGet[int](container))
		assert.Equal(t, 2, di.MustGet[int](container, "two"))
	})
}

func BenchmarkSingleton(b *testing.B) {
	container, cleanup := di.New()
	defer cleanup()
	type Foo struct {
		Name string
	}
	type Bar struct {
		Name string
	}
	type FooBar struct {
		Foo
		*Bar
	}
	container.Bind(di.DefineValue(1))
	container.Bind(di.DefineValue(2, di.WithTag("two")))
	container.Bind(di.Define[Foo](func() Foo { return Foo{Name: "foo"} }))
	container.Bind(di.Define[Foo](func() Foo { return Foo{Name: "foo-tagged"} }, di.WithTag("tagged")))
	container.Load()

	b.Run("resolve value", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			container.Resolve(tid.From[int]("two"))
		}
	})

	b.Run("resolve", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			container.Resolve(tid.From[Foo]())
		}
	})

	b.Run("resolve tagged", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			container.Resolve(tid.From[Foo]("tagged"))
		}
	})

	b.Run("invoke", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			container.Invoke(func(i int) {
			})
		}
	})

	b.Run("invoke in struct", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			container.Invoke(func(s struct {
				I int `tag:"two"`
			}) {
			})
		}
	})

}
