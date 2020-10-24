package assertion_test

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/elethoughts-code/goasserts/assertion"
)

func Test_BodyToString_Http_record_transformation_passes(t *testing.T) {
	// Given
	assert := assertion.New(t)

	handler := func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, "this is working ...")
	}

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	// When / Then
	assert.That(w).BodyToString().IsEq("this is working ...")
}

func Test_JSONBodyToMap_Http_record_transformation_passes(t *testing.T) {
	// Given
	assert := assertion.New(t)

	handler := func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, `{"a": 1, "b": [1,2,3], "c": "123"}`)
	}

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	// When / Then
	assert.That(w).JSONBodyToMap().IsDeepEq(map[string]interface{}{
		"a": 1.0,
		"b": []interface{}{1.0, 2.0, 3.0},
		"c": "123",
	})
}

func Test_JSONBodyToSlice_Http_record_transformation_passes(t *testing.T) {
	// Given
	assert := assertion.New(t)

	handler := func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, `["a", "b", "c"]`)
	}

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	// When / Then
	assert.That(w).JSONBodyToSlice().IsDeepEq([]interface{}{"a", "b", "c"})
}

func Test_DecodeBody_Http_record_transformation_passes(t *testing.T) {
	// Given
	assert := assertion.New(t)

	handler := func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, `not 123`)
	}

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	// When / Then
	assert.That(w).DecodeBody(func(body *bytes.Buffer) (interface{}, error) {
		return 123, nil
	}).IsEq(123)
}

func Test_DecodeBody_Http_record_transformation_panic_on_error(t *testing.T) {
	// Given
	assert := assertion.New(t)
	err := errors.New("some error")
	handler := func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, `not 123`)
	}

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	// When / Then
	defer func() {
		r := recover()
		assert.That(r).IsEq(err)
	}()
	assert.That(w).DecodeBody(func(body *bytes.Buffer) (interface{}, error) {
		return 123, err
	})
}

func Test_Record_transformations_should_panic_when_no_record_passed(t *testing.T) {
	// Given
	assert := assertion.New(t)
	panicFunc := func(expect func(e assertion.Expectation)) {
		defer func() {
			r := recover()
			assert.That(r).IsEq(assertion.ErrNotOfResponseRecorderType)
		}()
		expect(assert.That(123))
	}

	// When / Then
	panicFunc(func(e assertion.Expectation) { e.BodyToString() })
	panicFunc(func(e assertion.Expectation) { e.JSONBodyToSlice() })
	panicFunc(func(e assertion.Expectation) { e.JSONBodyToSlice() })
	panicFunc(func(e assertion.Expectation) {
		e.DecodeBody(func(body *bytes.Buffer) (interface{}, error) {
			return nil, nil
		})
	})
	panicFunc(func(e assertion.Expectation) { e.Response() })
}

func Test_Multiple_Http_record_transformation_passes(t *testing.T) {
	// Given
	assert := assertion.New(t)

	handler := func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{Name: "my-cookie", Value: "654321", Domain: "my-domain"})
		w.Header().Add("X-SOME-HEADER", "123456")
		w.WriteHeader(201)

		_, _ = io.WriteString(w, `["a", "b", "c"]`)
	}

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	// When / Then
	assert.That(w).Status().IsEq(201)
	assert.That(w).Status().Not().IsEq(200)
	assert.That(w).BodyToString().IsEq(`["a", "b", "c"]`)
	assert.That(w).JSONBodyToSlice().IsDeepEq([]interface{}{"a", "b", "c"})
	assert.That(w).Headers().HasLen(2)
	assert.That(w).Cookies().HasLen(1)
	assert.That(w).Header("X-SOME-HEADER").IsEq("123456")
	assert.That(w).Cookie("my-cookie").Attr("Name").IsEq("my-cookie")
	assert.That(w).Cookie("my-cookie").Attr("Value").IsEq("654321")
	assert.That(w).Cookie("my-cookie-2").IsNil()
}
