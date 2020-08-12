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
	return fmt.Sprintf("%x", ToInt(clr))
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
