package core

import (
	"time"

	"github.com/zekroTJA/shinpuru/internal/util"
)

type LCHandler func(now time.Time)

type LCTimer struct {
	ticker   *time.Ticker
	handlers map[string]LCHandler

	stopChan chan bool
}

func NewLTCTimer(each time.Duration) *LCTimer {
	return &LCTimer{
		ticker:   time.NewTicker(each),
		handlers: make(map[string]LCHandler),
		stopChan: make(chan bool, 1),
	}
}

func (t *LCTimer) OnTick(handler LCHandler) func() {
	uid := util.NodeLCHandler.Generate().String()
	t.handlers[uid] = handler
	return func() {
		delete(t.handlers, uid)
	}
}

func (t *LCTimer) OnTickOnce(handler LCHandler) func() {
	var unreg func()
	unreg = t.OnTick(func(now time.Time) {
		handler(now)
		unreg()
	})
	return unreg
}

func (t *LCTimer) AfterTime(after time.Time, handler LCHandler) func() {
	return t.OnTickOnce(func(now time.Time) {
		if now.After(after) {
			handler(now)
		}
	})
}

func (t *LCTimer) AfterDuration(after time.Duration, handler LCHandler) func() {
	afterTime := time.Now().Add(after)
	return t.AfterTime(afterTime, handler)
}

func (t *LCTimer) Start() {
	go func() {
		for {
			select {

			case now := <-t.ticker.C:
				for _, h := range t.handlers {
					h(now)
				}

			case <-t.stopChan:
				t.ticker.Stop()
				return

			}
		}
	}()
}

func (t *LCTimer) Stop() {
	go func() {
		t.stopChan <- true
	}()
}
