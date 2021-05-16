package guildlog

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/util/static"
)

type loggerImpl struct {
	s      *discordgo.Session
	db     database.Database
	module string
}

func NewLogger(container di.Container) *loggerImpl {
	return &loggerImpl{
		s:  container.Get(static.DiDiscordSession).(*discordgo.Session),
		db: container.Get(static.DiDatabase).(database.Database),
	}
}
