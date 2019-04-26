package vibrant

import (
	"image"
	"image/color"
	"math"

	"github.com/nfnt/resize"
)

// type bitmap is a simple wrapper for an image.Image
type bitmap struct {
	Width  int
	Height int
	Source image.Image
}

func newBitmap(input image.Image) *bitmap {
	bounds := input.Bounds()
	return &bitmap{bounds.Dx(), bounds.Dy(), input}
}

// Scales input image.Image by aspect ratio using https://github.com/nfnt/resize
func newScaledBitmap(input image.Image, ratio float64) *bitmap {
	bounds := input.Bounds()
	w := math.Ceil(float64(bounds.Dx()) * ratio)
	h := math.Ceil(float64(bounds.Dy()) * ratio)
	return &bitmap{int(w), int(h), resize.Resize(uint(w), uint(h), input, resize.Bilinear)}
}

// Returns all of the pixels of this bitmap.Source as a 1D array of image/color.Color
func (b *bitmap) Pixels() []color.Color {
	c := make([]color.Color, 0)
	for y := 0; y < b.Height; y++ {
		for x := 0; x < b.Width; x++ {
			c = append(c, b.Source.At(x, y))
		}
	}
	return c
}
