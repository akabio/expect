package expect_test

import (
	"bytes"
	"image"
	"os"
	"testing"

	"github.com/akabio/expect"
	"github.com/akabio/expect/internal/test"
)

var sampleImage = func() image.Image {
	sp, err := os.ReadFile("./testdata/sample.png")
	if err != nil {
		panic(err)
	}

	img, _, err := image.Decode(bytes.NewBuffer(sp))
	if err != nil {
		panic(err)
	}

	return img
}()

func TestCreateSnapshot(t *testing.T) {
	cleanTestData(t)
	expect.Value(t, "content", "we are all crazy").ToBeSnapshot("testdata/volatile/ss1.txt")

	data, err := os.ReadFile("testdata/volatile/ss1.txt")
	expect.Error(t, err).ToBe(nil)
	expect.Value(t, "content", string(data)).ToBe("we are all crazy")
}

func TestMismatchSnapshot(t *testing.T) {
	cleanTestData(t)
	expect.Value(t, "content", "we are all crazy").ToBeSnapshot("testdata/volatile/ss1.txt")
	l := test.New(t, func(t expect.Test) {
		expect.Value(t, "content", "we are all nuts").ToBeSnapshot("testdata/volatile/ss1.txt")
	})
	l.ExpectMessage(0).ToBe("snapshot for testdata/volatile/ss1.txt does not match current output")

	data, err := os.ReadFile("testdata/volatile/ss1.txt")
	expect.Error(t, err).ToBe(nil)
	expect.Value(t, "content", string(data)).ToBe("we are all crazy")

	data, err = os.ReadFile("testdata/volatile/ss1.txt.current")
	expect.Error(t, err).ToBe(nil)
	expect.Value(t, "content", string(data)).ToBe("we are all nuts")
}

func TestMatchAfterMismatchSnapshot(t *testing.T) {
	cleanTestData(t)
	expect.Value(t, "content", "we are all crazy").ToBeSnapshot("testdata/volatile/ss1.txt")
	test.New(t, func(t expect.Test) {
		expect.Value(t, "content", "we are all nuts").ToBeSnapshot("testdata/volatile/ss1.txt")
	})
	expect.Value(t, "content", "we are all crazy").ToBeSnapshot("testdata/volatile/ss1.txt")

	data, err := os.ReadFile("testdata/volatile/ss1.txt")
	expect.Error(t, err).ToBe(nil)
	expect.Value(t, "content", string(data)).ToBe("we are all crazy")

	data, err = os.ReadFile("testdata/volatile/ss1.txt.current")
	expect.Error(t, err).Message().ToBe("open testdata/volatile/ss1.txt.current: no such file or directory")
}

func TestCreateSnapshotFromBytes(t *testing.T) {
	cleanTestData(t)
	expect.Value(t, "content", []byte{1, 2, 3}).ToBeSnapshot("testdata/volatile/ss1.bin")

	data, err := os.ReadFile("testdata/volatile/ss1.bin")
	expect.Error(t, err).ToBe(nil)
	expect.Value(t, "content", data).ToBe([]byte{1, 2, 3})
}

func cleanTestData(t *testing.T) {
	err := os.RemoveAll("testdata/volatile")
	if err != nil {
		t.Fatal("Failed to clear testdata folder")
	}
}

func TestCreateSnapshotImage(t *testing.T) {
	cleanTestData(t)
	expect.Value(t, "content", sampleImage).ToBeSnapshotImage("testdata/sample.png")
}

func TestCreateSnapshotImageSize(t *testing.T) {
	cleanTestData(t)

	l := test.New(t, func(t expect.Test) {
		expect.Value(t, "content", sampleImage).ToBeSnapshotImage("testdata/snapshots/sample-small.png")
	})
	l.ExpectMessages().ToCount(1)
	l.ExpectMessage(0).ToBe("expected image size to be (64,64) but it is (128,128)")
}

func TestCreateSnapshotImageDifferent(t *testing.T) {
	cleanTestData(t)

	l := test.New(t, func(t expect.Test) {
		expect.Value(t, "content", sampleImage).ToBeSnapshotImage("testdata/snapshots/sample-dirty.png")
	})
	l.ExpectMessages().ToCount(1)
	l.ExpectMessage(0).ToBe("images are not the same at 76, 20")
}
