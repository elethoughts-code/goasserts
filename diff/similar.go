package diff

import (
	"fmt"
	"reflect"
)

// Similar function returns all extracted dissimilarities between two variables a and b.
// Similar uses field based equality check :
//
// - It do not check types.
// - It compares structs to maps.
// - Empty slices and maps are equal to nils.
func Similar(a, b interface{}, checkUnordered bool) (diffs []Diff) {
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

	findSimilarityDiffs(path, va, vb, &diffs, make(map[similarVisit]bool), checkUnordered)

	return diffs
}

func dereference(v reflect.Value) reflect.Value {
	switch v.Kind() {
	case reflect.Interface, reflect.Ptr:
		if v.IsNil() {
			return v
		}
		return dereference(v.Elem())
	default:
		return v
	}
}

func isIndexed(v reflect.Value, k reflect.Kind) (int, bool) {
	switch k {
	case reflect.Array, reflect.Slice:
		return v.Len(), true
	default:
		return 0, false
	}
}

func isFielded(v reflect.Value, k reflect.Kind) (map[string]reflect.Value, bool) {
	switch {
	case k == reflect.Map && v.Type().Key().Kind() == reflect.String:
		fields := make(map[string]reflect.Value, v.Len())
		for _, k := range v.MapKeys() {
			fields[k.String()] = v.MapIndex(k)
		}
		return fields, true
	case k == reflect.Struct:
		t := v.Type()
		nbFields := t.NumField()
		fields := make(map[string]reflect.Value, nbFields)
		for i := 0; i < nbFields; i++ {
			fName := t.Field(i).Name
			fields[fName] = v.FieldByName(fName)
		}
		return fields, true
	default:
		return nil, false
	}
}

func isNil(v reflect.Value, k reflect.Kind) bool {
	switch k {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return v.IsNil()
	default:
		return false
	}
}

func slicesOfComparable(v reflect.Value, len int) ([]interface{}, reflect.Type, bool) {
	var typ reflect.Type = nil
	s := make([]interface{}, len)
	for i := 0; i < len; i++ {
		cVal := v.Index(i)
		cdVal := dereference(cVal)
		cTyp := cdVal.Type()
		if !cTyp.Comparable() {
			return nil, nil, false
		}

		if typ == nil {
			typ = cTyp
		} else if typ.Kind() != cTyp.Kind() {
			return nil, nil, false
		}
		s[i] = cdVal.Interface()
	}
	return s, typ, true
}

