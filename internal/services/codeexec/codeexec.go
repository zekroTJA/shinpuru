package codeexec

import (
	"time"

	"github.com/ranna-go/ranna/pkg/models"
)

var AvailableFactories = []string{"ranna", "jdoodle"}

type Payload struct {
	Language    string
	Code        string
	Args        []string
	Environment map[string]string
	Inline      bool
}

type Response struct {
	StdOut   string
	StdErr   string
	ExecTime time.Duration
	MemUsed  string
	CpuUsed  string
}

type Factory interface {
	Name() string
	Specs() (models.SpecMap, error)
	NewExecutor(guildID string) (Executor, error)
}

type Executor interface {
	Exec(Payload) (Response, error)
}
