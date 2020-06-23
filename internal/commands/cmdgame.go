package commands

import (
	"strings"
	"time"

	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/presence"
	"github.com/zekroTJA/shinpuru/internal/util/static"
)

type CmdGame struct {
}

func (c *CmdGame) GetInvokes() []string {
	return []string{"game", "presence", "botmsg"}
}

func (c *CmdGame) GetDescription() string {
	return "set the presence of the bot"
}

func (c *CmdGame) GetHelp() string {
	return "`game msg <displayMessage>` - set the presence game text\n" +
		"`game status <online|dnd|idle>` - set the status"
}

func (c *CmdGame) GetGroup() string {
	return GroupGlobalAdmin
}

func (c *CmdGame) GetDomainName() string {
	return "sp.game"
}

func (c *CmdGame) GetSubPermissionRules() []SubPermission {
	return nil
}

func (c *CmdGame) Exec(args *CommandArgs) error {

	if len(args.Args) < 2 {
		return util.SendEmbedError(args.Session, args.Channel.ID,
			"Use the sub command `msg` to change the game text and `status` to update the status.").
			DeleteAfter(8 * time.Second).Error()
	}

	rawPresence, err := args.CmdHandler.db.GetSetting(static.SettingPresence)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	defPresence := &presence.Presence{
		Game:   args.CmdHandler.config.Discord.GeneralPrefix + "help | zekro.de",
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

	switch strings.ToLower(args.Args[0]) {

	case "msg":
		pre.Game = strings.Join(args.Args[1:], " ")

	case "status":
		pre.Status = strings.ToLower(args.Args[1])

	default:
		return util.SendEmbedError(args.Session, args.Channel.ID,
			"Use the sub command `msg` to change the game text and `status` to update the status.").
			DeleteAfter(8 * time.Second).Error()
	}

	if err = pre.Validate(); err != nil {
		return util.SendEmbedError(args.Session, args.Channel.ID, err.Error()).
			DeleteAfter(8 * time.Second).Error()
	}

	err = args.Session.UpdateStatusComplex(pre.ToUpdateStatusData())
	if err != nil {
		return err
	}

	err = args.CmdHandler.db.SetSetting(static.SettingPresence, pre.Marshal())
	if err != nil {
		return err
	}

	return util.SendEmbed(args.Session, args.Channel.ID,
		"Presence updated.", "", static.ColorEmbedUpdated).
		DeleteAfter(5 * time.Second).Error()
}
