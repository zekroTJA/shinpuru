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

type CmdAutovc struct {
}

func (c *CmdAutovc) GetInvokes() []string {
	return []string{"autovoice", "autovc", "avc", "avcs"}
}

func (c *CmdAutovc) GetDescription() string {
	return "Set the auto voicechannel for the current guild."
}

func (c *CmdAutovc) GetHelp() string {
	return "`autovc` - display currently set auto voicechannel(s)\n" +
		"`autovc <channelResolvable> (<channelResolvable> (<channelResolvable> (...)))` - set an auto voicechannel for the current guild\n" +
		"`autovc reset` - disable auto voicechannels"
}

func (c *CmdAutovc) GetGroup() string {
	return shireikan.GroupGuildConfig
}

func (c *CmdAutovc) GetDomainName() string {
	return "sp.guild.config.autovc"
}

func (c *CmdAutovc) GetSubPermissionRules() []shireikan.SubPermission {
	return nil
}

func (c *CmdAutovc) IsExecutableInDMChannels() bool {
	return false
}

func (c *CmdAutovc) Exec(ctx shireikan.Context) error {
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

func (c *CmdAutovc) list(ctx shireikan.Context, db database.Database, st *dgrs.State) (err error) {
	autoVCIDs, err := db.GetGuildAutoVC(ctx.GetGuild().ID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return
	}

	if len(autoVCIDs) == 0 {
		return util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
			"No auto voicechannels are set.", "", 0).DeleteAfter(10 * time.Second).Error()
	}

	guildChannels, err := st.Channels(ctx.GetGuild().ID, true)
	if err != nil {
		return err
	}
	guildChannelIDs := make([]string, len(guildChannels))
	for i, channel := range guildChannels {
		guildChannelIDs[i] = channel.ID
	}

	if nc := stringutil.NotContained(autoVCIDs, guildChannelIDs); len(nc) > 0 {
		autoVCIDs = stringutil.Contained(autoVCIDs, guildChannelIDs)
		am, err := acceptmsg.New().
			WithSession(ctx.GetSession()).
			DeleteAfterAnswer().
			LockOnUser(ctx.GetUser().ID).
			WithContent(fmt.Sprintf(
				"%d %s are not existent anymore. "+
					"Do you want to remove them now from the list of auto voicechannels?",
				len(nc), util.Pluralize(len(nc), "autovoicechannel"))).
			DoOnAccept(func(_ *discordgo.Message) error {
				return db.SetGuildAutoVC(ctx.GetGuild().ID, autoVCIDs)
			}).
			Send(ctx.GetChannel().ID)
		if err != nil {
			return err
		}
		if err = am.Error(); err != nil {
			return err
		}
	}

	var vcNames strings.Builder
	vcNames.WriteString("Following auto voicechannel(s) are set:\n")
	for _, vcid := range autoVCIDs {
		vcNames.WriteString(fmt.Sprintf(" - <@&%s> (%s)\n", vcid, vcid))
	}

	return util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
		vcNames.String(), "", 0).DeleteAfter(12 * time.Second).Error()
}

func (c *CmdAutovc) set(ctx shireikan.Context, db database.Database, st *dgrs.State) (err error) {
	autoVCIDs := make([]string, 0, len(ctx.GetArgs()))

	for _, arg := range ctx.GetArgs() {
		if len(arg) == 0 {
			continue
		}
		vc, err := fetch.FetchChannel(fetch.WrapDrgs(st), ctx.GetGuild().ID, arg)
		if err != nil {
			return err
		}
		autoVCIDs = append(autoVCIDs, vc.ID)
	}

	if err = db.SetGuildAutoVC(ctx.GetGuild().ID, autoVCIDs); err != nil {
		return err
	}

	return util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
		"Auto voicechannels set.", "", 0).DeleteAfter(5 * time.Second).Error()
}

func (c *CmdAutovc) reset(ctx shireikan.Context, db database.Database) (err error) {
	if err = db.SetGuildAutoVC(ctx.GetGuild().ID, []string{}); err != nil {
		return err
	}

	return util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
		"Auto voicechannels reseted.", "", 0).DeleteAfter(5 * time.Second).Error()
}
