package assertion_test

import (
	"testing"

	"github.com/elethoughts-code/goasserts/assertion"
)

func Test_Values_should_pass_assertions(t *testing.T) {
	// Given
	assert := assertion.New(t)

	// When
	assert.That(nil).Values().IsNil()
	assert.That(map[int]string{0: "a", 1: "b", 2: "c"}).Values().Unordered([]string{"a", "b", "c"})
	assert.That(map[int]string{0: "a", 1: "b", 2: "c"}).Values().Not().Unordered([]string{"a", "b", "c", "d"})
	assert.That(map[int]string{0: "a", 1: "b", 2: "c"}).Values().Contains("c")
	assert.That(map[int]string{0: "a", 1: "b", 2: "c"}).Not().Values().Contains("d")

	// Then nothing
}

func Test_Values_panic_if_not_map(t *testing.T) {
	// Given
	assert := assertion.New(t)

	// When / Then
	defer func() {
		r := recover()
		assert.That(r).IsEq("[type error] value should be a map")
	}()
	assert.That("abc").Values().Unordered([]string{"a", "b", "c"})
}

func Test_Keys_should_pass_assertions(t *testing.T) {
	// Given
	assert := assertion.New(t)

	// When
	assert.That(nil).Keys().IsNil()
	assert.That(map[int]string{0: "a", 1: "b", 2: "c"}).Keys().Unordered([]int{0, 1, 2})
	assert.That(map[int]string{0: "a", 1: "b", 2: "c"}).Keys().Not().Unordered([]string{"a", "b", "c", "d"})
	assert.That(map[int]string{0: "a", 1: "b", 2: "c"}).Keys().Contains(1)
	assert.That(map[int]string{0: "a", 1: "b", 2: "c"}).Not().Keys().Contains(4)

	// Then nothing
}

func Test_Keys_panic_if_not_map(t *testing.T) {
	// Given
	assert := assertion.New(t)

	// When / Then
	defer func() {
		r := recover()
		assert.That(r).IsEq("[type error] value should be a map")
	}()
	assert.That("abc").Keys().Unordered([]string{"a", "b", "c"})
}

func Test_Entries_should_pass_assertions(t *testing.T) {
	// Given
	assert := assertion.New(t)

	// When
	assert.That(nil).Entries().IsNil()
	assert.That(map[int]string{0: "a", 1: "b", 2: "c"}).Entries().Unordered([]assertion.MapEntry{
		{Key: 1, Value: "b"},
		{Key: 0, Value: "a"},
		{Key: 2, Value: "c"},
	})
	assert.That(map[int]string{0: "a", 1: "b", 2: "c"}).Entries().Not().Unordered([]assertion.MapEntry{
		{Key: 1, Value: "a"},
		{Key: 0, Value: "b"},
		{Key: 2, Value: "c"},
	})
	assert.That(map[int]string{0: "a", 1: "b", 2: "c"}).Entries().Contains(assertion.MapEntry{Key: 1, Value: "b"})
	assert.That(map[int]string{0: "a", 1: "b", 2: "c"}).Not().Entries().Contains(4)

	// Then nothing
}

func Test_Entries_panic_if_not_map(t *testing.T) {
	// Given
	assert := assertion.New(t)

	// When / Then
	defer func() {
		r := recover()
		assert.That(r).IsEq("[type error] value should be a map")
	}()
	assert.That("abc").Entries().Unordered([]string{"a", "b", "c"})
}
