package assertion

// StringExpectation interface encloses string related expectations.
//
// IsBlank() check if the value is a blank string.
//
// MatchRe(reg string) applies regex on the value string.
type StringExpectation interface {
	IsBlank()
	MatchRe(reg string)
	HasPrefix(prefix string)
	HasSuffix(suffix string)
}

func (exp *expectation) IsBlank() {
	exp.t.Helper()
	exp.Matches(IsBlank())
}

func (exp *expectation) MatchRe(reg string) {
	exp.t.Helper()
	exp.Matches(MatchRe(reg))
}

func (exp *expectation) HasPrefix(prefix string) {
	exp.t.Helper()
	exp.Matches(HasPrefix(prefix))
}

func (exp *expectation) HasSuffix(suffix string) {
	exp.t.Helper()
	exp.Matches(HasSuffix(suffix))
}
