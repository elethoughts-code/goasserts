package assertion

import (
	"fmt"
	"reflect"
)

func toSlice(v interface{}) ([]interface{}, bool) {
	if v == nil {
		return nil, true
	}
	switch reflect.TypeOf(v).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(v)
		is := make([]interface{}, s.Len())
		for i := 0; i < s.Len(); i++ {
			element := s.Index(i).Interface()
			is[i] = element
		}
		return is, true
	default:
		return nil, false
	}
}

func Contains(e interface{}) Matcher {
	return func(v interface{}) (MatchResult, error) {
		iv, isSlice := toSlice(v)
		if !isSlice {
			return errored(ErrNotOfSliceType)
		}
		for _, item := range iv {
			if item == e {
				return truthy(fmt.Sprintf("\nValue should not contains element : %v", e))
			}
		}
		return falsy(fmt.Sprintf("\nValue should contains element : %v", e))
	}
}

func Unordered(e interface{}) Matcher {
	return func(v interface{}) (MatchResult, error) {
		iv, isSlice := toSlice(v)
		if !isSlice {
			return errored(ErrNotOfSliceType)
		}
		ie, isSlice := toSlice(e)
		if !isSlice {
			return errored(ErrNotOfSliceType)
		}

		if len(ie) != len(iv) {
			return falsy(fmt.Sprintf("\nValue should contains all elements : %v", e))
		}

		for _, expectedItem := range ie {
			found := false
			for _, sliceItem := range iv {
				if sliceItem == expectedItem {
					found = true
					break
				}
			}
			if !found {
				return falsy(fmt.Sprintf("\nElement [%v] not found", expectedItem))
			}
		}
		return truthy(fmt.Sprintf("\nValue should not contain all elements : %v", e))
	}
}
