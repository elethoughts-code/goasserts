package diff_test

import (
	"errors"
	"reflect"
	"strings"
	"testing"
	"unsafe"

	"github.com/elethoughts-code/goasserts/diff"
)

func Test_Similar_simple_one_level_differences_functions(t *testing.T) {
	// Given

	// When
	d1 := diff.Similar(func() {}, func() {}, false)
	d2 := diff.Similar(nil, nil, false)

	// Then
	if !errors.As(d1[0].Value, &diff.FuncDiff{}) || d1[0].Value.Error() != "functions cannot be compared" {
		t.Fail()
	}
	if len(d2) != 0 {
		t.Fail()
	}
}

func Test_Similar_simple_one_level_differences(t *testing.T) {
	// Given

	rootLevelDiffs := func(values ...error) []diff.Diff {
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
			a: [3]string{"a", "b", "c"},
			b: [2]string{"a", "b"},
			result: rootLevelDiffs(diff.LenDiff{
				CommonDiff: diff.CommonDiff{A: [3]string{"a", "b", "c"}, B: [2]string{"a", "b"}},
				Value:      1,
			}),
		},
		{
			a: []string{"a", "b", "c"},
			b: [2]string{"a", "b"},
			result: rootLevelDiffs(diff.LenDiff{
				CommonDiff: diff.CommonDiff{A: []string{"a", "b", "c"}, B: [2]string{"a", "b"}},
				Value:      1,
			}),
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
			a: float64(12),
			b: float64(15),
			result: rootLevelDiffs(diff.CommonDiff{
				A: float64(12),
				B: float64(15),
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
		{
			a: []string{"A", "B", "C"},
			b: struct {
				A string
			}{"A"},
			result: rootLevelDiffs(diff.TypeDiff{
				A: []string{"A", "B", "C"},
				B: struct {
					A string
				}{"A"},
			}),
		},
		{
			a: map[string]int{"A": 1, "B": 2, "C": 3},
			b: 123,
			result: rootLevelDiffs(diff.TypeDiff{
				A: map[string]int{"A": 1, "B": 2, "C": 3},
				B: 123,
			}),
		},
	}

	for i, tc := range testCases {
		// When
		d := diff.Similar(tc.a, tc.b, false)
		// Then
		if !reflect.DeepEqual(d, tc.result) {
			t.Logf("%v", i)
			t.Fail()
		}
	}
}

func Test_Similar_simple_second_level_differences(t *testing.T) {
	path := func(p ...string) []string { return p }
	diffs := func(d ...diff.Diff) []diff.Diff { return d }
	d := func(path []string, value error) diff.Diff {
		return diff.Diff{
			Path:  path,
			Value: value,
		}
	}

	type KeyStruct struct {
		Key string
	}

	keyA := KeyStruct{Key: "A"}
	keyA2 := KeyStruct{Key: "A"}

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
			a: map[string]int{"A": 1, "B": 2},
			b: struct {
				A int
				B int
			}{1, 3},
			result: diffs(
				d(path("[B]"), diff.CommonDiff{A: int64(2), B: int64(3)}),
			),
		},
		{
			a: map[string]int{"a": 1, "b": 2},
			b: map[string]int{"a": 1, "c": 2},
			result: diffs(
				d(path("[b]"), diff.KeyNotFoundDiff{Key: "b", A: true, B: false}),
				d(path("[c]"), diff.KeyNotFoundDiff{Key: "c", A: false, B: true}),
			),
		},
		{
			a: map[string]int{"A": 1, "B": 2},
			b: struct {
				A int
				C int
			}{1, 3},
			result: diffs(
				d(path("[B]"), diff.KeyNotFoundDiff{Key: "B", A: true, B: false}),
				d(path("[C]"), diff.KeyNotFoundDiff{Key: "C", A: false, B: true}),
			),
		},
		{
			a: map[KeyStruct]int{{"a"}: 1, {"b"}: 2},
			b: map[KeyStruct]int{{"a"}: 1, {"b"}: 3},
			result: diffs(
				d(path("[{b}]"), diff.CommonDiff{A: int64(2), B: int64(3)}),
			),
		},
		{
			a: map[KeyStruct]int{{"a"}: 1, {"b"}: 2},
			b: map[KeyStruct]int{{"a"}: 1, {"c"}: 2},
			result: diffs(
				d(path("[{b}]"), diff.KeyNotFoundDiff{Key: "{b}", A: true, B: false}),
				d(path("[{c}]"), diff.KeyNotFoundDiff{Key: "{c}", A: false, B: true}),
			),
		},

		{
			a: SampleStruct{A: 1},
			b: SampleStruct{A: 2},
			result: diffs(
				d(path("[A]"), diff.CommonDiff{A: int64(1), B: int64(2)}),
			),
		},
		{
			a: map[string]interface{}{
				"A": KeyStruct{Key: "A"},
				"B": KeyStruct{Key: "A"},
				"C": KeyStruct{Key: "C"},
			},
			b: struct {
				A KeyStruct
				B KeyStruct
				C KeyStruct
			}{
				A: KeyStruct{Key: "A"},
				B: KeyStruct{Key: "A"},
				C: KeyStruct{Key: "A"},
			},
			result: diffs(
				d(path("[C]", "[Key]"), diff.CommonDiff{A: "C", B: "A"}),
			),
		},
		{
			a: map[string]interface{}{
				"A": &keyA,
				"B": &keyA,
				"C": KeyStruct{Key: "C"},
			},
			b: struct {
				A *KeyStruct
				B *KeyStruct
				C map[string]string
			}{
				A: &keyA,
				B: &keyA,
				C: map[string]string{"Key": "A"},
			},
			result: diffs(
				d(path("[C]", "[Key]"), diff.CommonDiff{A: "C", B: "A"}),
			),
		},
		{
			a: struct {
				A *KeyStruct
				B *KeyStruct
				C map[string]string
			}{
				A: &keyA2,
				B: &keyA,
				C: map[string]string{"Key": "A"},
			},
			b: map[string]interface{}{
				"A": &keyA,
				"B": &keyA2,
				"C": KeyStruct{Key: "C"},
			},
			result: diffs(
				d(path("[C]", "[Key]"), diff.CommonDiff{A: "A", B: "C"}),
			),
		},
	}

	for _, tc := range testCases {
		// When
		d := diff.Similar(tc.a, tc.b, false)
		// Then
		if !reflect.DeepEqual(d, tc.result) {
			t.Fail()
		}
	}
}

