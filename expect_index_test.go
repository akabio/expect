package expect_test

import (
	"testing"

	"gitlab.com/akabio/expect"
)

func TestExpectFirstInSlice(t *testing.T) {
	expect.Value(t, "int slice", []int{1, 2, 3}).First().ToBe(1)
}
func TestExpectFirstInString(t *testing.T) {
	expect.Value(t, "string", "Alabama").First().ToBe(byte(65))
}
func TestExpectLastInSlice(t *testing.T) {
	expect.Value(t, "int slice", []int{1, 2, 3}).Last().ToBe(3)
}
