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
			assertFunc: func(assert assertion.Assert) { assert.That("").Not().IsBlank() },
			errLog:     "\nValue should not be blank string",
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
