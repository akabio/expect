package expect

import (
	"reflect"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

type output string

var PlainOutput = output("plain")
var ColoredDiffOutput = output("coloredDiffOutput")

var Output = PlainOutput

// Value wraps a value and provides expectations for this value.
func Value(t Test, name string, val interface{}) Val {
	return Val{
		name:  name,
		t:     t,
		value: val,
	}
}

// Val to call expectations on.
type Val struct {
	name  string
	t     Test
	value interface{}
}

// ToBe asserts that the value is deeply equals to expected value.
func (e Val) ToBe(expected interface{}) Val {
	if !reflect.DeepEqual(e.value, expected) {
		x, xp := f(expected)
		v, vp := f(e.value)
		if Output == ColoredDiffOutput && (len(x) > 30 || len(v) > 30) {
			dmp := diffmatchpatch.New()
			diffs := dmp.DiffMainRunes([]rune(x), []rune(v), false)
			dmp.DiffCleanupSemantic(diffs)
			e.t.Error(dmp.DiffPrettyText(diffs))
		} else {
			e.t.Errorf("expected %v to be%v%v%vbut it is%v%v", e.name, xp, x, xp, vp, v)
		}
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
		x, xp := f(unExpected)
		e.t.Errorf("expected %v to NOT be%v%v%vbut it is", e.name, xp, x, xp)
	}
	return e
}

// ToBeAbout asserts that the number is in deltas range of expected value.
// Only works for numbers.
func (e Val) ToBeAbout(expected, delta float64) Val {
	val := 0.0
	switch t := e.value.(type) {
	case float32:
		val = float64(t)
	case float64:
		val = t
	case int:
		val = float64(t)
	case uint:
		val = float64(t)
	case int8:
		val = float64(t)
	case uint8:
		val = float64(t)
	case int16:
		val = float64(t)
	case uint16:
		val = float64(t)
	case int32:
		val = float64(t)
	case uint32:
		val = float64(t)
	case int64:
		val = float64(t)
	case uint64:
		val = float64(t)
	default:
		e.t.Fatalf("ToBeAbout() can only work on number values but it's called on type %T", e.value)
	}
	if val < expected-delta || val > expected+delta {
		e.t.Errorf("expected %v to be %v±%v but it is %v", e.name, expected, delta, e.value)
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
