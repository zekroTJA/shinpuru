package inits

import (
	"strings"

	"github.com/sarulabs/di/v2"
	"github.com/sirupsen/logrus"
	"github.com/zekroTJA/shinpuru/internal/config"
	"github.com/zekroTJA/shinpuru/internal/services/storage"
	"github.com/zekroTJA/shinpuru/internal/util/static"
)

func InitStorage(container di.Container) storage.Storage {
	var st storage.Storage
	var err error

	cfg := container.Get(static.DiConfig).(*config.Config)

	switch strings.ToLower(cfg.Storage.Type) {
	case "minio", "s3", "googlecloud":
		st = new(storage.Minio)
	case "file":
		st = new(storage.File)
	}

	if err = st.Connect(cfg); err != nil {
		logrus.WithError(err).Fatal("Failed connecting to storage device")
	}

	logrus.Info("Connected to storage device")

	return st
}
