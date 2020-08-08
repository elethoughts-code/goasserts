package assertion

import (
	"reflect"
)

// PublicTB ...
type PublicTB interface {
	Cleanup(func())
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fail()
	FailNow()
	Failed() bool
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Helper()
	Log(args ...interface{})
	Logf(format string, args ...interface{})
	Name() string
	Skip(args ...interface{})
	SkipNow()
	Skipf(format string, args ...interface{})
	Skipped() bool
}

type Assert interface {
	That(v interface{}) Expectation
}

type Expectation interface {
	Not() Expectation
	OrFatal() Expectation
	Silent() Expectation
	Logf(format string, args ...interface{}) Expectation
	IsEq(e interface{})
	IsDeepEq(e interface{})
}

type assert struct {
	t PublicTB
}

type expectation struct {
	*assert
	v        interface{}
	log      *logStruct
	negation bool
	isFatal  bool
	silent   bool
}

func (a *assert) That(v interface{}) Expectation {
	return &expectation{
		assert:   a,
		v:        v,
		log:      nil,
		negation: false,
		isFatal:  false,
		silent:   false,
	}
}

// New ...
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

type logStruct struct {
	format string
	args   []interface{}
}

func logf(format string, args ...interface{}) logStruct {
	return logStruct{
		format: format,
		args:   args,
	}
}

func (exp *expectation) Logf(format string, args ...interface{}) Expectation {
	log := logf(format, args...)
	exp.log = &log
	return exp
}

func (exp *expectation) failWithNegation(failCondition bool) bool {
	return exp.negation != failCondition
}

func (exp *expectation) defaultLog(straightLog, negatedLog logStruct) {
	if exp.log == nil && !exp.negation {
		exp.Logf(straightLog.format, straightLog.args...)
	} else if exp.log == nil && exp.negation {
		exp.Logf(negatedLog.format, negatedLog.args...)
	}
}

func (exp *expectation) handleFailure() {
	exp.t.Helper()
	switch {
	case !exp.silent && exp.isFatal && exp.log != nil:
		exp.t.Fatalf(exp.log.format, exp.log.args...)
	case !exp.silent && exp.log != nil:
		exp.t.Errorf(exp.log.format, exp.log.args...)
	case exp.isFatal:
		exp.t.FailNow()
	default:
		exp.t.Fail()
	}
}

func (exp *expectation) IsEq(e interface{}) {
	exp.t.Helper()
	fail := e != exp.v
	if exp.failWithNegation(fail) {
		exp.defaultLog(
			logf("\nValue is not equal to expectation.\nExpected : %v\nGot : %v", e, exp.v),
			logf("\nValue should not be equal to : %v", e))
		exp.handleFailure()
	}
}

func (exp *expectation) IsDeepEq(e interface{}) {
	exp.t.Helper()
	fail := !reflect.DeepEqual(e, exp.v)
	if exp.failWithNegation(fail) {
		exp.defaultLog(
			logf("\nValue is not deep equal to expectation.\nExpected : %v\nGot : %v", e, exp.v),
			logf("\nValue should not be deep equal to : %v", e))
		exp.handleFailure()
	}
}
