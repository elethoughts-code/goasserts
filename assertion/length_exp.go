package assertion

// LengthExpectation apply on Array, Slice, Map and String values.
// It adds expectations related to the value length.
//
// HasLen(len int) check if the value length is equal to the len parameter.
//
// HasMaxLen(len int) check if the value length have not exceeded the len parameter.
//
// HasMinLen(len int) check if the value length is greater than the len parameter.
//
// IsEmpty() check if the value is empty.
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
