package expect_test

import (
	"testing"

	"gitlab.com/testle/expect"
	"gitlab.com/testle/expect/internal/test"
)

func TestExample(t *testing.T) {
	l := &test.Logger{}
	expect.Value(l, "the guy", "Peter").ToBe("Steven")
	expect.Value(t, "error", l.Messages[0]).ToBe("expected the guy to be 'Steven' but it is 'Peter'")
}

func TestToBeString(t *testing.T) {
	expect.Value(t, "foo", "xxx").ToBe("xxx")
}

func TestToFailToBeString(t *testing.T) {
	l := &test.Logger{}
	expect.Value(l, "foo", "xxx").ToBe("yyy")
	expect.Value(t, "errors", l.Messages).ToCount(1)
	expect.Value(t, "error", l.Messages[0]).ToBe("expected foo to be 'yyy' but it is 'xxx'")
}

func TestToBeFloat64(t *testing.T) {
	expect.Value(t, "liters", 3.45).ToBe(3.45)
}

func TestFailToBeFloat64(t *testing.T) {
	l := &test.Logger{}
	expect.Value(l, "liters", 3.45).ToBe(3.45002)
	expect.Value(t, "error", l.Messages[0]).ToBe("expected liters to be '3.45002' but it is '3.45'")
}

func TestFailToBeMap(t *testing.T) {
	l := &test.Logger{}
	expect.Value(l, "names", map[string]int{"peter": 3, "johan": 2}).ToBe(map[string]int{"peter": 3, "johan": 1})
	expect.Value(t, "error", l.Messages[0]).ToBe(`expected names to be:
johan: 1
peter: 3

but it is:
johan: 2
peter: 3
`)
}

func TestToCountString(t *testing.T) {
	expect.Value(t, "foo", "xxx").ToCount(3)
}

func TestFailToCountString(t *testing.T) {
	l := &test.Logger{}
	expect.Value(l, "foo", "xxx").ToCount(1)
	expect.Value(t, "errors", l.Messages).ToCount(1)
	expect.Value(t, "error", l.Messages[0]).ToBe("expected foo to have 1 elements but it has 3 elements")
}

func TestErrorToCountInt(t *testing.T) {
	l := &test.Logger{}
	expect.Value(l, "foo", 2).ToCount(2)
	expect.Value(t, "errors", l.Messages).ToCount(1)
	expect.Value(t, "error", l.Messages[0]).ToBe("foo is not a datatype with a length (array, slice, map, chan, string)")
}

func TestToHavePrefix(t *testing.T) {
	expect.Value(t, "statement", "we are all crazy").ToHavePrefix("we are")
}

func TestFailToHavePrefix(t *testing.T) {
	l := &test.Logger{}
	expect.Value(l, "statement", "we are all crazy").ToHavePrefix("i am")
	expect.Value(t, "error", l.Messages[0]).ToBe("expected statement to have prefix 'i am' but it is 'we are all crazy'")
}

func TestErrorToHavePrefixOnInt(t *testing.T) {
	l := &test.Logger{}
	expect.Value(l, "number", 7).ToHavePrefix("i am")
	expect.Value(t, "error", l.Messages[0]).ToBe("ToHavePrefix must only be called on a string value")
}

func TestNotToBe(t *testing.T) {
	expect.Value(t, "number", 7).NotToBe(8)
}

func TestFailNotToBe(t *testing.T) {
	l := &test.Logger{}
	expect.Value(l, "number", 7).NotToBe(7)
	expect.Value(t, "error", l.Messages[0]).ToBe("expected number to NOT be '7' but it is")
}

func TestFailNotToBeSlice(t *testing.T) {
	l := &test.Logger{}
	expect.Value(l, "numbers", []int{3, 2, 1}).NotToBe([]int{3, 2, 1})
	expect.Value(t, "error", l.Messages[0]).ToBe(`expected numbers to NOT be:
- 3
- 2
- 1

but it is`)
}
