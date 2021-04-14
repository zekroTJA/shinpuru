package inits

import (
	"strings"

	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/core/config"
	"github.com/zekroTJA/shinpuru/internal/core/storage"
	"github.com/zekroTJA/shinpuru/internal/util"
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
		util.Log.Fatal("Failed connecting to storage device:", err)
	}

	util.Log.Info("Connected to storage device")

	return st
}
