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

func Unordered(e interface{}, areEq func(v, e interface{}) bool) Matcher {
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

		for i, expectedItem := range ie {
			found := false
			for _, sliceItem := range iv {
				if areEq(sliceItem, expectedItem) {
					found = true
					break
				}
			}
			if !found {
				return falsy(fmt.Sprintf("\nElement [%d]=%v not found", i, expectedItem))
			}
		}
		return truthy(fmt.Sprintf("\nValue should not contain all elements : %v", e))
	}
}

func applyMatcherToAllElements(v interface{}, matcher func(v interface{}) bool) (int, []int, []int, error) {
	iv, isSlice := toSlice(v)
	if !isSlice {
		return 0, nil, nil, ErrNotOfSliceType
	}
	truthyIndex := make([]int, 0)
	falsyIndex := make([]int, 0)
	for i, item := range iv {
		mr := matcher(item)
		if mr {
			truthyIndex = append(truthyIndex, i)
		} else {
			falsyIndex = append(falsyIndex, i)
		}
	}
	return len(iv), truthyIndex, falsyIndex, nil
}

func All(matcher func(v interface{}) bool) Matcher {
	return func(v interface{}) (MatchResult, error) {
		vLen, truthyIndex, falsyIndex, err := applyMatcherToAllElements(v, matcher)
		if err != nil {
			return MatchResult{}, err
		}
		if len(truthyIndex) != vLen {
			return falsy(fmt.Sprintf("\nMatcher dont apply to all values. Non matching indexes : %v", falsyIndex))
		}

		return truthy("\nMatcher should not apply to all elements")
	}
}

func AtLeast(n int, matcher func(v interface{}) bool) Matcher {
	return func(v interface{}) (MatchResult, error) {
		_, truthyIndex, falsyIndex, err := applyMatcherToAllElements(v, matcher)
		if err != nil {
			return MatchResult{}, err
		}
		if len(truthyIndex) < n {
			return falsy(fmt.Sprintf("\nAt least %d element(s) should match. Non matching indexes : %v", n, falsyIndex))
		}

		return truthy(fmt.Sprintf("\nMatcher should not apply to %d element(s) or more", n))
	}
}
