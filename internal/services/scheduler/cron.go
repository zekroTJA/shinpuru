package scheduler

import (
	"errors"

	"github.com/robfig/cron/v3"
)

// CronScheduler implements LifeCycleTimer
// wrapping a cron.Cron instance.
type CronScheduler struct {
	C *cron.Cron
}

func (lct *CronScheduler) Schedule(spec interface{}, job func()) (id interface{}, err error) {
	specStr, ok := spec.(string)
	if !ok {
		return nil, errors.New("invalid spec type: must be a string")
	}
	return lct.C.AddFunc(specStr, job)
}

func (lct *CronScheduler) Unschedule(id interface{}) error {
	cid, ok := id.(cron.EntryID)
	if !ok {
		return errors.New("invalid id type")
	}
	lct.C.Remove(cid)
	return nil
}

func (lct *CronScheduler) Start() {
	lct.C.Start()
}

func (lct *CronScheduler) Stop() {
	lct.C.Stop()
}
