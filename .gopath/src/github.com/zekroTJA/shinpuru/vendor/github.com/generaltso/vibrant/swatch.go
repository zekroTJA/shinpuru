package vibrant

import (
	"fmt"
	"strings"
)

type Swatch struct {
	Color      Color
	Population int
	Name       string // might be empty
}

// Convenience method that returns CSS e.g.
//	.vibrant{background-color:#bada55;color:#ffffff;}
func (sw *Swatch) String() string {
	return fmt.Sprintf(
		".%s{background-color:%s;color:%s;}",
		strings.ToLower(sw.Name),
		sw.Color.RGBHex(),
		sw.Color.BodyTextColor(),
	)
}
