package assertion

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

type HTTPRecorderParser interface {
	BodyToString() Expectation
	JSONBodyToMap() Expectation
	JSONBodyToSlice() Expectation
	DecodeBody(decoder func(*bytes.Buffer) (interface{}, error)) Expectation
	Response() Expectation
	Status() Expectation
	Headers() Expectation
	Header(header string) Expectation
	Cookies() Expectation
	Cookie(cookie string) Expectation
}

func (exp *expectation) BodyToString() Expectation {
	return exp.DecodeBody(func(body *bytes.Buffer) (interface{}, error) {
		parsed, err := exp.assert.br.ReadAll(body)
		return string(parsed), err
	})
}

func (exp *expectation) JSONBodyToMap() Expectation {
	return exp.DecodeBody(func(body *bytes.Buffer) (interface{}, error) {
		parsed := make(map[string]interface{})
		err := json.NewDecoder(body).Decode(&parsed)
		return parsed, err
	})
}

func (exp *expectation) JSONBodyToSlice() Expectation {
	return exp.DecodeBody(func(body *bytes.Buffer) (interface{}, error) {
		parsed := make([]interface{}, 0)
		err := json.NewDecoder(body).Decode(&parsed)
		return parsed, err
	})
}

func (exp *expectation) DecodeBody(decoder func(*bytes.Buffer) (interface{}, error)) Expectation {
	recorder, ok := exp.v.(*httptest.ResponseRecorder)
	if !ok {
		panic(ErrNotOfResponseRecorderType)
	}
	// Body content is read and preserved
	originalBody, err := exp.assert.br.ReadAll(recorder.Body)
	if err != nil {
		panic(err)
	}
	body, err := decoder(bytes.NewBuffer(originalBody))
	if err != nil {
		panic(err)
	}
	recorder.Body = bytes.NewBuffer(originalBody)
	exp.v = body
	return exp
}

func (exp *expectation) Response() Expectation {
	recorder, ok := exp.v.(*httptest.ResponseRecorder)
	if !ok {
		panic(ErrNotOfResponseRecorderType)
	}
	exp.v = recorder.Result()
	return exp
}

func (exp *expectation) Status() Expectation {
	exp.Response()
	response := exp.v.(*http.Response)
	exp.v = response.StatusCode
	return exp
}

func (exp *expectation) Headers() Expectation {
	exp.Response()
	response := exp.v.(*http.Response)
	exp.v = response.Header
	return exp
}

func (exp *expectation) Header(header string) Expectation {
	exp.Response()
	response := exp.v.(*http.Response)
	exp.v = response.Header.Get(header)
	return exp
}

func (exp *expectation) Cookies() Expectation {
	exp.Response()
	response := exp.v.(*http.Response)
	exp.v = response.Cookies()
	return exp
}

func (exp *expectation) Cookie(cookie string) Expectation {
	exp.Response()
	response := exp.v.(*http.Response)
	cookies := response.Cookies()
	for _, c := range cookies {
		if c.Name == cookie {
			exp.v = c
			return exp
		}
	}
	exp.v = nil
	return exp
}
