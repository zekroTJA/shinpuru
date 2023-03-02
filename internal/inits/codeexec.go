package inits

import (
	"strings"

	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/services/codeexec"
	"github.com/zekroTJA/shinpuru/internal/services/config"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekrotja/rogu/log"
)

func InitCodeExec(container di.Container) codeexec.Factory {
	cfg := container.Get(static.DiConfig).(config.Provider)

	log := log.Tagged("CodeExec")
	log.Info().Msg("Initializing code execution ...")

	switch strings.ToLower(cfg.Config().CodeExec.Type) {

	case "ranna":
		exec, err := codeexec.NewRannaFactory(container)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed setting up ranna factroy")
		}
		return exec

	default:
		return codeexec.NewJdoodleFactory(container)
	}
}
