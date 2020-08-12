package assertion_test

import (
	"testing"

	"github.com/elethoughts-code/goasserts/assertion"
	mocks "github.com/elethoughts-code/goasserts/mocks/assertion"
	"github.com/golang/mock/gomock"
)

func Test_ContainsValue_should_pass(t *testing.T) {
	// Given
	assert := assertion.New(t)

	// When

	assert.That(map[int]string{0: "a", 1: "b", 2: "c"}).ContainsValue("b")
	assert.That(map[int]string{0: "a", 1: "b", 2: "c"}).Not().ContainsValue("d")

	assert.That(map[int]string{}).Not().ContainsValue("d")
	assert.That("abcd").Not().ContainsValue("d")

	assert.That(map[int]struct{ a string }{0: {"a"}, 1: {"b"}, 2: {"c"}}).ContainsValue(struct{ a string }{"b"})

	// Then nothing
}

func Test_ContainsKey_should_pass(t *testing.T) {
	// Given
	assert := assertion.New(t)

	// When

	assert.That(map[int]string{0: "a", 1: "b", 2: "c"}).ContainsKey(1)
	assert.That(map[int]string{0: "a", 1: "b", 2: "c"}).Not().ContainsKey(4)

	assert.That(map[int]string{}).Not().ContainsKey(1)
	assert.That("abcd").Not().ContainsKey(1)

	assert.That(map[struct{ a string }]int{{"a"}: 0, {"b"}: 1, {"c"}: 2}).ContainsKey(struct{ a string }{"b"})

	// Then nothing
}

func Test_Map_Matchers_should_fail(t *testing.T) {
	// Mock preparation
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testEntries := []struct {
		assertFunc func(assert assertion.Assert)
		errLog     string
	}{
		{
			assertFunc: func(assert assertion.Assert) {
				assert.That(map[int]string{0: "a", 1: "b", 2: "c"}).ContainsValue("d")
			},
			errLog: "\nValue should contains element : d",
		},
		{
			assertFunc: func(assert assertion.Assert) {
				assert.That(map[int]string{0: "a", 1: "b", 2: "c"}).Not().ContainsValue("b")
			},
			errLog: "\nValue should not contains element : b",
		},
		{
			assertFunc: func(assert assertion.Assert) {
				assert.That(map[int]string{}).ContainsValue("b")
			},
			errLog: "\nValue should contains element : b",
		},
		{
			assertFunc: func(assert assertion.Assert) { assert.That("abcd").ContainsValue("d") },
			errLog:     "\nValue should be a map",
		},

		{
			assertFunc: func(assert assertion.Assert) {
				assert.That(map[int]string{0: "a", 1: "b", 2: "c"}).ContainsKey(4)
			},
			errLog: "\nValue should contains key : 4",
		},
		{
			assertFunc: func(assert assertion.Assert) {
				assert.That(map[int]string{0: "a", 1: "b", 2: "c"}).Not().ContainsKey(2)
			},
			errLog: "\nValue should not contains key : 2",
		},
		{
			assertFunc: func(assert assertion.Assert) {
				assert.That("abcd").ContainsKey(1)
			},
			errLog: "\nValue should be a map",
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
