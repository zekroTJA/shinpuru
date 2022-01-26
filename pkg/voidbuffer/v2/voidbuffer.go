// Package voidbuffer provides a simple, concurrency
// proof push buffer with a fixed size which "removes"
// firstly pushed values when fully filled.
package voidbuffer

import "sync"

// VoidBuffer wraps a push buffer with a fixed size
// which will scrap first pushed values when pusing
// into a fully filled buffer.
type VoidBuffer[T comparable] struct {
	m *sync.RWMutex

	buf  []T
	last int
	size int
}

// New initializes a new VoidBuffer with
// the passed size.
func New[T comparable](size int) *VoidBuffer[T] {
	return &VoidBuffer[T]{
		m:    &sync.RWMutex{},
		buf:  make([]T, size),
		last: -1,
		size: size,
	}
}

// Push adds the passed value into the buffer.
// If the buffer is full, the first input value
// to the buffer will be "pushed out".
func (vb *VoidBuffer[T]) Push(v T) {
	vb.m.Lock()
	defer vb.m.Unlock()

	vb.last++
	if vb.last >= vb.size {
		vb.last = 0
	}

	vb.buf[vb.last] = v
}

// Get returns the value in the buffer at the
// specified position. Order is determinated
// by push order.
func (vb *VoidBuffer[T]) Get(i int) (v T) {
	vb.m.RLock()
	defer vb.m.RUnlock()

	if i < 0 || i >= vb.size {
		panic("index out of range")
	}

	if vb.last == -1 {
		return def[T]()
	}

	pos := vb.last - i
	if pos < 0 {
		pos += vb.size
	}

	return vb.buf[pos]
}

// Contains is true if the passed value is
// currently contained in the buffer at any
// position.
func (vb *VoidBuffer[T]) Contains(v T) bool {
	vb.m.RLock()
	defer vb.m.RUnlock()

	for _, e := range vb.buf {
		if e == v {
			return true
		}
	}

	return false
}

// Flush sets all values of the buffer to
// nil.
func (vb *VoidBuffer[T]) Flush() {
	vb.m.Lock()
	defer vb.m.Unlock()

	for i := range vb.buf {
		vb.buf[i] = def[T]()
	}
}

// Size returns the fixed predefined size
// of the buffer.
func (vb *VoidBuffer[T]) Size() int {
	return vb.size
}

func (vb *VoidBuffer[T]) Snapshot() (r []T) {
	vb.m.RLock()
	defer vb.m.RUnlock()

	r = make([]T, vb.Size())
	copy(r, vb.buf)

	return
}
