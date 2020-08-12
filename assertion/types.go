package assertion

import "errors"

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

var ErrNotOfErrorType = errors.New("value is not of type error")
var ErrNotOfLenType = errors.New("value type should be Array, Slice, String or Map")
var ErrNotOfSliceType = errors.New("value should be a slice")
var ErrNotOfMapType = errors.New("value should be a map")
var ErrNotOfStringType = errors.New("value should be a string")
