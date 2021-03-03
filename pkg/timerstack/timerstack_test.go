package timerstack

import (
	"testing"
	"time"
)

func TestAfter(t *testing.T) {
	ts := New().
		After(1*time.Second, func() bool { return true }).
		After(2*time.Second, func() bool { return true }).
		After(3*time.Second, func() bool { return true })

	if len(ts.stack) < 3 {
		t.Error("not all timers have been registered")
	}

	for i, timer := range ts.stack {
		if timer.delay != time.Duration(i+1)*time.Second {
			t.Error("invalid registered timer duration")
		}
	}
}

func TestRunBlocking(t *testing.T) {
	curr := 0

	ts := New().
		After(1*time.Second, func() bool {
			curr = 1
			return true
		}).
		After(2*time.Second, func() bool {
			curr = 2
			return true
		}).
		After(3*time.Second, func() bool {
			curr = 3
			return true
		})

	go ts.RunBlocking()

	if curr != 0 {
		t.Errorf("curr should have been 0 but was %d", curr)
	}

	time.Sleep(1010 * time.Millisecond)
	if curr != 1 {
		t.Errorf("curr should have been 1 but was %d", curr)
	}

	time.Sleep(3010 * time.Millisecond)
	if curr != 2 {
		t.Errorf("curr should have been 2 but was %d", curr)
	}

	time.Sleep(6010 * time.Millisecond)
	if curr != 3 {
		t.Errorf("curr should have been 3 but was %d", curr)
	}
}

func TestRunBlocking_ActionFalse(t *testing.T) {
	curr := 0

	ts := New().
		After(1*time.Second, func() bool {
			curr = 1
			return true
		}).
		After(2*time.Second, func() bool {
			curr = 2
			return false
		}).
		After(3*time.Second, func() bool {
			curr = 3
			return true
		})

	go ts.RunBlocking()

	if curr != 0 {
		t.Errorf("curr should have been 0 but was %d", curr)
	}

	time.Sleep(1010 * time.Millisecond)
	if curr != 1 {
		t.Errorf("curr should have been 1 but was %d", curr)
	}

	time.Sleep(3010 * time.Millisecond)
	if curr != 2 {
		t.Errorf("curr should have been 2 but was %d", curr)
	}

	time.Sleep(6010 * time.Millisecond)
	if curr != 2 {
		t.Errorf("curr should have been 2 but was %d", curr)
	}
}

func TestStop(t *testing.T) {
	curr := 0

	ts := New().
		After(1*time.Second, func() bool {
			curr = 1
			return true
		}).
		After(2*time.Second, func() bool {
			curr = 2
			return true
		})

	go ts.RunBlocking()

	if curr != 0 {
		t.Errorf("curr should have been 0 but was %d", curr)
	}

	time.Sleep(1010 * time.Millisecond)
	if curr != 1 {
		t.Errorf("curr should have been 1 but was %d", curr)
	}

	ts.Stop()

	time.Sleep(3010 * time.Millisecond)
	if curr != 1 {
		t.Errorf("curr should have been 1 but was %d", curr)
	}
}
