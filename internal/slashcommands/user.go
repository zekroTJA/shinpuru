package slashcommands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/services/config"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
	"github.com/zekroTJA/shinpuru/pkg/stringutil"
	"github.com/zekrotja/dgrs"
	"github.com/zekrotja/ken"
)

type User struct{}

var (
	_ ken.SlashCommand        = (*User)(nil)
	_ permissions.PermCommand = (*User)(nil)
)

func (c *User) Name() string {
	return "user"
}

func (c *User) Description() string {
	return "Get information about a user."
}

func (c *User) Version() string {
	return "1.0.0"
}

func (c *User) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *User) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionUser,
			Name:        "user",
			Description: "The user to be displayed.",
			Required:    true,
		},
	}
}

func (c *User) Domain() string {
	return "sp.chat.profile"
}

func (c *User) SubDomains() []permissions.SubPermission {
	return nil
}

func (c *User) Run(ctx *ken.Ctx) (err error) {
	if err = ctx.Defer(); err != nil {
		return
	}

	st := ctx.Get(static.DiState).(*dgrs.State)
	cfg := ctx.Get(static.DiConfig).(config.Provider)
	db := ctx.Get(static.DiDatabase).(database.Database)
	pmw := ctx.Get(static.DiPermissions).(*permissions.Permissions)

	var user *discordgo.User

	for _, user = range ctx.Event.ApplicationCommandData().Resolved.Users {
		break
	}

	if user == nil {
		user = ctx.Options().GetByName("user").UserValue(ctx)
	}

	member, err := st.Member(ctx.Event.GuildID, user.ID, true)
	if err != nil {
		return
	}

	guild, err := st.Guild(ctx.Event.GuildID, true)
	if err != nil {
		return
	}

	membRoleIDs := make(map[string]struct{})
	for _, rID := range member.Roles {
		membRoleIDs[rID] = struct{}{}
	}

	maxPos := len(guild.Roles)
	roleColor := static.ColorEmbedGray
	for _, guildRole := range guild.Roles {
		if _, ok := membRoleIDs[guildRole.ID]; ok && guildRole.Position < maxPos && guildRole.Color != 0 {
			maxPos = guildRole.Position
			roleColor = guildRole.Color
		}
	}

	createdTime, err := discordutil.GetDiscordSnowflakeCreationTime(member.User.ID)
	if err != nil {
		return err
	}

	perms, _, err := pmw.GetPermissions(ctx.Session, guild.ID, member.User.ID)
	if err != nil {
		return err
	}

	guildReps, err := db.GetReportsFiltered(guild.ID, member.User.ID, -1, 0, 1000)
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

	karma, err := db.GetKarma(member.User.ID, guild.ID)
	if !database.IsErrDatabaseNotFound(err) && err != nil {
		return err
	}

	karmaTotal, err := db.GetKarmaSum(member.User.ID)
	if !database.IsErrDatabaseNotFound(err) && err != nil {
		return err
	}

	embed := &discordgo.MessageEmbed{
		Color: roleColor,
		Title: fmt.Sprintf("Info about member %s#%s", member.User.Username, member.User.Discriminator),
		Description: fmt.Sprintf("[**Here**](%s/guilds/%s/%s) you can find this users profile in the web interface.",
			cfg.Config().WebServer.PublicAddr, guild.ID, member.User.ID),
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
				Value: stringutil.EnsureNotEmpty(member.JoinedAt.Format(time.RFC1123),
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

	return ctx.FollowUpEmbed(embed).Error
}
