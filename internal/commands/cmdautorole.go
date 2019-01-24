package commands

import (
	"fmt"
	"time"

	"github.com/zekroTJA/shinpuru/internal/core"
	"github.com/zekroTJA/shinpuru/internal/util"
)

type CmdAutorole struct {
	PermLvl int
}

func (c *CmdAutorole) GetInvokes() []string {
	return []string{"autorole", "arole"}
}

func (c *CmdAutorole) GetDescription() string {
	return "set the autorole for the current guild"
}

func (c *CmdAutorole) GetHelp() string {
	return "`autorole` - display currently set autorole\n" +
		"`autorole <roleResolvable>` - set an auto role for the current guild"
}

func (c *CmdAutorole) GetGroup() string {
	return GroupGuildConfig
}

func (c *CmdAutorole) GetPermission() int {
	return c.PermLvl
}

func (c *CmdAutorole) SetPermission(permLvl int) {
	c.PermLvl = permLvl
}

func (c *CmdAutorole) Exec(args *CommandArgs) error {
	if len(args.Args) < 1 {
		currAutoRoleID, err := args.CmdHandler.db.GetGuildAutoRole(args.Guild.ID)
		if err != nil && !core.IsErrDatabaseNotFound(err) {
			return err
		}
		if currAutoRoleID == "" {
			_, err := util.SendEmbed(args.Session, args.Channel.ID,
				"There is no autorole set on this guild currently.", "", 0)
			return err
		}
		_, err = util.FetchRole(args.Session, args.Guild.ID, currAutoRoleID)
		if err != nil {
			_, err := util.SendEmbedError(args.Session, args.Channel.ID,
				"**ATTENTION:** The set auto role is no more existent on the guild!")
			return err
		}
		_, err = util.SendEmbed(args.Session, args.Channel.ID,
			fmt.Sprintf("Currently, <@&%s> is set as auto role.", currAutoRoleID), "", 0)
		return err
	}

	newAutoRole, err := util.FetchRole(args.Session, args.Guild.ID, args.Args[0])
	if err != nil {
		msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
			"Sorry, but the entered role could not be fetched :(")
		util.DeleteMessageLater(args.Session, msg, 5*time.Second)
		return err
	}
	err = args.CmdHandler.db.SetGuildAutoRole(args.Guild.ID, newAutoRole.ID)
	if err != nil {
		return err
	}
	_, err = util.SendEmbed(args.Session, args.Channel.ID,
		fmt.Sprintf("Autorole set to <@&%s>.", newAutoRole.ID), "", util.ColorEmbedUpdated)
	return err
}
