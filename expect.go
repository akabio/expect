package expect

import (
	"encoding/json"
	"reflect"
)

// NamedValue creates an expectation for the provided value with given
// name for reporting.
func NamedValue(t Test, n string, val interface{}) Value {
	return Value{
		name:  n,
		t:     t,
		value: val,
	}
}

// Value to assert expectations on.
type Value struct {
	name  string
	t     Test
	value interface{}
}

// ToBe asserts that the value is deeply equals to expected value.
func (e Value) ToBe(expected interface{}) Value {
	if !reflect.DeepEqual(e.value, expected) {
		if needsFormating(e.value) {
			// if it's a "complex" type we try to print the value as formated yaml
			exp, erre := json.MarshalIndent(expected, "--", "  ")
			val, errv := json.MarshalIndent(e.value, "--", "  ")
			if erre != nil || errv != nil {
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
func (e Value) ToCount(c int) Value {
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

// NotToBe asserts that the value is not deeply equals to expected value.
func (e Value) NotToBe(unExpected interface{}) Value {
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
