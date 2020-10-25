package assertion

import "reflect"

// CommonExpectation interface hold commonly used expectations.
// It also exposes Matches(m Matcher) as a generic expectation method.
// Nearly all expectations interfaces are shortcuts for Matches(m Matcher).
//
// IsEq(e interface{}) expects simple "==" equality.
//
// IsDeepEq(e interface{}) expects value and expectation to be deep equal (reflect.DeepEq).
//
// IsNil() expects value to be nil.
//
// HaveKind(k reflect.Kind) expects value to be of the given reflection Kind.
//
// IsError(target error) and AsError(target interface{}) apply errors.IsError and errors.AsError checks/
//
// NoDiff(e interface{}) uses  diff.Diffs(v, e) to check equality. When the expectation fails,
// it log the deltas detected between the value and the expectation.
type CommonExpectation interface {
	Matches(m Matcher)
	IsEq(e interface{})
	IsDeepEq(e interface{})
	NoDiff(e interface{})
	IsNil()
	HaveKind(k reflect.Kind)
	IsError(target error)
	AsError(target interface{})
}

func (exp *expectation) IsEq(e interface{}) {
	exp.t.Helper()
	exp.Matches(IsEq(e))
}

func (exp *expectation) IsDeepEq(e interface{}) {
	exp.t.Helper()
	exp.Matches(IsDeepEq(e))
}

func (exp *expectation) NoDiff(e interface{}) {
	exp.t.Helper()
	exp.Matches(NoDiff(e))
}

func (exp *expectation) IsNil() {
	exp.t.Helper()
	exp.Matches(IsNil())
}

func (exp *expectation) HaveKind(k reflect.Kind) {
	exp.t.Helper()
	exp.Matches(HaveKind(k))
}

func (exp *expectation) IsError(target error) {
	exp.t.Helper()
	exp.Matches(IsError(target))
}

func (exp *expectation) AsError(target interface{}) {
	exp.t.Helper()
	exp.Matches(AsError(target))
}
