package assertion

import (
	"reflect"
)

type ReflectTransformer interface {
	Dereference() Expectation
}

func (exp *expectation) Dereference() Expectation {
	exp.t.Helper()
	v := reflect.ValueOf(exp.v)
	switch v.Kind() {
	case reflect.Ptr:
		exp.v = v.Elem().Interface()
	default:
		panic("value is not a pointer")
	}
	return exp
}
