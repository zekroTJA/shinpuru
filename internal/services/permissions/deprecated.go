package permissions

import (
	"time"

	"github.com/zekroTJA/shinpuru/internal/services/config"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
	"github.com/zekroTJA/shireikan"
)

func (m *Permissions) Handle(
	cmd shireikan.Command, ctx shireikan.Context, layer shireikan.MiddlewareLayer) (next bool, err error) {

	if m.db == nil {
		m.db, _ = ctx.GetObject(static.DiDatabase).(database.Database)
	}

	if m.cfg == nil {
		m.cfg, _ = ctx.GetObject(static.DiConfig).(config.Provider)
	}

	var guildID string
	if ctx.GetGuild() != nil {
		guildID = ctx.GetGuild().ID
	}

	ok, _, err := m.CheckPermissions(ctx.GetSession(), guildID, ctx.GetUser().ID, cmd.GetDomainName())

	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return false, err
	}

	if !ok {
		msg, _ := ctx.ReplyEmbedError("You are not permitted to use this command!", "Missing Permission")
		discordutil.DeleteMessageLater(ctx.GetSession(), msg, 8*time.Second)
		return false, nil
	}

	return true, nil
}

func (m *Permissions) GetLayer() shireikan.MiddlewareLayer {
	return shireikan.LayerBeforeCommand
}
