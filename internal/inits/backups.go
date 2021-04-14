package inits

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/core/backup"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/core/storage"
	"github.com/zekroTJA/shinpuru/internal/util/static"
)

func InitBackupHandler(container di.Container) *backup.GuildBackups {
	session := container.Get(static.DiDiscordSession).(*discordgo.Session)
	db := container.Get(static.DiDatabase).(database.Database)
	storage := container.Get(static.DiObjectStorage).(storage.Storage)

	return backup.New(session, db, storage)
}
