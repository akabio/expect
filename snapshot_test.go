package expect_test

import (
	"os"
	"testing"

	"github.com/akabio/expect"
	"github.com/akabio/expect/internal/test"
)

func TestCreateSnapshot(t *testing.T) {
	cleanTestData(t)
	expect.Value(t, "content", "we are all crazy").ToBeSnapshot("testdata/ss1.txt")

	data, err := os.ReadFile("testdata/ss1.txt")
	expect.Error(t, err).ToBe(nil)
	expect.Value(t, "content", string(data)).ToBe("we are all crazy")
}

func TestMismatchSnapshot(t *testing.T) {
	cleanTestData(t)
	expect.Value(t, "content", "we are all crazy").ToBeSnapshot("testdata/ss1.txt")
	l := test.New(t, func(t expect.Test) {
		expect.Value(t, "content", "we are all nuts").ToBeSnapshot("testdata/ss1.txt")
	})
	l.ExpectMessage(0).ToBe("snapshot for testdata/ss1.txt does not match current output")

	data, err := os.ReadFile("testdata/ss1.txt")
	expect.Error(t, err).ToBe(nil)
	expect.Value(t, "content", string(data)).ToBe("we are all crazy")

	data, err = os.ReadFile("testdata/ss1.txt.current")
	expect.Error(t, err).ToBe(nil)
	expect.Value(t, "content", string(data)).ToBe("we are all nuts")
}

func TestMatchAfterMismatchSnapshot(t *testing.T) {
	cleanTestData(t)
	expect.Value(t, "content", "we are all crazy").ToBeSnapshot("testdata/ss1.txt")
	test.New(t, func(t expect.Test) {
		expect.Value(t, "content", "we are all nuts").ToBeSnapshot("testdata/ss1.txt")
	})
	expect.Value(t, "content", "we are all crazy").ToBeSnapshot("testdata/ss1.txt")

	data, err := os.ReadFile("testdata/ss1.txt")
	expect.Error(t, err).ToBe(nil)
	expect.Value(t, "content", string(data)).ToBe("we are all crazy")

	data, err = os.ReadFile("testdata/ss1.txt.current")
	expect.Error(t, err).Message().ToBe("open testdata/ss1.txt.current: no such file or directory")
}

func TestCreateSnapshotFromBytes(t *testing.T) {
	cleanTestData(t)
	expect.Value(t, "content", []byte{1, 2, 3}).ToBeSnapshot("testdata/ss1.bin")

	data, err := os.ReadFile("testdata/ss1.bin")
	expect.Error(t, err).ToBe(nil)
	expect.Value(t, "content", data).ToBe([]byte{1, 2, 3})
}

func cleanTestData(t *testing.T) {
	err := os.RemoveAll("testdata")
	if err != nil {
		t.Fatal("Failed to clear testdata folder")
	}
}
