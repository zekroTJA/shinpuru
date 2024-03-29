// Package timerstack provides a timer which can
// execute multiple delayed functions one after one.
package timerstack

import (
	"sync"
	"time"
)

// Action is the function being executed.
// If the function returns false, the timer
// stack stops executing after this function.
type Action func() bool

type timer struct {
	delay  time.Duration
	action Action
}

// TimerStack allows to stack multiple timers on
// top to be executed one after one.
type TimerStack struct {
	stack []*timer

	mtx       sync.Mutex
	currTimer *time.Timer
	stopNext  bool
}

// New returns a new empty TimerStack.
func New() *TimerStack {
	return &TimerStack{
		stack: make([]*timer, 0),
	}
}

// After adds a new timer to the stack which is executed
// after the given time delay after the last executed
// timer. On execution, a is executed. If this function
// returns false, the execution stops after this function.
func (ts *TimerStack) After(d time.Duration, a Action) *TimerStack {
	ts.stack = append(ts.stack, &timer{d, a})
	return ts
}

// RunBlocking starts the timer queue blocking the
// current go-routine until all timers on the stack
// are executed or until the timer has been stoped.
func (ts *TimerStack) RunBlocking() {
	for _, t := range ts.stack {
		ts.mtx.Lock()
		ts.currTimer = time.NewTimer(t.delay)
		ts.mtx.Unlock()

		<-ts.currTimer.C

		ts.mtx.Lock()
		exit := !t.action() || ts.stopNext
		ts.mtx.Unlock()

		if exit {
			break
		}
	}
}

// Stop stops the timer stack execution.
func (ts *TimerStack) Stop() {
	ts.mtx.Lock()
	defer ts.mtx.Unlock()

	if ts.currTimer == nil {
		return
	}
	ts.stopNext = true
	ts.currTimer.Stop()
}
