package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/zekrotja/discordgo"

	"github.com/zekroTJA/shinpuru/internal/config"
	"github.com/zekroTJA/shinpuru/internal/middleware"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
	"github.com/zekroTJA/shinpuru/pkg/fetch"
	"github.com/zekroTJA/shinpuru/pkg/stringutil"
	"github.com/zekroTJA/shireikan"
)

type CmdProfile struct {
	PermLvl int
}

func (c *CmdProfile) GetInvokes() []string {
	return []string{"user", "u", "profile"}
}

func (c *CmdProfile) GetDescription() string {
	return "Get information about a user."
}

func (c *CmdProfile) GetHelp() string {
	return "`profile (<userResolvable>)` - get user info"
}

func (c *CmdProfile) GetGroup() string {
	return shireikan.GroupChat
}

func (c *CmdProfile) GetDomainName() string {
	return "sp.chat.profile"
}

func (c *CmdProfile) GetSubPermissionRules() []shireikan.SubPermission {
	return nil
}

func (c *CmdProfile) IsExecutableInDMChannels() bool {
	// TODO: Change to true; Required some
	// modification to the command
	return false
}

func (c *CmdProfile) Exec(ctx shireikan.Context) error {
	member, err := ctx.GetSession().GuildMember(ctx.GetGuild().ID, ctx.GetUser().ID)
	if err != nil {
		return err
	}
	if len(ctx.GetArgs()) > 0 {
		member, err = fetch.FetchMember(ctx.GetSession(), ctx.GetGuild().ID, strings.Join(ctx.GetArgs(), " "))
		if err != nil || member == nil {
			return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
				"Could not fetch any member by the passed resolvable.").
				DeleteAfter(8 * time.Second).Error()
		}
	}

	membRoleIDs := make(map[string]struct{})
	for _, rID := range member.Roles {
		membRoleIDs[rID] = struct{}{}
	}

	maxPos := len(ctx.GetGuild().Roles)
	roleColor := static.ColorEmbedGray
	for _, guildRole := range ctx.GetGuild().Roles {
		if _, ok := membRoleIDs[guildRole.ID]; ok && guildRole.Position < maxPos && guildRole.Color != 0 {
			maxPos = guildRole.Position
			roleColor = guildRole.Color
		}
	}

	joinedTime, err := member.JoinedAt.Parse()
	if err != nil {
		return err
	}
	createdTime, err := discordutil.GetDiscordSnowflakeCreationTime(member.User.ID)
	if err != nil {
		return err
	}

	pmw, _ := ctx.GetObject(static.DiPermissionMiddleware).(*middleware.PermissionsMiddleware)
	perms, _, err := pmw.GetPermissions(ctx.GetSession(), ctx.GetGuild().ID, member.User.ID)
	if err != nil {
		return err
	}

	db, _ := ctx.GetObject(static.DiDatabase).(database.Database)

	guildReps, err := db.GetReportsFiltered(ctx.GetGuild().ID, member.User.ID, -1)
	if err != nil {
		return err
	}
	repsOnRecord := len(guildReps)
	repsOnRecordStr := "This user has a white vest :ok_hand:"
	if repsOnRecord > 0 {
		repsOnRecordStr = fmt.Sprintf("This user has **%d** reports on record on this guild.", repsOnRecord)
	}

	roles := make([]string, len(member.Roles))
	for i, rID := range member.Roles {
		roles[i] = "<@&" + rID + ">"
	}

	karma, err := db.GetKarma(member.User.ID, ctx.GetGuild().ID)
	if !database.IsErrDatabaseNotFound(err) && err != nil {
		return err
	}

	karmaTotal, err := db.GetKarmaSum(member.User.ID)
	if !database.IsErrDatabaseNotFound(err) && err != nil {
		return err
	}

	cfg, _ := ctx.GetObject(static.DiConfig).(*config.Config)

	embed := &discordgo.MessageEmbed{
		Color: roleColor,
		Title: fmt.Sprintf("Info about member %s#%s", member.User.Username, member.User.Discriminator),
		Description: fmt.Sprintf("[**Here**](%s/guilds/%s/%s) you can find this users profile in the web interface.",
			cfg.WebServer.PublicAddr, ctx.GetGuild().ID, member.User.ID),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: member.User.AvatarURL(""),
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Inline: true,
				Name:   "Tag",
				Value:  member.User.Username + "#" + member.User.Discriminator,
			},
			{
				Inline: true,
				Name:   "Nickname",
				Value:  stringutil.EnsureNotEmpty(member.Nick, "*no nick*"),
			},
			{
				Name:  "ID",
				Value: "```\n" + member.User.ID + "\n```",
			},
			{
				Name: "Guild Joined",
				Value: stringutil.EnsureNotEmpty(joinedTime.Format(time.RFC1123),
					"*failed parsing timestamp*"),
			},
			{
				Name: "Account Created",
				Value: stringutil.EnsureNotEmpty(createdTime.Format(time.RFC1123),
					"*failed parsing timestamp*"),
			},
			{
				Name: "Karma",
				Value: fmt.Sprintf("On this Guild: **%d**\nTotal: **%d**",
					karma, karmaTotal),
			},
			{
				Name:  "Permissions",
				Value: stringutil.EnsureNotEmpty(strings.Join(perms, "\n"), "*no permissions defined*"),
			},
			{
				Name:  "Reports",
				Value: repsOnRecordStr + "\n*Use `rep <user>` to list all reports of this user.*",
			},
			{
				Name:  "Roles",
				Value: stringutil.EnsureNotEmpty(strings.Join(roles, ", "), "*no roles assigned*"),
			},
		},
	}

	if member.User.Bot {
		embed.Description = ":robot:  **This is a bot account**"
	}

	_, err = ctx.GetSession().ChannelMessageSendEmbed(ctx.GetChannel().ID, embed)
	return err
}
