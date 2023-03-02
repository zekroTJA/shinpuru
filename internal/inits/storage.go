package inits

import (
	"strings"

	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/services/config"
	"github.com/zekroTJA/shinpuru/internal/services/storage"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekrotja/rogu/log"
)

func InitStorage(container di.Container) storage.Storage {
	var st storage.Storage
	var err error

	cfg := container.Get(static.DiConfig).(config.Provider)

	log := log.Tagged("Storage")

	drv := strings.ToLower(cfg.Config().Storage.Type)
	log.Info().Field("driver", drv).Msg("Initializing storage ...")

	switch drv {
	case "minio", "s3", "googlecloud":
		st = new(storage.Minio)
	case "file":
		st = new(storage.File)
	default:
		log.Fatal().Field("driver", drv).Msg("Invalid or unsupported storage driver")
	}

	if err = st.Connect(cfg); err != nil {
		log.Fatal().Err(err).Msg("Failed connecting to storage device")
	}

	log.Info().Msg("Connected to storage device")

	return st
}
