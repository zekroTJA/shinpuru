package limiter

import (
	"testing"
	"time"

	"github.com/zekroTJA/ratelimit"
)

func TestRetieve(t *testing.T) {
	var m *manager
	var rl, rl2 *ratelimit.Limiter

	m = newManager(1*time.Second, 500*time.Millisecond, 1)
	rl = m.retrieve("test")
	rl2 = m.retrieve("test")
	if rl != rl2 {
		t.Error("retrieved rl was not the same ")
	}

	m = newManager(1*time.Second, 500*time.Millisecond, 1)
	rl = m.retrieve("test1")
	rl2 = m.retrieve("test2")
	if rl == rl2 {
		t.Error("valid rl instance was reused")
	}
}
