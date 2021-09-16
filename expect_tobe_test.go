package expect_test

import (
	"testing"

	"gitlab.com/akabio/expect"
	"gitlab.com/akabio/expect/internal/test"
)

func TestToBeString(t *testing.T) {
	expect.Value(t, "foo", "xxx").ToBe("xxx")
}

func TestFailToBeString(t *testing.T) {
	l := test.New(t, func(t expect.Test) {
		expect.Value(t, "foo", "xxx").ToBe("yyy")
	})
	l.ExpectMessages().ToCount(1)
	l.ExpectMessage(0).ToBe("expected foo to be 'yyy' but it is 'xxx'")
}

func TestFailToBeMultilineString(t *testing.T) {
	l := test.New(t, func(t expect.Test) {
		expect.Value(t, "foo", "A\nB\nC").ToBe("a\nb\nc")
	})
	l.ExpectMessages().ToCount(1)
	l.ExpectMessage(0).ToBe(`expected foo to be
    a
    b
    c
but it is
    A
    B
    C`)
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

func TestFailToBeFloat32Type(t *testing.T) {
	l := test.New(t, func(t expect.Test) {
		expect.Value(t, "liters", 3.45).ToBe(float32(3.45))
	})
	l.ExpectMessage(0).ToBe("expected liters to be of type float32 but it is of type float64")
}

func TestFailTypeCheckMap(t *testing.T) {
	l := test.New(t, func(t expect.Test) {
		expect.Value(t, "inventory", map[string][]string{}).ToBe([]int{})
	})
	l.ExpectMessage(0).ToBe("expected inventory to be of type []int but it is of type map[string][]string")
}

func TestFailTypeCheckNonPointer(t *testing.T) {
	l := test.New(t, func(t expect.Test) {
		expect.Value(t, "ref", &[]int{}).ToBe([]int{})
	})
	l.ExpectMessage(0).ToBe("expected ref to be of type []int but it is of type *[]int")
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

func TestToBeArray(t *testing.T) {
	a := [3]string{"a", "b", "c"}
	b := [3]string{"a", "b", "c"}
	expect.Value(t, "array", a).ToBe(b)
}

func TestFailToBeArray(t *testing.T) {
	a := [3]string{"a", "b", "c"}
	b := [3]string{"a", "b", "s"}
	l := test.New(t, func(t expect.Test) {
		expect.Value(t, "array", a).ToBe(b)
	})
	l.ExpectMessage(0).ToBe(`expected array to be
    - a
    - b
    - s
but it is
    - a
    - b
    - c`)
}

func TestNilTypeToBeNil(t *testing.T) {
	type vs struct{}
	var vsv *vs
	expect.Value(t, "vsv", vsv).ToBe(nil)
}

func TestNilValueToBeNil(t *testing.T) {
	expect.Value(t, "vsv", nil).ToBe(nil)
}
