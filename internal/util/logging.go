package util

import (
	logging "github.com/op/go-logging"
)

var Log = GetLogger()

func GetLogger() *logging.Logger {
	logger := logging.MustGetLogger("main")
	format := logging.MustStringFormatter(`%{color}▶  %{level:.4s} %{id:05d}%{color:reset} %{message}`)
	logging.SetFormatter(format)
	logging.SetLevel(logging.INFO, "main")
	return logger
}

func SetLogLevel(level int) {
	logging.SetLevel(logging.Level(level), "main")
}
