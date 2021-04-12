package thumbnail

import (
	"image"
	"testing"
)

func TestMake(t *testing.T) {
	check(500, 1000, 2000, 250, 500, t)
	check(500, 500, 500, 500, 500, t)
	check(5, 2, 1, 2, 1, t)
	check(0, 2, 1, 0, 0, t)
	check(500, 1337, 6969, 95, 500, t)
}

func check(max, isW, isH, mustW, mustH int, t *testing.T) {
	t.Helper()

	in := image.NewRGBA(image.Rect(0, 0, isW, isH))
	out := Make(in, max)

	is := out.Bounds().Dx()
	if is != mustW {
		t.Errorf("width was %d - must %d", is, mustW)
	}

	is = out.Bounds().Dy()
	if is != mustH {
		t.Errorf("height was %d - must %d", is, mustH)
	}
}
