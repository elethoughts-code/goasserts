package assertion_test

import (
	"testing"

	"github.com/elethoughts-code/goasserts/assertion"
	mocks "github.com/elethoughts-code/goasserts/mocks/assertion"
	"github.com/golang/mock/gomock"
)

func Test_IsBlank_should_pass(t *testing.T) {
	// Given
	assert := assertion.New(t)

	// When

	assert.That("").IsBlank()

	assert.That([]string{" "}).Not().IsBlank()
	assert.That("123").Not().IsBlank()

	// Incompatible types (these will be treated with feature #2)
	assert.That(666).Not().IsBlank()
	assert.That(struct{}{}).Not().IsBlank()
	assert.That(false).Not().IsBlank()

	// Then nothing
}

func Test_MatchRe_should_pass(t *testing.T) {
	// Given
	assert := assertion.New(t)

	// When

	assert.That("123456").MatchRe("^\\d+$")
	assert.That("123456a").Not().MatchRe("^\\d+$")

	assert.That([]string{" "}).Not().MatchRe("^\\d+$")

	// Incompatible types (these will be treated with feature #2)
	assert.That(666).Not().MatchRe("^\\d+$")
	assert.That(struct{}{}).Not().MatchRe("^\\d+$")
	assert.That(false).Not().MatchRe("^\\d+$")

	// Then nothing
}

func Test_String_Matchers_should_fail(t *testing.T) {
	// Mock preparation
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testEntries := []struct {
		assertFunc func(assert assertion.Assert)
		errLog     string
	}{
		{
			assertFunc: func(assert assertion.Assert) { assert.That("abc").IsBlank() },
			errLog:     "\nValue is not a blank string : abc",
		},
		{
			assertFunc: func(assert assertion.Assert) { assert.That("123456").Not().MatchRe("^\\d+$") },
			errLog:     "\nValue should not match regexp : ^\\d+$",
		},
		{
			assertFunc: func(assert assertion.Assert) { assert.That("123456a").MatchRe("^\\d+$") },
			errLog:     "\nValue do not match regexp : ^\\d+$",
		},
		{
			assertFunc: func(assert assertion.Assert) { assert.That(123456).MatchRe("^\\d+$") },
			errLog:     "\nValue type is not a string : 123456",
		},
		{
			assertFunc: func(assert assertion.Assert) { assert.That("123456").MatchRe("^[\\d+$") },
			errLog:     "\nCannot match for regexp : ^[\\d+$",
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
