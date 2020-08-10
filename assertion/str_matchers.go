package assertion

import "fmt"

func IsBlank() Matcher {
	return func(v interface{}) MatchResult {
		if v == "" {
			return truthy("\nValue should not be blank string")
		}
		return falsy(fmt.Sprintf("\nValue is not a blank string : %v", v))
	}
}