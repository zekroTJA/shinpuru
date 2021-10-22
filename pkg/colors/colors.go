// Package color provides general utilities for
// image/color objects and color codes.
package colors

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"

	"github.com/generaltso/vibrant"
	"github.com/zekroTJA/shinpuru/pkg/httpreq"
)

// FromHex returns a color.RGBA object reference
// from the passed hexVal HEX RGBA color code.
//
// When the passed color code is malformed, an
// error is returned.
func FromHex(hexVal string) (*color.RGBA, error) {
	if hexVal == "" {
		return nil, errors.New("invalid color format")
	}

	if hexVal[0] == '#' {
		return FromHex(hexVal[1:])
	}

	v, err := hex.DecodeString(hexVal)
	if err != nil {
		return nil, err
	}

	if len(v) < 4 {
		v = append(v, 255)
	}

	return &color.RGBA{v[0], v[1], v[2], v[3]}, nil
}

// ToInt returns an integer color value from
// the oassed color.RGBA object reference.
func ToInt(clr *color.RGBA) int {
	return int(clr.B) | int(clr.G)<<8 | int(clr.R)<<16
}

// ToHex returns a HEX RBGA color string from
// the passed color.RGBA object reference.
func ToHex(clr *color.RGBA) string {
	return fmt.Sprintf("%06X", ToInt(clr))
}

// CreateImage generates a PNG image filled with
// the passed color in the size of the passed
// xSize and ySize dimensions.
//
// The generated image is returned as bytes.Buffer
// reference. When the image generation fails, an
// error is returned.
func CreateImage(clr *color.RGBA, xSize, ySize int) (*bytes.Buffer, error) {
	// Create image and fill it with the color
	// of the clr color object.
	img := image.NewRGBA(image.Rect(0, 0, xSize, ySize))
	draw.Draw(img, img.Bounds(), &image.Uniform{*clr}, image.Point{}, draw.Src)

	// Encode image object to image data using
	// the png encoder
	buff := bytes.NewBuffer([]byte{})
	if err := png.Encode(buff, img); err != nil {
		return nil, err
	}

	return buff, nil
}

// GetVibrantColorFromImage returns the vribrant accent
// color of an image passed.
func GetVibrantColorFromImage(img image.Image) (clr int, err error) {
	palette, err := vibrant.NewPaletteFromImage(img)
	if err != nil {
		return
	}

	for name, swatch := range palette.ExtractAwesome() {
		if name == "Vibrant" {
			clr = int(swatch.Color)
			break
		}
	}

	return
}

// GetVibrantColorFromImage requests the image from the given
// URL and returns the vribrant accent color of an image passed.
func GetVibrantColorFromImageUrl(url string) (clr int, err error) {
	body, _, err := httpreq.GetFile(url, nil)
	if err != nil {
		return
	}

	img, _, err := image.Decode(body)
	if err != nil {
		return
	}

	return GetVibrantColorFromImage(img)
}
