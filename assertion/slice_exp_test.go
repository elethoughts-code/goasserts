package assertion_test

import (
	"errors"
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

func Test_UnorderedDeepEq_should_pass(t *testing.T) {
	// Given
	assert := assertion.New(t)

	// When
	assert.That(nil).Unordered(nil)

	assert.That([]string{"a", "b", "c"}).UnorderedDeepEq([]string{"a", "b", "c"})
	assert.That([]string{"a", "b", "c"}).UnorderedDeepEq([]string{"b", "a", "c"})
	assert.That([]string{"a", "b", "c"}).UnorderedDeepEq([]string{"c", "a", "b"})

	assert.That([]string{"a", "b", "c"}).Not().UnorderedDeepEq([]string{"a", "b"})
	assert.That([]string{"a", "b", "c"}).Not().UnorderedDeepEq([]string{"b", "a", "c", "d"})
	assert.That([]string{"a", "b", "c"}).Not().UnorderedDeepEq([]int{1, 2, 3})

	// Then nothing
}

func Test_UnorderedNoDiff_should_pass(t *testing.T) {
	// Given
	assert := assertion.New(t)

	// When
	assert.That(nil).Unordered(nil)

	assert.That([]string{"a", "b", "c"}).UnorderedNoDiff([]string{"a", "b", "c"})
	assert.That([]string{"a", "b", "c"}).UnorderedNoDiff([]string{"b", "a", "c"})
	assert.That([]string{"a", "b", "c"}).UnorderedNoDiff([]string{"c", "a", "b"})

	assert.That([]string{"a", "b", "c"}).Not().UnorderedNoDiff([]string{"a", "b"})
	assert.That([]string{"a", "b", "c"}).Not().UnorderedNoDiff([]string{"b", "a", "c", "d"})
	assert.That([]string{"a", "b", "c"}).Not().UnorderedNoDiff([]int{1, 2, 3})

	// Then nothing
}

func Test_All_should_pass(t *testing.T) {
	// Given
	assert := assertion.New(t)

	// When

	assert.That([]string{"a", "b", "c"}).All(func(v interface{}) (assertion.MatchResult, error) {
		return assertion.MatchResult{Matches: len(v.(string)) == 1}, nil
	})
	assert.That([]string{}).All(func(v interface{}) (assertion.MatchResult, error) {
		return assertion.MatchResult{Matches: len(v.(string)) == 1}, nil
	})

	assert.That([]string{"a", "bb", "c"}).Not().All(func(v interface{}) (assertion.MatchResult, error) {
		return assertion.MatchResult{Matches: len(v.(string)) == 1}, nil
	})
	// Then nothing
}

func Test_AtLeast_should_pass(t *testing.T) {
	// Given
	assert := assertion.New(t)

	// When

	assert.That([]string{"a", "b", "c"}).AtLeast(2, func(v interface{}) (assertion.MatchResult, error) {
		return assertion.MatchResult{Matches: len(v.(string)) == 1}, nil
	})

	assert.That([]string{"a", "bb", "c"}).AtLeast(2, func(v interface{}) (assertion.MatchResult, error) {
		return assertion.MatchResult{Matches: len(v.(string)) == 1}, nil
	})

	assert.That([]string{}).Not().AtLeast(2, func(v interface{}) (assertion.MatchResult, error) {
		return assertion.MatchResult{Matches: len(v.(string)) == 1}, nil
	})

	assert.That([]string{"aa", "bb", "c"}).Not().AtLeast(2, func(v interface{}) (assertion.MatchResult, error) {
		return assertion.MatchResult{Matches: len(v.(string)) == 1}, nil
	})
	// Then nothing
}

func Test_Any_should_pass(t *testing.T) {
	// Given
	assert := assertion.New(t)

	// When

	assert.That([]string{"a", "b", "c"}).Any(func(v interface{}) (assertion.MatchResult, error) {
		return assertion.MatchResult{Matches: len(v.(string)) == 1}, nil
	})

	assert.That([]string{"aa", "b", "cc"}).Any(func(v interface{}) (assertion.MatchResult, error) {
		return assertion.MatchResult{Matches: len(v.(string)) == 1}, nil
	})
	assert.That([]string{"aa", "bb", "c"}).Any(func(v interface{}) (assertion.MatchResult, error) {
		return assertion.MatchResult{Matches: len(v.(string)) == 1}, nil
	})

	assert.That([]string{}).Not().Any(func(v interface{}) (assertion.MatchResult, error) {
		return assertion.MatchResult{Matches: len(v.(string)) == 1}, nil
	})

	assert.That([]string{"aa", "bb", "cc"}).Not().Any(func(v interface{}) (assertion.MatchResult, error) {
		return assertion.MatchResult{Matches: len(v.(string)) == 1}, nil
	})
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
			assertFunc: func(assert assertion.Assert) { assert.That([]string{}).UnorderedDeepEq([]string{"b"}) },
			errLog:     fmt.Sprintf("\nValue should contains all elements : %v", []string{"b"}),
		},
		{
			assertFunc: func(assert assertion.Assert) { assert.That([]string{}).UnorderedNoDiff([]string{"b"}) },
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
		{
			assertFunc: func(assert assertion.Assert) {
				assert.That([]string{"a", "bb", "c"}).All(func(v interface{}) (assertion.MatchResult, error) {
					return assertion.MatchResult{Matches: len(v.(string)) == 1}, nil
				})
			},
			errLog: fmt.Sprintf("\nMatcher dont apply to all values. Non matching indexes : %v", []int{1}),
		},
		{
			assertFunc: func(assert assertion.Assert) {
				assert.That([]string{"a", "b", "c"}).Not().All(func(v interface{}) (assertion.MatchResult, error) {
					return assertion.MatchResult{Matches: len(v.(string)) == 1}, nil
				})
			},
			errLog: "\nMatcher should not apply to all elements",
		},
		{
			assertFunc: func(assert assertion.Assert) {
				assert.That([]string{"aa", "bb", "cc"}).AtLeast(2, func(v interface{}) (assertion.MatchResult, error) {
					return assertion.MatchResult{Matches: len(v.(string)) == 1}, nil
				})
			},
			errLog: fmt.Sprintf("\nAt least %d element(s) should match. Non matching indexes : %v", 2, []int{0, 1, 2}),
		},
		{
			assertFunc: func(assert assertion.Assert) {
				assert.That([]string{"a", "b", "cc"}).Not().AtLeast(2, func(v interface{}) (assertion.MatchResult, error) {
					return assertion.MatchResult{Matches: len(v.(string)) == 1}, nil
				})
			},
			errLog: fmt.Sprintf("\nMatcher should not apply to %d element(s) or more", 2),
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
	err := errors.New("matcher proper error")

	testEntries := []struct {
		assertFunc func(assert assertion.Assert)
		err        error
		times      int
	}{
		{
			assertFunc: func(assert assertion.Assert) {
				assert.That("abcd").Contains("d")
				assert.That("abcd").Not().Contains("d")
				assert.That("abcd").All(func(v interface{}) (assertion.MatchResult, error) {
					return assertion.MatchResult{Matches: len(v.(string)) == 1}, nil
				})
				assert.That("abcd").AtLeast(2, func(v interface{}) (assertion.MatchResult, error) {
					return assertion.MatchResult{Matches: len(v.(string)) == 1}, nil
				})
			},
			err:   assertion.ErrNotOfSliceType,
			times: 4,
		},
		{
			assertFunc: func(assert assertion.Assert) {
				assert.That("abcd").Unordered([]string{"a"})
				assert.That([]string{"a"}).Not().Unordered("d")
			},
			err:   assertion.ErrNotOfSliceType,
			times: 2,
		},
		{
			assertFunc: func(assert assertion.Assert) {
				assert.That([]string{"a", "b", "c"}).All(func(v interface{}) (assertion.MatchResult, error) {
					return assertion.MatchResult{}, err
				})
			},
			err:   err,
			times: 1,
		},
		{
			assertFunc: func(assert assertion.Assert) {
				assert.That([]string{"a", "b", "c"}).AtLeast(2, func(v interface{}) (assertion.MatchResult, error) {
					return assertion.MatchResult{}, err
				})
			},
			err:   err,
			times: 1,
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
