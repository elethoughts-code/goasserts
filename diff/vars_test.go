package diff_test

import (
	"github.com/elethoughts-code/goasserts/assertion"

	"testing"
	"unsafe"

	"github.com/elethoughts-code/goasserts/diff"
)

func Test_simple_one_level_differences(t *testing.T) {
	// Given
	assert := assertion.New(t)

	rootLevelDiffs := func(values ...interface{}) []diff.Diff {
		diffs := make([]diff.Diff, len(values))
		for i, v := range values {
			diffs[i] = diff.Diff{
				Path:  []string{},
				Value: v,
			}
		}
		return diffs
	}
	a := 1
	b := 1
	upa := unsafe.Pointer(&a)
	upb := unsafe.Pointer(&b)

	c1 := make(chan int)
	c2 := make(chan string)
	c3 := make(chan int)

	testCases := []struct {
		a      interface{}
		b      interface{}
		result []diff.Diff
	}{
		{
			a:      nil,
			b:      "this is a string",
			result: rootLevelDiffs(diff.CommonDiff{B: "this is a string"}),
		},
		{
			a:      123456,
			b:      nil,
			result: rootLevelDiffs(diff.CommonDiff{A: 123456}),
		},
		{
			a:      123456,
			b:      "this is a string",
			result: rootLevelDiffs(diff.TypeDiff{A: 123456, B: "this is a string"}),
		},
		{
			a:      "this is a string",
			b:      "this is another string",
			result: rootLevelDiffs(diff.CommonDiff{A: "this is a string", B: "this is another string"}),
		},
		{
			a:      [3]string{"a", "b", "c"},
			b:      [2]string{"a", "b"},
			result: rootLevelDiffs(diff.TypeDiff{A: [3]string{"a", "b", "c"}, B: [2]string{"a", "b"}}),
		},
		{
			a: []string{"a", "b"},
			b: []string{"a", "b", "c"},
			result: rootLevelDiffs(diff.LenDiff{
				CommonDiff: diff.CommonDiff{A: []string{"a", "b"}, B: []string{"a", "b", "c"}},
				Value:      -1,
			}),
		},
		{
			a: map[int]string{1: "a", 2: "b", 3: "c"},
			b: map[int]string{1: "a", 2: "b"},
			result: rootLevelDiffs(diff.LenDiff{
				CommonDiff: diff.CommonDiff{
					A: map[int]string{1: "a", 2: "b", 3: "c"},
					B: map[int]string{1: "a", 2: "b"},
				},
				Value: 1,
			}),
		},
		{
			a: complex(12, 12),
			b: complex(12, 15),
			result: rootLevelDiffs(diff.CommonDiff{
				A: complex(12, 12),
				B: complex(12, 15),
			}),
		},
		{
			a: upa,
			b: upb,
			result: rootLevelDiffs(diff.CommonDiff{
				A: upa,
				B: upb,
			}),
		},
		{
			a: c1,
			b: c2,
			result: rootLevelDiffs(diff.TypeDiff{
				A: c1,
				B: c2,
			}),
		},
		{
			a: c1,
			b: c3,
			result: rootLevelDiffs(diff.CommonDiff{
				A: c1,
				B: c3,
			}),
		},
	}

	for _, tc := range testCases {
		// When
		d := diff.Diffs(tc.a, tc.b)
		// Then
		assert.That(d).IsDeepEq(tc.result)
	}
}

type SampleStruct struct {
	A int
	B string
	C OtherStruct
	D *OtherStruct
}

type OtherStruct struct {
	A []OtherStruct
	B map[string][]int
	C interface{}
	D uint8
	E bool
}

