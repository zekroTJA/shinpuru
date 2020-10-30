package bytecount

import "testing"

func TestFormat(t *testing.T) {
	if r := Format(615); r != "615 B" {
		t.Errorf("wrong format: was %s", r)
	}

	if r := Format(5623); r != "5.491 kiB" {
		t.Errorf("wrong format: was %s", r)
	}

	if r := Format(4425623); r != "4.221 MiB" {
		t.Errorf("wrong format: was %s", r)
	}

	if r := Format(9134426623); r != "8.507 GiB" {
		t.Errorf("wrong format: was %s", r)
	}

	if r := Format(4534425386623); r != "4.124 TiB" {
		t.Errorf("wrong format: was %s", r)
	}
}
