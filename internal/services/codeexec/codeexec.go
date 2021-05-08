package codeexec

import "time"

type Payload struct {
	Language    string
	Code        string
	Args        []string
	Environment map[string]string
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
	Languages() ([]string, error)
	NewExecutor(guildID string) (Executor, error)
}

type Executor interface {
	Exec(Payload) (Response, error)
}
