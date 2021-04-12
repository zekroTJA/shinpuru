// Package thumbnail provides simple functionalities to
// generate thumbnails from images with a max witdh or
// height.
package thumbnail

import (
	"image"
	"math"

	"golang.org/x/image/draw"
)

// Make returns a scaled image from the given in image with
// the given maxSize width or height, maintaining the input
// images aspect ratio.
//
// Optionally, you can also pass a custom scaler
// implementation, if desired.
func Make(in image.Image, maxSize int, scaler ...draw.Scaler) (out image.Image) {
	out = in

	var sc draw.Scaler = draw.BiLinear
	if len(scaler) > 0 {
		sc = scaler[0]
	}

	width, height := in.Bounds().Dx(), in.Bounds().Dy()
	if width > maxSize || height > maxSize {
		var scale float64
		if width > height {
			scale = float64(maxSize) / float64(width)
			width = int(maxSize)
			height = int(math.Floor(float64(height) * scale))
		} else {
			scale = float64(maxSize) / float64(height)
			height = int(maxSize)
			width = int(math.Floor(float64(width) * scale))
		}
		outRect := image.Rect(0, 0, width, height)
		outImg := image.NewRGBA(outRect)
		sc.Scale(outImg, outRect, in, in.Bounds(), draw.Over, nil)
		out = outImg
	}

	return
}
