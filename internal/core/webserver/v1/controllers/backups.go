package controllers

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"

	"github.com/bwmarrin/discordgo"
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/core/middleware"
	"github.com/zekroTJA/shinpuru/internal/core/storage"
	"github.com/zekroTJA/shinpuru/internal/core/webserver/v1/models"
	"github.com/zekroTJA/shinpuru/internal/util/static"
)

type GuildBackupsController struct {
	db database.Database
	st storage.Storage
}

func (c *GuildBackupsController) Setup(container di.Container, router fiber.Router) {
	c.db = container.Get(static.DiDatabase).(database.Database)
	c.st = container.Get(static.DiObjectStorage).(storage.Storage)

	session := container.Get(static.DiDiscordSession).(*discordgo.Session)
	pmw := container.Get(static.DiPermissionMiddleware).(*middleware.PermissionsMiddleware)

	router.Get("", c.getBackups)
	router.Get("/toggle", pmw.HandleWs(session, "sp.guild.admin.backup"), c.postToggleBackups)
	router.Get("/:backupid/download", pmw.HandleWs(session, "sp.guild.admin.backup"), c.downloadBackup)
}

func (c *GuildBackupsController) getBackups(ctx *fiber.Ctx) error {
	guildID := ctx.Params("guildid")

	backupEntries, err := c.db.GetBackups(guildID)
	if database.IsErrDatabaseNotFound(err) {
		return fiber.ErrNotFound
	} else if err != nil {
		return err
	}

	return ctx.JSON(&models.ListResponse{N: len(backupEntries), Data: backupEntries})
}

func (c *GuildBackupsController) downloadBackup(ctx *fiber.Ctx) error {
	guildID := ctx.Params("guildid")
	backupID := ctx.Params("backupid")

	backupEntries, err := c.db.GetBackups(guildID)
	if database.IsErrDatabaseNotFound(err) {
		return fiber.ErrNotFound
	} else if err != nil {
		return err
	}

	var found bool
	for _, e := range backupEntries {
		if e.FileID == backupID {
			found = true
			break
		}
	}

	if !found {
		return fiber.ErrNotFound
	}

	f, size, err := c.st.GetObject(static.StorageBucketBackups, backupID)
	if err != nil {
		return err
	}
	defer f.Close()

	buff := bytes.NewBuffer([]byte{})
	zf := gzip.NewWriter(buff)
	zf.Name = fmt.Sprintf("backup_%s_%s.json", guildID, backupID)

	_, err = io.CopyN(zf, f, size)
	if err != nil {
		return err
	}
	defer zf.Close()

	// 24 hours browser caching
	ctx.Set("Cache-Control", "public, max-age=86400‬, immutable")
	ctx.Set("Content-Type", "application/gzip")
	return ctx.SendStream(buff)
}

func (c *GuildBackupsController) postToggleBackups(ctx *fiber.Ctx) error {
	guildID := ctx.Params("guildid")

	var data struct {
		Enabled bool `json:"enabled"`
	}

	if err := ctx.BodyParser(&data); err != nil {
		return err
	}

	if err := c.db.SetGuildBackup(guildID, data.Enabled); err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusOK)
}
