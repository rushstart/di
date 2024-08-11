package di

import "github.com/rushstart/tid"

func Get[T any](s Scope, tag ...string) (service T, err error) {
	val, err := s.Resolve(tid.From[T](tag...))
	if err != nil {
		return service, err
	}
	return val.Interface().(T), nil
}

func MustGet[T any](s Scope, tag ...string) T {
	return must1(Get[T](s, tag...))
}
