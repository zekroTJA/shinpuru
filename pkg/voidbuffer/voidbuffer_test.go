package voidbuffer

import (
	"testing"
)

func TestNewVoidBuffer(t *testing.T) {
	const size = 5

	vb := New(size)

	if vb == nil {
		t.Error("new voidbuffer was nil")
	}
	if len(vb.buf) != size {
		t.Errorf("buffer size does not match (must %d, was %d)", size, len(vb.buf))
	}
}

func TestPush(t *testing.T) {
	vb := getPreFilled()

	checkPosVal(t, vb, 0, 6)
	checkPosVal(t, vb, 1, 7)
	checkPosVal(t, vb, 2, 8)
	checkPosVal(t, vb, 3, 4)
	checkPosVal(t, vb, 4, 5)
}

func TestGet(t *testing.T) {
	vb := getPreFilled()

	// 8 7 6 5 4

	// 6 7 8 4 5
	//     ^
	// Here is next

	v := vb.Get(1).(int)
	if v != 7 {
		t.Errorf("recovered value was invalid (was %d, must %d)", v, 7)
	}

	v = vb.Get(3).(int)
	if v != 5 {
		t.Errorf("recovered value was invalid (was %d, must %d)", v, 5)
	}

	vb = New(5)
	vi := vb.Get(1)
	if vi != nil {
		t.Errorf("nil value was not nil: %v", vi)
	}
}

func TestContains(t *testing.T) {
	vb := getPreFilled()

	if !vb.Contains(6) {
		t.Errorf("did not detect contained value %d", 6)
	}

	if vb.Contains(12) {
		t.Errorf("falsely detected non-contained value %d", 12)
	}
}

// --- HELPERS ---------------------------------------

func getPreFilled() (vb *VoidBuffer) {
	vb = New(5)

	// 1 2 3 4 5
	// 6 7 8 4 5

	for i := 1; i < 9; i++ {
		vb.Push(i)
	}

	return
}

func checkPosVal(t *testing.T, vb *VoidBuffer, i int, must interface{}) {
	t.Helper()

	if vb.buf[i] != must {
		t.Errorf("value invalid at [%d]: must %v, was %v",
			i, must, vb.buf[i])
	}
}
