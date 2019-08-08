package expect

import (
	"fmt"
	"testing"
)

var cases = map[interface{}][]string{
	"foo":   {"'foo'", " "},
	7:       {"7", " "},
	uint(7): {"7", " "},
	12.1:    {"12.1", " "},
	// &[]string{"a", "b"}: []string{"", ""},
}

func runTest(t *testing.T, i interface{}, x, xp string) {
	actual, aPad := f(i)
	Value(t, fmt.Sprintf("f(%v %T)", i, i), actual).ToBe(x)
	Value(t, fmt.Sprintf("f(%v %T) pad", i, i), aPad).ToBe(xp)
}

func TestFormat(t *testing.T) {
	for i, expected := range cases {
		t.Run("", func(t *testing.T) {
			runTest(t, i, expected[0], expected[1])
		})
	}
}

func TestFormatSlice(t *testing.T) {
	runTest(t, []string{"a", "b"}, "  > - a\n  > - b", "\n")
}
