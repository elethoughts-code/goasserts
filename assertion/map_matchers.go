package assertion

import (
	"fmt"
	"reflect"
)

func ContainsValue(e interface{}) Matcher {
	return func(v interface{}) MatchResult {
		switch reflect.TypeOf(v).Kind() {
		case reflect.Map:
			vm := reflect.ValueOf(v)
			mk := vm.MapKeys()
			for _, k := range mk {
				element := vm.MapIndex(k).Interface()
				if element == e {
					return truthy(fmt.Sprintf("\nValue should not contains element : %v", e))
				}
			}
			return falsy(fmt.Sprintf("\nValue should contains element : %v", e))

		default:
			return falsy("\nValue should be a map")
		}
	}
}

func ContainsKey(e interface{}) Matcher {
	return func(v interface{}) MatchResult {
		switch reflect.TypeOf(v).Kind() {
		case reflect.Map:
			vm := reflect.ValueOf(v)
			mk := vm.MapKeys()
			for _, k := range mk {
				if k.Interface() == e {
					return truthy(fmt.Sprintf("\nValue should not contains key : %v", e))
				}
			}
			return falsy(fmt.Sprintf("\nValue should contains key : %v", e))

		default:
			return falsy("\nValue should be a map")
		}
	}
}
