package assertion

type LengthExpectation interface {
	HasLen(len int)
	HasMaxLen(len int)
	HasMinLen(len int)
	IsEmpty()
}

func (exp *expectation) HasLen(len int) {
	exp.t.Helper()
	exp.Matches(HasLen(len))
}

func (exp *expectation) HasMaxLen(len int) {
	exp.t.Helper()
	exp.Matches(HasMaxLen(len))
}

func (exp *expectation) HasMinLen(len int) {
	exp.t.Helper()
	exp.Matches(HasMinLen(len))
}

func (exp *expectation) IsEmpty() {
	exp.t.Helper()
	exp.Matches(IsEmpty())
}