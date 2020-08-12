package assertion

import (
	"fmt"
	"regexp"
)

func IsBlank() Matcher {
	return func(v interface{}) MatchResult {
		if v == "" {
			return truthy("\nValue should not be blank string")
		}
		return falsy(fmt.Sprintf("\nValue is not a blank string : %v", v))
	}
}

func MatchRe(reg string) Matcher {
	return func(v interface{}) MatchResult {
		s, ok := v.(string)
		if !ok {
			return falsy(fmt.Sprintf("\nValue type is not a string : %v", v))
		}
		match, err := regexp.Match(reg, []byte(s))
		if err != nil {
			return falsy(fmt.Sprintf("\nCannot match for regexp : %s", reg))
		}
		if match {
			return truthy(fmt.Sprintf("\nValue should not match regexp : %s", reg))
		}
		return falsy(fmt.Sprintf("\nValue do not match regexp : %s", reg))
	}
}
