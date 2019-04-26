// Copyright 2012 The Gorilla Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package color

import (
	"image/color"
	"math"
)

// HSVModel converts any Color to a HSV color.
var HSVModel = color.ModelFunc(hsvModel)

// HSV represents a cylindrical coordinate of points in an RGB color model.
//
// Values are in the range 0 to 1.
type HSV struct {
	H, S, V float64
}

// RGBA returns the alpha-premultiplied red, green, blue and alpha values
// for the HSV.
func (c HSV) RGBA() (uint32, uint32, uint32, uint32) {
	r, g, b := HSVToRGB(c.H, c.S, c.V)
	return uint32(r) * 0x101, uint32(g) * 0x101, uint32(b) * 0x101, 0xffff
}

// hsvModel converts a Color to HSV.
func hsvModel(c color.Color) color.Color {
	if _, ok := c.(HSV); ok {
		return c
	}
	r, g, b, _ := c.RGBA()
	h, s, v := RGBToHSV(uint8(r>>8), uint8(g>>8), uint8(b>>8))
	return HSV{h, s, v}
}

// RGBToHSV converts an RGB triple to a HSV triple.
//
// Ported from http://goo.gl/Vg1h9
func RGBToHSV(r, g, b uint8) (h, s, v float64) {
	fR := float64(r) / 255
	fG := float64(g) / 255
	fB := float64(b) / 255
	max := math.Max(math.Max(fR, fG), fB)
	min := math.Min(math.Min(fR, fG), fB)
	d := max - min
	s, v = 0, max
	if max > 0 {
		s = d / max
	}
	if max == min {
		// Achromatic.
		h = 0
	} else {
		// Chromatic.
		switch max {
		case fR:
			h = (fG - fB) / d
			if fG < fB {
				h += 6
			}
		case fG:
			h = (fB-fR)/d + 2
		case fB:
			h = (fR-fG)/d + 4
		}
		h /= 6
	}
	return
}

// HSVToRGB converts an HSV triple to a RGB triple.
//
// Ported from http://goo.gl/Vg1h9
func HSVToRGB(h, s, v float64) (r, g, b uint8) {
	var fR, fG, fB float64
	i := math.Floor(h * 6)
	f := h*6 - i
	p := v * (1.0 - s)
	q := v * (1.0 - f*s)
	t := v * (1.0 - (1.0-f)*s)
	switch int(i) % 6 {
	case 0:
		fR, fG, fB = v, t, p
	case 1:
		fR, fG, fB = q, v, p
	case 2:
		fR, fG, fB = p, v, t
	case 3:
		fR, fG, fB = p, q, v
	case 4:
		fR, fG, fB = t, p, v
	case 5:
		fR, fG, fB = v, p, q
	}
	r, g, b = float64ToUint8(fR), float64ToUint8(fG), float64ToUint8(fB)
	return
}
