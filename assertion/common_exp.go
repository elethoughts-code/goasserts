package assertion

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