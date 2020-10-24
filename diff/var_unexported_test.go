package diff //nolint: testpackage

import (
	"reflect"
	"testing"
)

func Test_invalid_values_diff(t *testing.T) {
	// given

	va := reflect.Value{}
	vb := reflect.ValueOf(3)
	diffs := make([]Diff, 0)

	// When

	findDiffs([]string{"a", "b", "c"}, va, vb, &diffs, map[visit]bool{})

	// Then
	if !reflect.DeepEqual(diffs, []Diff{
		{
			Path: []string{"a", "b", "c"},
			Value: InvalidDiff{
				A: false,
				B: true,
			},
		},
	}) {
		t.Fail()
	}
}
