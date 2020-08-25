package assertion

type SliceExpectation interface {
	Contains(e interface{})
	Unordered(e interface{})
}

func (exp *expectation) Contains(e interface{}) {
	exp.t.Helper()
	exp.Matches(Contains(e))
}

func (exp *expectation) Unordered(e interface{}) {
	exp.t.Helper()
	exp.Matches(Unordered(e))
}
