package expect

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

var gmtM1 = func() *time.Location {
	l, err := time.LoadLocation("Etc/GMT-1")
	if err != nil {
		panic(err)
	}
	return l
}()

func runTest(t *testing.T, i interface{}, x, xp string) {
	actual, p, d := format(i)
	aPad := presentations[p]
	Value(t, fmt.Sprintf("f(%v %T)", i, i), del(actual, d)).ToBe(x)
	Value(t, fmt.Sprintf("f(%v %T) pad", i, i), aPad).ToBe(xp)
}

func TestFormatPrimitives(t *testing.T) {
	runTest(t, "foo", "'foo'", " ")
	runTest(t, 7, "7", " ")
	runTest(t, uint(7), "7", " ")
	runTest(t, 12.1, "12.1", " ")
	runTest(t, true, "true", " ")
	runTest(t, nil, "nil", " ")
	runTest(t, time.Date(2020, 3, 4, 8, 32, 34, 0, time.UTC), "2020-03-04T08:32:34Z", " ")
	runTest(t, time.Date(2020, 3, 4, 8, 32, 34, 0, gmtM1), "2020-03-04T08:32:34+01:00", " ")
	runTest(t, errors.New("xyz"), "'xyz'", " ")
}

func TestFormatSlice(t *testing.T) {
	runTest(t, []string{"a", "b"}, "- a\n- b", "\n")
}

func TestFormatMap(t *testing.T) {
	runTest(t, map[string]int{"a": 2, "b": 15}, "a: 2\nb: 15", "\n")
}

type Struct struct {
	Foo   string
	Count int
}

func TestFormatStruct(t *testing.T) {
	runTest(t, Struct{Foo: "Bar"}, "Count: 0\nFoo: Bar", "\n")
}

func TestFormatStructPtr(t *testing.T) {
	runTest(t, &Struct{Foo: "Bar"}, "Count: 0\nFoo: Bar", "\n")
}

type unm struct {
	F func()
	X string
}

func TestFormatUnmarshalable(t *testing.T) {
	// a struct wit a public func field can not me marsheled into yaml
	// so we print the normal string representation instead
	runTest(t, &unm{X: "Foo"}, "{<nil> Foo}", " ")
}
