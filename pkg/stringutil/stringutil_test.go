package stringutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsInteger(t *testing.T) {
	if !IsInteger("123123") {
		t.Error("number was not detected as number")
	}

	if !IsInteger("-23896472386472863") {
		t.Error("negative number was not detected as number")
	}

	if IsInteger("") {
		t.Error("empty string was detected as number")
	}

	if IsInteger("123123a123") {
		t.Error("not-number string was detected as number")
	}
}

func TestEnsureNotEmpty(t *testing.T) {
	const (
		val = "val"
		def = "def"
	)

	if EnsureNotEmpty("", def) != def {
		t.Error("did not return default string")
	}

	if EnsureNotEmpty(val, def) != val {
		t.Error("did not return value string")
	}

	if EnsureNotEmpty("", "") != "" {
		t.Error("return value was not empty")
	}
}

func TestFromBool(t *testing.T) {
	const (
		tr = "true"
		fa = "false"
	)

	if FromBool(true, tr, fa) != tr {
		t.Error("true does not return true string")
	}

	if FromBool(false, tr, fa) != fa {
		t.Error("false does not return false string")
	}
}

func TestIndexOf(t *testing.T) {
	arr := []string{"0", "1", "2", "3", "4", "5"}

	if i := IndexOf("2", arr); i != 2 {
		t.Errorf("got index %d instead of 2", i)
	}

	if i := IndexOf("6", arr); i != -1 {
		t.Errorf("got index %d instead of -1", i)
	}
}

func TestContainsAny(t *testing.T) {
	arr := []string{"0", "1", "2", "3", "4", "5"}

	if !ContainsAny("1", arr) {
		t.Error("contained value was detected as not contained")
	}

	if ContainsAny("6", arr) {
		t.Error("not contained value was detected as contained")
	}

	if ContainsAny("", arr) {
		t.Error("not contained value was detected as contained")
	}
}

func TestContained(t *testing.T) {
	var arr, subset, ct []string

	arr = []string{"0", "1", "2", "3", "4", "5"}
	subset = nil
	ct = Contained(subset, arr)
	assert.ElementsMatch(t, ct, []string{})

	arr = nil
	subset = []string{"0", "1"}
	ct = Contained(subset, arr)
	assert.ElementsMatch(t, ct, []string{})

	arr = []string{"0", "1", "2", "3", "4", "5"}
	subset = []string{"0", "1", "2"}
	ct = Contained(subset, arr)
	assert.ElementsMatch(t, ct, []string{"0", "1", "2"})

	arr = []string{"0", "1", "2", "3", "4", "5"}
	subset = []string{"0", "1", "2", "6", "7", "8"}
	ct = Contained(subset, arr)
	assert.ElementsMatch(t, ct, []string{"0", "1", "2"})

	arr = []string{"0", "1", "2", "3", "4", "5"}
	subset = []string{"6", "7", "8"}
	ct = Contained(subset, arr)
	assert.ElementsMatch(t, ct, []string{})
}

func TestNotContained(t *testing.T) {
	var arr, must, nc []string

	arr = []string{"0", "1", "2", "3", "4", "5"}
	must = nil
	nc = NotContained(must, arr)
	assert.ElementsMatch(t, nc, []string{})

	arr = nil
	must = []string{"0", "1", "2"}
	nc = NotContained(must, arr)
	assert.ElementsMatch(t, nc, must)

	arr = []string{"0", "1", "2", "3", "4", "5"}
	must = []string{"0", "1", "2"}
	nc = NotContained(must, arr)
	assert.ElementsMatch(t, nc, []string{})

	arr = []string{"0", "1", "2", "3", "4", "5"}
	must = []string{"0", "1", "2", "6", "7"}
	nc = NotContained(must, arr)
	assert.ElementsMatch(t, nc, []string{"6", "7"})
}

func TestHasPrefixAny(t *testing.T) {
	if !HasPrefixAny("kekw", "p", "ke") {
		t.Error("falsely detected has no prefix")
	}

	if !HasPrefixAny("kekw", "p", "u", "k") {
		t.Error("falsely detected has no prefix")
	}

	if HasPrefixAny("kekw", "p", "u") {
		t.Error("falsely detected has prefix")
	}

	if HasPrefixAny("kekw") {
		t.Error("falsely detected has prefix")
	}

	if HasPrefixAny("", "p", "u") {
		t.Error("falsely detected has prefix")
	}
}

func TestHasSuffixAny(t *testing.T) {
	if !HasSuffixAny("kekw", "p", "kw") {
		t.Error("falsely detected has no prefix")
	}

	if !HasSuffixAny("kekw", "p", "u", "w") {
		t.Error("falsely detected has no prefix")
	}

	if HasSuffixAny("kekw", "p", "u") {
		t.Error("falsely detected has prefix")
	}

	if HasSuffixAny("kekw") {
		t.Error("falsely detected has prefix")
	}

	if HasSuffixAny("", "p", "u") {
		t.Error("falsely detected has prefix")
	}
}

func TestSpice(t *testing.T) {
	assert.Equal(t,
		Splice([]string{"a", "b", "c"}, -1),
		[]string{"a", "b", "c"})
	assert.Equal(t,
		Splice([]string{"a", "b", "c"}, 3),
		[]string{"a", "b", "c"})
	assert.Equal(t,
		Splice([]string{"a", "b", "c"}, 0),
		[]string{"b", "c"})
	assert.Equal(t,
		Splice([]string{"a", "b", "c"}, 2),
		[]string{"a", "b"})
	assert.Equal(t,
		Splice([]string{"a", "b", "c"}, 1),
		[]string{"a", "c"})
}

func TestCapitalize(t *testing.T) {
	assert.Equal(t, "", Capitalize("", false))
	assert.Equal(t, "", Capitalize("", true))

	assert.Equal(t, "H", Capitalize("h", false))
	assert.Equal(t, "H", Capitalize("h", true))

	assert.Equal(t, "Hey", Capitalize("hey", false))
	assert.Equal(t, "Hey", Capitalize("hey", true))

	assert.Equal(t, "Hey was geht ab", Capitalize("hey was geht ab", false))
	assert.Equal(t, "Hey Was Geht Ab", Capitalize("hey was geht ab", true))

	assert.Equal(t, "Hey Was Geht Ab", Capitalize("Hey Was Geht Ab", false))
	assert.Equal(t, "Hey Was Geht Ab", Capitalize("Hey Was Geht Ab", true))
}
