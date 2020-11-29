package assertion

import (
	"encoding/json"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type FsTransformer interface {
	FileAsString() Expectation
	FileAsJSON(content interface{}) Expectation
	FileAsYAML(content interface{}) Expectation
}

func (exp *expectation) FileAsString() Expectation {
	exp.t.Helper()
	content, err := ioutil.ReadFile(exp.v.(string))
	if err != nil {
		panic(err)
	}
	exp.v = string(content)
	return exp
}

func (exp *expectation) decode(decoder func(b []byte) (interface{}, error)) Expectation {
	exp.t.Helper()
	b, err := ioutil.ReadFile(exp.v.(string))
	if err != nil {
		panic(err)
	}

	decoded, err := decoder(b)
	if err != nil {
		panic(err)
	}
	exp.v = decoded
	return exp
}

func (exp *expectation) FileAsJSON(content interface{}) Expectation {
	exp.t.Helper()
	return exp.decode(func(b []byte) (interface{}, error) {
		err := json.Unmarshal(b, content)
		return content, err
	})
}

func (exp *expectation) FileAsYAML(content interface{}) Expectation {
	exp.t.Helper()
	return exp.decode(func(b []byte) (interface{}, error) {
		err := yaml.Unmarshal(b, content)
		return content, err
	})
}
