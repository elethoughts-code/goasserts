package assertion

import (
	"fmt"
)

type Assert interface {
	That(v interface{}) Expectation
}

type CommonExpectation interface {
	Matches(m Matcher)
	IsEq(e interface{})
	IsDeepEq(e interface{})
	IsNil()
}

type Expectation interface {
	Not() Expectation
	OrFatal() Expectation
	Silent() Expectation
	Logf(format string, args ...interface{}) Expectation
	Log(log string) Expectation
	CommonExpectation
	LengthExpectation
	StringExpectation
}

type assert struct {
	t PublicTB
}

type expectation struct {
	*assert
	v        interface{}
	log      string
	negation bool
	isFatal  bool
	silent   bool
}

func (a *assert) That(v interface{}) Expectation {
	return &expectation{
		assert:   a,
		v:        v,
		log:      "",
		negation: false,
		isFatal:  false,
		silent:   false,
	}
}

func New(t PublicTB) Assert {
	return &assert{
		t: t,
	}
}

func (exp *expectation) Not() Expectation {
	exp.negation = true
	return exp
}

func (exp *expectation) OrFatal() Expectation {
	exp.isFatal = true
	return exp
}

func (exp *expectation) Silent() Expectation {
	exp.silent = true
	return exp
}

func (exp *expectation) Logf(format string, args ...interface{}) Expectation {
	exp.log = fmt.Sprintf(format, args...)
	return exp
}

func (exp *expectation) Log(log string) Expectation {
	exp.log = log
	return exp
}

func (exp *expectation) handleFailure() {
	exp.t.Helper()
	switch {
	case !exp.silent && exp.isFatal && exp.log != "":
		exp.t.Fatal(exp.log)
	case !exp.silent && exp.log != "":
		exp.t.Error(exp.log)
	case exp.isFatal:
		exp.t.FailNow()
	default:
		exp.t.Fail()
	}
}

func (exp *expectation) Matches(m Matcher) {
	exp.t.Helper()
	mr := m(exp.v)
	fail := exp.negation == mr.Matches
	if fail {
		if exp.log == "" && exp.negation {
			exp.Log(mr.NLog)
		} else if exp.log == "" {
			exp.Log(mr.Log)
		}
		exp.handleFailure()
	}
}