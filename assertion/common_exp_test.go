package assertion_test

import (
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

func Test_Common_Matchers_should_fail(t *testing.T) {
	// Mock preparation
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

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
