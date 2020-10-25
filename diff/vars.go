package diff

import (
	"fmt"
	"reflect"
)

// Diff is a structure pointing a difference on a variable Path (Deep navigation on attributes and indexes).
// A difference is reported in the Value attribute as un error.
type Diff struct {
	Path  []string
	Value error
}

func newDiff(path []string, value error) Diff {
	cPath := make([]string, len(path))
	copy(cPath, path)
	return Diff{
		Path:  cPath,
		Value: value,
	}
}

// CommonDiff is a simple difference error between two values A and B.
type CommonDiff struct {
	A interface{}
	B interface{}
}

func (cd CommonDiff) Error() string {
	return fmt.Sprintf("values diff\nA=%v\nB=%v", cd.A, cd.B)
}

// TypeDiff is a type difference error between two values A and B.
type TypeDiff CommonDiff

func (td TypeDiff) Error() string {
	return fmt.Sprintf("value types diff\nType of A=%v\nType of B=%v", reflect.TypeOf(td.A), reflect.TypeOf(td.B))
}

// TypeDiff is a function difference error between two values A and B. (all functions are considered different).
type FuncDiff CommonDiff

func (fd FuncDiff) Error() string {
	return "functions cannot be compared"
}

// LenDiff reports a length different between A and B.
type LenDiff struct {
	CommonDiff
	Value int
}

func (ld LenDiff) Error() string {
	return fmt.Sprintf("value length diff = %v", ld.Value)
}

// KeyNotFoundDiff reports that a Key exists only on A or B.
type KeyNotFoundDiff struct {
	Key string
	A   bool
	B   bool
}

func (kd KeyNotFoundDiff) Error() string {
	return fmt.Sprintf("key [%s] not found", kd.Key)
}

// InvalidDiff reports that one or both of the values A and B are in an Invalid State.
type InvalidDiff struct {
	A bool
	B bool
}

func (id InvalidDiff) Error() string {
	return "invalid value"
}

// Diffs function returns all extracted differences between to variables a and b.
// It copies most of the standard reflect.DeepEq algorithm (getting around some unexported capabilities).
func Diffs(a, b interface{}) (diffs []Diff) {
	diffs = make([]Diff, 0)
	path := make([]string, 0)
	if a == nil && b == nil {
		return diffs
	}

	if a == nil {
		diffs = append(diffs, Diff{Path: path, Value: CommonDiff{a, b}})
		return diffs
	}

	if b == nil {
		diffs = append(diffs, Diff{Path: path, Value: CommonDiff{a, b}})
		return diffs
	}

	va := reflect.ValueOf(a)
	vb := reflect.ValueOf(b)

	findDiffs(path, va, vb, &diffs, make(map[visit]bool))

	return diffs
}

type visit struct {
	lowerAddr  uintptr
	higherAddr uintptr
	typ        reflect.Type
}

func checkVisited(va, vb reflect.Value, typ reflect.Type, visited map[visit]bool) bool {
	// see reflect.DeepEqual since this code is partially copied from it
	// NOTE: Interface is not treated directly as possible Pointer value but we wait
	// to unwrap it from " case reflect.Interface: " line of the findDiffs function
	// this is because reflect.DeepEqual uses internal and non exported methods
	hard := func(va, vb reflect.Value) bool {
		switch va.Kind() {
		case reflect.Map, reflect.Slice, reflect.Ptr:
			return !va.IsNil() && !vb.IsNil()
		}
		return false
	}

	if hard(va, vb) {
		lowerAddr := va.Pointer()
		higherAddr := vb.Pointer()
		if lowerAddr > higherAddr {
			// *** reflect.DeepEqual comments (as is) ***
			// Canonicalize order to reduce number of entries in visited.
			// Assumes non-moving garbage collector.
			lowerAddr, higherAddr = higherAddr, lowerAddr
		}

		v := visit{lowerAddr, higherAddr, typ}
		if visited[v] {
			return true
		}

		visited[v] = true
	}

	return false
}

const (
	areNil int = iota
	aNil
	bNil
	noNil
)

func checkNilValue(va, vb reflect.Value, currentPath []string, diffs *[]Diff) int {
	if va.IsNil() && vb.IsNil() {
		return areNil
	}
	if va.IsNil() {
		*diffs = append(*diffs, newDiff(currentPath, CommonDiff{nil, vb.Interface()}))
		return aNil
	}
	if vb.IsNil() {
		*diffs = append(*diffs, newDiff(currentPath, CommonDiff{va.Interface(), nil}))
		return bNil
	}
	return noNil
}

func simpleEqDiff(a, b interface{}, currentPath []string, diffs *[]Diff) {
	if a != b {
		*diffs = append(*diffs, newDiff(currentPath, CommonDiff{a, b}))
	}
}

func checkArrays(currentPath []string, va, vb reflect.Value, diffs *[]Diff, visited map[visit]bool) {
	for i := 0; i < va.Len(); i++ {
		iKey := fmt.Sprintf("[%d]", i)
		findDiffs(append(currentPath, iKey), va.Index(i), vb.Index(i), diffs, visited)
	}
}

