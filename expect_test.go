package expect_test

import (
	"errors"
	"testing"

	"gitlab.com/akabio/expect"
	"gitlab.com/akabio/expect/internal/test"
)

func TestExample(t *testing.T) {
	l := test.New(t, func(t expect.Test) {
		expect.Value(t, "the guy", "Peter").ToBe("Steven")
	})
	l.ExpectMessage(0).ToBe("expected the guy to be 'Steven' but it is 'Peter'")
}

func TestToBeString(t *testing.T) {
	expect.Value(t, "foo", "xxx").ToBe("xxx")
}

func TestToFailToBeString(t *testing.T) {
	l := test.New(t, func(t expect.Test) {
		expect.Value(t, "foo", "xxx").ToBe("yyy")
	})
	l.ExpectMessages().ToCount(1)
	l.ExpectMessage(0).ToBe("expected foo to be 'yyy' but it is 'xxx'")
}

func TestToBeFloat64(t *testing.T) {
	expect.Value(t, "liters", 3.45).ToBe(3.45)
}

func TestFailToBeFloat64(t *testing.T) {
	l := test.New(t, func(t expect.Test) {
		expect.Value(t, "liters", 3.45).ToBe(3.45002)
	})
	l.ExpectMessage(0).ToBe("expected liters to be 3.45002 but it is 3.45")
}

func TestFailToBeMap(t *testing.T) {
	l := test.New(t, func(t expect.Test) {
		expect.Value(t, "names", map[string]int{"peter": 3, "johan": 2}).ToBe(map[string]int{"peter": 3, "johan": 1})
	})
	l.ExpectMessage(0).ToBe(`expected names to be
johan: 1
peter: 3
but it is
johan: 2
peter: 3`)
}

func TestToCountString(t *testing.T) {
	expect.Value(t, "foo", "xxx").ToCount(3)
}

func TestFailToCountString(t *testing.T) {
	l := test.New(t, func(t expect.Test) {
		expect.Value(t, "foo", "xxx").ToCount(1)
	})

	l.ExpectMessages().ToCount(1)
	l.ExpectMessage(0).ToBe("expected foo to have 1 elements but it has 3 elements")
}

func TestErrorToCountInt(t *testing.T) {
	l := test.New(t, func(t expect.Test) {
		expect.Value(t, "foo", 2).ToCount(2)
	})
	l.ExpectMessages().ToCount(1)
	l.ExpectMessage(0).ToBe("foo is not a datatype with a length (array, slice, map, chan, string)")
}

func TestToHavePrefix(t *testing.T) {
	expect.Value(t, "statement", "we are all crazy").ToHavePrefix("we are")
}

func TestFailToHavePrefix(t *testing.T) {
	l := test.New(t, func(t expect.Test) {
		expect.Value(t, "statement", "we are all crazy").ToHavePrefix("i am")
	})
	l.ExpectMessage(0).ToBe("expected statement to have prefix 'i am' but it is 'we are all crazy'")
}

func TestErrorToHavePrefixOnInt(t *testing.T) {
	l := test.New(t, func(t expect.Test) {
		expect.Value(t, "number", 7).ToHavePrefix("i am")
	})
	l.ExpectMessage(0).ToBe("ToHavePrefix must only be called on a string value")
}

func TestToHaveSuffix(t *testing.T) {
	expect.Value(t, "statement", "we are all crazy").ToHaveSuffix("all crazy")
}

func TestFailToHaveSufix(t *testing.T) {
	l := test.New(t, func(t expect.Test) {
		expect.Value(t, "statement", "we are all crazy").ToHaveSuffix("all nuts")
	})
	l.ExpectMessage(0).ToBe("expected statement to have suffix 'all nuts' but it is 'we are all crazy'")
}

func TestErrorToHaveSufixOnInt(t *testing.T) {
	l := test.New(t, func(t expect.Test) {
		expect.Value(t, "number", 7).ToHaveSuffix("i am")
	})
	l.ExpectMessage(0).ToBe("ToHaveSuffix must only be called on a string value")
}

func TestNotToBe(t *testing.T) {
	expect.Value(t, "number", 7).NotToBe(8)
}

func TestFailNotToBe(t *testing.T) {
	l := test.New(t, func(t expect.Test) {
		expect.Value(t, "number", 7).NotToBe(7)
	})
	l.ExpectMessage(0).ToBe("expected number to NOT be 7 but it is")
}

func TestFailNotToBeSlice(t *testing.T) {
	l := test.New(t, func(t expect.Test) {
		expect.Value(t, "numbers", []int{3, 2, 1}).NotToBe([]int{3, 2, 1})
	})
	l.ExpectMessage(0).ToBe(`expected numbers to NOT be
- 3
- 2
- 1
but it is`)
}

func TestToBeAbout(t *testing.T) {
	expect.Value(t, "liters", 1.92).ToBeAbout(2, 0.1)
}

func TestFailToBeAbout(t *testing.T) {
	l := test.New(t, func(t expect.Test) {
		expect.Value(t, "liters", 1.92).ToBeAbout(2, 0.01)
	})
	l.ExpectMessage(0).ToBe("expected liters to be 2±0.01 but it is 1.92")
}

func TestErrorToHaveMessage(t *testing.T) {
	expect.Error(t, errors.New("I am the error message")).Message().ToBe("I am the error message")
	expect.Value(t, "error", errors.New("I am the error message")).Message().ToBe("I am the error message")
}

func TestNilErrorMessage(t *testing.T) {
	expect.Error(t, nil).Message().ToBe("")
}

func TestIntToNotAllowMessageMethod(t *testing.T) {
	l := test.New(t, func(t expect.Test) {
		expect.Value(t, "int", 0).Message().ToBe("0")
	})
	l.ExpectMessage(0).ToBe("Message must only be called on a error value")
}

func TestExpectFirstInSlice(t *testing.T) {
	expect.Value(t, "int slice", []int{1, 2, 3}).First().ToBe(1)
}
func TestExpectFirstInString(t *testing.T) {
	expect.Value(t, "string", "Alabama").First().ToBe(byte(65))
}
func TestExpectLastInSlice(t *testing.T) {
	expect.Value(t, "int slice", []int{1, 2, 3}).Last().ToBe(3)
}

func TestColoredOutputSpaceChars(t *testing.T) {
	l := test.New(t, func(t expect.Test) {
		o := expect.Default.Output
		expect.Default.Output = expect.ColoredDiffOutput
		expect.Value(t, "spaces", " 	\n----------------------").ToBe("----------------------")
		expect.Default.Output = o
	})
	l.ExpectMessage(0).ToBe("'[31m․↦↵\n[0m----------------------'")
}
