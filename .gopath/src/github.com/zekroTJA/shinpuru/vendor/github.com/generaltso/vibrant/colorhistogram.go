package vibrant

import (
	"image/color"
	"sort"
)

// Colors and ColorCounts are transformed into a map[int]int by colorCutQuantizer
// where key is the 24-bit color and value is the count
type colorHistogram struct {
	Colors       []int // 24-bit packed int colors
	ColorCounts  []int // index refers to above color
	NumberColors int
}

// See colorCutQuantizer source for how this is used.
func newColorHistogram(colorPixels []color.Color) *colorHistogram {
	// Transform []color.Color into array of 24-bit ints
	pixels := make([]int, len(colorPixels))
	for _, px := range colorPixels {
		pixels = append(pixels, packColor(colorToRgb(px)))
	}

	// Sort the pixels to enable counting
	sort.Ints(pixels)

	numColors := countDistinctColors(pixels)
	colors := make([]int, numColors)
	colorCounts := make([]int, numColors)

	if numColors > 0 {
		curIndex := 0
		curColor := pixels[0]
		colors[0] = curColor
		colorCounts[0] = 1

		for _, px := range pixels {
			if px == curColor {
				// same color, increase population
				colorCounts[curIndex]++
			} else {
				// new color, increase index
				curColor = px
				curIndex++
				colors[curIndex] = curColor
				colorCounts[curIndex] = 1
			}
		}
	}

	return &colorHistogram{colors, colorCounts, numColors}
}

func countDistinctColors(pixels []int) int {
	if len(pixels) < 2 {
		return len(pixels)
	}
	count := 1
	current := pixels[0]
	for _, px := range pixels {
		if px != current {
			// new color, increase population
			current = px
			count++
		}
	}
	return count
}
