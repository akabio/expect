package expect

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"strings"
)

type Option interface {
	apply(o *snapshotImageOptions)
}

type snapshotImageOptions struct {
	pixelTolerance float64
	matchTolerance float64
}

type options struct {
	pixelTolerance *float64
	matchTolerance *float64
}

func WithExact() Option {
	return &options{
		pixelTolerance: fptr(0),
		matchTolerance: fptr(0),
	}
}

// WithPixelTolerance sets the maximum allowed color difference for a pixel to be considered a match.
// 0 means exact match only. 1 means all pixels matched.
func WithPixelTolerance(t float64) Option {
	return &options{
		pixelTolerance: fptr(t),
	}
}

// WithMatchTolerance sets the maximum fraction of pixels that may differ while still accepting the image.
// 0 means no mismatches allowed. 1 means all mismatches allowed.
func WithMatchTolerance(t float64) Option {
	return &options{
		matchTolerance: fptr(t),
	}
}

func (s *options) apply(o *snapshotImageOptions) {
	if s.pixelTolerance != nil {
		o.pixelTolerance = *s.pixelTolerance
	}

	if s.matchTolerance != nil {
		o.matchTolerance = *s.matchTolerance
	}
}

func fptr(f float64) *float64 {
	return &f
}

// ToBeSnapshotImage saves the image in the first run, in later runs, compares the image to the saved one.
// If they are not the same it will write a .current.pn and .diff.png version of the image.
// The images match by default when 99% of the pixels colors are by less than 10% off.
// The Parameter SnapshotImageOptionExact forces the images to be exactly the same.
func (e Val) ToBeSnapshotImage(path string, opts ...Option) Val {
	e.t.Helper()

	if !strings.HasSuffix(path, ".png") {
		e.t.Fatalf("only png format is supported, pleas add a .png extension to the snapshot path")
	}

	optOb := &snapshotImageOptions{
		pixelTolerance: 0.1,
		matchTolerance: 0.01,
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
		os.RemoveAll(currentPath(path))
		os.RemoveAll(diffPath(path))
		return e
	}

	// not the same image, write current output
	current := bytes.NewBuffer(nil)
	err = png.Encode(current, img)
	if err != nil {
		e.t.Fatalf("failed encode snapshot image %v", currentPath(path))
	}

	err = os.WriteFile(currentPath(path), current.Bytes(), 0o644)
	if err != nil {
		e.t.Fatalf("failed to write snapshot %v", currentPath(path))
	}

	if diffImg != nil {
		diff := bytes.NewBuffer(nil)
		err = png.Encode(diff, diffImg)
		if err != nil {
			e.t.Fatalf("failed encode diff image %v", diffPath(path))
		}

		err = os.WriteFile(diffPath(path), diff.Bytes(), 0o644)
		if err != nil {
			e.t.Fatalf("failed to diff snapshot %v", diffPath(path))
		}
	}

	return e
}

func currentPath(i string) string {
	return strings.TrimSuffix(i, ".png") + ".current.png"
}

func diffPath(i string) string {
	return strings.TrimSuffix(i, ".png") + ".diff.png"
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
			if avg > opts.pixelTolerance {
				mismatches++
			}
		}
	}

	m := float64(mismatches) / float64(snapshotSize.X*snapshotSize.Y)

	if m > opts.matchTolerance {
		t.Errorf("expected image does not match snapshot, %.1f%% of pixels do not match", m*100)
		return false, diffImg
	}

	return true, nil
}

func getDiffFor(rs, rc uint32) float64 {
	return math.Abs(float64(rs)-float64(rc)) / (256*256 - 1)
}
