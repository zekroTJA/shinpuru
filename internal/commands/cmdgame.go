package commands

import (
	"strings"
	"time"

	"github.com/zekroTJA/shinpuru/internal/config"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/presence"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shireikan"
)

type CmdGame struct {
}

func (c *CmdGame) GetInvokes() []string {
	return []string{"game", "presence", "botmsg"}
}

func (c *CmdGame) GetDescription() string {
	return "Set the presence of the bot."
}

func (c *CmdGame) GetHelp() string {
	return "`game msg <displayMessage>` - set the presence game text\n" +
		"`game status <online|dnd|idle>` - set the status"
}

func (c *CmdGame) GetGroup() string {
	return shireikan.GroupGlobalAdmin
}

func (c *CmdGame) GetDomainName() string {
	return "sp.game"
}

func (c *CmdGame) GetSubPermissionRules() []shireikan.SubPermission {
	return nil
}

func (c *CmdGame) IsExecutableInDMChannels() bool {
	return true
}

func (c *CmdGame) Exec(ctx shireikan.Context) error {

	if len(ctx.GetArgs()) < 2 {
		return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
			"Use the sub command `msg` to change the game text and `status` to update the status.").
			DeleteAfter(8 * time.Second).Error()
	}

	db, _ := ctx.GetObject(static.DiDatabase).(database.Database)
	rawPresence, err := db.GetSetting(static.SettingPresence)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	cfg, _ := ctx.GetObject(static.DiConfig).(*config.Config)
	defPresence := &presence.Presence{
		Game:   cfg.Discord.GeneralPrefix + "help | zekro.de",
		Status: "online",
	}

	var pre *presence.Presence
	if rawPresence == "" {
		pre = defPresence
	} else {
		pre, err = presence.Unmarshal(rawPresence)
		if err != nil {
			pre = defPresence
		}
	}

	switch strings.ToLower(ctx.GetArgs().Get(0).AsString()) {

	case "msg":
		pre.Game = strings.Join(ctx.GetArgs()[1:], " ")

	case "status":
		pre.Status = strings.ToLower(ctx.GetArgs().Get(1).AsString())

	default:
		return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
			"Use the sub command `msg` to change the game text and `status` to update the status.").
			DeleteAfter(8 * time.Second).Error()
	}

	if err = pre.Validate(); err != nil {
		return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID, err.Error()).
			DeleteAfter(8 * time.Second).Error()
	}

	err = ctx.GetSession().UpdateStatusComplex(pre.ToUpdateStatusData())
	if err != nil {
		return err
	}

	err = db.SetSetting(static.SettingPresence, pre.Marshal())
	if err != nil {
		return err
	}

	return util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
		"Presence updated.", "", static.ColorEmbedUpdated).
		DeleteAfter(5 * time.Second).Error()
}
