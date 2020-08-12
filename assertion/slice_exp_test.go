package assertion_test

import (
	"testing"

	"github.com/elethoughts-code/goasserts/assertion"
	mocks "github.com/elethoughts-code/goasserts/mocks/assertion"
	"github.com/golang/mock/gomock"
)

func Test_Contains_should_pass(t *testing.T) {
	// Given
	assert := assertion.New(t)

	// When

	assert.That([]string{"a", "b", "c"}).Contains("b")
	assert.That([]string{"a", "b", "c"}).Not().Contains("d")
	assert.That([]string{}).Not().Contains("d")
	assert.That("abcd").Not().Contains("d")

	assert.That([]struct{ a string }{{"a"}, {"b"}, {"c"}}).Contains(struct{ a string }{"b"})

	// Then nothing
}

func Test_Slice_Matchers_should_fail(t *testing.T) {
	// Mock preparation
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testEntries := []struct {
		assertFunc func(assert assertion.Assert)
		errLog     string
	}{
		{
			assertFunc: func(assert assertion.Assert) { assert.That([]string{"a", "b", "c"}).Contains("d") },
			errLog:     "\nValue should contains element : d",
		},
		{
			assertFunc: func(assert assertion.Assert) { assert.That([]string{"a", "b", "c"}).Not().Contains("b") },
			errLog:     "\nValue should not contains element : b",
		},
		{
			assertFunc: func(assert assertion.Assert) { assert.That([]string{}).Contains("b") },
			errLog:     "\nValue should contains element : b",
		},
		{
			assertFunc: func(assert assertion.Assert) { assert.That("abc").Contains("b") },
			errLog:     "\nValue should be a slice",
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
