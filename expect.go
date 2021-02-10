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
// This is a shortcut for value using "error" as name.
func (e *Expect) Error(t Test, val interface{}) Val {
	return Val{
		ex:    e,
		name:  "error",
		t:     t,
		value: val,
	}
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
	if !reflect.DeepEqual(e.value, expected) {
		x, delimiterX := format(expected)
		v, delimiterV := format(e.value)
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
			e.t.Errorf("expected %v to be%v%v%vbut it is%v%v", e.name, delimiterX, x, delimiterX, delimiterV, v)
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
		x, delimiter := format(unExpected)
		e.t.Errorf("expected %v to NOT be%v%v%vbut it is", e.name, delimiter, x, delimiter)
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

// ToHaveSuffix asserts that the string value ends with the provided sufix.
func (e Val) ToHaveSuffix(suffix string) Val {
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

func (e Val) index(i int) Val {
	if !isIndexable(e.value) {
		e.t.Fatalf("%v is not a indexable datatype (array, slice, string)", e.name)
		return e
	}

	l := reflect.ValueOf(e.value).Len()
	if i < 0 {
		i = l + i
	}
	if i >= l || i < 0 {
		e.t.Fatalf("%v has length of %, index %v is out of bonds", e.name, l, i)
	}

	v := reflect.ValueOf(e.value).Index(i)

	return Val{
		ex:    e.ex,
		name:  "element at index " + strconv.Itoa(i) + " of " + e.name,
		t:     e.t,
		value: v.Interface(),
	}
}

func (e Val) First() Val {
	return e.index(0)
}

func (e Val) Last() Val {
	return e.index(-1)
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
