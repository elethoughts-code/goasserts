package assertion

import (
	"fmt"
	"reflect"
)

type MatchResult struct {
	Matches bool
	Log     string
	NLog    string
}

func truthy(nlog string) MatchResult {
	return MatchResult{
		Matches: true,
		Log:     "",
		NLog:    nlog,
	}
}

func falsy(log string) MatchResult {
	return MatchResult{
		Matches: false,
		Log:     log,
		NLog:    "",
	}
}

type Matcher func(v interface{}) MatchResult

func IsEq(e interface{}) Matcher {
	return func(v interface{}) MatchResult {
		if e == v {
			return truthy(fmt.Sprintf("\nValue should not be equal to : %v", e))
		}
		return falsy(fmt.Sprintf("\nValue is not equal to expectation.\nExpected : %v\nGot : %v", e, v))
	}
}

func IsDeepEq(e interface{}) Matcher {
	return func(v interface{}) MatchResult {
		if reflect.DeepEqual(v, e) {
			return truthy(fmt.Sprintf("\nValue should not be deep equal to : %v", e))
		}
		return falsy(fmt.Sprintf("\nValue is not deep equal to expectation.\nExpected : %v\nGot : %v", e, v))
	}
}
