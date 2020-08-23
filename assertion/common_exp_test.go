package assertion_test

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/elethoughts-code/goasserts/assertion"
	mocks "github.com/elethoughts-code/goasserts/mocks/assertion"
	"github.com/golang/mock/gomock"
)

func Test_IsNil_should_pass(t *testing.T) {
	// Given
	assert := assertion.New(t)

	// When

	assert.That(nil).IsNil()

	assert.That([]string{" "}).Not().IsNil()
	assert.That("123").Not().IsNil()

	// Then nothing
}

func Test_HaveKind_should_pass(t *testing.T) {
	// Given
	assert := assertion.New(t)

	// When

	assert.That("").HaveKind(reflect.String)
	assert.That(123).HaveKind(reflect.Int)
	assert.That([]string{" "}).HaveKind(reflect.Slice)

	assert.That("123").Not().HaveKind(reflect.Slice)

	// Then nothing
}

type myError struct {
	msg     string
	wrapped error
}

func (myError *myError) Error() string { return myError.msg }
func (myError *myError) Unwrap() error { return myError.wrapped }

func Test_IsError_should_pass(t *testing.T) {
	// Given
	assert := assertion.New(t)

	err1 := errors.New("some error")
	err2 := fmt.Errorf("some other error > %w", err1)

	err3 := &myError{
		msg:     "my error",
		wrapped: err2,
	}

	// When
	assert.That(err1).IsError(err1)
	assert.That(err2).IsError(err2)
	assert.That(err3).IsError(err3)

	assert.That(err3).IsError(err2)
	assert.That(err2).IsError(err1)
	assert.That(err3).IsError(err1)

	assert.That(err1).Not().IsError(err2)
	assert.That(err2).Not().IsError(err3)
	assert.That(err3).Not().IsError(fmt.Errorf("error 4"))

	// Then nothing
}

func Test_AsError_should_pass(t *testing.T) {
	// Given
	assert := assertion.New(t)

	err1 := errors.New("some error")
	err2 := fmt.Errorf("some other error > %w", err1)

	err3 := &myError{
		msg:     "my error",
		wrapped: err2,
	}

	// When
	var e *myError

	assert.That(err3).AsError(&e)
	assert.That(err1).Not().AsError(&e)
	assert.That(err2).Not().AsError(&e)

	// Then nothing
}

func Test_Common_Matchers_should_fail(t *testing.T) {
	// Mock preparation
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	var e *myError

	testEntries := []struct {
		assertFunc func(assert assertion.Assert)
		errLog     string
	}{
		{
			assertFunc: func(assert assertion.Assert) { assert.That("abc").IsNil() },
			errLog:     "\nValue is not nil : abc",
		},
		{
			assertFunc: func(assert assertion.Assert) { assert.That(nil).Not().IsNil() },
			errLog:     "\nValue should not be nil",
		},

		{
			assertFunc: func(assert assertion.Assert) { assert.That("abc").HaveKind(reflect.Bool) },
			errLog:     "\nValue is not of the expected Kind.\nExpected : bool\nGot : string",
		},
		{
			assertFunc: func(assert assertion.Assert) { assert.That(nil).HaveKind(reflect.Func) },
			errLog:     "\nValue is not of the expected Kind.\nExpected : func\nGot : invalid",
		},
		{
			assertFunc: func(assert assertion.Assert) { assert.That("abc").Not().HaveKind(reflect.String) },
			errLog:     "\nValue should not of Kind : string",
		},

		{
			assertFunc: func(assert assertion.Assert) {
				assert.That(errors.New("some error")).AsError(&e)
			},
			errLog: "\nError Value is not as the expected type",
		},
		{
			assertFunc: func(assert assertion.Assert) {
				err := &myError{}
				assert.That(err).Not().AsError(&e)
			},
			errLog: "\nError value should not be as expected type",
		},

		{
			assertFunc: func(assert assertion.Assert) {
				assert.That(errors.New("error 1")).IsError(errors.New("error 2"))
			},
			errLog: fmt.Sprintf("\nError Value is not of the expected type.\nExpected : %v\nGot : %v",
				errors.New("error 2"),
				errors.New("error 1")),
		},

		{
			assertFunc: func(assert assertion.Assert) {
				err1 := errors.New("error 1")
				err2 := fmt.Errorf("%w", err1)
				assert.That(err2).Not().IsError(err1)
			},
			errLog: fmt.Sprintf("\nError value should not be : %v", errors.New("error 1")),
		},
	}

	for _, entry := range testEntries {
		// Given
		tMock := mocks.NewMockPublicTB(ctrl)
		assert := assertion.New(tMock)

		// Expectation
		tMock.EXPECT().Helper().AnyTimes()
		tMock.EXPECT().Error(entry.errLog)

		// When
		entry.assertFunc(assert)
	}
}

