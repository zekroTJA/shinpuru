package guildlog

import (
	"fmt"
	"time"

	"github.com/sarulabs/di/v2"
	"github.com/sirupsen/logrus"
	"github.com/zekroTJA/shinpuru/internal/models"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/util/snowflakenodes"
	"github.com/zekroTJA/shinpuru/internal/util/static"
)

type loggerImpl struct {
	db     database.Database
	module string
}

func New(container di.Container) Logger {
	return &loggerImpl{
		db: container.Get(static.DiDatabase).(database.Database),
	}
}

func (l *loggerImpl) Section(module string) Logger {
	return &loggerImpl{
		db:     l.db,
		module: module,
	}
}

func (l *loggerImpl) log(severity models.GuildLogSeverity, guildID, message string, data ...interface{}) (err error) {
	defer func() {
		if err != nil {
			logrus.WithError(err).WithFields(logrus.Fields{
				"gid":      guildID,
				"message":  message,
				"severity": severity,
			}).Error("failed logging guildlog")
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
		Timestamp: time.Now(),
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
