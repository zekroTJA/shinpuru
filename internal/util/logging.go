package util

import (
	logging "github.com/op/go-logging"
)

// Log is the default logger.
var Log = GetLogger()

// GetLogger returns a defautly configured logger.
func GetLogger() *logging.Logger {
	logger := logging.MustGetLogger("main")
	format := logging.MustStringFormatter(`%{color}▶  %{level:.4s} %{id:05d}%{color:reset} %{message}`)
	logging.SetFormatter(format)
	logging.SetLevel(logging.INFO, "main")
	return logger
}

// SetLogLevel sets the log level.
func SetLogLevel(level int) {
	logging.SetLevel(logging.Level(level), "main")
}