func Test_Common_Matchers_should_fail_with_error(t *testing.T) {
	// Mock preparation
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	var e *myError

	testEntries := []struct {
		assertFunc func(assert assertion.Assert)
		err        error
		times      int
	}{
		{
			assertFunc: func(assert assertion.Assert) {
				assert.That(123).AsError(&e)
				assert.That(123).IsError(errors.New("my error"))
				assert.That(123).Not().AsError(&e)
				assert.That(123).Not().IsError(errors.New("my error"))
			},
			err:   assertion.ErrNotOfErrorType,
			times: 4,
		},
		{
			assertFunc: func(assert assertion.Assert) {
				assert.That([]string{}).IsEq([]string{})
				assert.That([]string{}).Not().IsEq([]string{})
			},
			err:   errors.New("[panic error occurred] runtime error: comparing uncomparable type []string"),
			times: 2,
		},
		{
			assertFunc: func(assert assertion.Assert) {
				assert.That([]string{}).Matches(func(v interface{}) (assertion.MatchResult, error) {
					panic("for some reason")
				})
				assert.That([]string{}).Not().Matches(func(v interface{}) (assertion.MatchResult, error) {
					panic("for some reason")
				})
			},
			err:   errors.New("[panic error occurred] for some reason"),
			times: 2,
		},
	}

	for _, entry := range testEntries {
		// Given
		tMock := mocks.NewMockPublicTB(ctrl)
		assert := assertion.New(tMock)

		// Expectation
		tMock.EXPECT().Helper().AnyTimes()
		tMock.EXPECT().Fatalf("\n%s", entry.err.Error()).Times(entry.times)

		// When
		entry.assertFunc(assert)
	}
}

type SampleStruct struct {
	A int
	B string
	C OtherStruct
	D *OtherStruct
}

type OtherStruct struct {
	A []OtherStruct
	B map[string][]int
	C interface{}
	D uint8
	E bool
}

func Test_NoDiff_should_pass(t *testing.T) {
	// Given
	assert := assertion.New(t)
	os := OtherStruct{D: 0, E: true}
	os2 := OtherStruct{D: 0, E: true}
	a := SampleStruct{
		A: 1,
		B: "b1",
		C: OtherStruct{
			A: []OtherStruct{{D: 1}, {D: 2}, {D: 3}},
			B: nil,
			C: nil,
			D: 0,
			E: false,
		},
		D: &os,
	}
	b := SampleStruct{
		A: 1,
		B: "b1",
		C: OtherStruct{
			A: []OtherStruct{{D: 1}, {D: 2}, {D: 3}},
			B: nil,
			C: nil,
			D: 0,
			E: false,
		},
		D: &os,
	}
	c := SampleStruct{
		A: 2,
		B: "b2",
		C: OtherStruct{
			A: []OtherStruct{{D: 1}, {D: 2}, {D: 4}},
			B: nil,
			C: nil,
			D: 0,
			E: false,
		},
		D: &os2,
	}
	// When
	assert.That(a).NoDiff(b)
	assert.That(a).Not().NoDiff(c)
}

func Test_NoDiff_should_not_pass_1(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	os := OtherStruct{D: 0, E: true}
	a := SampleStruct{
		A: 1,
		B: "b1",
		C: OtherStruct{
			A: []OtherStruct{{D: 1}, {D: 2}, {D: 3}},
			B: nil,
			C: nil,
			D: 0,
			E: false,
		},
		D: &os,
	}
	b := SampleStruct{
		A: 1,
		B: "b1",
		C: OtherStruct{
			A: []OtherStruct{{D: 1}, {D: 2}, {D: 3}},
			B: nil,
			C: nil,
			D: 0,
			E: false,
		},
		D: &os,
	}

	// Expectation
	tMock := mocks.NewMockPublicTB(ctrl)
	assert := assertion.New(tMock)
	tMock.EXPECT().Helper().AnyTimes()
	tMock.EXPECT().Error("Value should have diffs with expectation")

	// When
	assert.That(a).Not().NoDiff(b)
}

func Test_NoDiff_should_not_pass_2(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	os := OtherStruct{D: 0, E: true}
	os2 := OtherStruct{D: 0, E: true}
	a := SampleStruct{
		A: 1,
		B: "b1",
		C: OtherStruct{
			A: []OtherStruct{{D: 1}, {D: 2}, {D: 3}},
			B: nil,
			C: nil,
			D: 0,
			E: false,
		},
		D: &os,
	}
	b := SampleStruct{
		A: 2,
		B: "b2",
		C: OtherStruct{
			A: []OtherStruct{{D: 1}, {D: 2}, {D: 4}},
			B: nil,
			C: nil,
			D: 0,
			E: false,
		},
		D: &os2,
	}

	// Expectation
	tMock := mocks.NewMockPublicTB(ctrl)
	assert := assertion.New(tMock)
	tMock.EXPECT().Helper().AnyTimes()
	tMock.EXPECT().Error("Value have following diffs with expectation :\n" +
		"Path [[A]] : values diff\n" +
		"A=1\n" +
		"B=2\n" +
		"Path [[B]] : values diff\n" +
		"A=b1\n" +
		"B=b2\n" +
		"Path [[C] [A] [2] [D]] : values diff\n" +
		"A=3\n" +
		"B=4")

	// When
	assert.That(a).NoDiff(b)
}
