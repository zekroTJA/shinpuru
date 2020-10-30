package stringutil

import "testing"

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
