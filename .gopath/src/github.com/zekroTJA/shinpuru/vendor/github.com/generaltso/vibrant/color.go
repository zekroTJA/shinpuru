package vibrant

import (
	"fmt"
	"image/color"
	"math"

	colorconv "github.com/generaltso/sadbox/color" // by rodrigo moraes, exported from google code
)

type Color int

// Same as RGBHex()
func (c Color) String() string {
	return c.RGBHex()
}

func (c Color) RGB() (r, g, b int) {
	return unpackColor(int(c))
}

// e.g. "#bada55"
func (c Color) RGBHex() string {
	r, g, b := unpackColor(int(c))
	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

// Returns either black or white based on contrastRatio.
func (c Color) TextColor(contrastRatio float64) Color {
	if contrast(0xffffff, int(c)) >= contrastRatio {
		return Color(0xffffff)
	}
	return Color(0)
}

// Returns either black or white based on MIN_CONTRAST_TITLE_TEXT
func (c Color) TitleTextColor() Color {
	return c.TextColor(MIN_CONTRAST_TITLE_TEXT)
}

// Returns either black or white based on MIN_CONTRAST_BODY_TEXT
func (c Color) BodyTextColor() Color {
	return c.TextColor(MIN_CONTRAST_BODY_TEXT)
}

// takes an image/color.Color and normalizes it into
// r, g, b components in the range of 0-255
func colorToRgb(c color.Color) (int, int, int) {
	r, g, b, _ := c.RGBA()
	return int(r >> 8), int(g >> 8), int(b >> 8)
}

// takes r, g, b components in the range of 0-255 and packs them into
// a 24-bit int
func packColor(r, g, b int) int {
	return (r << 16) | (g << 8) | b
}

// inverse of packColor
func unpackColor(color int) (r, g, b int) {
	r = color >> 16 & 0xff
	g = color >> 8 & 0xff
	b = color >> 0 & 0xff
	return r, g, b
}

// floating point version of unpackColor
func unpackColorFloat(color int) (r, g, b float64) {
	ir, ig, ib := unpackColor(color)
	r = float64(ir)
	g = float64(ig)
	b = float64(ib)
	return r, g, b
}

// given a 24-bit int color (aka HTML hex aka #FFFFFF = 0xFFFFFF = white)
// returns Hue, Saturation, and Lightness components
// uses github.com/generaltso/sadbox/color for conversion because math is hard
// by rodrigo moraes, exported from google code
func rgbToHsl(color int) (h, s, l float64) {
	r, g, b := unpackColor(color)
	h, s, l = colorconv.RGBToHSL(uint8(r), uint8(g), uint8(b))
	return
}

// given Hue, Saturation, and Lightness components, returns a 24-bit int color
// uses github.com/generaltso/sadbox/color for conversion because math is hard
// by rodrigo moraes, exported from google code
func hslToRgb(h, s, l float64) (rgb int) {
	r, g, b := colorconv.HSLToRGB(h, s, l)
	return packColor(int(r), int(g), int(b))
}

// returns the contrast ratio of 24-bit int colors fg and bg (foreground and background)
func contrast(fg, bg int) float64 {
	lum1 := luminance(unpackColorFloat(fg))
	lum2 := luminance(unpackColorFloat(bg))
	return math.Max(lum1, lum2) / math.Min(lum1, lum2)
}

// http://www.w3.org/TR/2008/REC-WCAG20-20081211/#relativeluminancedef
func luminance(red, green, blue float64) float64 {
	red /= 255.0
	if red < 0.03928 {
		red /= 12.92
	} else {
		red = math.Pow((red+0.055)/1.055, 2.4)
	}
	green /= 255.0
	if green < 0.03928 {
		green /= 12.92
	} else {
		green = math.Pow((green+0.055)/1.055, 2.4)
	}
	blue /= 255.0
	if blue < 0.03928 {
		blue /= 12.92
	} else {
		blue = math.Pow((blue+0.055)/1.055, 2.4)
	}
	return (0.2126 * red) + (0.7152 * green) + (0.0722 * blue)
}
