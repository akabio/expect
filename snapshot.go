package expect

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"path/filepath"

	_ "image/gif"
	_ "image/jpeg"

	"github.com/ghodss/yaml"
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

	current, err := asBytes(e.value)
	if err != nil {
		e.t.Error(err)
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

func asBytes(in any) ([]byte, error) {
	switch t := in.(type) {
	case []byte:
		return t, nil
	case string:
		return []byte(t), nil
	default:
		return yaml.Marshal(in)
	}
}

type SnapshotImageOption interface {
	apply(o *snapshotImageOptions)
}

type snapshotImageOptions struct {
	colorMatchFactor float64
	matchRatio       float64
}

type snapshotImageOptionExact int

var SnapshotImageOptionExact snapshotImageOptionExact

func (s snapshotImageOptionExact) apply(o *snapshotImageOptions) {
	o.colorMatchFactor = 1
	o.matchRatio = 1
}

// ToBeSnapshotImage saves the image in the first run, in later runs, compares the image to the saved one.
// If they are not the same it will write a .current.pn and .diff.png version of the image.
// The images match by default when 90% of the pixels colors are by less than 10% off.
// The Parameter SnapshotImageOptionExact forces the images to be exactly the same.
func (e Val) ToBeSnapshotImage(path string, opts ...SnapshotImageOption) Val {
	e.t.Helper()

	optOb := &snapshotImageOptions{
		colorMatchFactor: 0.9,
		matchRatio:       0.9,
	}

	for _, opt := range opts {
		opt.apply(optOb)
	}

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

	var img image.Image

	switch t := e.value.(type) {
	case image.Image:
		img = t
	case []byte:
		img, _, err = image.Decode(bytes.NewReader(t))
		if err != nil {
			e.t.Fatalf("[]byte value of .ToBeSnapshotImage is not an image, %v", err)
		}

	default:
		e.t.Fatalf("value of .ToBeSnapshotImage must be of type image or []byte but it is %T", e.value)
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

	isSame, diffImg := isSameImage(e.t, snapshotImage, img, optOb)
	if isSame {
		// all is well, snapshot is matched, remove a possible current version
		os.RemoveAll(path + ".current.png")
		os.RemoveAll(path + ".diff.png")
		return e
	}

	// not the same image, write current output
	current := bytes.NewBuffer(nil)
	err = png.Encode(current, img)
	if err != nil {
		e.t.Fatalf("failed encode snapshot image %v", path+".current.png")
	}

	err = os.WriteFile(path+".current.png", current.Bytes(), 0o644)
	if err != nil {
		e.t.Fatalf("failed to write snapshot %v", path+".current.png")
	}

	diff := bytes.NewBuffer(nil)
	err = png.Encode(diff, diffImg)
	if err != nil {
		e.t.Fatalf("failed encode diff image %v", path+".diff.png")
	}

	err = os.WriteFile(path+".diff.png", diff.Bytes(), 0o644)
	if err != nil {
		e.t.Fatalf("failed to diff snapshot %v", path+".diff.png")
	}

	return e
}

func isSameImage(t Test, snapshot, current image.Image, opts *snapshotImageOptions) (bool, image.Image) {
	snapshotSize := snapshot.Bounds().Size()
	currentSize := current.Bounds().Size()
	if snapshotSize != currentSize {
		t.Errorf("expected image size to be %v but it is %v", snapshotSize, currentSize)
		return false, nil
	}

	diffImg := image.NewRGBA(snapshot.Bounds())

	mismatches := 0

	for y := 0; y < snapshotSize.Y; y++ {
		for x := 0; x < snapshotSize.X; x++ {
			rs, gs, bs, as := snapshot.At(x, y).RGBA()
			rc, gc, bc, ac := current.At(x, y).RGBA()

			rd := getDiffFor(rs, rc)
			gd := getDiffFor(gs, gc)
			bd := getDiffFor(bs, bc)
			ad := getDiffFor(as, ac)

			rb := math.Min(255, rd*3*255)
			gb := math.Min(255, gd*3*255)
			bb := math.Min(255, bd*3*255)

			diffImg.SetRGBA(x, y, color.RGBA{R: uint8(rb), G: uint8(gb), B: uint8(bb), A: 255})

			avg := (rd + gd + bd + ad) / 4
			if avg > (1 - opts.colorMatchFactor) {
				mismatches++
			}
		}
	}

	m := float64(mismatches) / float64(snapshotSize.X*snapshotSize.Y)

	if m > (1 - opts.matchRatio) {
		t.Errorf("expected image does not match snapshot")
		return false, diffImg
	}

	return true, nil
}

func getDiffFor(rs, rc uint32) float64 {
	return math.Abs(float64(rs)-float64(rc)) / (256*256 - 1)
}
