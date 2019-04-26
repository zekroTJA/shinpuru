package vibrant

import "container/heap"

const (
	blackMaxLightness float64 = 0.05
	whiteMinLightness float64 = 0.95
)

// A color quantizer based on the Median-cut algorithm, optimized for
// picking out distinct colors rather than representation colors.
//
// The color space is represented as a 3-dimensional cube with each
// dimension being an RGB component. The cube is then repeatedly divided
// until we have reduced the color space to the requested number of colors.
// An average color is then generated from each cube.
//
// Whereas median-cut divides cubes so they all have roughly the same
// population, this quantizer divides boxes based on their color volume.
type colorCutQuantizer struct {
	Colors           []int
	ColorPopulations map[int]int
	QuantizedColors  []*Swatch
}

// true if the color is close to pure black, pure white, or
// "the red side of the I line" which I believe is Google-speak for
// "that particular shade of red which occurs in the red-eye effect"
// see enwp.org/Red-eye_effect
func shouldIgnoreColor(color int) bool {
	h, s, l := rgbToHsl(color)
	return l <= blackMaxLightness || l >= whiteMinLightness || (h >= 0.0278 && h <= 0.1028 && s <= 0.82)
}

func shouldIgnoreColorSwatch(sw *Swatch) bool {
	return shouldIgnoreColor(int(sw.Color))
}

func newColorCutQuantizer(bitmap bitmap, maxColors int) *colorCutQuantizer {
	pixels := bitmap.Pixels()
	histo := newColorHistogram(pixels)
	colorPopulations := make(map[int]int, histo.NumberColors)
	for i, c := range histo.Colors {
		colorPopulations[c] = histo.ColorCounts[i]
	}
	validColors := make([]int, 0)
	i := 0
	for _, c := range histo.Colors {
		if !shouldIgnoreColor(c) {
			validColors = append(validColors, c)
			i++
		}
	}
	validCount := len(validColors)
	ccq := &colorCutQuantizer{Colors: validColors, ColorPopulations: colorPopulations}
	if validCount <= maxColors {
		// note: no quantization actually occurs
		for _, c := range validColors {
			ccq.QuantizedColors = append(ccq.QuantizedColors, &Swatch{Color: Color(c), Population: colorPopulations[c]})
		}
	} else {
		ccq.quantizePixels(validCount-1, maxColors)
	}
	return ccq
}

// see also vbox.go
func (ccq *colorCutQuantizer) quantizePixels(maxColorIndex, maxColors int) {
	pq := make(priorityQueue, 0)
	heap.Init(&pq)
	heap.Push(&pq, newVbox(0, maxColorIndex, ccq.Colors, ccq.ColorPopulations))
	for pq.Len() < maxColors {
		v := heap.Pop(&pq).(*vbox)
		if v.CanSplit() {
			heap.Push(&pq, v.Split())
			heap.Push(&pq, v)
		} else {
			break
		}
	}
	for pq.Len() > 0 {
		v := heap.Pop(&pq).(*vbox)
		swatch := v.AverageColor()
		if !shouldIgnoreColorSwatch(swatch) {
			ccq.QuantizedColors = append(ccq.QuantizedColors, swatch)
		}
	}
}
