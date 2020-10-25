package assertion

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
