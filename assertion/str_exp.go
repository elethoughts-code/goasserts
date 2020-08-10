package assertion

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