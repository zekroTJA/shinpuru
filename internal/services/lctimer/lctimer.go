package lctimer

// LifeCycleTimer provides a job scheduler
// which executes given jobs at a specified
// time or in a given location.
type LifeCycleTimer interface {

	// Schedule adds a job to the scheduler which is
	// executed according to the given spec.
	//
	// Returns an id of the scheduled job and n error
	// when the job could not be scheduled.
	Schedule(spec interface{}, job func()) (id interface{}, err error)

	// Unschedule removes the given job by its id
	// from the scheduler so it will not be executed
	// anymore.
	Unschedule(id interface{}) error

	// Start runs the scheduler cycle.
	Start()

	// Stop cancels the scheduler cycle.
	Stop()
}
