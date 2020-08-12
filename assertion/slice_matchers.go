package assertion

import (
	"fmt"
	"reflect"
)

func Contains(e interface{}) Matcher {
	return func(v interface{}) (MatchResult, error) {
		switch reflect.TypeOf(v).Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(v)

			for i := 0; i < s.Len(); i++ {
				element := s.Index(i).Interface()
				if element == e {
					return truthy(fmt.Sprintf("\nValue should not contains element : %v", e))
				}
			}
			return falsy(fmt.Sprintf("\nValue should contains element : %v", e))

		default:
			return falsy("\nValue should be a slice")
		}
	}
}
