package shared

type LifeCycleTimer interface {
	Schedule(spec string, job func()) (id interface{}, err error)
	Unschedule(id interface{}) error
	Start()
	Stop()
}
