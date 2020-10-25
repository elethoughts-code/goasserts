package assertion

import "reflect"

// MapEntry is the struct used by Entries() transformation to change a map to a slice of entries.
type MapEntry struct {
	Key   interface{}
	Value interface{}
}

// MapTransformer interface encloses map related transformations.
// All transformations return the same expectation interface to pile in calls (Fluent API).
//
// Values() Transform the assert value from map to a slice of its values.
//
// Keys() Transform the assert value from map to a slice of its keys.
//
// Entries() Transform the assert value from map to a slice of its entries (MapEntry).
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
