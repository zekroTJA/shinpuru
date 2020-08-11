package colors

import (
	"encoding/hex"
	"errors"
	"image/color"
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
