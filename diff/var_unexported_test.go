package diff

import (
	"reflect"
	"testing"

	"github.com/elethoughts-code/goasserts/assertion"
)

func Test_invalid_values_diff(t *testing.T) {
	// given
	assert := assertion.New(t)

	va := reflect.Value{}
	vb := reflect.ValueOf(3)
	diffs := make([]Diff, 0)

	// When

	findDiffs([]string{"a", "b", "c"}, va, vb, &diffs, map[visit]bool{})

	// Then

	assert.That(diffs).IsDeepEq([]Diff{
		{
			Path: []string{"a", "b", "c"},
			Value: InvalidDiff{
				A: false,
				B: true,
			},
		},
	})
}
