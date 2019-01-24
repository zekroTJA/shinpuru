package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/zekroTJA/shinpuru/internal/util"
)

type CmdProfile struct {
	PermLvl int
}

func (c *CmdProfile) GetInvokes() []string {
	return []string{"user", "u", "profile"}
}

func (c *CmdProfile) GetDescription() string {
	return "Get information about a user"
}

func (c *CmdProfile) GetHelp() string {
	return "`profile (<userResolvable>)` - get user info"
}

func (c *CmdProfile) GetGroup() string {
	return GroupEtc
}

func (c *CmdProfile) GetPermission() int {
	return c.PermLvl
}

func (c *CmdProfile) SetPermission(permLvl int) {
	c.PermLvl = permLvl
}

func (c *CmdProfile) Exec(args *CommandArgs) error {
	member, err := args.Session.GuildMember(args.Guild.ID, args.User.ID)
	if err != nil {
		return err
	}
	if len(args.Args) > 0 {
		member, err = util.FetchMember(args.Session, args.Guild.ID, strings.Join(args.Args, " "))
		if err != nil || member == nil {
			msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
				"Could not fetch any member by the passed resolvable.")
			util.DeleteMessageLater(args.Session, msg, 6*time.Second)
			return err
		}
	}

	membRoleIDs := make(map[string]struct{})
	for _, rID := range member.Roles {
		membRoleIDs[rID] = struct{}{}
	}

	maxPos := len(args.Guild.Roles)
	roleColor := util.ColorEmbedGray
	for _, guildRole := range args.Guild.Roles {
		if _, ok := membRoleIDs[guildRole.ID]; ok && guildRole.Position < maxPos && guildRole.Color != 0 {
			maxPos = guildRole.Position
			roleColor = guildRole.Color
		}
	}

	joinedTime, err := member.JoinedAt.Parse()
	if err != nil {
		return err
	}
	createdTime, err := util.GetDiscordSnowflakeCreationTime(member.User.ID)
	if err != nil {
		return err
	}

	var permLvl int
	if member.User.ID == args.CmdHandler.config.Discord.OwnerID {
		permLvl = util.PermLvlBotOwner
	} else if member.User.ID == args.Guild.OwnerID {
		permLvl = util.PermLvlGuildOwner
	} else {
		permLvl, err = args.CmdHandler.db.GetMemberPermissionLevel(args.Session, args.Guild.ID, member.User.ID)
		if err != nil {
			return err
		}
	}

	var repsOnRecord int
	guildReps, err := args.CmdHandler.db.GetReportsGuild(args.Guild.ID)
	if err != nil {
		return err
	}
	for _, grep := range guildReps {
		if grep.VictimID == member.User.ID {
			repsOnRecord++
		}
	}
	repsOnRecordStr := "This user has a white vest :ok_hand:"
	if repsOnRecord > 0 {
		repsOnRecordStr = fmt.Sprintf("This user has **%d** reports on record on this guild.", repsOnRecord)
	}

	roles := make([]string, len(member.Roles))
	for i, rID := range member.Roles {
		roles[i] = "<@&" + rID + ">"
	}

	embed := &discordgo.MessageEmbed{
		Color: roleColor,
		Title: fmt.Sprintf("Info about member %s#%s", member.User.Username, member.User.Discriminator),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: member.User.AvatarURL(""),
		},
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Inline: true,
				Name:   "Tag",
				Value:  member.User.Username + "#" + member.User.Discriminator,
			},
			&discordgo.MessageEmbedField{
				Inline: true,
				Name:   "Nickname",
				Value:  util.EnsureNotEmpty(member.Nick, "*no nick*"),
			},
			&discordgo.MessageEmbedField{
				Name:  "ID",
				Value: "```\n" + member.User.ID + "\n```",
			},
			&discordgo.MessageEmbedField{
				Name: "Guild Joined",
				Value: util.EnsureNotEmpty(joinedTime.Format(time.RFC1123),
					"*failed parsing timestamp*"),
			},
			&discordgo.MessageEmbedField{
				Name: "Account Created",
				Value: util.EnsureNotEmpty(createdTime.Format(time.RFC1123),
					"*failed parsing timestamp*"),
			},
			&discordgo.MessageEmbedField{
				Name:  "Permission Level",
				Value: fmt.Sprintf("`%d`", permLvl),
			},
			&discordgo.MessageEmbedField{
				Name:  "Reports",
				Value: repsOnRecordStr + "\n*Use `rep <user>` to list all reports of this user.*",
			},
			&discordgo.MessageEmbedField{
				Name:  "Roles",
				Value: util.EnsureNotEmpty(strings.Join(roles, ", "), "*no roles assigned*"),
			},
		},
	}

	if member.User.Bot {
		embed.Description = ":robot:  **This is a bot account**"
	}

	_, err = args.Session.ChannelMessageSendEmbed(args.Channel.ID, embed)
	return err
}
