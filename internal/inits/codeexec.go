package inits

import (
	"strings"

	"github.com/sarulabs/di/v2"
	"github.com/sirupsen/logrus"
	"github.com/zekroTJA/shinpuru/internal/config"
	"github.com/zekroTJA/shinpuru/internal/services/codeexec"
	"github.com/zekroTJA/shinpuru/internal/util/static"
)

func InitCodeExec(container di.Container) codeexec.Factory {
	cfg := container.Get(static.DiConfig).(*config.Config)

	if cfg.CodeExec == nil {
		cfg.CodeExec = cfg.Defaults.CodeExec
	}

	switch strings.ToLower(cfg.CodeExec.Type) {

	case "ranna":
		exec, err := codeexec.NewRannaFactory(container)
		if err != nil {
			logrus.WithError(err).Fatal("failed setting up ranna factroy")
		}
		return exec

	default:
		return codeexec.NewJdoodleFactory(container)
	}
}
