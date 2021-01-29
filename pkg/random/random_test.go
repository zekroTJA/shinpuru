package random

import "testing"

func TestGetRandBase64Str(t *testing.T) {
	v, err := GetRandBase64Str(0)
	if v != "" || err != ErrInvalidLen {
		t.Error("invalid length not detected")
	}

	v, err = GetRandBase64Str(-1)
	if v != "" || err != ErrInvalidLen {
		t.Error("invalid length not detected")
	}

	v, err = GetRandBase64Str(10)
	if err != nil {
		t.Error("valid errored: ", err)
	}
	if len(v) != 10 {
		t.Error("invalid result length: ", len(v))
	}

	v, err = GetRandBase64Str(16)
	if err != nil {
		t.Error("valid errored: ", err)
	}
	if len(v) != 16 {
		t.Error("invalid result length: ", len(v))
	}
}

func TestGetRandBase64StrUniqueness(t *testing.T) {
	const n = 100000

	var err error
	arr := make([]string, n)
	for i := range arr {
		if arr[i], err = GetRandBase64Str(10); err != nil {
			t.Error(err)
		}
	}

	if hasDuplicatesStrings(arr) {
		t.Error("has duplicates")
	}
}

func TestGetRandByteArray(t *testing.T) {
	v, err := GetRandByteArray(0)
	if v != nil || err != ErrInvalidLen {
		t.Error("invalid length not detected")
	}

	v, err = GetRandByteArray(-1)
	if v != nil || err != ErrInvalidLen {
		t.Error("invalid length not detected")
	}

	v, err = GetRandByteArray(10)
	if err != nil {
		t.Error("valid errored: ", err)
	}
	if len(v) != 10 {
		t.Error("invalid result length: ", len(v))
	}

	v, err = GetRandByteArray(16)
	if err != nil {
		t.Error("valid errored: ", err)
	}
	if len(v) != 16 {
		t.Error("invalid result length: ", len(v))
	}
}

func TestGetRandByteArrayUniqueness(t *testing.T) {
	const n = 100000

	var err error
	arr := make([][]byte, n)
	for i := range arr {
		if arr[i], err = GetRandByteArray(10); err != nil {
			t.Error(err)
		}
	}

	if hasDuplicatesByteArrays(arr) {
		t.Error("has duplicates")
	}
}

// --- HELPER ---------

func hasDuplicatesStrings(arr []string) bool {
	m := make(map[string]*struct{})

	for _, v := range arr {
		if _, ok := m[v]; ok {
			return true
		}
		m[v] = nil
	}

	return false
}

func hasDuplicatesByteArrays(arr [][]byte) bool {
	m := make(map[string]*struct{})

	for _, v := range arr {
		if _, ok := m[string(v)]; ok {
			return true
		}
		m[string(v)] = nil
	}

	return false
}
