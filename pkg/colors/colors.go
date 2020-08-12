package colors

import (
	"bytes"
	"encoding/hex"
	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"strings"
)

func FromHex(hexVal string) (*color.RGBA, error) {
	if hexVal == "" {
		return nil, errors.New("invalid color format")
	}

	if hexVal[1] == '#' {
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

func CreateImage(hexVal string, xSize, ySize int) (*bytes.Buffer, error) {
	// Remove # when color code starts with it.
	if strings.HasPrefix(hexVal, "#") {
		hexVal = hexVal[1:]
	}

	// Trim and lowercase color code
	hexVal = strings.Trim(
		strings.ToLower(hexVal), " ")

	// Get color.RGBA object from color code
	clr, err := FromHex(hexVal)
	if err != nil {
		return nil, err
	}

	// Create image and fill it with the color
	// of the clr color object.
	img := image.NewRGBA(image.Rect(0, 0, xSize, ySize))
	draw.Draw(img, img.Bounds(), &image.Uniform{*clr}, image.Point{}, draw.Src)

	// Encode image object to image data using
	// the png encoder
	buff := bytes.NewBuffer([]byte{})
	if err = png.Encode(buff, img); err != nil {
		return nil, err
	}

	return buff, nil
}
