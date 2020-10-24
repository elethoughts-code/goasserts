package assertion

import (
	"fmt"
	"reflect"
)

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
