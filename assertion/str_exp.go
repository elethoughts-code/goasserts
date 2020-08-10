package assertion

func (exp *expectation) IsBlank() {
	exp.t.Helper()
	exp.Matches(IsBlank())
}