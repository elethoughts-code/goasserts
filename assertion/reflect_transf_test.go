package assertion_test

import (
	"testing"

	"github.com/elethoughts-code/goasserts/assertion"
)

type sample struct {
	a string
	b int
}

func Test_should_dereference_pointer(t *testing.T) {
	assert := assertion.New(t)
	// Then
	assert.That(&sample{
		a: "hello world",
		b: 10,
	}).Dereference().IsEq(sample{
		a: "hello world",
		b: 10,
	})
}

func Test_should_panic_when_non_pointer_value_is_de_referenced(t *testing.T) {
	assert := assertion.New(t)
	defer func() {
		r := recover()
		assert.That(r).IsEq("value is not a pointer")
	}()

	assert.That(sample{
		a: "hello world",
		b: 10,
	}).Dereference().IsEq(sample{
		a: "hello world",
		b: 10,
	})
}
