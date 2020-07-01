// Package bytecount provides functionalities to
// format byte counts.
package bytecount

import "fmt"

// Format returns a human readable string from the
// passed byte count.
//
// Example:
//   size := uint64(3371549327)
//   hSize := bytecount.Format(size)
//   fmt.Println(hSize)
//   // -> "3.140 GiB"
func Format(bc uint64) string {
	f1k := float64(1024)
	if bc < 1024 {
		return fmt.Sprintf("%d B", bc)
	}
	if bc < 1024*1024 {
		return fmt.Sprintf("%.3f kiB", float64(bc)/f1k)
	}
	if bc < 1024*1024*1024 {
		return fmt.Sprintf("%.3f MiB", float64(bc)/f1k/f1k)
	}
	if bc < 1024*1024*1024*1024 {
		return fmt.Sprintf("%.3f GiB", float64(bc)/f1k/f1k/f1k)
	}
	return fmt.Sprintf("%.3f TiB", float64(bc)/f1k/f1k/f1k/f1k)
}
