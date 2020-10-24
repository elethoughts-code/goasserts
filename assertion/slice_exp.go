package assertion

type SliceExpectation interface {
	Contains(e interface{})
	Unordered(e interface{})
	All(m Matcher)
	AtLeast(n int, m Matcher)
	Any(m Matcher)
}

func (exp *expectation) Contains(e interface{}) {
	exp.t.Helper()
	exp.Matches(Contains(e))
}

func (exp *expectation) Unordered(e interface{}) {
	exp.t.Helper()
	exp.Matches(Unordered(e))
}

func (exp *expectation) All(m Matcher) {
	exp.t.Helper()
	exp.Matches(All(m))
}

func (exp *expectation) AtLeast(n int, m Matcher) {
	exp.t.Helper()
	exp.Matches(AtLeast(n, m))
}

func (exp *expectation) Any(m Matcher) {
	exp.t.Helper()
	exp.Matches(AtLeast(1, m))
}
