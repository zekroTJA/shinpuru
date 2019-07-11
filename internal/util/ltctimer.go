package util

import "time"

type LTCHandler func(now time.Time)

type LTCTimer struct {
	ticker   *time.Ticker
	handlers map[string]LTCHandler

	stopChan chan bool
}

func NewLTCTimer(each time.Duration) *LTCTimer {
	return &LTCTimer{
		ticker:   time.NewTicker(each),
		handlers: make(map[string]LTCHandler),
		stopChan: make(chan bool, 1),
	}
}

func (t *LTCTimer) OnTick(handler LTCHandler) func() {
	uid := NodeLTCHandler.Generate().String()
	t.handlers[uid] = handler
	return func() {
		delete(t.handlers, uid)
	}
}

func (t *LTCTimer) OnTickOnce(handler LTCHandler) func() {
	var unreg func()
	unreg = t.OnTick(func(now time.Time) {
		handler(now)
		unreg()
	})
	return unreg
}

func (t *LTCTimer) AfterTime(after time.Time, handler LTCHandler) func() {
	return t.OnTickOnce(func(now time.Time) {
		if now.After(after) {
			handler(now)
		}
	})
}

func (t *LTCTimer) AfterDuration(after time.Duration, handler LTCHandler) func() {
	afterTime := time.Now().Add(after)
	return t.AfterTime(afterTime, handler)
}

func (t *LTCTimer) Start() {
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

func (t *LTCTimer) Stop() {
	go func() {
		t.stopChan <- true
	}()
}
