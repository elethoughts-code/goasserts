package assertion_test

import (
	"fmt"
	"testing"

	"github.com/elethoughts-code/goasserts/assertion"
	mocks "github.com/elethoughts-code/goasserts/mocks/assertion"
	"github.com/golang/mock/gomock"
)

func Test_IsEmpty_should_pass(t *testing.T) {
	// Given
	assert := assertion.New(t)

	// When

	assert.That([]string{}).IsEmpty()
	assert.That("").IsEmpty()
	assert.That(map[int]string{}).IsEmpty()

	assert.That([]string{""}).Not().IsEmpty()
	assert.That("123").Not().IsEmpty()
	assert.That(map[int]string{0: ""}).Not().IsEmpty()

	// Then nothing
}

func Test_HasLen_should_pass(t *testing.T) {
	// Given
	assert := assertion.New(t)

	// When

	assert.That([]string{"", "", ""}).HasLen(3)
	assert.That("abc").HasLen(3)
	assert.That(map[int]string{0: "", 1: "", 2: ""}).HasLen(3)

	assert.That([]string{"", ""}).Not().HasLen(3)
	assert.That("abcd").Not().HasLen(3)
	assert.That(map[int]string{0: ""}).Not().HasLen(3)

	// Then nothing
}

func Test_HasMaxLen_should_pass(t *testing.T) {
	// Given
	assert := assertion.New(t)

	// When

	assert.That([]string{"", "", ""}).HasMaxLen(3)
	assert.That("ab").HasMaxLen(3)
	assert.That(map[int]string{}).HasMaxLen(3)

	assert.That([]string{"", "", "", ""}).Not().HasMaxLen(3)
	assert.That("abcd").Not().HasMaxLen(3)
	assert.That(map[int]string{0: "", 1: "", 2: "", 3: "abc"}).Not().HasMaxLen(3)

	// Then nothing
}

func Test_HasMinLen_should_pass(t *testing.T) {
	// Given
	assert := assertion.New(t)

	// When

	assert.That([]string{"", "", "", ""}).HasMinLen(3)
	assert.That("abcd").HasMinLen(3)
	assert.That(map[int]string{0: "", 1: "", 2: "", 3: "abc"}).HasMinLen(3)

	assert.That([]string{"", ""}).Not().HasMinLen(3)
	assert.That("ab").Not().HasMinLen(3)
	assert.That(map[int]string{}).Not().HasMinLen(3)

	// Then nothing
}

func Test_Length_Matchers_should_fail(t *testing.T) {
	// Mock preparation
	ctrl := gomock.NewController(t)

	testEntries := []struct {
		assertFunc func(assert assertion.Assert)
		errLog     string
	}{
		{
			assertFunc: func(assert assertion.Assert) { assert.That([]string{""}).IsEmpty() },
			errLog:     fmt.Sprintf("\nValue length should be empty. Actual length is equal to : %d", 1),
		},
		{
			assertFunc: func(assert assertion.Assert) { assert.That([]string{}).Not().IsEmpty() },
			errLog:     "\nValue length should not be empty",
		},

		{
			assertFunc: func(assert assertion.Assert) { assert.That([]string{""}).HasLen(2) },
			errLog:     fmt.Sprintf("\nValue length is not equal to expectation.\nExpected : %d\nGot : %d", 2, 1),
		},
		{
			assertFunc: func(assert assertion.Assert) { assert.That([]string{"", ""}).Not().HasLen(2) },
			errLog:     fmt.Sprintf("\nValue length should not be equal to : %d", 2),
		},

		{
			assertFunc: func(assert assertion.Assert) { assert.That([]string{""}).HasMinLen(2) },
			errLog:     fmt.Sprintf("\nValue length is higher than : %d", 2),
		},
		{
			assertFunc: func(assert assertion.Assert) { assert.That([]string{"", "", ""}).Not().HasMinLen(2) },
			errLog:     fmt.Sprintf("\nValue length should not be higher than : %d. \n Actual value : %d", 2, 3),
		},

		{
			assertFunc: func(assert assertion.Assert) { assert.That([]string{"", "", ""}).HasMaxLen(2) },
			errLog:     fmt.Sprintf("\nValue length is lower than : %d", 2),
		},
		{
			assertFunc: func(assert assertion.Assert) { assert.That([]string{""}).Not().HasMaxLen(2) },
			errLog:     fmt.Sprintf("\nValue length should not be lower than : %d. \n Actual value : %d", 2, 1),
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

func Test_Length_Matchers_should_fail_with_error(t *testing.T) {
	// Mock preparation
	ctrl := gomock.NewController(t)

	testEntries := []struct {
		assertFunc func(assert assertion.Assert)
		err        error
		times      int
	}{
		{
			assertFunc: func(assert assertion.Assert) {
				assert.That(666).Not().IsEmpty()
				assert.That(struct{}{}).Not().IsEmpty()
				assert.That(false).Not().IsEmpty()
				assert.That(666).Not().HasLen(3)
				assert.That(struct{}{}).Not().HasLen(3)
				assert.That(false).Not().HasLen(3)
				assert.That(666).Not().HasMaxLen(3)
				assert.That(struct{}{}).Not().HasMaxLen(3)
				assert.That(false).Not().HasMaxLen(3)
				assert.That(666).Not().HasMinLen(3)
				assert.That(struct{}{}).Not().HasMinLen(3)
				assert.That(false).Not().HasMinLen(3)
			},
			err:   assertion.ErrNotOfLenType,
			times: 12,
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
