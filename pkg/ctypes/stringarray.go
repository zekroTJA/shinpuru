package ctypes

// StringArray extends a string slice by some
// useful functionalities.
type StringArray []string

// IndexOf returns the index of v in arr.
// If not found, the returned index is -1.
func (arr StringArray) IndexOf(v string) int {
	for i, s := range arr {
		if v == s {
			return i
		}
	}

	return -1
}

// Contains returns true when v is included
// in arr.
func (arr StringArray) Contains(v string) bool {
	return arr.IndexOf(v) > -1
}

// Slice returns a new array sliced at i by
// the range of r.
func (arr StringArray) Splice(i, r int) StringArray {
	l := len(arr)
	if i >= l {
		return arr
	}
	if i+r >= l {
		return arr[:i]
	}

	return append(arr[:i], arr[i+r:]...)
}