func checkSlices(currentPath []string, va, vb reflect.Value, diffs *[]Diff, visited map[visit]bool) {
	if checkNilValue(va, vb, currentPath, diffs) != noNil {
		return
	}
	if va.Pointer() == vb.Pointer() {
		return
	}
	lenVa := va.Len()
	lenDiff := lenVa - vb.Len()
	if lenDiff != 0 {
		*diffs = append(*diffs, newDiff(currentPath, LenDiff{CommonDiff{va.Interface(), vb.Interface()}, lenDiff}))
		return
	}

	for i := 0; i < lenVa; i++ {
		iKey := fmt.Sprintf("[%d]", i)
		findDiffs(append(currentPath, iKey), va.Index(i), vb.Index(i), diffs, visited)
	}
}

func checkMaps(currentPath []string, va, vb reflect.Value, diffs *[]Diff, visited map[visit]bool) {
	if checkNilValue(va, vb, currentPath, diffs) != noNil {
		return
	}
	if va.Pointer() == vb.Pointer() {
		return
	}
	lenVa := va.Len()
	lenDiff := lenVa - vb.Len()
	if lenDiff != 0 {
		*diffs = append(*diffs, newDiff(currentPath, LenDiff{CommonDiff{va.Interface(), vb.Interface()}, lenDiff}))
		return
	}
	for _, k := range va.MapKeys() {
		fieldName := fmt.Sprintf("[%v]", k)
		bValue := vb.MapIndex(k)
		if !bValue.IsValid() || bValue.IsZero() {
			*diffs = append(*diffs, newDiff(append(currentPath, fieldName),
				KeyNotFoundDiff{Key: fmt.Sprintf("%v", k), A: true, B: false}))
		} else {
			findDiffs(append(currentPath, fieldName), va.MapIndex(k), bValue, diffs, visited)
		}
	}
	for _, k := range vb.MapKeys() {
		fieldName := fmt.Sprintf("[%v]", k)
		aValue := va.MapIndex(k)
		if !aValue.IsValid() || aValue.IsZero() {
			*diffs = append(*diffs, newDiff(append(currentPath, fieldName),
				KeyNotFoundDiff{Key: fmt.Sprintf("%v", k), A: false, B: true}))
		}
	}
}

func checkStructs(currentPath []string, va, vb reflect.Value, diffs *[]Diff, visited map[visit]bool) {
	t := va.Type()
	nbFields := t.NumField()
	for i := 0; i < nbFields; i++ {
		fName := t.Field(i).Name
		ffName := fmt.Sprintf("[%s]", fName)
		findDiffs(append(currentPath, ffName), va.FieldByName(fName), vb.FieldByName(fName), diffs, visited)
	}
}

func findDiffs(currentPath []string, va, vb reflect.Value, diffs *[]Diff, visited map[visit]bool) {
	if !va.IsValid() || !vb.IsValid() {
		*diffs = append(*diffs, newDiff(currentPath, InvalidDiff{va.IsValid(), vb.IsValid()}))
		return
	}
	ta, tb := va.Type(), vb.Type()
	if ta != tb {
		*diffs = append(*diffs, newDiff(currentPath, TypeDiff{va.Interface(), vb.Interface()}))
		return
	}

	if checkVisited(va, vb, ta, visited) {
		return
	}

	switch va.Kind() {
	case reflect.Array:
		// Array len is part of the Type()
		checkArrays(currentPath, va, vb, diffs, visited)
	case reflect.Slice:
		checkSlices(currentPath, va, vb, diffs, visited)
	case reflect.Interface:
		if checkNilValue(va, vb, currentPath, diffs) != noNil {
			return
		}
		findDiffs(append(currentPath, "[interface{}]"), va.Elem(), vb.Elem(), diffs, visited)
	case reflect.Ptr:
		if va.Pointer() == vb.Pointer() {
			return
		}
		if checkNilValue(va, vb, currentPath, diffs) != noNil {
			return
		}
		findDiffs(append(currentPath, "[&]"), va.Elem(), vb.Elem(), diffs, visited)
	case reflect.Struct:
		checkStructs(currentPath, va, vb, diffs, visited)
	case reflect.Map:
		checkMaps(currentPath, va, vb, diffs, visited)
	case reflect.Func:
		if checkNilValue(va, vb, currentPath, diffs) != noNil {
			return
		}
		// Can't do better than this:
		*diffs = append(*diffs, newDiff(currentPath, FuncDiff{va.Interface(), vb.Interface()}))
	// Note: since Value.Interface() do not return unexported attributes
	// Continue case by case
	case reflect.Bool:
		simpleEqDiff(va.Bool(), vb.Bool(), currentPath, diffs)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		simpleEqDiff(va.Int(), vb.Int(), currentPath, diffs)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		simpleEqDiff(va.Uint(), vb.Uint(), currentPath, diffs)
	case reflect.Float32, reflect.Float64:
		simpleEqDiff(va.Float(), vb.Float(), currentPath, diffs)
	case reflect.Complex64, reflect.Complex128:
		simpleEqDiff(va.Complex(), vb.Complex(), currentPath, diffs)
	case reflect.UnsafePointer:
		simpleEqDiff(va.Interface(), vb.Interface(), currentPath, diffs)
	case reflect.String:
		simpleEqDiff(va.String(), vb.String(), currentPath, diffs)
	case reflect.Chan:
		simpleEqDiff(va.Interface(), vb.Interface(), currentPath, diffs)
	}
}
