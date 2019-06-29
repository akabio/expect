package expect_test

import (
	"testing"

	"gitlab.com/testle/expect"
	"gitlab.com/testle/expect/internal/test"
)

func TestExample(t *testing.T) {
	l := &test.Logger{}
	expect.NamedValue(l, "the guy", "Peter").ToBe("Steven")
	expect.NamedValue(t, "error", l.Messages[0]).ToBe("expected the guy to be 'Steven' but it is 'Peter'")
}

func TestToBeString(t *testing.T) {
	expect.NamedValue(t, "foo", "xxx").ToBe("xxx")
}

func TestToFailToBeString(t *testing.T) {
	l := &test.Logger{}
	expect.NamedValue(l, "foo", "xxx").ToBe("yyy")
	expect.NamedValue(t, "errors", l.Messages).ToCount(1)
	expect.NamedValue(t, "error", l.Messages[0]).ToBe("expected foo to be 'yyy' but it is 'xxx'")
}

func TestToCountString(t *testing.T) {
	expect.NamedValue(t, "foo", "xxx").ToCount(3)
}

func TestToFailToCountString(t *testing.T) {
	l := &test.Logger{}
	expect.NamedValue(l, "foo", "xxx").ToCount(1)
	expect.NamedValue(t, "errors", l.Messages).ToCount(1)
	expect.NamedValue(t, "error", l.Messages[0]).ToBe("expected foo to have 1 elements but it has 3 elements")
}

func TestToHavePrefix(t *testing.T) {
	expect.NamedValue(t, "statement", "we are all crazy").ToHavePrefix("we are")
}

func TestFailToHavePrefix(t *testing.T) {
	l := &test.Logger{}
	expect.NamedValue(l, "statement", "we are all crazy").ToHavePrefix("i am")
	expect.NamedValue(t, "error", l.Messages[0]).ToBe("expected statement to have prefix 'i am' but it is 'we are all crazy'")
}
