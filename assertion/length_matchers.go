package assertion

import (
	"fmt"
	"reflect"
)

func lenCondition(cond func(int) bool, logFunc func(int) string, nlogFunc func(int) string) Matcher {
	return func(v interface{}) MatchResult {
		vv := reflect.ValueOf(v)
		switch vv.Kind() {
		case reflect.Array, reflect.String, reflect.Slice, reflect.Map:
			actualLen := vv.Len()
			if cond(actualLen) {
				return truthy(nlogFunc(actualLen))
			}
			return falsy(logFunc(actualLen))
		default:
			return falsy("\nValue type should be Array, Slice, String or Map")
		}
	}
}

func HasLen(len int) Matcher {
	return lenCondition(
		func(al int) bool { return al == len },
		func(al int) string {
			return fmt.Sprintf("\nValue length is not equal to expectation.\nExpected : %d\nGot : %d", len, al)
		},
		func(al int) string { return fmt.Sprintf("\nValue length should not be equal to : %d", len) },
	)
}

func HasMaxLen(len int) Matcher {
	return lenCondition(
		func(al int) bool { return al <= len },
		func(al int) string {
			return fmt.Sprintf("\nValue length is lower than : %d", len)
		},
		func(al int) string {
			return fmt.Sprintf("\nValue length should not be lower than : %d. \n Actual value : %d", len, al)
		})
}

func HasMinLen(len int) Matcher {
	return lenCondition(
		func(al int) bool { return al >= len },
		func(al int) string {
			return fmt.Sprintf("\nValue length is higher than : %d", len)
		},
		func(al int) string {
			return fmt.Sprintf("\nValue length should not be higher than : %d. \n Actual value : %d", len, al)
		})
}

func IsEmpty() Matcher {
	return lenCondition(
		func(al int) bool { return al == 0 },
		func(al int) string {
			return fmt.Sprintf("\nValue length should be empty. Actual length is equal to : %d", al)
		},
		func(al int) string {
			return "\nValue length should not be empty"
		})
}
