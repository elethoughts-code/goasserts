package assertion

import (
	"errors"
	"fmt"
	"reflect"
)

type MatchResult struct {
	Matches bool
	Log     string
	NLog    string
}

var zeroMatchResult = MatchResult{}

func errored(err error) (MatchResult, error) {
	return zeroMatchResult, err
}

func truthy(nlog string) (MatchResult, error) {
	return MatchResult{
		Matches: true,
		Log:     "",
		NLog:    nlog,
	}, nil
}

func falsy(log string) (MatchResult, error) {
	return MatchResult{
		Matches: false,
		Log:     log,
		NLog:    "",
	}, nil
}

type Matcher func(v interface{}) (MatchResult, error)

func IsEq(e interface{}) Matcher {
	return func(v interface{}) (MatchResult, error) {
		if e == v {
			return truthy(fmt.Sprintf("\nValue should not be equal to : %v", e))
		}
		return falsy(fmt.Sprintf("\nValue is not equal to expectation.\nExpected : %v\nGot : %v", e, v))
	}
}

func IsNil() Matcher {
	return func(v interface{}) (MatchResult, error) {
		if v == nil {
			return truthy("\nValue should not be nil")
		}
		return falsy(fmt.Sprintf("\nValue is not nil : %v", v))
	}
}

func IsDeepEq(e interface{}) Matcher {
	return func(v interface{}) (MatchResult, error) {
		if reflect.DeepEqual(v, e) {
			return truthy(fmt.Sprintf("\nValue should not be deep equal to : %v", e))
		}
		return falsy(fmt.Sprintf("\nValue is not deep equal to expectation.\nExpected : %v\nGot : %v", e, v))
	}
}

func HaveKind(k reflect.Kind) Matcher {
	return func(v interface{}) (MatchResult, error) {
		vv := reflect.ValueOf(v)
		if vv.Kind() == k {
			return truthy(fmt.Sprintf("\nValue should not of Kind : %v", k))
		}
		return falsy(fmt.Sprintf("\nValue is not of the expected Kind.\nExpected : %v\nGot : %v", k, vv.Kind()))
	}
}

func IsError(target error) Matcher {
	return func(v interface{}) (MatchResult, error) {
		ve, ok := v.(error)
		if !ok {
			return errored(ErrNotOfErrorType)
		}

		if errors.Is(ve, target) {
			return truthy(fmt.Sprintf("\nError value should not be : %v", target))
		}

		return falsy(fmt.Sprintf("\nError Value is not of the expected type.\nExpected : %v\nGot : %v", target, ve))
	}
}

func AsError(target interface{}) Matcher {
	return func(v interface{}) (MatchResult, error) {
		ve, ok := v.(error)
		if !ok {
			return errored(ErrNotOfErrorType)
		}

		if errors.As(ve, target) {
			return truthy("\nError value should not be as expected type")
		}

		return falsy("\nError Value is not as the expected type")
	}
}
