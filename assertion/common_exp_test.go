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

func (me *myError) Error() string { return me.msg }
func (me *myError) Unwrap() error { return me.wrapped }

func Test_IsError_should_pass(t *testing.T) {
	// Given
	assert := assertion.New(t)

	err1 := errors.New("some error")
	err2 := fmt.Errorf("some other error > %w", err1)

	err3 := &myError{
		msg:     "My error",
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
	assert.That(err3).Not().IsError(fmt.Errorf("Error 4"))

	// Then nothing
}

func Test_AsError_should_pass(t *testing.T) {
	// Given
	assert := assertion.New(t)

	err1 := errors.New("some error")
	err2 := fmt.Errorf("some other error > %w", err1)

	err3 := &myError{
		msg:     "My error",
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
				assert.That(123).AsError(&e)
			},
			errLog: "\nValue is not of type error.",
		},
		{
			assertFunc: func(assert assertion.Assert) {
				assert.That(123).IsError(errors.New("My error"))
			},
			errLog: "\nValue is not of type error.",
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
				assert.That(errors.New("Error 1")).IsError(errors.New("Error 2"))
			},
			errLog: fmt.Sprintf("\nError Value is not of the expected type.\nExpected : %v\nGot : %v",
				errors.New("Error 2"),
				errors.New("Error 1")),
		},

		{
			assertFunc: func(assert assertion.Assert) {
				err1 := errors.New("Error 1")
				err2 := fmt.Errorf("%w", err1)
				assert.That(err2).Not().IsError(err1)
			},
			errLog: fmt.Sprintf("\nError value should not be : %v", errors.New("Error 1")),
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
