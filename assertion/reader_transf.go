package assertion

import (
	"bytes"
	"io"
)

type ReaderTransformer interface {
	ReadCloserToBytes() Expectation
	ReaderToBytes() Expectation
	BytesToString() Expectation
}

func (exp *expectation) ReadCloserToBytes() Expectation {
	if exp.v == nil {
		return exp
	}
	rc, isReadClosed := exp.v.(io.ReadCloser)

	if !isReadClosed {
		panic("[type error] value should be of type io.ReadCloser")
	}

	defer func() {
		_ = rc.Close()
	}()
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(rc); err != nil {
		panic(err)
	}
	exp.v = buf.Bytes()
	return exp
}

func (exp *expectation) ReaderToBytes() Expectation {
	if exp.v == nil {
		return exp
	}
	r, isReader := exp.v.(io.Reader)

	if !isReader {
		panic("[type error] value should be of type io.Reader")
	}

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(r); err != nil {
		panic(err)
	}
	exp.v = buf.Bytes()
	return exp
}

func (exp *expectation) BytesToString() Expectation {
	if exp.v == nil {
		return exp
	}
	b, isBytes := exp.v.([]byte)

	if !isBytes {
		panic("[type error] value should be of type []byte")
	}

	exp.v = string(b)
	return exp
}
