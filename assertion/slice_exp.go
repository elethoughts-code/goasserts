package assertion

type SliceExpectation interface {
	Contains(e interface{})
}

func (exp *expectation) Contains(e interface{}) {
	exp.t.Helper()
	exp.Matches(Contains(e))
}