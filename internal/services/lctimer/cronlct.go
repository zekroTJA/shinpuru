package lctimer

import (
	"errors"

	"github.com/robfig/cron/v3"
)

// CronLifeCycleTimer implements LifeCycleTimer
// wrapping a cron.Cron instance.
type CronLifeCycleTimer struct {
	C *cron.Cron
}

func (lct *CronLifeCycleTimer) Schedule(spec interface{}, job func()) (id interface{}, err error) {
	specStr, ok := spec.(string)
	if !ok {
		return nil, errors.New("invalid spec type: must be a string")
	}
	return lct.C.AddFunc(specStr, job)
}

func (lct *CronLifeCycleTimer) Unschedule(id interface{}) error {
	cid, ok := id.(cron.EntryID)
	if !ok {
		return errors.New("invalid id type")
	}
	lct.C.Remove(cid)
	return nil
}

func (lct *CronLifeCycleTimer) Start() {
	lct.C.Start()
}

func (lct *CronLifeCycleTimer) Stop() {
	lct.C.Stop()
}
