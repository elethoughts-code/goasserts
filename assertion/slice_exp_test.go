package assertion_test

import (
	"fmt"
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

	assert.That([]struct{ a string }{{"a"}, {"b"}, {"c"}}).Contains(struct{ a string }{"b"})

	// Then nothing
}

func Test_Unordered_should_pass(t *testing.T) {
	// Given
	assert := assertion.New(t)

	// When
	assert.That(nil).Unordered(nil)

	assert.That([]string{"a", "b", "c"}).Unordered([]string{"a", "b", "c"})
	assert.That([]string{"a", "b", "c"}).Unordered([]string{"b", "a", "c"})
	assert.That([]string{"a", "b", "c"}).Unordered([]string{"c", "a", "b"})

	assert.That([]string{"a", "b", "c"}).Not().Unordered([]string{"a", "b"})
	assert.That([]string{"a", "b", "c"}).Not().Unordered([]string{"b", "a", "c", "d"})
	assert.That([]string{"a", "b", "c"}).Not().Unordered([]int{1, 2, 3})

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
			assertFunc: func(assert assertion.Assert) { assert.That([]string{}).Unordered([]string{"b"}) },
			errLog:     fmt.Sprintf("\nValue should contains all elements : %v", []string{"b"}),
		},
		{
			assertFunc: func(assert assertion.Assert) { assert.That([]string{"b", "c"}).Not().Unordered([]string{"b", "c"}) },
			errLog:     fmt.Sprintf("\nValue should not contain all elements : %v", []string{"b", "c"}),
		},
		{
			assertFunc: func(assert assertion.Assert) { assert.That([]string{"b", "c"}).Unordered([]string{"a", "b"}) },
			errLog:     fmt.Sprintf("\nElement [%v] not found", "a"),
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

func Test_Slice_Matchers_should_fail_with_error(t *testing.T) {
	// Mock preparation
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testEntries := []struct {
		assertFunc func(assert assertion.Assert)
		err        error
		times      int
	}{
		{
			assertFunc: func(assert assertion.Assert) {
				assert.That("abcd").Contains("d")
				assert.That("abcd").Not().Contains("d")
			},
			err:   assertion.ErrNotOfSliceType,
			times: 2,
		},
		{
			assertFunc: func(assert assertion.Assert) {
				assert.That("abcd").Unordered([]string{"a"})
				assert.That([]string{"a"}).Not().Unordered("d")
			},
			err:   assertion.ErrNotOfSliceType,
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
