package expect

import (
	"bytes"
	"image"
	"image/png"
	"os"
	"path/filepath"

	_ "image/gif"
	_ "image/jpeg"

	"golang.org/x/exp/slices"
)

func (e Val) ToBeSnapshot(path string) Val {
	e.t.Helper()

	folder := filepath.Dir(path)
	if folder != "." {
		err := os.MkdirAll(folder, 0o755)
		if err != nil {
			e.t.Fatalf("failed to create target folder %v", folder)
		}
	}

	existing, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		e.t.Fatalf("failed to read snaphsot %v: %v", path, err)
	}

	currentString, isString := e.value.(string)
	current, isBytes := e.value.([]byte)
	if !isString && !isBytes {
		e.t.Fatalf("value of .ToBeSnaphsot must be of type string or []byte but it is %T", e.value)
	}

	if isString {
		current = []byte(currentString)
	}

	if existing == nil {
		// snapshot does not exist, create it
		err = os.WriteFile(path, current, 0o644)
		if err != nil {
			e.t.Fatalf("failed to write snapshot %v", path)
		}
	} else {
		if slices.Equal(current, existing) {
			// all is well, snapshot is matched, remove a possible current version
			os.RemoveAll(path + ".current")
		} else {
			e.t.Errorf("snapshot for %v does not match current output", path)
			err = os.WriteFile(path+".current", current, 0o644)
			if err != nil {
				e.t.Fatalf("failed to write snapshot %v", path)
			}
		}
	}

	return e
}

func (e Val) ToBeSnapshotImage(path string) Val {
	e.t.Helper()

	folder := filepath.Dir(path)
	if folder != "." {
		err := os.MkdirAll(folder, 0o755)
		if err != nil {
			e.t.Fatalf("failed to create target folder %v", folder)
		}
	}

	existing, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		e.t.Fatalf("failed to read snaphsot %v: %v", path, err)
	}

	img, isImage := e.value.(image.Image)
	if !isImage {
		e.t.Fatalf("value of .ToBeSnaphsotImage must be of type image but it is %T", e.value)
	}

	// snapshot does not exist, create it
	if existing == nil {
		encoded := bytes.NewBuffer(nil)

		err = png.Encode(encoded, img)
		if err != nil {
			e.t.Fatalf("failed encode snapshot image %v", path)
		}

		err = os.WriteFile(path, encoded.Bytes(), 0o644)
		if err != nil {
			e.t.Fatalf("failed to write snapshot %v", path)
		}

		return e
	}

	// snapshot exists, compare it
	encoded := bytes.NewBuffer(existing)
	snapshotImage, _, err := image.Decode(encoded)
	if err != nil {
		e.t.Fatalf("failed to read snapshot %v", err)
	}

	if isSameImage(e.t, snapshotImage, img) {
		// all is well, snapshot is matched, remove a possible current version
		os.RemoveAll(path + ".current")
		return e
	}

	// not the same image, write current output
	err = png.Encode(encoded, img)
	if err != nil {
		e.t.Fatalf("failed encode snapshot image %v", path+".current")
	}

	err = os.WriteFile(path+".current", encoded.Bytes(), 0o644)
	if err != nil {
		e.t.Fatalf("failed to write snapshot %v", path+".current")
	}

	return e
}

func isSameImage(t Test, snapshot, current image.Image) bool {
	snapshotSize := snapshot.Bounds().Size()
	currentSize := current.Bounds().Size()
	if snapshotSize != currentSize {
		t.Errorf("expected image size to be %v but it is %v", snapshotSize, currentSize)
		return false
	}

	for y := 0; y < snapshotSize.Y; y++ {
		for x := 0; x < snapshotSize.X; x++ {
			rs, gs, bs, as := snapshot.At(x, y).RGBA()
			rc, gc, bc, ac := current.At(x, y).RGBA()

			if rs != rc || gs != gc || bs != bc || as != ac {
				t.Errorf("images are not the same at %v, %v", x, y)
				return false
			}
		}
	}

	return true
}
