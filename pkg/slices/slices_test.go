package slices

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndexOf(t *testing.T) {
	assert.Equal(t,
		1,
		IndexOf([]int{1, 2, 3}, 2))
	assert.Equal(t,
		2,
		IndexOf([]int{1, 2, 3}, 3))
	assert.Equal(t,
		-1,
		IndexOf([]int{1, 2, 3}, 9))
	assert.Equal(t,
		-1,
		IndexOf([]int{}, 9))

	var s []int
	assert.Equal(t,
		-1,
		IndexOf(s, 9))
}

func TestContains(t *testing.T) {
	assert.True(t, Contains([]int{1, 2, 3}, 2))
	assert.False(t, Contains([]int{1, 2, 3}, 4))
	assert.False(t, Contains([]int{}, 4))

	var s []int
	assert.False(t, Contains(s, 4))
}

func TestSplice(t *testing.T) {
	var s, ns, rest []int

	s = []int{1, 2, 3, 4, 5}
	ns, rest = Splice(s, 1, 2)
	assert.Equal(t, []int{1, 4, 5}, ns)
	assert.Equal(t, []int{2, 3}, rest)

	s = []int{1, 2, 3, 4, 5}
	ns, rest = Splice(s, 0, 2)
	assert.Equal(t, []int{3, 4, 5}, ns)
	assert.Equal(t, []int{1, 2}, rest)

	s = []int{1, 2, 3, 4, 5}
	ns, rest = Splice(s, -1, 2)
	assert.Equal(t, []int{3, 4, 5}, ns)
	assert.Equal(t, []int{1, 2}, rest)

	s = []int{1, 2, 3, 4, 5}
	ns, rest = Splice(s, 3, 6)
	assert.Equal(t, []int{1, 2, 3}, ns)
	assert.Equal(t, []int{4, 5}, rest)

	s = []int{1, 2, 3, 4, 5}
	ns, rest = Splice(s, 1, 1)
	assert.Equal(t, []int{1, 3, 4, 5}, ns)
	assert.Equal(t, []int{2}, rest)
}
