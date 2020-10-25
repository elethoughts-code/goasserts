// Assertion is handled by two main interfaces that act as a Fluent API.
//
// Assert interface for value holding.
//
// Expectation interface to apply transformations, expectations, matches. etc on the held value.
//
package assertion

import (
	"fmt"
)

// Assert is the assertion entry point.
// It takes a variable as a value and hold it to apply transformations and expectations to it.
//
// That should takes a value and return a wrapped Expectation around it.
type Assert interface {
	That(v interface{}) Expectation
}

// Expectation interface have three roles. It changes Expectation state by setting negation, changing
// expectation result message, etc. It apply transformations on assertion value. It executes expectations.
//
// Not() sets current expectation state to negation.
// thus `assert.That("123").Not().IsNil()` will pass
// and assert.That(1).Not().IsEq(1) will fail.
//
// OrFatal() tell the expectation execution to fail with t.Fatal instead of t.Error.
//
// Silent() tell the expectation to fail without logging its message.
//
// Logf() and Log() set custom failing message.
type Expectation interface {
	Not() Expectation
	OrFatal() Expectation
	Silent() Expectation
	Logf(format string, args ...interface{}) Expectation
	Log(log string) Expectation
	CommonExpectation
	LengthExpectation
	StringExpectation
	SliceExpectation
	MapExpectation
	MapTransformer
	AttributeParser
	HTTPRecorderParser
}

type assert struct {
	t  PublicTB
	br bytesReader
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

// New is an Assert builder. it takes a *testing.T like variable (uses PublicTB interface).
func New(t PublicTB) Assert {
	return newWithBr(t, stdBytesReader{})
}

func newWithBr(t PublicTB, br bytesReader) Assert {
	return &assert{
		t:  t,
		br: br,
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

func runMatcher(m Matcher, v interface{}) (mr MatchResult, err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = fmt.Errorf("[panic error occurred] %w", e)
			} else {
				err = fmt.Errorf("[panic error occurred] %v", r) //nolint:goerr113
			}
		}
	}()
	mr, err = m(v)
	return mr, err
}

func (exp *expectation) Matches(m Matcher) {
	exp.t.Helper()
	mr, err := runMatcher(m, exp.v)
	if err != nil {
		exp.t.Fatalf("\n%s", err.Error())
		return
	}
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
