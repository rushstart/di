package di

import "reflect"

type InStruct struct{}

func (InStruct) inStructGuard() {}

func IsInStruct(t reflect.Type) bool {
	if t.Kind() == reflect.Struct && t.PkgPath() == "" {
		return true
	}
	_, ok := reflect.New(t).Elem().Interface().(interface{ inStructGuard() })
	return ok
}