func unorderedEq(a, b []interface{}) bool {
	for _, expectedItem := range a {
		found := false
		for _, sliceItem := range b {
			if sliceItem == expectedItem {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// nolint:gocognit,gocyclo,nestif
func findSimilarityDiffs(currentPath []string, va, vb reflect.Value, diffs *[]Diff,
	visited map[similarVisit]bool, checkUnordered bool) {
	if !va.IsValid() || !vb.IsValid() {
		*diffs = append(*diffs, newDiff(currentPath, InvalidDiff{va.IsValid(), vb.IsValid()}))
		return
	}

	if checkSimilarVisited(va, vb, visited) {
		return
	}

	// Same pointer
	if va.Kind() == reflect.Ptr && vb.Kind() == reflect.Ptr && va.Pointer() == vb.Pointer() {
		return
	}

	va = dereference(va)
	vb = dereference(vb)

	ka, kb := va.Kind(), vb.Kind()

	// Check func
	if ka == reflect.Func {
		*diffs = append(*diffs, newDiff(currentPath, FuncDiff{va.Interface(), vb.Interface()}))
		return
	}

	// Check nils
	aIsNil, bIsNil := isNil(va, ka), isNil(vb, kb)

	if aIsNil && bIsNil {
		return
	}

	// Check indexed (array or slice)
	lenA, aIsIndexed := isIndexed(va, ka)
	lenB, bIsIndexed := isIndexed(vb, kb)

	if (aIsIndexed || aIsNil) && (bIsIndexed || bIsNil) {
		lenDiff := lenA - lenB
		if lenDiff != 0 {
			*diffs = append(*diffs, newDiff(currentPath, LenDiff{CommonDiff{va.Interface(), vb.Interface()}, lenDiff}))
			return
		}
		if checkUnordered {
			scA, typA, isSCA := slicesOfComparable(va, lenA)
			scB, typB, isSCB := slicesOfComparable(vb, lenB)
			if isSCA && isSCB && typA == typB {
				if !unorderedEq(scA, scB) {
					*diffs = append(*diffs, newDiff(currentPath, CommonDiff{scA, scB}))
				}
				// If slice of comparable check unordered
				return
			}
		}
		for i := 0; i < lenA; i++ {
			iKey := fmt.Sprintf("[%d]", i)
			findSimilarityDiffs(append(currentPath, iKey), va.Index(i), vb.Index(i), diffs, visited, checkUnordered)
		}
		return
	}

	if aIsIndexed != bIsIndexed {
		*diffs = append(*diffs, newDiff(currentPath, TypeDiff{va.Interface(), vb.Interface()}))
		return
	}

	// Check fielded
	aFields, aIsFielded := isFielded(va, ka)
	bFields, bIsFielded := isFielded(vb, kb)

	if (aIsFielded || aIsNil) && (bIsFielded || bIsNil) {
		for k, aValue := range aFields {
			fieldName := fmt.Sprintf("[%v]", k)

			if bValue, exists := bFields[k]; !exists {
				*diffs = append(*diffs, newDiff(append(currentPath, fieldName),
					KeyNotFoundDiff{Key: fmt.Sprintf("%v", k), A: true, B: false}))
			} else {
				findSimilarityDiffs(append(currentPath, fieldName), aValue, bValue, diffs, visited, checkUnordered)
			}
		}
		for k := range bFields {
			fieldName := fmt.Sprintf("[%v]", k)
			if _, exists := aFields[k]; !exists {
				*diffs = append(*diffs, newDiff(append(currentPath, fieldName),
					KeyNotFoundDiff{Key: fmt.Sprintf("%v", k), A: false, B: true}))
			}
		}
		return
	}

	if aIsFielded != bIsFielded {
		*diffs = append(*diffs, newDiff(currentPath, TypeDiff{va.Interface(), vb.Interface()}))
		return
	}

	// Check simple types
	checkSimpleTypes(currentPath, va, vb, ka, diffs, visited, checkUnordered)
}

func checkSimpleTypes(currentPath []string, va, vb reflect.Value,
	ka reflect.Kind, diffs *[]Diff, visited map[similarVisit]bool, checkUnordered bool) {
	ta, tb := va.Type(), vb.Type()
	if ta != tb {
		*diffs = append(*diffs, newDiff(currentPath, TypeDiff{va.Interface(), vb.Interface()}))
		return
	}
	switch ka {
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
	case reflect.Map:
		if ta.Key().Kind() == reflect.String {
			panic(fmt.Errorf("should not have a string keyed map")) //nolint:goerr113
		}
		checkSimilarMaps(currentPath, va, vb, diffs, visited, checkUnordered)
	default:
		panic(fmt.Errorf("should not have kind : %v", ka)) //nolint:goerr113
	}
}

func checkSimilarMaps(currentPath []string, va, vb reflect.Value, diffs *[]Diff,
	visited map[similarVisit]bool, checkUnordered bool) {
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
			findSimilarityDiffs(append(currentPath, fieldName), va.MapIndex(k), bValue, diffs, visited, checkUnordered)
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

type similarVisit struct {
	lowerAddr  uintptr
	higherAddr uintptr
	lowerTyp   reflect.Type
	higherTyp  reflect.Type
}

func interfaceDereference(v reflect.Value) reflect.Value {
	switch v.Kind() {
	case reflect.Interface:
		if v.IsNil() {
			return v
		}
		return interfaceDereference(v.Elem())
	default:
		return v
	}
}

func checkSimilarVisited(va, vb reflect.Value, visited map[similarVisit]bool) bool {
	// see reflect.DeepEqual since this code is partially copied from it
	// NOTE: Interface is not treated directly as possible Pointer value but we wait
	// to unwrap it from " case reflect.Interface: " line of the findDiffs function
	// this is because reflect.DeepEqual uses internal and non exported methods

	hard := func(va, vb reflect.Value) bool {
		aN, bN := false, false
		switch va.Kind() {
		case reflect.Map, reflect.Slice, reflect.Ptr:
			aN = !va.IsNil()
		}
		switch vb.Kind() {
		case reflect.Map, reflect.Slice, reflect.Ptr:
			bN = !vb.IsNil()
		}
		return aN && bN
	}

	va = interfaceDereference(va)
	vb = interfaceDereference(vb)

	if hard(va, vb) {
		lowerAddr := va.Pointer()
		higherAddr := vb.Pointer()
		lowerTyp := va.Type()
		higherTyp := vb.Type()
		if lowerAddr > higherAddr {
			// *** reflect.DeepEqual comments (as is) ***
			// Canonicalize order to reduce number of entries in visited.
			// Assumes non-moving garbage collector.
			lowerAddr, higherAddr = higherAddr, lowerAddr
			lowerTyp, higherTyp = higherTyp, lowerTyp
		}

		v := similarVisit{lowerAddr, higherAddr, lowerTyp, higherTyp}
		if visited[v] {
			return true
		}

		visited[v] = true
	}

	return false
}
