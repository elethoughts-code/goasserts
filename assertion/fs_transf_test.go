package assertion_test

import (
	"os"
	"testing"

	"github.com/elethoughts-code/goasserts/assertion"
	fsBuilder "github.com/elethoughts-code/goasserts/fs_builder"
)

func Test_should_read_file_into_string(t *testing.T) {
	assert := assertion.New(t)

	// When
	tmpFolder := fsBuilder.TmpDir("", "my_folder_1")

	file := tmpFolder.File("file", os.O_CREATE, 0755).WriteString("hello world")

	// Then
	assert.That(file.Name()).FileAsString().IsEq("hello world")

	// Clean up
	t.Cleanup(tmpFolder.Root().RemoveAll)
}

type sampleStruct struct {
	Value1 string   `json:"value1" yaml:"value1"`
	Value2 int      `json:"value2" yaml:"value2"`
	Value3 []string `json:"value3" yaml:"value3"`
}

func Test_should_read_file_into_struct_as_JSON(t *testing.T) {
	assert := assertion.New(t)

	// When
	tmpFolder := fsBuilder.TmpDir("", "my_folder_1")

	file := tmpFolder.File("file", os.O_CREATE, 0755).WriteString(`
		{
			"value1": "hello world",
			"value2": 10,
			"value3": ["a", "b", "c"]
		}
	`)

	// Then
	assert.That(file.Name()).FileAsJSON(&sampleStruct{}).Dereference().IsDeepEq(sampleStruct{
		Value1: "hello world",
		Value2: 10,
		Value3: []string{"a", "b", "c"},
	})

	// Clean up
	t.Cleanup(tmpFolder.Root().RemoveAll)
}

func Test_should_read_file_into_struct_as_YAML(t *testing.T) {
	assert := assertion.New(t)

	// When
	tmpFolder := fsBuilder.TmpDir("", "my_folder_1")

	file := tmpFolder.File("file", os.O_CREATE, 0755).WriteStringDedent(`
			value1: hello world
			value2: 10
			value3:
			  - a
			  - b
			  - c`)

	// Then
	assert.That(file.Name()).FileAsYAML(&sampleStruct{}).Dereference().IsDeepEq(sampleStruct{
		Value1: "hello world",
		Value2: 10,
		Value3: []string{"a", "b", "c"},
	})

	// Clean up
	t.Cleanup(tmpFolder.Root().RemoveAll)
}

func Test_should_panic_when_reading_file_as_string_and_dont_exists(t *testing.T) {
	assert := assertion.New(t)
	defer func() {
		r := recover()
		assert.That(r).IsError(os.ErrNotExist)
	}()

	assert.That("/file").FileAsString()
}

func Test_should_panic_when_reading_file_as_json_and_dont_exists(t *testing.T) {
	assert := assertion.New(t)
	defer func() {
		r := recover()
		assert.That(r).IsError(os.ErrNotExist)
	}()

	assert.That("/file").FileAsJSON([]string{})
}

func Test_should_panic_when_reading_file_as_yaml_and_dont_exists(t *testing.T) {
	assert := assertion.New(t)
	defer func() {
		r := recover()
		assert.That(r).IsError(os.ErrNotExist)
	}()

	assert.That("/file").FileAsYAML([]string{})
}

func Test_should_panic_when_reading_file_as_yaml_and_syntax_error(t *testing.T) {
	assert := assertion.New(t)
	defer func() {
		r := recover()
		assert.That(r.(error).Error()).HasPrefix("yaml: line")
	}()

	// When
	tmpFolder := fsBuilder.TmpDir("", "my_folder_1")

	file := tmpFolder.File("file", os.O_CREATE, 0755).WriteStringDedent(`
			value1: hello world
			value2: 10
			value3:
			  - a
			  - b
			- c`)

	// Then
	assert.That(file.Name()).FileAsYAML(&sampleStruct{}).Dereference().IsDeepEq(sampleStruct{
		Value1: "hello world",
		Value2: 10,
		Value3: []string{"a", "b", "c"},
	})

	// Clean up
	t.Cleanup(tmpFolder.Root().RemoveAll)
}
