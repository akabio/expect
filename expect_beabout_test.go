package expect_test

import (
	"testing"

	"gitlab.com/akabio/expect"
	"gitlab.com/akabio/expect/internal/test"
)

func TestToBeAbout(t *testing.T) {
	expect.Value(t, "liters", 1.92).ToBeAbout(2, 0.1)
}

func TestFailToBeAbout(t *testing.T) {
	l := test.New(t, func(t expect.Test) {
		expect.Value(t, "liters", 1.92).ToBeAbout(2, 0.01)
	})
	l.ExpectMessageNoLoc(0).ToBe("expected liters to be 2±0.01 but it is 1.92")
}
