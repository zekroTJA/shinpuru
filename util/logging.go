package util

import (
	logging "github.com/op/go-logging"
)

var Log = GetLogger()

func GetLogger() *logging.Logger {
	logger := logging.MustGetLogger("main")
	format := logging.MustStringFormatter(`%{color}▶  %{level:.4s} %{id:05d}%{color:reset} %{message}`)
	logging.SetFormatter(format)
	logging.SetLevel(logging.DEBUG, "main")
	return logger
}
