package controllers

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/middleware"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/storage"
	"github.com/zekroTJA/shinpuru/internal/services/webserver/v1/models"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/onetimeauth/v2"
	"github.com/zekrotja/discordgo"
)

type GuildBackupsController struct {
	db  database.Database
	st  storage.Storage
	ota onetimeauth.OneTimeAuth
}

func (c *GuildBackupsController) Setup(container di.Container, router fiber.Router) {
	c.db = container.Get(static.DiDatabase).(database.Database)
	c.st = container.Get(static.DiObjectStorage).(storage.Storage)
	c.ota = container.Get(static.DiOneTimeAuth).(onetimeauth.OneTimeAuth)

	session := container.Get(static.DiDiscordSession).(*discordgo.Session)
	pmw := container.Get(static.DiPermissionMiddleware).(*middleware.PermissionsMiddleware)

	router.Get("", c.getBackups)
	router.Post("/toggle", pmw.HandleWs(session, "sp.guild.admin.backup"), c.postToggleBackups)
	router.Post("/:backupid/download", pmw.HandleWs(session, "sp.guild.admin.backup"), c.postDownloadBackup)
	router.Get("/:backupid/download", c.getDownloadBackup)
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

func (c *GuildBackupsController) postDownloadBackup(ctx *fiber.Ctx) error {
	guildID := ctx.Params("guildid")
	backupID := ctx.Params("backupid")

	ident := getBackupIdent(guildID, backupID)

	token, expires, err := c.ota.GetKey(ident)
	if err != nil {
		return err
	}

	return ctx.JSON(&models.AccessTokenResponse{
		Token:   token,
		Expires: expires,
	})
}

func (c *GuildBackupsController) getDownloadBackup(ctx *fiber.Ctx) error {
	guildID := ctx.Params("guildid")
	backupID := ctx.Params("backupid")

	ident, _ := ctx.Locals("uid").(string)
	if rGuildID, rBackupID := decodeBackupIdent(ident); rGuildID != guildID || rBackupID != backupID {
		return fiber.ErrForbidden
	}

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
	zf.Close()

	// 24 hours browser caching
	ctx.Set("Cache-Control", "public, max-age=86400‬, immutable")
	ctx.Set("Content-Type", "application/gzip")
	ctx.Set("Content-Disposition", fmt.Sprintf(`filename="backup_%s_%s.gz"`, guildID, backupID))
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

	return ctx.JSON(models.Ok)
}

// --- HELPERS ---

func getBackupIdent(guildID, backupID string) string {
	return fmt.Sprintf("%s#%s", guildID, backupID)
}

func decodeBackupIdent(ident string) (guildID, backupID string) {
	split := strings.Split(ident, "#")
	if len(split) == 2 {
		guildID = split[0]
		backupID = split[1]
	}
	return
}
