package codeexec

import (
	"errors"
	"strings"
	"time"

	ranna "github.com/ranna-go/ranna/pkg/client"
	"github.com/ranna-go/ranna/pkg/models"
	"github.com/sarulabs/di/v2"
	sharedmodels "github.com/zekroTJA/shinpuru/internal/models"
	"github.com/zekroTJA/shinpuru/internal/services/config"
	"github.com/zekroTJA/shinpuru/internal/util/static"
)

type RannaFactory struct {
	cfg *sharedmodels.CodeExecRanna
}

var _ Factory = (*RannaFactory)(nil)

func NewRannaFactory(container di.Container) (e *RannaFactory, err error) {
	e = &RannaFactory{}

	cfg := container.Get(static.DiConfig).(config.Provider)

	e.cfg = &cfg.Config().CodeExec.Ranna

	if e.cfg.Endpoint == "" {
		err = errors.New("no ranna endpoint provided")
		return
	}

	if e.cfg.Token != "" && !strings.Contains(e.cfg.Token, " ") {
		e.cfg.Token = "basic " + e.cfg.Token
	}

	return
}

func (e *RannaFactory) Name() string {
	return "ranna"
}

func (e *RannaFactory) Specs() (specs models.SpecMap, err error) {
	exec, err := e.NewExecutor("")
	if err != nil {
		return
	}

	client := exec.(*RannaExecutor).client

	specs, err = client.Spec()
	return
}

func (e *RannaFactory) NewExecutor(guildID string) (exec Executor, err error) {
	client, err := ranna.New(ranna.Options{
		Endpoint:      e.cfg.Endpoint,
		Version:       e.cfg.ApiVersion,
		Authorization: e.cfg.Token,
		UserAgent:     "shinpuru",
	})
	if err != nil {
		return
	}
	exec = &RannaExecutor{client}
	return
}

type RannaExecutor struct {
	client ranna.Client
}

func (e *RannaExecutor) Exec(p Payload) (res Response, err error) {
	r, err := e.client.Exec(models.ExecutionRequest{
		Language:         p.Language,
		Code:             p.Code,
		Arguments:        p.Args,
		Environment:      p.Environment,
		InlineExpression: p.Inline,
	})
	if err != nil {
		return
	}
	res.StdOut = r.StdOut
	res.StdErr = r.StdErr
	res.ExecTime = time.Duration(r.ExecTimeMS) * time.Millisecond
	return
}
