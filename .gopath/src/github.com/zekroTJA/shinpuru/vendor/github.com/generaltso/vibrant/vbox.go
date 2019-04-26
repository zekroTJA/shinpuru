package vibrant

import (
	"math"
	"sort"
)

const (
	componentRed   int = -3
	componentGreen int = -2
	componentBlue      = -1
)

// Represents a tightly fitting box around a color space.
type vbox struct {
	lowerIndex  int
	upperIndex  int
	minRed      int
	maxRed      int
	minGreen    int
	maxGreen    int
	minBlue     int
	maxBlue     int
	colors      []int
	populations map[int]int
}

func newVbox(lowerIndex, upperIndex int, colors []int, populations map[int]int) *vbox {
	v := &vbox{lowerIndex: lowerIndex, upperIndex: upperIndex, colors: colors, populations: populations}
	v.fitBox()
	return v
}

// Recomputes the boundaries of this box to tightly fit the colors within
func (v *vbox) fitBox() {
	// Reset the min and max to opposite values
	v.minRed = 255
	v.minGreen = 255
	v.minBlue = 255
	v.maxRed = 0
	v.maxGreen = 0
	v.maxBlue = 0

	for i := v.lowerIndex; i <= v.upperIndex; i++ {
		r, g, b := unpackColor(v.colors[i])
		if r > v.maxRed {
			v.maxRed = r
		}
		if r < v.minRed {
			v.minRed = r
		}
		if g > v.maxGreen {
			v.maxGreen = g
		}
		if g < v.minGreen {
			v.minGreen = g
		}
		if b > v.maxBlue {
			v.maxBlue = b
		}
		if b < v.minBlue {
			v.minBlue = b
		}
	}
}

func (v *vbox) Volume() int {
	return (v.maxRed - v.minRed + 1) * (v.maxGreen - v.minGreen + 1) * (v.maxBlue - v.minBlue + 1)
}

func (v *vbox) CanSplit() bool {
	return (v.upperIndex - v.lowerIndex + 1) > 1
}

// Split this color box at the mid-point along its longest dimension
func (v *vbox) Split() *vbox {
	if !v.CanSplit() {
		panic("Cannot split a box with only 1 color!")
	}

	lenRed := v.maxRed - v.minRed
	lenGreen := v.maxGreen - v.minGreen
	lenBlue := v.maxBlue - v.minBlue

	// Find the longest color dimension, and then sort the
	// sub-array based on that dimension value in each color.
	//
	// Rather than define the sort logic, we modify each color so that
	// its most significant is the desired dimension:
	// see modifySignificantOctet
	var longestDim, midPoint int
	switch {
	case lenRed >= lenGreen && lenRed >= lenBlue:
		longestDim = componentRed

		// Already in RGB, no need to do anything
		v.sortColors()

		midPoint = (v.minRed + v.maxRed) / 2
	case lenGreen >= lenRed && lenGreen >= lenBlue:
		longestDim = componentGreen

		// RGB to GRB swap
		v.modifySignificantOctet(longestDim)

		v.sortColors()

		// Now revert all of the colors so that they are RGB again
		v.modifySignificantOctet(longestDim)

		midPoint = (v.minGreen + v.maxGreen) / 2
	default:
		longestDim = componentBlue

		// RGB to BGR swap
		v.modifySignificantOctet(longestDim)

		v.sortColors()

		// Now revert all of the colors so that they are RGB again
		v.modifySignificantOctet(longestDim)
		midPoint = (v.minBlue + v.maxBlue) / 2
	}

	// Iterate over the colors until a color is found with at least the
	// midpoint of the whole box's dimension midpoint.
	splitPoint := v.lowerIndex
loop:
	for i := v.lowerIndex; i <= v.upperIndex; i++ {
		r, g, b := unpackColor(v.colors[i])
		switch longestDim {
		case componentRed:
			if r >= midPoint {
				splitPoint = i
				break loop
			}
		case componentGreen:
			if g >= midPoint {
				splitPoint = i
				break loop
			}
		case componentBlue:
			if b >= midPoint {
				splitPoint = i
				break loop
			}
		}
	}

	vbox := newVbox(splitPoint+1, v.upperIndex, v.colors, v.populations)

	// Now change this box's upperIndex and recompute the color boundaries
	v.upperIndex = splitPoint
	v.fitBox()
	return vbox
}

func (v *vbox) sortColors() {
	section := v.colors[v.lowerIndex : v.upperIndex+1]
	sort.Ints(section)
	i := v.lowerIndex
	for _, color := range section {
		v.colors[i] = color
		i++
	}
}

// Modify the significant octet in a packed color int.
// Allows sorting based on the value of a single color component.
func (v *vbox) modifySignificantOctet(dim int) {
	switch dim {
	case componentRed:
		// Already in RGB, no need to do anything
		return
	case componentGreen:
		// RGB to GRB swap
		for i := v.lowerIndex; i <= v.upperIndex; i++ {
			r, g, b := unpackColor(v.colors[i])
			v.colors[i] = packColor(g, r, b)
		}
	case componentBlue:
		// RGB to BGR swap
		for i := v.lowerIndex; i <= v.upperIndex; i++ {
			r, g, b := unpackColor(v.colors[i])
			v.colors[i] = packColor(b, g, r)
		}
	}
}

func (v *vbox) AverageColor() *Swatch {
	sumRed := 0
	sumGreen := 0
	sumBlue := 0
	sumPop := 0
	for i := v.lowerIndex; i <= v.upperIndex; i++ {
		color := v.colors[i]
		r, g, b := unpackColor(color)
		pop := v.populations[color]
		sumPop += pop
		sumRed += r * pop
		sumGreen += g * pop
		sumBlue += b * pop
	}
	avgRed := round(float64(sumRed) / float64(sumPop))
	avgGreen := round(float64(sumGreen) / float64(sumPop))
	avgBlue := round(float64(sumBlue) / float64(sumPop))

	return &Swatch{Color: Color(packColor(avgRed, avgGreen, avgBlue)), Population: sumPop}
}

// there is no math.Round ._.
func round(val float64) int {
	var ret float64
	pow := math.Pow(10, 14)
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= 0.5 {
		ret = math.Ceil(digit)
	} else {
		ret = math.Floor(digit)
	}
	return int(ret / pow)
}