func Test_Similar_multi_level_differences(t *testing.T) {
	path := func(p ...string) []string { return p }
	diffs := func(d ...diff.Diff) []diff.Diff { return d }
	d := func(path []string, value error) diff.Diff {
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
			b: map[string]interface{}{
				"A": 2,
				"B": "b2",
				"C": map[string]interface{}{
					"A": []OtherStruct{{D: 1}, {D: 2}, {D: 4}},
					"B": nil,
					"C": nil,
					"D": uint8(0),
					"E": false,
				},
				"D": nil,
			},
			result: diffs(
				d(path("[A]"), diff.CommonDiff{A: int64(1), B: int64(2)}),
				d(path("[B]"), diff.CommonDiff{A: "b1", B: "b2"}),
				d(path("[C]", "[A]", "[2]", "[D]"), diff.CommonDiff{A: uint64(3), B: uint64(4)}),
				d(path("[D]", "[A]"), diff.KeyNotFoundDiff{Key: "A", A: true, B: false}),
				d(path("[D]", "[B]"), diff.KeyNotFoundDiff{Key: "B", A: true, B: false}),
				d(path("[D]", "[C]"), diff.KeyNotFoundDiff{Key: "C", A: true, B: false}),
				d(path("[D]", "[D]"), diff.KeyNotFoundDiff{Key: "D", A: true, B: false}),
				d(path("[D]", "[E]"), diff.KeyNotFoundDiff{Key: "E", A: true, B: false}),
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
				D: nil,
			},
			result: diffs(
				d(path("[A]"), diff.CommonDiff{A: int64(1), B: int64(2)}),
				d(path("[B]"), diff.CommonDiff{A: "b1", B: "b2"}),
				d(path("[C]", "[A]", "[2]", "[D]"), diff.CommonDiff{A: uint64(3), B: uint64(4)}),
				d(path("[D]", "[A]"), diff.KeyNotFoundDiff{Key: "A", A: true, B: false}),
				d(path("[D]", "[B]"), diff.KeyNotFoundDiff{Key: "B", A: true, B: false}),
				d(path("[D]", "[C]"), diff.KeyNotFoundDiff{Key: "C", A: true, B: false}),
				d(path("[D]", "[D]"), diff.KeyNotFoundDiff{Key: "D", A: true, B: false}),
				d(path("[D]", "[E]"), diff.KeyNotFoundDiff{Key: "E", A: true, B: false}),
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
				d(path("[D]", "[B]", "[x]"), diff.KeyNotFoundDiff{Key: "x", A: true, B: false}),
				d(path("[D]", "[B]", "[y]"), diff.KeyNotFoundDiff{Key: "y", A: true, B: false}),
				d(path("[D]", "[B]", "[z]"), diff.KeyNotFoundDiff{Key: "z", A: true, B: false}),
				d(path("[D]", "[E]"), diff.CommonDiff{A: false, B: true}),
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
			b: map[string]interface{}{
				"A": 2,
				"B": "b2",
				"C": map[string]interface{}{
					"A": []interface{}{
						map[string]interface{}{"A": nil, "B": nil, "C": nil, "D": uint8(1), "E": false},
						OtherStruct{D: 2},
						OtherStruct{D: 4}},
					"B": nil,
					"C": nil,
					"D": uint8(0),
					"E": false,
				},
				"D": &os2,
			},
			result: diffs(
				d(path("[A]"), diff.CommonDiff{A: int64(1), B: int64(2)}),
				d(path("[B]"), diff.CommonDiff{A: "b1", B: "b2"}),
				d(path("[C]", "[A]", "[2]", "[D]"), diff.CommonDiff{A: uint64(3), B: uint64(4)}),
				d(path("[D]", "[B]", "[x]"), diff.KeyNotFoundDiff{Key: "x", A: true, B: false}),
				d(path("[D]", "[B]", "[y]"), diff.KeyNotFoundDiff{Key: "y", A: true, B: false}),
				d(path("[D]", "[B]", "[z]"), diff.KeyNotFoundDiff{Key: "z", A: true, B: false}),
				d(path("[D]", "[E]"), diff.CommonDiff{A: false, B: true}),
			),
		},
	}

	for _, tc := range testCases {
		// When
		d := diff.Similar(tc.a, tc.b, false)
		// Then
		if !reflect.DeepEqual(diffMap(d), diffMap(tc.result)) {
			t.Fail()
		}
	}
}

