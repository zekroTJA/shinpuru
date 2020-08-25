// Package angularservice provides bindings to
// start an Angular development server via the
// Angular CLI.
package angularservice

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"strconv"
)

// AngularService provides async bindings for  an
// Angular development server instance started by
// the Angular CLI.
type AngularService struct {
	options   *Options
	ctxCancel context.CancelFunc
	done      <-chan struct{}
}

// New returns a new instance of the AngularService
// with the passed options.
func New(options Options) *AngularService {
	return &AngularService{
		options: &options,
	}
}

// Start starts the Angular development server
// process.
func (s *AngularService) Start() (err error) {
	if s.options.Cd != "" {
		var currDir string
		if currDir, err = os.Getwd(); err != nil {
			return
		}
		if err = os.Chdir(s.options.Cd); err != nil {
			return
		}
		defer os.Chdir(currDir)
	}

	var ctx context.Context

	ctx, s.ctxCancel = context.WithCancel(context.Background())
	s.done = ctx.Done()

	cmdArgs := []string{"serve"}
	if s.options.Port > 0 {
		cmdArgs = append(cmdArgs, "--port", strconv.Itoa(s.options.Port))
	}

	if s.options.Args != nil {
		cmdArgs = append(cmdArgs, s.options.Args...)
	}

	cmd := exec.CommandContext(ctx, "ng", cmdArgs...)

	if s.options.Stdout != nil {
		cmd.Stdout = s.options.Stdout
	}
	if s.options.Stderr != nil {
		cmd.Stderr = s.options.Stderr
	}

	return cmd.Start()
}

// Stop stops a running Angular development server
// instance.
//
// Returns an error when no server instance is
// running on function call.
func (s *AngularService) Stop() error {
	if s.ctxCancel == nil {
		return errors.New("service was not started before")
	}

	s.ctxCancel()

	return nil
}

// Done returns the channel which is closed when the
// server process exited.
func (s *AngularService) Done() <-chan struct{} {
	return s.done
}
