package assertion

import (
	"reflect"

	"github.com/elethoughts-code/goasserts/diff"
)

// SliceExpectation interface encloses slice related expectations.
//
// Contains(e interface{}) check if e parameter is into the slice.
//
// Unordered(e interface{}) check if all elements into the e parameter (should be a slice)
// are into the slice regardless of the order.
//
// All(m Matcher) check if all elements of the value slice matches.
//
// AtLeast(n int, m Matcher) check if at least n elements from the value slice matches.
//
// Any(m Matcher) check if at any element from the value slice matches.
type SliceExpectation interface {
	Contains(e interface{})
	Unordered(e interface{})
	UnorderedDeepEq(e interface{})
	UnorderedNoDiff(e interface{})
	All(m func(v interface{}) bool)
	AtLeast(n int, m func(v interface{}) bool)
	Any(m func(v interface{}) bool)
	Every(matchers []func(v interface{}) bool)
}

func (exp *expectation) Contains(e interface{}) {
	exp.t.Helper()
	exp.Matches(Contains(e))
}

func (exp *expectation) Unordered(e interface{}) {
	exp.t.Helper()
	exp.Matches(Unordered(e, func(v, e interface{}) bool {
		return v == e
	}))
}

func (exp *expectation) UnorderedDeepEq(e interface{}) {
	exp.t.Helper()
	exp.Matches(Unordered(e, reflect.DeepEqual))
}

func (exp *expectation) UnorderedNoDiff(e interface{}) {
	exp.t.Helper()
	exp.Matches(Unordered(e, func(v, e interface{}) bool {
		diffs := diff.Diffs(v, e)
		return len(diffs) == 0
	}))
}

func (exp *expectation) All(m func(v interface{}) bool) {
	exp.t.Helper()
	exp.Matches(All(m))
}

func (exp *expectation) AtLeast(n int, m func(v interface{}) bool) {
	exp.t.Helper()
	exp.Matches(AtLeast(n, m))
}

func (exp *expectation) Any(m func(v interface{}) bool) {
	exp.t.Helper()
	exp.Matches(AtLeast(1, m))
}

func (exp *expectation) Every(matchers []func(v interface{}) bool) {
	exp.t.Helper()
	exp.Matches(Unordered(matchers, func(v, e interface{}) bool {
		currentM, ok := e.(func(v interface{}) bool)
		if !ok {
			panic("expectation variable e should be of type 'func(v interface{}) bool'")
		}
		return currentM(v)
	}))
}
