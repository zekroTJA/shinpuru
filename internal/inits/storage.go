package inits

import (
	"strings"

	"github.com/zekroTJA/shinpuru/internal/core/config"
	"github.com/zekroTJA/shinpuru/internal/core/storage"
	"github.com/zekroTJA/shinpuru/internal/util"
)

func InitStorage(cfg *config.Config) storage.Storage {
	var st storage.Storage
	var err error

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
