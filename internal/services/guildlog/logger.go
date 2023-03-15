package guildlog

import (
	"fmt"

	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/models"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/timeprovider"
	"github.com/zekroTJA/shinpuru/internal/util/snowflakenodes"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekrotja/rogu"
	"github.com/zekrotja/rogu/log"
)

type loggerImpl struct {
	db     database.Database
	module string
	tp     timeprovider.Provider
	l      rogu.Logger
}

func New(container di.Container) Logger {
	return &loggerImpl{
		db: container.Get(static.DiDatabase).(database.Database),
		tp: container.Get(static.DiTimeProvider).(timeprovider.Provider),
		l:  log.Tagged("GuildLog"),
	}
}

func (l *loggerImpl) Section(module string) Logger {
	return &loggerImpl{
		db:     l.db,
		tp:     l.tp,
		module: module,
	}
}

func (l *loggerImpl) log(severity models.GuildLogSeverity, guildID, message string, data ...interface{}) (err error) {
	defer func() {
		if err != nil {
			l.l.Error().Err(err).Fields(
				"gid", guildID,
				"message", message,
				"severity", severity,
			).Msg("Failed creating guildlog entry")
		}
	}()

	message = fmt.Sprintf(message, data...)

	module := l.module
	if module == "" {
		module = "global"
	}

	ok, err := l.db.GetGuildLogDisable(guildID)
	if ok || err != nil && !database.IsErrDatabaseNotFound(err) {
		return
	}

	err = l.db.AddGuildLogEntry(models.GuildLogEntry{
		ID:        snowflakenodes.NodeGuildLog.Generate(),
		GuildID:   guildID,
		Module:    module,
		Message:   message,
		Severity:  severity,
		Timestamp: l.tp.Now(),
	})

	return
}

func (l *loggerImpl) Debugf(guildID, message string, data ...interface{}) error {
	return l.log(models.GLDebug, guildID, message, data...)
}

func (l *loggerImpl) Infof(guildID, message string, data ...interface{}) error {
	return l.log(models.GLInfo, guildID, message, data...)
}

func (l *loggerImpl) Warnf(guildID, message string, data ...interface{}) error {
	return l.log(models.GLWarn, guildID, message, data...)
}

func (l *loggerImpl) Errorf(guildID, message string, data ...interface{}) error {
	return l.log(models.GLError, guildID, message, data...)
}

func (l *loggerImpl) Fatalf(guildID, message string, data ...interface{}) error {
	return l.log(models.GLFatal, guildID, message, data...)
}
