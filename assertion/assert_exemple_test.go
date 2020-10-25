package assertion_test

import (
	"net/http/httptest"
	"testing"

	"github.com/elethoughts-code/goasserts/assertion"
)

func Example() {
	// Testing variable
	t := &testing.T{}

	// Examples
	assert := assertion.New(t)
	assert.That(1).IsEq(1)
	assert.That("123").Not().IsNil()
	assert.That(httptest.NewRecorder()).Cookie("my-cookie").IsNil()
}
