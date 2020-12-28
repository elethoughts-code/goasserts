package assertion_test

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/elethoughts-code/goasserts/assertion"
	fsBuilder "github.com/elethoughts-code/goasserts/fs_builder"
)

func Test_Reader_transformer_should_read_bytes(t *testing.T) {
	// Given
	assert := assertion.New(t)

	// When
	assert.That(nil).ReaderToBytes().IsNil()
	assert.That(nil).BytesToString().IsNil()

	assert.That(strings.NewReader("ABCDEFGH 123456789")).ReaderToBytes().BytesToString().IsEq("ABCDEFGH 123456789")

	// Then nothing
}

func Test_ReadCloser_transformer_should_read_bytes_and_close_reader(t *testing.T) {
	// Given
	assert := assertion.New(t)

	file := fsBuilder.TmpDir("", "root").File("file_1", os.O_CREATE, 0755).WriteString("ABCDEFGH 123456789")

	// When
	assert.That(nil).ReadCloserToBytes().IsNil()
	f, _ := os.OpenFile(file.Name(), os.O_RDONLY, 0)
	assert.That(f).ReadCloserToBytes().BytesToString().IsEq("ABCDEFGH 123456789")

	// Then nothing
}

func Test_ReadCloser_panic_if_not_ReadCloser(t *testing.T) {
	// Given
	assert := assertion.New(t)

	// When / Then
	defer func() {
		r := recover()
		assert.That(r).IsEq("[type error] value should be of type io.ReadCloser")
	}()
	assert.That("abc").ReadCloserToBytes()
}

func Test_ReadCloser_panic_if_not_Reader(t *testing.T) {
	// Given
	assert := assertion.New(t)

	// When / Then
	defer func() {
		r := recover()
		assert.That(r).IsEq("[type error] value should be of type io.Reader")
	}()
	assert.That("abc").ReaderToBytes()
}

func Test_ReadCloser_panic_if_not_bytes(t *testing.T) {
	// Given
	assert := assertion.New(t)

	// When / Then
	defer func() {
		r := recover()
		assert.That(r).IsEq("[type error] value should be of type []byte")
	}()
	assert.That("abc").BytesToString()
}

func Test_ReadCloser_panic_if_reading_is_errored(t *testing.T) {
	// Given
	assert := assertion.New(t)
	err := errors.New("some error")

	// When / Then
	defer func() {
		r := recover()
		assert.That(r).IsError(err)
	}()
	assert.That(&testRc{err: err}).ReadCloserToBytes()
}

func Test_Reader_panic_if_reading_is_errored(t *testing.T) {
	// Given
	assert := assertion.New(t)
	err := errors.New("some error")

	// When / Then
	defer func() {
		r := recover()
		assert.That(r).IsError(err)
	}()
	assert.That(&testRc{err: err}).ReaderToBytes()
}

type testRc struct {
	err error
}

func (r *testRc) Read(_ []byte) (n int, err error) {
	return 0, r.err
}

func (r *testRc) Close() error {
	return nil
}
