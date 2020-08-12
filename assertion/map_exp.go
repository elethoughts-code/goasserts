package assertion

type MapExpectation interface {
	ContainsValue(e interface{})
	ContainsKey(e interface{})
}

func (exp *expectation) ContainsValue(e interface{}) {
	exp.t.Helper()
	exp.Matches(ContainsValue(e))
}

func (exp *expectation) ContainsKey(e interface{}) {
	exp.t.Helper()
	exp.Matches(ContainsKey(e))
}
