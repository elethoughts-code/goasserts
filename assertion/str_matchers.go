package assertion

import (
	"fmt"
	"regexp"
	"strings"
)

func IsBlank() Matcher {
	return func(v interface{}) (MatchResult, error) {
		if v == "" {
			return truthy("\nValue should not be blank string")
		}
		return falsy(fmt.Sprintf("\nValue is not a blank string : %v", v))
	}
}

func MatchRe(reg string) Matcher {
	return func(v interface{}) (MatchResult, error) {
		s, ok := v.(string)
		if !ok {
			return errored(ErrNotOfStringType)
		}
		match, err := regexp.Match(reg, []byte(s))
		if err != nil {
			return errored(err)
		}
		if match {
			return truthy(fmt.Sprintf("\nValue should not match regexp : %s", reg))
		}
		return falsy(fmt.Sprintf("\nValue do not match regexp : %s", reg))
	}
}

func HasPrefix(prefix string) Matcher {
	return func(v interface{}) (MatchResult, error) {
		if strings.HasPrefix(v.(string), prefix) {
			return truthy(fmt.Sprintf("\nValue should have prefix = %s", prefix))
		}
		return falsy(fmt.Sprintf("\nValue should not have prefix = %s", prefix))
	}
}

func HasSuffix(suffix string) Matcher {
	return func(v interface{}) (MatchResult, error) {
		if strings.HasSuffix(v.(string), suffix) {
			return truthy(fmt.Sprintf("\nValue should have suffix = %s", suffix))
		}
		return falsy(fmt.Sprintf("\nValue should not have suffix = %s", suffix))
	}
}
