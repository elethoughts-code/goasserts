package fsbuilder //nolint: testpackage

import (
	"os"
	"testing"

	"github.com/elethoughts-code/goasserts/assertion"
)

func Test_RemoveAll_should_panic(t *testing.T) {
	assert := assertion.New(t)
	defer func() {
		r := recover()
		err := &os.PathError{}
		assert.That(r).AsError(&err)
	}()
	dir := dirBuilder{
		name:   ".",
		parent: nil,
	}
	dir.RemoveAll()
}
