package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/fetch"
	"github.com/zekroTJA/shireikan"
)

type CmdAutorole struct {
}

func (c *CmdAutorole) GetInvokes() []string {
	return []string{"autorole", "arole"}
}

func (c *CmdAutorole) GetDescription() string {
	return "Set the autorole for the current guild."
}

func (c *CmdAutorole) GetHelp() string {
	return "`autorole` - display currently set autorole\n" +
		"`autorole <roleResolvable>` - set an auto role for the current guild\n" +
		"`autorole reset` - disable autorole"
}

func (c *CmdAutorole) GetGroup() string {
	return shireikan.GroupGuildConfig
}

func (c *CmdAutorole) GetDomainName() string {
	return "sp.guild.config.autorole"
}

func (c *CmdAutorole) GetSubPermissionRules() []shireikan.SubPermission {
	return nil
}

func (c *CmdAutorole) IsExecutableInDMChannels() bool {
	return false
}

func (c *CmdAutorole) Exec(ctx shireikan.Context) error {
	db, _ := ctx.GetObject("db").(database.Database)

	if len(ctx.GetArgs()) < 1 {
		currAutoRoleID, err := db.GetGuildAutoRole(ctx.GetGuild().ID)
		if err != nil && !database.IsErrDatabaseNotFound(err) {
			return err
		}
		if currAutoRoleID == "" {
			return util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
				"There is no autorole set on this guild currently.", "", 0).Error()
		}
		_, err = fetch.FetchRole(ctx.GetSession(), ctx.GetGuild().ID, currAutoRoleID)
		if err != nil {
			return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
				"**ATTENTION:** The set auto role is no more existent on the guild!").Error()
		}
		return util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
			fmt.Sprintf("Currently, <@&%s> is set as auto role.", currAutoRoleID), "", 0).Error()
	}

	if strings.ToLower(ctx.GetArgs().Get(0).AsString()) == "reset" {
		err := db.SetGuildAutoRole(ctx.GetGuild().ID, "")
		if err != nil {
			return err
		}
		return util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
			"Autorole reseted.", "", static.ColorEmbedUpdated).Error()
	}

	newAutoRole, err := fetch.FetchRole(ctx.GetSession(), ctx.GetGuild().ID, ctx.GetArgs().Get(0).AsString())
	if err != nil {
		return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
			"Sorry, but the entered role could not be fetched :(").
			DeleteAfter(5 * time.Second).Error()
	}
	err = db.SetGuildAutoRole(ctx.GetGuild().ID, newAutoRole.ID)
	if err != nil {
		return err
	}
	return util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
		fmt.Sprintf("Autorole set to <@&%s>.", newAutoRole.ID), "", static.ColorEmbedUpdated).
		Error()
}
