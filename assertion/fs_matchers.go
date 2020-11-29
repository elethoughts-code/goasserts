package assertion

import (
	"fmt"
	"os"
)

func FileExists() Matcher {
	return func(v interface{}) (MatchResult, error) {
		if _, err := os.Stat(v.(string)); !os.IsNotExist(err) {
			return truthy(fmt.Sprintf("\nFile %v should not exists", v))
		}
		return falsy(fmt.Sprintf("\nFile %v do not exists", v))
	}
}
