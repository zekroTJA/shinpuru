package webserver

import (
	"compress/gzip"
	"fmt"
	"io"

	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/util/static"
)

// ---------------------------------------------------------------------------
// - GET /api/guilds/:guildid/backups

func (ws *WebServer) handlerGetGuildBackups(ctx *routing.Context) error {
	guildID := ctx.Param("guildid")

	backupEntries, err := ws.db.GetBackups(guildID)
	if database.IsErrDatabaseNotFound(err) {
		return jsonError(ctx, nil, fasthttp.StatusNotFound)
	} else if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	return jsonResponse(ctx, ListResponse{N: len(backupEntries), Data: backupEntries}, fasthttp.StatusOK)
}

// ---------------------------------------------------------------------------
// - GET /api/guilds/:guildid/backups/:backupid/download

func (ws *WebServer) handlerGetGuildBackupDownload(ctx *routing.Context) error {
	guildID := ctx.Param("guildid")
	backupID := ctx.Param("backupid")

	backupEntries, err := ws.db.GetBackups(guildID)
	if database.IsErrDatabaseNotFound(err) {
		return jsonError(ctx, nil, fasthttp.StatusNotFound)
	} else if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	var found bool
backupEntriesLoop:
	for _, e := range backupEntries {
		if e.FileID == backupID {
			found = true
			break backupEntriesLoop
		}
	}

	if !found {
		return jsonError(ctx, nil, fasthttp.StatusNotFound)
	}

	f, size, err := ws.st.GetObject(static.StorageBucketBackups, backupID)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}
	defer f.Close()

	zf := gzip.NewWriter(ctx.Response.BodyWriter())
	zf.Name = fmt.Sprintf("backup_%s_%s.json", guildID, backupID)

	_, err = io.CopyN(zf, f, size)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}
	defer zf.Close()

	// 24 hours browser caching
	ctx.Response.Header.Set("Cache-Control", "public, max-age=86400‬, immutable")
	ctx.Response.Header.SetContentType("application/gzip")

	return nil
}

// ---------------------------------------------------------------------------
// - POST /api/guilds/:guildid/backups/toggle

func (ws *WebServer) handlerPostGuildBackupsToggle(ctx *routing.Context) error {
	userID := ctx.Get("uid").(string)

	guildID := ctx.Param("guildid")

	if ok, _, err := ws.pmw.CheckPermissions(ws.session, guildID, userID, "sp.guild.admin.backup"); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	} else if !ok {
		return jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized)
	}

	var data struct {
		Enabled bool `json:"enabled"`
	}

	if err := parseJSONBody(ctx, &data); err != nil {
		return jsonError(ctx, err, fasthttp.StatusBadRequest)
	}

	if err := ws.db.SetGuildBackup(guildID, data.Enabled); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	return jsonResponse(ctx, nil, fasthttp.StatusOK)
}
