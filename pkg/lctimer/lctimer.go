// Package lctimer provides a life cycle timer which
// calls registered callback handlers on timer elapse.
package lctimer

import (
	"sync"
	"sync/atomic"
	"time"
)

// Handler is a lifecycle timer
// elapse callback function
type Handler func(now time.Time)

// LifeCycleTimer provides functionalities to
// execute registered handlers on lifetime
// timer elapse.
type LifeCycleTimer struct {
	ticker   *time.Ticker
	handlers *sync.Map

	rid      int32
	stopChan chan bool
}

// New initializes a new LifeCycleTimer instance
// with the given elapse duration.
//
// This function does not start the actual timer.
func New(each time.Duration) *LifeCycleTimer {
	return &LifeCycleTimer{
		ticker:   time.NewTicker(each),
		handlers: &sync.Map{},
		stopChan: make(chan bool, 1),
	}
}

// OnTick executes the passed handler function on
// each life cycle timer elapse.
//
// Returned function removes the handler on call.
func (t *LifeCycleTimer) OnTick(handler Handler) func() {
	uid := atomic.LoadInt32(&t.rid)
	atomic.AddInt32(&t.rid, 1)
	t.handlers.Store(uid, handler)
	return func() {
		t.handlers.Delete(uid)
	}
}

// OnTickOnce executes the passed handler once at
// next life time cycle elapse.
//
// Returned function removes the handler on call.
func (t *LifeCycleTimer) OnTickOnce(handler Handler) (unreg func()) {
	unreg = t.OnTick(func(now time.Time) {
		handler(now)
		unreg()
	})

	return
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
				t.handlers.Range(func(_, value interface{}) bool {
					h, ok := value.(Handler)
					if ok {
						h(now)
					}
					return true
				})

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
