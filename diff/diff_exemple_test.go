package diff_test

import (
	"fmt"

	"github.com/elethoughts-code/goasserts/diff"
)

func Example() {
	a := SampleStruct{
		A: 1,
		B: "b1",
		C: OtherStruct{
			A: []OtherStruct{{D: 1}, {D: 2}, {D: 3}},
			D: 0,
			E: false,
		},
		D: &OtherStruct{},
	}

	b := SampleStruct{
		A: 2,
		B: "b2",
		C: OtherStruct{
			A: []OtherStruct{{D: 1}, {D: 2}, {D: 4}},
			D: 0,
			E: false,
		},
		D: nil,
	}

	diffs := diff.Diffs(a, b)
	for _, d := range diffs {
		fmt.Printf("%v\n\n", d)
	}
}
