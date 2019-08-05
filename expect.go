package expect

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/ghodss/yaml"
)

// Value creates an expectation for the provided value with given
// name for reporting.
func Value(t Test, name string, val interface{}) Val {
	return Val{
		name:  name,
		t:     t,
		value: val,
	}
}

// Val to test expectations against.
type Val struct {
	name  string
	t     Test
	value interface{}
}

// ToBe asserts that the value is deeply equals to expected value.
func (e Val) ToBe(expected interface{}) Val {
	if !reflect.DeepEqual(e.value, expected) {
		if needsFormating(e.value) {
			// if it's a "complex" type we try to print the value as formated yaml
			exp, erre := yaml.Marshal(expected)
			val, errv := yaml.Marshal(e.value)
			if erre == nil && errv == nil {
				e.t.Errorf("expected %v to be:\n%v\nbut it is:\n%v", e.name, string(exp), string(val))
				return e
			}
		}
		// otherwise or if serialisation failed print it as it is
		e.t.Errorf("expected %v to be '%v' but it is '%v'", e.name, expected, e.value)
	}
	return e
}

// ToCount asserts that the list/map/chan/string has c elements.
func (e Val) ToCount(c int) Val {
	if !hasLen(e.value) {
		e.t.Fatalf("%v is not a datatype with a length (array, slice, map, chan, string)", e.name)
		return e
	}

	l := reflect.ValueOf(e.value).Len()
	if l != c {
		e.t.Errorf("expected %v to have %v elements but it has %v elements", e.name, c, l)
	}

	return e
}

// NotToBe asserts that the value is not deeply equals to expected value.
func (e Val) NotToBe(unExpected interface{}) Val {
	if reflect.DeepEqual(e.value, unExpected) {
		exp, err := json.MarshalIndent(unExpected, "--", "  ")
		if err != nil {
			e.t.Errorf("expected %v to NOT be '%v' but it is", e.name, unExpected)
		} else {
			e.t.Errorf("expected %v to NOT be:\n%v\nbut it is", e.name, string(exp))
		}
	}
	return e
}

// ToHavePrefix asserts that the string value starts with the provided prefix.
func (e Val) ToHavePrefix(prefix string) Val {
	actual, is := e.value.(string)
	if !is {
		e.t.Fatalf("ToHavePrefix must only be called on a string value")
	}
	if !strings.HasPrefix(actual, prefix) {
		e.t.Errorf("expected %v to have prefix '%v' but it is '%v'", e.name, prefix, actual)
	}
	return e
}

func hasLen(v interface{}) bool {
	switch reflect.TypeOf(v).Kind() {
	case reflect.Array:
		return true
	case reflect.Chan:
		return true
	case reflect.Map:
		return true
	case reflect.Slice:
		return true
	case reflect.String:
		return true
	}
	return false
}

func needsFormating(v interface{}) bool {
	switch reflect.TypeOf(v).Kind() {
	case reflect.Array:
		return true
	case reflect.Map:
		return true
	case reflect.Slice:
		return true
	case reflect.Struct:
		return true
	}
	return false
}
