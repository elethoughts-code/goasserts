package assertion

import "reflect"

type MapEntry struct {
	Key   interface{}
	Value interface{}
}

type MapTransformer interface {
	Values() Expectation
	Keys() Expectation
	Entries() Expectation
}

func (exp *expectation) Values() Expectation {
	if exp.v == nil {
		return exp
	}
	switch reflect.TypeOf(exp.v).Kind() {
	case reflect.Map:
		m := reflect.ValueOf(exp.v)
		values := make([]interface{}, m.Len())
		i := 0
		iter := m.MapRange()
		for iter.Next() {
			values[i] = iter.Value().Interface()
			i++
		}
		exp.v = values
	default:
		panic("[type error] value should be a map")
	}
	return exp
}

func (exp *expectation) Keys() Expectation {
	if exp.v == nil {
		return exp
	}
	switch reflect.TypeOf(exp.v).Kind() {
	case reflect.Map:
		m := reflect.ValueOf(exp.v)
		values := make([]interface{}, m.Len())
		i := 0
		iter := m.MapRange()
		for iter.Next() {
			values[i] = iter.Key().Interface()
			i++
		}
		exp.v = values
	default:
		panic("[type error] value should be a map")
	}
	return exp
}

func (exp *expectation) Entries() Expectation {
	if exp.v == nil {
		return exp
	}
	switch reflect.TypeOf(exp.v).Kind() {
	case reflect.Map:
		m := reflect.ValueOf(exp.v)
		values := make([]interface{}, m.Len())
		i := 0
		iter := m.MapRange()
		for iter.Next() {
			values[i] = MapEntry{
				Key:   iter.Key().Interface(),
				Value: iter.Value().Interface(),
			}
			i++
		}
		exp.v = values
	default:
		panic("[type error] value should be a map")
	}
	return exp
}
