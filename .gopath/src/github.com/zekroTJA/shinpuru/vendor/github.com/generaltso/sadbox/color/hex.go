// Copyright 2012 The Gorilla Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package color

import (
	"fmt"
	"image/color"
	"strconv"
)

// HexModel converts any Color to an Hex color.
var HexModel = color.ModelFunc(hexModel)

// Hex represents an RGB color in hexadecimal format.
//
// The length must be 3 or 6 characters, preceded or not by a '#'.
type Hex string

// RGBA returns the alpha-premultiplied red, green, blue and alpha values
// for the Hex.
func (c Hex) RGBA() (uint32, uint32, uint32, uint32) {
	r, g, b := HexToRGB(c)
	return uint32(r) * 0x101, uint32(g) * 0x101, uint32(b) * 0x101, 0xffff
}

// hexModel converts a Color to Hex.
func hexModel(c color.Color) color.Color {
	if _, ok := c.(Hex); ok {
		return c
	}
	r, g, b, _ := c.RGBA()
	return RGBToHex(uint8(r>>8), uint8(g>>8), uint8(b>>8))
}

// RGBToHex converts an RGB triple to an Hex string.
func RGBToHex(r, g, b uint8) Hex {
	return Hex(fmt.Sprintf("#%02X%02X%02X", r, g, b))
}

// HexToRGB converts an Hex string to a RGB triple.
func HexToRGB(h Hex) (uint8, uint8, uint8) {
	if len(h) > 0 && h[0] == '#' {
		h = h[1:]
	}
	if len(h) == 3 {
		h = h[:1] + h[:1] + h[1:2] + h[1:2] + h[2:] + h[2:]
	}
	if len(h) == 6 {
		if rgb, err := strconv.ParseUint(string(h), 16, 32); err == nil {
			return uint8(rgb >> 16), uint8((rgb >> 8) & 0xFF), uint8(rgb & 0xFF)
		}
	}
	return 0, 0, 0
}
