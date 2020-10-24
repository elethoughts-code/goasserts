package assertion_test

import (
	"testing"

	"github.com/elethoughts-code/goasserts/assertion"
)

func Test_Attr_should_pass_assertions(t *testing.T) {
	// Given
	assert := assertion.New(t)

	// When
	assert.That(struct {
		A int
	}{33}).Attr("A").IsEq(33)

	assert.That(struct {
		A int
		B string
	}{33, "123"}).Attr("B").IsEq("123")

	assert.That(map[string]int{
		"a": 33,
		"b": 34,
	}).Attr("a").IsEq(33)

	type someType struct {
		a int
		B string
	}

	assert.That(&someType{
		a: 10,
		B: "123",
	}).Attr("B").IsEq("123")

	var b interface{} = &someType{
		a: 10,
		B: "456",
	}
	assert.That(b).Attr("B").IsEq("456")

	assert.That(map[someType]int{
		someType{1, "a"}: 123, //nolint:gofmt
		someType{1, "b"}: 456,
		someType{2, "c"}: 789,
	}).Attr(someType{1, "b"}).IsEq(456)
}

func Test_Index_should_pass_assertions(t *testing.T) {
	// Given
	assert := assertion.New(t)

	// When
	assert.That([]int{10, 11, 12}).Index(1).IsEq(11)
	assert.That([3]int{10, 11, 12}).Index(1).IsEq(11)
	assert.That(&[]int{10, 11, 12}).Index(1).IsEq(11)
	assert.That(&[3]int{10, 11, 12}).Index(1).IsEq(11)
}

func Test_Should_allow_complex_struct_navigation(t *testing.T) {
	// Given
	assert := assertion.New(t)
	type someType struct {
		A int
		B string
		C []someType
	}

	// When
	v := someType{
		A: 0,
		B: "",
		C: []someType{
			{},
			{A: 10, B: "123"},
			{A: 20, B: "456"},
			{A: 30, B: "789"},
		},
	}
	assert.That(v).Attr("C").Index(2).Attr("B").IsEq("456")
	assert.That(&v).Attr("C").Index(2).IsDeepEq(someType{A: 20, B: "456"})
}

func Test_should_panic_when_non_attribute_type_passed(t *testing.T) {
	assert := assertion.New(t)
	defer func() {
		r := recover()
		assert.That(r).IsEq("value should be an attribute type (struct, map, interface, ptr)")
	}()
	assert.That(33).Attr("C")
}

func Test_should_panic_when_non_indexed_type_passed(t *testing.T) {
	assert := assertion.New(t)
	defer func() {
		r := recover()
		assert.That(r).IsEq("value should be an indexed type (array, slice, interface, ptr)")
	}()
	assert.That(33).Index(2)
}

func Test_should_panic_when_attribute_not_found_in_struct(t *testing.T) {
	assert := assertion.New(t)
	defer func() {
		r := recover()
		assert.That(r).IsEq("attribute C not found")
	}()
	type someType struct {
		A int
		B string
	}
	assert.That(someType{}).Attr("C")
}

func Test_should_panic_when_attribute_not_found_in_map(t *testing.T) {
	assert := assertion.New(t)
	defer func() {
		r := recover()
		assert.That(r).IsEq("attribute C not found")
	}()
	assert.That(map[string]int{}).Attr("C")
}

func Test_attr_should_panic_when_value_is_nil(t *testing.T) {
	assert := assertion.New(t)
	defer func() {
		r := recover()
		assert.That(r).IsEq("value is invalid")
	}()
	assert.That(nil).Attr("C")
}

func Test_index_should_panic_when_value_is_nil(t *testing.T) {
	assert := assertion.New(t)
	defer func() {
		r := recover()
		assert.That(r).IsEq("value is invalid")
	}()
	assert.That(nil).Index(1)
}

func Test_index_should_panic_when_out_of_bound(t *testing.T) {
	assert := assertion.New(t)
	defer func() {
		r := recover()
		assert.That(r).IsEq("index 3 out of bound")
	}()
	assert.That([]int{1, 2}).Index(3)
}

func Test_attr_should_panic_when_value_ptr_is_nil(t *testing.T) {
	assert := assertion.New(t)
	defer func() {
		r := recover()
		assert.That(r).IsEq("value is nil")
	}()
	var p *struct {
		A int
	}
	assert.That(p).Attr("A")
}

func Test_index_should_panic_when_value_ptr_is_nil(t *testing.T) {
	assert := assertion.New(t)
	defer func() {
		r := recover()
		assert.That(r).IsEq("value is nil")
	}()
	var p *[]int
	assert.That(p).Index(1)
}

func Test_index_should_panic_when_map_is_nil(t *testing.T) {
	assert := assertion.New(t)
	defer func() {
		r := recover()
		assert.That(r).IsEq("value is nil")
	}()
	var p map[string]int
	assert.That(p).Attr("C")
}

func Test_index_should_panic_when_key_is_not_string_for_struct_access(t *testing.T) {
	assert := assertion.New(t)
	defer func() {
		r := recover()
		assert.That(r).IsEq("attribute key not of string type")
	}()
	assert.That(struct{}{}).Attr(123)
}
