package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
)

type CmdAutorole struct {
}

func (c *CmdAutorole) GetInvokes() []string {
	return []string{"autorole", "arole"}
}

func (c *CmdAutorole) GetDescription() string {
	return "set the autorole for the current guild"
}

func (c *CmdAutorole) GetHelp() string {
	return "`autorole` - display currently set autorole\n" +
		"`autorole <roleResolvable>` - set an auto role for the current guild\n" +
		"`autorole reset` - disable autorole"
}

func (c *CmdAutorole) GetGroup() string {
	return GroupGuildConfig
}

func (c *CmdAutorole) GetDomainName() string {
	return "sp.guild.config.autorole"
}

func (c *CmdAutorole) GetSubPermissionRules() []SubPermission {
	return nil
}

func (c *CmdAutorole) Exec(args *CommandArgs) error {
	if len(args.Args) < 1 {
		currAutoRoleID, err := args.CmdHandler.db.GetGuildAutoRole(args.Guild.ID)
		if err != nil && !database.IsErrDatabaseNotFound(err) {
			return err
		}
		if currAutoRoleID == "" {
			return util.SendEmbed(args.Session, args.Channel.ID,
				"There is no autorole set on this guild currently.", "", 0).Error()
		}
		_, err = util.FetchRole(args.Session, args.Guild.ID, currAutoRoleID)
		if err != nil {
			return util.SendEmbedError(args.Session, args.Channel.ID,
				"**ATTENTION:** The set auto role is no more existent on the guild!").Error()
		}
		return util.SendEmbed(args.Session, args.Channel.ID,
			fmt.Sprintf("Currently, <@&%s> is set as auto role.", currAutoRoleID), "", 0).Error()
	}

	if strings.ToLower(args.Args[0]) == "reset" {
		err := args.CmdHandler.db.SetGuildAutoRole(args.Guild.ID, "")
		if err != nil {
			return err
		}
		return util.SendEmbed(args.Session, args.Channel.ID,
			"Autorole reseted.", "", static.ColorEmbedUpdated).Error()
	}

	newAutoRole, err := util.FetchRole(args.Session, args.Guild.ID, args.Args[0])
	if err != nil {
		return util.SendEmbedError(args.Session, args.Channel.ID,
			"Sorry, but the entered role could not be fetched :(").
			DeleteAfter(5 * time.Second).Error()
	}
	err = args.CmdHandler.db.SetGuildAutoRole(args.Guild.ID, newAutoRole.ID)
	if err != nil {
		return err
	}
	return util.SendEmbed(args.Session, args.Channel.ID,
		fmt.Sprintf("Autorole set to <@&%s>.", newAutoRole.ID), "", static.ColorEmbedUpdated).
		Error()
}
