package lctimer

import (
	"time"

	"github.com/zekroTJA/shinpuru/internal/util/snowflakenodes"
)

// Handler is a lifecycle timer
// elapse callback function
type Handler func(now time.Time)

// LifeCycleTimer provides functionalities to
// execute registered handlers on lifetime
// timer elapse.
type LifeCycleTimer struct {
	ticker   *time.Ticker
	handlers map[string]Handler

	stopChan chan bool
}

// New initializes a new LifeCycleTimer instance
// with the given elapse duration.
//
// This function does not start the actual timer.
func New(each time.Duration) *LifeCycleTimer {
	return &LifeCycleTimer{
		ticker:   time.NewTicker(each),
		handlers: make(map[string]Handler),
		stopChan: make(chan bool, 1),
	}
}

// OnTick executes the passed handler function on
// each life cycle timer elapse.
//
// Returned function removes the handler on call.
func (t *LifeCycleTimer) OnTick(handler Handler) func() {
	uid := snowflakenodes.NodeLCHandler.Generate().String()
	t.handlers[uid] = handler
	return func() {
		delete(t.handlers, uid)
	}
}

// OnTickOnce executes the passed handler once at
// next life time cycle elapse.
//
// Returned function removes the handler on call.
func (t *LifeCycleTimer) OnTickOnce(handler Handler) func() {
	var unreg func()
	unreg = t.OnTick(func(now time.Time) {
		handler(now)
		unreg()
	})
	return unreg
}

// AfterTimeOnce is shorthand for OnTickOnce, but only
// executes the handler on timer elapse after specified
// timestamp.
//
// Returned function removes the handler on call.
func (t *LifeCycleTimer) AfterTimeOnce(after time.Time, handler Handler) func() {
	return t.OnTickOnce(func(now time.Time) {
		if now.After(after) {
			handler(now)
		}
	})
}

// AfterDurationOnce is shorthand for OnTickOnce, but only
// executes the handler on timer elapse after specified
// duration has elapsed.
//
// Returned function removes the handler on call.
func (t *LifeCycleTimer) AfterDurationOnce(after time.Duration, handler Handler) func() {
	afterTime := time.Now().Add(after)
	return t.AfterTimeOnce(afterTime, handler)
}

// Start starts the life cycle timer loop.
func (t *LifeCycleTimer) Start() {
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

// Stop stops the life cycle timer loop.
func (t *LifeCycleTimer) Stop() {
	go func() {
		t.stopChan <- true
	}()
}
