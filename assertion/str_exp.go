package assertion

// StringExpectation interface encloses string related expectations.
//
// IsBlank() check if the value is a blank string.
//
// MatchRe(reg string) applies regex on the value string.
type StringExpectation interface {
	IsBlank()
	MatchRe(reg string)
}

func (exp *expectation) IsBlank() {
	exp.t.Helper()
	exp.Matches(IsBlank())
}

func (exp *expectation) MatchRe(reg string) {
	exp.t.Helper()
	exp.Matches(MatchRe(reg))
}
