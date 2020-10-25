package assertion

import (
	"fmt"
	"reflect"
)

// AttributeParser interface encloses attribute/index access transformations.
// All transformations return the same expectation interface to pile in calls (Fluent API).
//
// Attr(key interface{}) change value to the corresponding attribute.
// Value should be a struct, a map or a ptr/interface to struct or map.
//
// Index(i int) change value to the corresponding index.
// Value should be a slice, an array or a ptr/interface to slice or array.
type AttributeParser interface {
	Attr(key interface{}) Expectation
	Index(i int) Expectation
}

func (exp *expectation) Attr(key interface{}) Expectation {
	value := reflect.ValueOf(exp.v)
	attrValue := elemFromKey(value, key)
	exp.v = attrValue
	return exp
}

func (exp *expectation) Index(i int) Expectation {
	value := reflect.ValueOf(exp.v)
	elementValue := elemFromIndex(value, i)
	exp.v = elementValue
	return exp
}

func elemFromKey(v reflect.Value, key interface{}) interface{} {
	if !v.IsValid() {
		panic("value is invalid")
	}
	switch v.Kind() {
	case reflect.Interface, reflect.Ptr:
		if v.IsNil() {
			panic("value is nil")
		}
		return elemFromKey(v.Elem(), key)
	case reflect.Struct:
		if _, ok := key.(string); !ok {
			panic("attribute key not of string type")
		}
		field := v.FieldByName(key.(string))
		if reflect.ValueOf(field).IsZero() {
			panic(fmt.Sprintf("attribute %s not found", key))
		}
		return field.Interface()
	case reflect.Map:
		if v.IsNil() {
			panic("value is nil")
		}
		field := v.MapIndex(reflect.ValueOf(key))
		if reflect.ValueOf(field).IsZero() {
			panic(fmt.Sprintf("attribute %s not found", key))
		}
		return field.Interface()
	default:
		panic("value should be an attribute type (struct, map, interface, ptr)")
	}
}

func elemFromIndex(v reflect.Value, i int) interface{} {
	if !v.IsValid() {
		panic("value is invalid")
	}
	switch v.Kind() {
	case reflect.Interface, reflect.Ptr:
		if v.IsNil() {
			panic("value is nil")
		}
		return elemFromIndex(v.Elem(), i)
	case reflect.Array, reflect.Slice:
		if i >= v.Len() {
			panic(fmt.Sprintf("index %d out of bound", i))
		}
		element := v.Index(i)
		return element.Interface()
	default:
		panic("value should be an indexed type (array, slice, interface, ptr)")
	}
}