func Test_simple_second_level_differences(t *testing.T) {
	assert := assertion.New(t)

	path := func(p ...string) []string { return p }
	diffs := func(d ...diff.Diff) []diff.Diff { return d }
	d := func(path []string, value interface{}) diff.Diff {
		return diff.Diff{
			Path:  path,
			Value: value,
		}
	}

	testCases := []struct {
		a      interface{}
		b      interface{}
		result []diff.Diff
	}{
		{
			a: []string{"a", "b"},
			b: []string{"a", "c"},
			result: diffs(
				d(path("[1]"), diff.CommonDiff{A: "b", B: "c"}),
			),
		},
		{
			a: map[string]int{"a": 1, "b": 2},
			b: map[string]int{"a": 1, "b": 3},
			result: diffs(
				d(path("[b]"), diff.CommonDiff{A: int64(2), B: int64(3)}),
			),
		},
		{
			a: map[string]int{"a": 1, "b": 2},
			b: map[string]int{"a": 1, "c": 2},
			result: diffs(
				d(path("[b]"), diff.KeyNotFoundDiff{A: true, B: false}),
				d(path("[c]"), diff.KeyNotFoundDiff{A: false, B: true}),
			),
		},
		{
			a: SampleStruct{A: 1},
			b: SampleStruct{A: 2},
			result: diffs(
				d(path("[A]"), diff.CommonDiff{A: int64(1), B: int64(2)}),
			),
		},
	}

	for _, tc := range testCases {
		// When
		d := diff.Diffs(tc.a, tc.b)
		// Then
		assert.That(d).IsDeepEq(tc.result)
	}
}

func Test_multi_level_differences(t *testing.T) {
	assert := assertion.New(t)

	path := func(p ...string) []string { return p }
	diffs := func(d ...diff.Diff) []diff.Diff { return d }
	d := func(path []string, value interface{}) diff.Diff {
		return diff.Diff{
			Path:  path,
			Value: value,
		}
	}
	os := OtherStruct{
		A: nil,
		B: map[string][]int{
			"x": {1},
			"y": {1, 2},
			"z": {1, 2, 3},
		},
		C: nil,
		D: 0,
		E: false,
	}
	os2 := OtherStruct{
		A: nil,
		B: nil,
		C: nil,
		D: 0,
		E: true,
	}
	testCases := []struct {
		a      interface{}
		b      interface{}
		result []diff.Diff
	}{
		{
			a: SampleStruct{
				A: 1,
				B: "b1",
				C: OtherStruct{
					A: []OtherStruct{{D: 1}, {D: 2}, {D: 3}},
					B: nil,
					C: nil,
					D: 0,
					E: false,
				},
				D: &os,
			},
			b: SampleStruct{
				A: 2,
				B: "b2",
				C: OtherStruct{
					A: []OtherStruct{{D: 1}, {D: 2}, {D: 4}},
					B: nil,
					C: nil,
					D: 0,
					E: false,
				},
				D: nil,
			},
			result: diffs(
				d(path("[A]"), diff.CommonDiff{A: int64(1), B: int64(2)}),
				d(path("[B]"), diff.CommonDiff{A: "b1", B: "b2"}),
				d(path("[C]", "[A]", "[2]", "[D]"), diff.CommonDiff{A: uint64(3), B: uint64(4)}),
				d(path("[D]"), diff.CommonDiff{A: &os, B: nil}),
			),
		},
		{
			a: SampleStruct{
				A: 1,
				B: "b1",
				C: OtherStruct{
					A: []OtherStruct{{D: 1}, {D: 2}, {D: 3}},
					B: nil,
					C: nil,
					D: 0,
					E: false,
				},
				D: &os,
			},
			b: SampleStruct{
				A: 2,
				B: "b2",
				C: OtherStruct{
					A: []OtherStruct{{D: 1}, {D: 2}, {D: 4}},
					B: nil,
					C: nil,
					D: 0,
					E: false,
				},
				D: &os2,
			},
			result: diffs(
				d(path("[A]"), diff.CommonDiff{A: int64(1), B: int64(2)}),
				d(path("[B]"), diff.CommonDiff{A: "b1", B: "b2"}),
				d(path("[C]", "[A]", "[2]", "[D]"), diff.CommonDiff{A: uint64(3), B: uint64(4)}),
				d(path("[D]", "[&]", "[B]"), diff.CommonDiff{A: map[string][]int{
					"x": {1},
					"y": {1, 2},
					"z": {1, 2, 3},
				}, B: nil}),
				d(path("[D]", "[&]", "[E]"), diff.CommonDiff{A: false, B: true}),
			),
		},
	}

	for _, tc := range testCases {
		// When
		d := diff.Diffs(tc.a, tc.b)
		// Then
		assert.That(d).IsDeepEq(tc.result)
	}
}
