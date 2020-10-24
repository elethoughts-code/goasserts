package assertion //nolint: testpackage

import (
	"bytes"
	"errors"
	"io"
	"net/http/httptest"
	"testing"
)

type testBytesReader struct {
	result []byte
	err    error
}

func (br testBytesReader) ReadAll(r io.Reader) ([]byte, error) {
	return br.result, br.err
}

func Test_Byte_reading_error_should_be_reported(t *testing.T) {
	// Given
	err := errors.New("some error")
	assert := newWithBr(t, testBytesReader{
		result: nil,
		err:    err,
	})

	// When / then
	defer func() {
		r := recover()
		assert.That(r).IsEq(err)
	}()
	assert.That(&httptest.ResponseRecorder{
		Body: bytes.NewBuffer([]byte{1}),
	}).BodyToString()
}
