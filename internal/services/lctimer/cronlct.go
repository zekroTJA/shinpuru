package lctimer

import (
	"errors"

	"github.com/robfig/cron/v3"
)

type CronLifeCycleTimer struct {
	C *cron.Cron
}

func (lct *CronLifeCycleTimer) Schedule(spec string, job func()) (id interface{}, err error) {
	return lct.C.AddFunc(spec, job)
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