func diffMap(diffs []diff.Diff) map[string]diff.Diff {
	m := make(map[string]diff.Diff)
	for _, d := range diffs {
		m[strings.Join(d.Path, ".")] = d
	}
	if len(m) != len(diffs) {
		panic("same path found")
	}
	return m
}

func Test_Similar_Unordered_simple_one_level_differences(t *testing.T) {
	// Given

	rootLevelDiffs := func(values ...error) []diff.Diff {
		diffs := make([]diff.Diff, len(values))
		for i, v := range values {
			diffs[i] = diff.Diff{
				Path:  []string{},
				Value: v,
			}
		}
		return diffs
	}
	path := func(p ...string) []string { return p }
	diffs := func(d ...diff.Diff) []diff.Diff { return d }
	d := func(path []string, value error) diff.Diff {
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

	type ComparableStr struct {
		A string
	}

	type ComparableStr2 struct {
		A string
	}
	testCases := []struct {
		a      interface{}
		b      interface{}
		result []diff.Diff
	}{
		{
			a:      []string{"a", "b", "c"},
			b:      []string{"a", "d", "c"},
			result: rootLevelDiffs(diff.CommonDiff{A: []interface{}{"a", "b", "c"}, B: []interface{}{"a", "d", "c"}}),
		},
		{
			a: []string{"a", "b", "c"},
			b: []string{"a", "b"},
			result: rootLevelDiffs(diff.LenDiff{
				CommonDiff: diff.CommonDiff{A: []string{"a", "b", "c"}, B: []string{"a", "b"}},
				Value:      1,
			}),
		},

		{
			a: []ComparableStr{{"a"}, {"b"}, {"c"}},
			b: []ComparableStr{{"a"}, {"b"}, {"d"}},
			result: rootLevelDiffs(diff.CommonDiff{
				A: []interface{}{ComparableStr{"a"}, ComparableStr{"b"}, ComparableStr{"c"}},
				B: []interface{}{ComparableStr{"a"}, ComparableStr{"b"}, ComparableStr{"d"}},
			}),
		},
		{
			a: []ComparableStr{{"a"}, {"b"}, {"c"}},
			b: []interface{}{ComparableStr{"a"}, &ComparableStr{"b"}, ComparableStr{"d"}},
			result: rootLevelDiffs(diff.CommonDiff{
				A: []interface{}{ComparableStr{"a"}, ComparableStr{"b"}, ComparableStr{"c"}},
				B: []interface{}{ComparableStr{"a"}, ComparableStr{"b"}, ComparableStr{"d"}},
			}),
		},
		{
			a:      []ComparableStr{{"a"}, {"b"}, {"c"}},
			b:      []ComparableStr2{{"a"}, {"b"}, {"d"}},
			result: diffs(d(path("[2]", "[A]"), diff.CommonDiff{A: "c", B: "d"})),
		},
		{
			a:      []interface{}{ComparableStr{"a"}, ComparableStr{"b"}, ComparableStr{"c"}},
			b:      []interface{}{ComparableStr{"a"}, "b", ComparableStr{"c"}},
			result: diffs(d(path("[1]"), diff.TypeDiff{A: ComparableStr{"b"}, B: "b"})),
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
			b: map[string]interface{}{
				"A": 2,
				"B": "b2",
				"C": map[string]interface{}{
					"A": []OtherStruct{{D: 1}, {D: 2}, {D: 4}},
					"B": nil,
					"C": nil,
					"D": uint8(0),
					"E": false,
				},
				"D": nil,
			},
			result: diffs(
				d(path("[A]"), diff.CommonDiff{A: int64(1), B: int64(2)}),
				d(path("[B]"), diff.CommonDiff{A: "b1", B: "b2"}),
				d(path("[C]", "[A]", "[2]", "[D]"), diff.CommonDiff{A: uint64(3), B: uint64(4)}),
				d(path("[D]", "[A]"), diff.KeyNotFoundDiff{Key: "A", A: true, B: false}),
				d(path("[D]", "[B]"), diff.KeyNotFoundDiff{Key: "B", A: true, B: false}),
				d(path("[D]", "[C]"), diff.KeyNotFoundDiff{Key: "C", A: true, B: false}),
				d(path("[D]", "[D]"), diff.KeyNotFoundDiff{Key: "D", A: true, B: false}),
				d(path("[D]", "[E]"), diff.KeyNotFoundDiff{Key: "E", A: true, B: false}),
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
				D: nil,
			},
			result: diffs(
				d(path("[A]"), diff.CommonDiff{A: int64(1), B: int64(2)}),
				d(path("[B]"), diff.CommonDiff{A: "b1", B: "b2"}),
				d(path("[C]", "[A]", "[2]", "[D]"), diff.CommonDiff{A: uint64(3), B: uint64(4)}),
				d(path("[D]", "[A]"), diff.KeyNotFoundDiff{Key: "A", A: true, B: false}),
				d(path("[D]", "[B]"), diff.KeyNotFoundDiff{Key: "B", A: true, B: false}),
				d(path("[D]", "[C]"), diff.KeyNotFoundDiff{Key: "C", A: true, B: false}),
				d(path("[D]", "[D]"), diff.KeyNotFoundDiff{Key: "D", A: true, B: false}),
				d(path("[D]", "[E]"), diff.KeyNotFoundDiff{Key: "E", A: true, B: false}),
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
				d(path("[D]", "[B]", "[x]"), diff.KeyNotFoundDiff{Key: "x", A: true, B: false}),
				d(path("[D]", "[B]", "[y]"), diff.KeyNotFoundDiff{Key: "y", A: true, B: false}),
				d(path("[D]", "[B]", "[z]"), diff.KeyNotFoundDiff{Key: "z", A: true, B: false}),
				d(path("[D]", "[E]"), diff.CommonDiff{A: false, B: true}),
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
			b: map[string]interface{}{
				"A": 2,
				"B": "b2",
				"C": map[string]interface{}{
					"A": []interface{}{
						map[string]interface{}{"A": nil, "B": nil, "C": nil, "D": uint8(1), "E": false},
						OtherStruct{D: 2},
						OtherStruct{D: 4}},
					"B": nil,
					"C": nil,
					"D": uint8(0),
					"E": false,
				},
				"D": &os2,
			},
			result: diffs(
				d(path("[A]"), diff.CommonDiff{A: int64(1), B: int64(2)}),
				d(path("[B]"), diff.CommonDiff{A: "b1", B: "b2"}),
				d(path("[C]", "[A]", "[2]", "[D]"), diff.CommonDiff{A: uint64(3), B: uint64(4)}),
				d(path("[D]", "[B]", "[x]"), diff.KeyNotFoundDiff{Key: "x", A: true, B: false}),
				d(path("[D]", "[B]", "[y]"), diff.KeyNotFoundDiff{Key: "y", A: true, B: false}),
				d(path("[D]", "[B]", "[z]"), diff.KeyNotFoundDiff{Key: "z", A: true, B: false}),
				d(path("[D]", "[E]"), diff.CommonDiff{A: false, B: true}),
			),
		},
	}

	for _, tc := range testCases {
		// When
		d := diff.Similar(tc.a, tc.b, true)
		// Then
		if !reflect.DeepEqual(diffMap(d), diffMap(tc.result)) {
			t.Logf("%v \n\n", diffMap(d))
			t.Logf("%v \n\n", diffMap(tc.result))
			t.Fail()
		}
	}
}
