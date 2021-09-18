package expect

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

type output string

var PlainOutput = output("plain")
var ColoredDiffOutput = output("coloredDiffOutput")

type Expect struct {
	Output output
}

var Default = &Expect{
	Output: PlainOutput,
}

// Value wraps a value and provides expectations for this value.
// It delegates to the default instance `Default`.
func Value(t Test, name string, val interface{}) Val {
	return Default.Value(t, name, val)
}

// Error wraps an error and provides expectations for this value.
// It delegates to the default instance `Default`.
func Error(t Test, val interface{}) Val {
	return Default.Error(t, val)
}

// Value wraps a value and provides expectations for this value.
func (e *Expect) Value(t Test, name string, val interface{}) Val {
	return Val{
		ex:    e,
		name:  name,
		t:     t,
		value: val,
	}
}

// Error wraps an error and provides expectations for this value.
// This is a shortcut for Value(t, "error", val).
func (e *Expect) Error(t Test, val interface{}) Val {
	return Value(t, "error", val)
}

// Val to call expectations on.
type Val struct {
	ex    *Expect
	name  string
	t     Test
	value interface{}
}

// ToBe asserts that the value is deeply equals to expected value.
func (e Val) ToBe(expected interface{}) Val {
	e.t.Helper()

	if !sameType(e.value, expected) {
		e.t.Errorf("expected %v to be of type %v but it is of type %v", e.name, typeName(expected), typeName(e.value))
		return e
	}

	// if both are some kind of nil we are fine
	if isNil(e.value) && isNil(expected) {
		return e
	}

	if !reflect.DeepEqual(e.value, expected) {
		x, v, del := formatBoth(expected, e.value)
		if e.ex.Output == ColoredDiffOutput && (len(x) > 20 || len(v) > 20) {
			dmp := diffmatchpatch.New()
			diffs := dmp.DiffMainRunes([]rune(v), []rune(x), false)
			diffs = dmp.DiffCleanupSemantic(diffs)
			txt := dmp.DiffPrettyText(diffs)
			txt = strings.ReplaceAll(txt, " ", "․")
			txt = strings.ReplaceAll(txt, "\t", "↦")
			txt = strings.ReplaceAll(txt, "\n", "↵\n")
			txt = strings.ReplaceAll(txt, "\r", "↵\n")
			e.t.Error(txt)
		} else {
			pres := presentations[del]
			e.t.Errorf("expected %v to be%v%v%vbut it is%v%v", e.name, pres, indent(x, del), pres, pres, indent(v, del))
		}
	}

	return e
}

// ToCount asserts that the list/map/chan/string has c elements. Strings use the number of unicode chars.
func (e Val) ToCount(c int) Val {
	e.t.Helper()

	if !hasLen(e.value) {
		e.t.Fatalf("%v is not a datatype with a length (array, slice, map, chan, string)", e.name)
		return e
	}

	l := reflect.ValueOf(e.value).Len()

	str, isStr := e.value.(string)
	if isStr {
		l = len([]rune(str))
	}

	if l != c {
		e.t.Errorf("expected %v to have %v elements but it has %v elements", e.name, c, l)
	}

	return e
}

// NotToBe asserts that the value is not deeply equals to expected value.
func (e Val) NotToBe(unExpected interface{}) Val {
	e.t.Helper()

	if reflect.DeepEqual(e.value, unExpected) {
		x, p := formatOne(unExpected)
		nl := presentations[p]
		e.t.Errorf("expected %v to NOT be%v%v%vbut it is", e.name, nl, x, nl)
	}

	return e
}

// ToBeAbout asserts that the number is in deltas range of expected value.
// Only works for numbers.
func (e Val) ToBeAbout(expected, delta float64) Val {
	e.t.Helper()

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
	e.t.Helper()

	actual, is := e.value.(string)
	if !is {
		e.t.Fatalf("ToHavePrefix must only be called on a string value")
	}

	if !strings.HasPrefix(actual, prefix) {
		e.t.Errorf("expected %v to have prefix '%v' but it is '%v'", e.name, prefix, actual)
	}

	return e
}

// ToHaveSuffix asserts that the string value ends with the provided sufix.
func (e Val) ToHaveSuffix(suffix string) Val {
	e.t.Helper()

	actual, is := e.value.(string)
	if !is {
		e.t.Fatalf("ToHaveSuffix must only be called on a string value")
	}

	if !strings.HasSuffix(actual, suffix) {
		e.t.Errorf("expected %v to have suffix '%v' but it is '%v'", e.name, suffix, actual)
	}

	return e
}

// Message creates a new value from the given errors message. If the error is nil the message
// wil be the empty string.
func (e Val) Message() Val {
	if e.value == nil {
		// nil always translates to empty string
		return Val{
			ex:    e.ex,
			name:  e.name + " message",
			t:     e.t,
			value: "",
		}
	}

	actual, is := e.value.(error)
	if !is {
		e.t.Fatalf("Message must only be called on a error value")
	}

	return Val{
		ex:    e.ex,
		name:  e.name + " message",
		t:     e.t,
		value: actual.Error(),
	}
}

func (e Val) index(t Test, i int) Val {
	t.Helper()

	calcIndex := func(l int) int {
		e.t.Helper()

		if l == 0 {
			if i == -1 {
				e.t.Fatalf("%v is empty, can not take last element", e.name)
			}

			if i == 0 {
				e.t.Fatalf("%v is empty, can not take first element", e.name)
			}
		}

		in := i
		if in < 0 {
			in = l + in
		}

		if in >= l || in < 0 {
			e.t.Fatalf("%v has length of %v, index %v is out of bounds", e.name, l, i)
		}

		return in
	}

	if !isIndexable(e.value) {
		e.t.Fatalf("%v is not an indexable datatype", e.name)
		return e
	}

	// strings are handled as rune slices
	str, isStr := e.value.(string)
	if isStr {
		runes := []rune(str)
		i = calcIndex(len(runes))

		return Val{
			ex:    e.ex,
			name:  "char at index " + strconv.Itoa(i) + " of " + e.name,
			t:     e.t,
			value: string(runes[i : i+1]),
		}
	}

	rVal := reflect.ValueOf(e.value)
	i = calcIndex(rVal.Len())

	v := rVal.Index(i)

	return Val{
		ex:    e.ex,
		name:  "element at index " + strconv.Itoa(i) + " of " + e.name,
		t:     e.t,
		value: v.Interface(),
	}
}

func (e Val) First() Val {
	e.t.Helper()
	return e.index(e.t, 0)
}

func (e Val) Last() Val {
	e.t.Helper()
	return e.index(e.t, -1)
}

func isIndexable(v interface{}) bool {
	switch reflect.TypeOf(v).Kind() {
	case reflect.Array:
		return true
	case reflect.Slice:
		return true
	case reflect.String:
		return true
	}

	return false
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
