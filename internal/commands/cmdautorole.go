package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/acceptmsg"
	"github.com/zekroTJA/shinpuru/pkg/fetch"
	"github.com/zekroTJA/shinpuru/pkg/stringutil"
	"github.com/zekroTJA/shireikan"
	"github.com/zekrotja/dgrs"
)

type CmdAutorole struct {
}

func (c *CmdAutorole) GetInvokes() []string {
	return []string{"autorole", "autoroles", "arole", "aroles"}
}

func (c *CmdAutorole) GetDescription() string {
	return "Set the autorole for the current guild."
}

func (c *CmdAutorole) GetHelp() string {
	return "`autorole` - display currently set autorole(s)\n" +
		"`autorole <roleResolvable> (<roleResolvable> (<roleResolvable> (...)))` - set an auto role for the current guild\n" +
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
	db, _ := ctx.GetObject(static.DiDatabase).(database.Database)
	st, _ := ctx.GetObject(static.DiState).(*dgrs.State)

	if len(ctx.GetArgs()) == 0 {
		return c.list(ctx, db, st)
	}
	if ctx.GetArgs().Get(0).AsString() == "reset" {
		return c.reset(ctx, db)
	}
	return c.set(ctx, db, st)
}

func (c *CmdAutorole) list(ctx shireikan.Context, db database.Database, st *dgrs.State) (err error) {
	autoRoleIDs, err := db.GetGuildAutoRole(ctx.GetGuild().ID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return
	}

	if len(autoRoleIDs) == 0 {
		return util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
			"No autoroles are set.", "", 0).DeleteAfter(10 * time.Second).Error()
	}

	guildRoles, err := st.Roles(ctx.GetGuild().ID, true)
	if err != nil {
		return err
	}
	guildRoleIDs := make([]string, len(guildRoles))
	for i, role := range guildRoles {
		guildRoleIDs[i] = role.ID
	}

	if nc := stringutil.NotContained(autoRoleIDs, guildRoleIDs); len(nc) > 0 {
		autoRoleIDs = stringutil.Contained(autoRoleIDs, guildRoleIDs)
		am, err := acceptmsg.New().
			WithSession(ctx.GetSession()).
			DeleteAfterAnswer().
			LockOnUser(ctx.GetUser().ID).
			WithContent(fmt.Sprintf(
				"%d %s are not existent anymore. "+
					"Do you want to remove them now from the list of autoroles?",
				len(nc), util.Pluralize(len(nc), "autorole"))).
			DoOnAccept(func(_ *discordgo.Message) error {
				return db.SetGuildAutoRole(ctx.GetGuild().ID, autoRoleIDs)
			}).
			Send(ctx.GetChannel().ID)
		if err != nil {
			return err
		}
		if err = am.Error(); err != nil {
			return err
		}
	}

	var roleNames strings.Builder
	roleNames.WriteString("Following autorole(s) are set:\n")
	for _, rid := range autoRoleIDs {
		roleNames.WriteString(fmt.Sprintf(" - <@&%s> (%s)\n", rid, rid))
	}

	return util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
		roleNames.String(), "", 0).DeleteAfter(12 * time.Second).Error()
}

func (c *CmdAutorole) set(ctx shireikan.Context, db database.Database, st *dgrs.State) (err error) {
	autoRoleIDs := make([]string, 0, len(ctx.GetArgs()))

	if stringutil.ContainsAny("@everyone", ctx.GetArgs()) {
		return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
			"`@everyone` can not be set as autorole.", "").DeleteAfter(10 * time.Second).Error()
	}

	for _, arg := range ctx.GetArgs() {
		if len(arg) == 0 {
			continue
		}
		role, err := fetch.FetchRole(fetch.WrapDrgs(st), ctx.GetGuild().ID, arg)
		if err != nil {
			return err
		}
		autoRoleIDs = append(autoRoleIDs, role.ID)
	}

	if err = db.SetGuildAutoRole(ctx.GetGuild().ID, autoRoleIDs); err != nil {
		return err
	}

	return util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
		"Autoroles set.", "", 0).DeleteAfter(5 * time.Second).Error()
}

func (c *CmdAutorole) reset(ctx shireikan.Context, db database.Database) (err error) {
	if err = db.SetGuildAutoRole(ctx.GetGuild().ID, []string{}); err != nil {
		return err
	}

	return util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
		"Autoroles reseted.", "", 0).DeleteAfter(5 * time.Second).Error()
}
