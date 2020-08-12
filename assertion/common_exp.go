package assertion

import "reflect"

type CommonExpectation interface {
	Matches(m Matcher)
	IsEq(e interface{})
	IsDeepEq(e interface{})
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
