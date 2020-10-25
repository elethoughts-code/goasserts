# Elethoughts Go Assertion Library

[![Go Report Card](https://goreportcard.com/badge/github.com/elethoughts-code/goasserts)](https://goreportcard.com/report/github.com/elethoughts-code/goasserts)
[![Coverage Status](https://coveralls.io/repos/github/elethoughts-code/goasserts/badge.svg?branch=master)](https://coveralls.io/github/elethoughts-code/goasserts?branch=master)
[![Build Status](https://travis-ci.org/elethoughts-code/goasserts.svg?branch=master)](https://travis-ci.org/elethoughts-code/goasserts)

Elethoughts Go Assertion library is what its name says : A simple assertion library for Go unit tests.

## Installation

```shell script
go get github.com/elethoughts-code/goasserts
```

## Usage

For detailed example, you can check the unit tests from the repository. Next are some code samples.

```go
func Test_Common_expectations(t *testing.T) {
	// Given
	assert := assertion.New(t)

	// When / Then

	assert.That(nil).IsNil()

	assert.That([]string{" "}).Not().IsNil()
	assert.That("123").Not().IsNil()
	assert.That([]string{""}).Not().IsEmpty()
	assert.That(map[int]string{0: "", 1: "", 2: ""}).HasLen(3)
}
```

```go
func Test_Map_expectations(t *testing.T) {
	// Given
	assert := assertion.New(t)

	// When / Then

	assert.That(map[int]string{0: "a", 1: "b", 2: "c"}).ContainsValue("b")
	assert.That(map[int]string{0: "a", 1: "b", 2: "c"}).Not().ContainsValue("d")

	assert.That(map[int]string{}).Not().ContainsValue("d")

	assert.That(map[int]struct{ a string }{0: {"a"}, 1: {"b"}, 2: {"c"}}).ContainsValue(struct{ a string }{"b"})
}
```

```go
func Test_Http_record_transformation(t *testing.T) {
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
```

## License

Elethoughts Go Assertion Library is licensed under the terms of the MIT license.
