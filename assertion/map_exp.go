package assertion

// MapExpectation interface encloses map related expectations.
//
// ContainsValue(e interface{}) check if the map have an equal to e parameter value.
//
// ContainsKey(e interface{}) check if the map have an equal to e parameter key.
type MapExpectation interface {
	ContainsValue(e interface{})
	ContainsKey(e interface{})
}

func (exp *expectation) ContainsValue(e interface{}) {
	exp.t.Helper()
	exp.Matches(ContainsValue(e))
}

func (exp *expectation) ContainsKey(e interface{}) {
	exp.t.Helper()
	exp.Matches(ContainsKey(e))
}
